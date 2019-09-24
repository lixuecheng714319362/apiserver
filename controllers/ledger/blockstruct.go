package ledger

import (
	"apiserver/models/gosdk/tool/blockdata"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

type Block struct {
	Number            uint64
	CurrentBlockHash  []byte
	PreviousHash      []byte
	DataHash          []byte
	TransactionNumber int
	Transactions      []*Transaction
}
type Transaction struct {
	CreatorMSPID   string
	CreateID  string
	Type      string
	Timestamp int64
	Nanos     int32
	ChannelId string
	TxId      string
	Actions   []*Action
}
type Action struct {
	CCID          string
	TxArgs        string
	NsRwSets      string
	ReponseStatus int32
}

func Getinfo(thisBlock *common.Block) (*Block, error) {
	b, err := blockdata.Getinfo(thisBlock)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Number:            b.Header.Number,
		CurrentBlockHash:  b.Header.CurrentBlockHash,
		PreviousHash:      b.Header.PreviousHash,
		DataHash:          b.Header.DataHash,
		TransactionNumber: b.TransactionNumber,
	}



	for _, v := range b.Transaction {
		t := &Transaction{
			CreatorMSPID:v.Header.SignatureHeader.Creator.MSPID,
			CreateID: LoadCertBytes(v.Header.SignatureHeader.Creator.IdBytes),
			Type:      v.Header.ChannelHeader.Type,
			Timestamp: v.Header.ChannelHeader.Timestamp,
			Nanos:     v.Header.ChannelHeader.Nanos,
			TxId:      v.Header.ChannelHeader.TxId,
			ChannelId: v.Header.ChannelHeader.ChannelId,
		}

		if v.Header.ChannelHeader.Type == "CONFIG" {
			block.Transactions = append(block.Transactions, t)
			continue
		}
		for _, act := range v.Actions {
			jtx, err := json.Marshal(act.CCProposalPayload.TxArgs)
			if err != nil {
				fmt.Println("tx args json marshal err", err)
				continue
			}
			jrw, err := json.Marshal(act.CCResponsePayload.NsRwSets)
			if err != nil {
				fmt.Println("tx args json marshal err", err)
				continue
			}
			a := &Action{
				CCID:          act.CCProposalPayload.CCID,
				TxArgs:        string(jtx),
				NsRwSets:      string(jrw),
				ReponseStatus: act.CCResponsePayload.ReponseStatus,
			}
			t.Actions = append(t.Actions, a)
		}

		block.Transactions = append(block.Transactions, t)
	}

	return block, nil

}
func LoadCertBytes(original []byte) (string) {
	certDerBlock,_:=pem.Decode(original)
	if certDerBlock ==nil{
		return ""
	}
	cert,err:=x509.ParseCertificate(certDerBlock.Bytes)
	if err != nil {
		return ""
	}
	return cert.Subject.CommonName
}