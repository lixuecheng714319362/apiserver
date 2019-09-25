package blockdata

import (
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
)

type Chain struct {
	Height int64
}

type Block struct {
	Header            *BlockHead
	TransactionNumber int
	Transaction       []*Transaction
	TransactionsFilter		[]uint8
}

type BlockHead struct {
	Number           uint64
	CurrentBlockHash []byte
	PreviousHash     []byte
	DataHash         []byte
}

type Transaction struct {
	Header *Header

	Actions []*Action
}
type Action struct {
	Header            *Header
	CCProposalPayload *CCProposalPayload
	CCResponsePayload *CCResponsePayload
	TransientMap      map[string][]byte
}
type CCResponsePayload struct {
	ProposalHash []byte

	NsRwSets []*NsRwSets

	ReponseStatus   int32
	ReponseMessage  string
	ReponsePlayload []byte
}
type NsRwSets struct {
	NameSpace        string
	Reads            []*kvrwset.KVRead
	Writes           []*kvrwset.KVWrite
	RangeQueriesInfo []*kvrwset.RangeQueryInfo
	MetadataWrites   []*kvrwset.KVMetadataWrite
}

type CCProposalPayload struct {
	CCtype      string
	CCPath      string
	CCID        string
	CCVersion   string
	TxArgs      []string
	Decorations map[string][]byte
	Timeout     int32
	Method      string
}
type Header struct {
	ChannelHeader   *ChannelHeader
	SignatureHeader *SignatureHeader
}
type SignatureHeader struct {
	Creator *Creater
	Nonce   []byte
}

type Creater struct {
	MSPID   string
	IdBytes []byte
}

type ChannelHeader struct {
	Type        string
	Version     int32
	Timestamp   int64
	Nanos       int32
	ChannelId   string
	TxId        string
	Epoch       uint64
	Extension   []byte
	TlsCertHash []byte
}

type ChainTxEvents struct {
	TxID, Chaincode, Name string
	Status                int
	Payload               []byte
}

type ChainBlock struct {
	Height       int64 `json:",string"`
	Hash         string
	TimeStamp    string
	Transactions []*Transaction
	TxEvents     []*ChainTxEvents
}

func Getinfo(thisBlock *common.Block) (*Block, error) {

	var txFilter []uint8
	if thisBlock.Metadata!=nil && len(thisBlock.Metadata.Metadata)>=2{
		txFilter=thisBlock.Metadata.Metadata[2]
	}
	newBlock := &Block{
		Header: &BlockHead{
			Number:       thisBlock.Header.Number,
			DataHash:     thisBlock.Header.DataHash,
			PreviousHash: thisBlock.Header.PreviousHash,

		},
		TransactionNumber: len(thisBlock.Data.Data),
		TransactionsFilter: txFilter,
	}

	//此处应该遍历block.Data.Data
	transaction := make([]*Transaction, 0)

	for _, data := range thisBlock.Data.Data {

		env, err := utils.GetEnvelopeFromBlock(data)
		if err != nil {
			return nil, err
		}

		chainTransaction, err := EnvelopeToTrasaction((*common.Envelope)(env))
		if err != nil {
			return nil, err
		}
		//if chainTransaction.Header.ChannelHeader.TxId != "" {
		transaction = append(transaction, chainTransaction)
		//}
	}

	newBlock.Transaction = transaction
	return newBlock, nil

}

func EnvelopeToTrasaction(env *common.Envelope) (*Transaction, error) {
	transaction := &Transaction{}

	var err error
	if env == nil {
		return nil, errors.New("<-common.Envelope is nil")
	}
	payl := &common.Payload{}
	err = proto.Unmarshal(env.Payload, payl)
	if err != nil {
		return nil, err
	}
	if payl.Header == nil {
		return nil, errors.New("<-  payl head nil")
	}
	header, err := getHander(payl.Header)
	if err != nil {
		return nil, err
	}
	transaction.Header = header

	tx := &pb.Transaction{}
	err = proto.Unmarshal(payl.Data, tx)
	if err != nil {
		return nil, err
	}

	for _, v := range tx.Actions {
		var action = &Action{
			CCProposalPayload: &CCProposalPayload{},
			CCResponsePayload: &CCResponsePayload{},
		}

		actionheader := &common.Header{}
		err = proto.Unmarshal(v.Header, actionheader)
		if err != nil {
			return nil, err
		}
		//header,err:=getHander(actionheader)
		//if err != nil {
		//	return nil,err
		//}
		//action.Header=header

		chaincodeActionPayload := &pb.ChaincodeActionPayload{}
		err = proto.Unmarshal(v.Payload, chaincodeActionPayload)
		if err != nil {
			return nil, err
		}

		if transaction.Header.ChannelHeader.Type == "CONFIG" {
			//block 0  区块进行下面解析会报错
			//提前退出
			return transaction, nil
		}

		if chaincodeActionPayload.Action != nil {
			proposalResponsePayload := &pb.ProposalResponsePayload{}
			err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
			if err != nil {
				//block 0  区块会报错
				return nil, err
			}
			action.CCResponsePayload.ProposalHash = proposalResponsePayload.ProposalHash

			repextension := &pb.ChaincodeAction{}
			err = proto.Unmarshal(proposalResponsePayload.Extension, repextension)
			if err != nil {
				return nil, err
			}

			if repextension.Response != nil {
				action.CCResponsePayload.ReponseStatus = repextension.Response.Status
				action.CCResponsePayload.ReponseMessage = repextension.Response.Message
				action.CCResponsePayload.ReponsePlayload = repextension.Response.Payload
			}

			resault := &rwset.TxReadWriteSet{}
			err = proto.Unmarshal(repextension.Results, resault)
			if err != nil {
				return nil, err
			}

			for _, kvrw := range resault.NsRwset {
				nsRwSet := &NsRwSets{}
				nsRwSet.NameSpace = kvrw.Namespace
				kv := &kvrwset.KVRWSet{}
				err = proto.Unmarshal(kvrw.Rwset, kv)
				if err != nil {
					return nil, err
				}
				nsRwSet.Reads = kv.Reads
				nsRwSet.Writes = kv.Writes
				nsRwSet.MetadataWrites = kv.MetadataWrites
				nsRwSet.RangeQueriesInfo = kv.RangeQueriesInfo
				action.CCResponsePayload.NsRwSets = append(action.CCResponsePayload.NsRwSets, nsRwSet)
				//for _, v := range kv.Writes {
				//	fmt.Println("__________________",string(v.Value))
				//}
			}

		}

		if chaincodeActionPayload.ChaincodeProposalPayload != nil {
			chaincodeProposalPayload := &pb.ChaincodeProposalPayload{}
			err = proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, chaincodeProposalPayload)
			if err != nil {
				return nil, err
			}
			invocation := &pb.ChaincodeInvocationSpec{}
			err = proto.Unmarshal(chaincodeProposalPayload.Input, invocation)
			if err != nil {
				return nil, err
			}
			spec := invocation.ChaincodeSpec

			if spec != nil {
				action.CCProposalPayload.CCtype = spec.CCType()
				if spec.ChaincodeId != nil {
					action.CCProposalPayload.CCPath = spec.ChaincodeId.Path
					action.CCProposalPayload.CCID = spec.ChaincodeId.Name
					action.CCProposalPayload.CCVersion = spec.ChaincodeId.Version
				}
				action.CCProposalPayload.Timeout = spec.Timeout
				if spec.Input != nil {
					if len(spec.GetInput().GetArgs()) != 0 {
						for _, v := range spec.GetInput().GetArgs() {
							action.CCProposalPayload.TxArgs = append(action.CCProposalPayload.TxArgs, string(v))
						}
						action.CCProposalPayload.Method = string(spec.GetInput().GetArgs()[0])
					}
					action.CCProposalPayload.Decorations = spec.Input.Decorations
				}

			}
			action.TransientMap = chaincodeProposalPayload.TransientMap

			transaction.Actions = append(transaction.Actions, action)
		}
	}

	return transaction, nil
}

func getHander(cheader *common.Header) (*Header, error) {
	header := &Header{
		ChannelHeader:   &ChannelHeader{},
		SignatureHeader: &SignatureHeader{},
	}
	channelHeader := &common.ChannelHeader{}
	sig := &common.SignatureHeader{}
	err := proto.Unmarshal(cheader.ChannelHeader, channelHeader)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(cheader.SignatureHeader, sig)
	if err != nil {
		return nil, err
	}
	header.ChannelHeader = &ChannelHeader{
		Type:    common.HeaderType(channelHeader.Type).String(),
		Version: channelHeader.Version,

		ChannelId:   channelHeader.ChannelId,
		TxId:        channelHeader.TxId,
		Epoch:       channelHeader.Epoch,
		Extension:   channelHeader.Extension,
		TlsCertHash: channelHeader.TlsCertHash,
	}

	if channelHeader.Timestamp != nil {
		header.ChannelHeader.Timestamp = channelHeader.Timestamp.Seconds
		header.ChannelHeader.Nanos = channelHeader.Timestamp.Nanos
	}

	var creater *Creater
	if sig.Creator != nil {
		serial := msp.SerializedIdentity{}
		err := proto.Unmarshal(sig.Creator, &serial)
		if err != nil {
			creater = nil
		} else {
			creater = &Creater{
				MSPID:   serial.Mspid,
				IdBytes: serial.IdBytes,
			}
			//
			//fmt.Println(serial.Mspid, "\n", string(serial.IdBytes))

		}
	}
	header.SignatureHeader = &SignatureHeader{
		Creator: creater,
		Nonce:   sig.Nonce,
	}
	return header, err
}

// blockToChainCodeEvents parses block events for chaincode events associated with individual transactions
func BlockToChainCodeEvents(block *common.Block) []*pb.ChaincodeEvent {
	if block == nil || block.Data == nil || block.Data.Data == nil || len(block.Data.Data) == 0 {
		return nil
	}
	events := make([]*pb.ChaincodeEvent, 0)
	//此处应该遍历block.Data.Data？
	for _, data := range block.Data.Data {
		event, err := GetChainCodeEventsByByte(data)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	return events
}

func GetChainCodeEventsByByte(data []byte) (*pb.ChaincodeEvent, error) {
	// env := &common.Envelope{}
	// if err := proto.Unmarshal(data, env); err != nil {
	// 	return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	// }

	env, err := utils.GetEnvelopeFromBlock(data)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	}
	// get the payload from the envelope
	payload, err := utils.GetPayload(env)
	if err != nil {
		return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
	}

	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not extract channel header from envelope, err %s", err)
	}

	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

		tx, err := utils.GetTransaction(payload.Data)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction payload for block event: %s", err)
		}
		//此处应该遍历tx.Actions？
		chaincodeActionPayload, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction action payload for block event: %s", err)
		}
		propRespPayload, err := utils.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling proposal response payload for block event: %s", err)
		}

		caPayload, err := utils.GetChaincodeAction(propRespPayload.Extension)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling chaincode action for block event: %s", err)
		}
		ccEvent, err := utils.GetChaincodeEvents(caPayload.Events)
		if ccEvent != nil {
			return ccEvent, nil
		}

	}
	return nil, errors.New("no HeaderType_ENDORSER_TRANSACTION type ")
}

func EventConvert(event *pb.ChaincodeEvent) *ChainTxEvents {
	if event == nil {
		return nil
	}
	clientEvent := &ChainTxEvents{}
	clientEvent.Chaincode = event.ChaincodeId
	clientEvent.Name = event.EventName
	clientEvent.Payload = event.Payload
	clientEvent.TxID = event.TxId
	return clientEvent
}
