package blockdata

import (
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric/protos/peer"
)

//block 结构
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
	Header            *SignatureHeader
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
	Writes           []*KVWrite
	RangeQueriesInfo []*kvrwset.RangeQueryInfo
	MetadataWrites   []*kvrwset.KVMetadataWrite
}

type KVWrite struct {
	Key string
	IsDelete bool
	Value string
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
	//随机值用来防止重放攻击
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


//背书节点返回交易响应
type ChannelResponse struct {
	Proposal         *TransactionProposal
	Responses        []*TransactionProposalResponse
	TransactionID string
	TxValidationCode int32
	ChaincodeStatus int32
	Payload string

}
type TransactionProposal struct {
	TxnID string
	Header *Header
	TxArgs []string
	TransientMap map[string]string
	ChainCodeSpec *peer.ChaincodeSpec
}
type TransactionProposalResponse struct {
	Endorser string
	Status int32
	ChaincodeStatus int32

	Version int32
	Timestamp   int64
	Nanos       int32
	Response *Response
	Payload  *ProposalResponsePayload
	Endorsement  *peer.Endorsement


}

type ProposalResponsePayload struct{
	ProposalHash []byte
	Extension  *ChaincodeAction
}
type ChaincodeAction struct {

	Results Results
	Events  *peer.ChaincodeEvent
	Response *Response
}



type Response struct {
	Status int32
	Message string
	Payload []byte
}



type Payload struct{
	Results *Results

}
type Results struct {
	NsRwSets []*NsRwSets
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


