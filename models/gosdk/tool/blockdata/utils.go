package blockdata

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	common2 "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
)

func GetProposal(P *fab.TransactionProposal)(*TransactionProposal,error){
	h,err:=utils.GetHeader(P.Header)
	if err != nil {
		return nil,err
	}
	header,err:=getHander(h)
	if err != nil {
		return nil,err
	}
	txProposal:=&TransactionProposal{
		Header:header,
	}

	txProposal.TxnID=string(P.TxnID)
	//Payload 内容
	p,err:=utils.GetChaincodeProposalPayload(P.Payload)
	if err != nil {
		return nil,err
	}
	transientMap:=make(map[string]string)
	for k, v := range p.TransientMap {
		transientMap[k]=string(v)
	}
	txProposal.TransientMap=transientMap

	invokeSpec,err:=utils.GetChaincodeInvocationSpec((*peer.Proposal)(P.Proposal))
	if err != nil {
		return nil,err
	}
	txProposal.ChainCodeSpec=invokeSpec.ChaincodeSpec
	return txProposal,nil
}

func GetProposalResponse(pr *fab.TransactionProposalResponse)(*TransactionProposalResponse,error)  {

	pl,err:=utils.GetProposalResponsePayload(pr.Payload)
	if err != nil {
		return nil,err
	}
	ca,err:=utils.GetChaincodeAction(pl.Extension)
	if err != nil {
		return nil,err
	}
	resault := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(ca.Results, resault)
	if err != nil {
		return nil, err
	}
	var rs Results
	for _, kvrw := range resault.NsRwset {
		nsRwSet := &NsRwSets{}
		nsRwSet.NameSpace = kvrw.Namespace
		kv := &kvrwset.KVRWSet{}
		err = proto.Unmarshal(kvrw.Rwset, kv)
		if err != nil {
			return nil, err
		}
		nsRwSet.Reads = kv.Reads
		writes:=[]*KVWrite{}
		for _, v := range kv.Writes {
			write:=&KVWrite{v.Key,v.IsDelete,string(v.Value)}
			writes=append(writes,write)
		}
		nsRwSet.Writes = writes
		nsRwSet.MetadataWrites = kv.MetadataWrites
		nsRwSet.RangeQueriesInfo = kv.RangeQueriesInfo
		rs.NsRwSets = append(rs.NsRwSets, nsRwSet)
	}

	env,err:=utils.GetChaincodeEvents(ca.Events)
	if err != nil {
		return nil,err
	}





	tpr:=&TransactionProposalResponse{
		Endorser:pr.Endorser,
		Status:pr.Status,
		ChaincodeStatus:pr.ChaincodeStatus,
		Version:pr.Version,


		Response:&Response{
			Status:pr.Response.Status,
			Message:pr.Response.Message,
			Payload:pr.Response.Payload,
		},
		Payload:&ProposalResponsePayload{
			ProposalHash:pl.ProposalHash,
			Extension:&ChaincodeAction{
				Results:rs,
				Events:env,
				Response:&Response{
					Status:ca.Response.Status,
					Message:ca.Response.Message,
					Payload:ca.Response.Payload,
				},
			},
		},
		Endorsement:(*peer.Endorsement)(pr.Endorsement),
	}
	if pr.Timestamp!=nil{
	tpr.Timestamp=pr.Timestamp.Seconds
	tpr.Nanos=pr.Timestamp.Nanos
	}

	return tpr,nil
}










func  GetCreater(b []byte) *Creater  {
	var creater *Creater
	if b != nil {
		serial := msp.SerializedIdentity{}
		err := proto.Unmarshal(b, &serial)
		if err != nil {
			creater = nil
		} else {
			creater = &Creater{
				MSPID:   serial.Mspid,
				IdBytes: serial.IdBytes,
			}
		}
	}
	return creater
}


func getHander(cheader *common2.Header) (*Header, error) {

	ch,err:=UnmarshalChannelHeader(cheader.ChannelHeader)
	if err != nil {
		return nil,err
	}
	sh,err:=UnmarshalSignatureHeader(cheader.SignatureHeader)
	if err != nil {
		return nil ,err
	}
	return &Header{ChannelHeader: ch, SignatureHeader: sh}, err
}

func UnmarshalChannelHeader(b []byte) (*ChannelHeader,error) {
	channelHeader,err:=utils.UnmarshalChannelHeader(b)
	if err != nil {
		return nil, err
	}
	//err = proto.Unmarshal(cheader.SignatureHeader, sig)

	ChannelHeader := &ChannelHeader{
		Type:    common.HeaderType(channelHeader.Type).String(),
		Version: channelHeader.Version,

		ChannelId:   channelHeader.ChannelId,
		TxId:        channelHeader.TxId,
		Epoch:       channelHeader.Epoch,
		Extension:   channelHeader.Extension,
		TlsCertHash: channelHeader.TlsCertHash,
	}

	if channelHeader.Timestamp != nil {
		ChannelHeader.Timestamp = channelHeader.Timestamp.Seconds
		ChannelHeader.Nanos = channelHeader.Timestamp.Nanos
	}
	return ChannelHeader,nil
}

func UnmarshalSignatureHeader(b []byte)(sh  *SignatureHeader,err error){
	sh=new(SignatureHeader)
	sig,err:=utils.GetSignatureHeader(b)
	if err != nil {
		return nil, err
	}
	creater:=new(Creater)
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
		}
	}
	sh.Creator=creater
	sh.Nonce=sig.Nonce
	return
}

