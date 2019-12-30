package blockdata

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func GetChannelResposeInfo(res *channel.Response)(*ChannelResponse,error){
	var responseData ChannelResponse
	responseData.TransactionID=string(res.TransactionID)
	responseData.TxValidationCode=int32(res.TxValidationCode)
	responseData.Payload=string(res.Payload)
	responseData.ChaincodeStatus= res.ChaincodeStatus
	//proposal 内容
	tp,err:=GetProposal(res.Proposal)
	if err != nil {
		return nil,err
	}
	responseData.Proposal=tp
	txPr:=make([]*TransactionProposalResponse,0)
	for _, v := range res.Responses {
		pr,err:=GetProposalResponse(v)
		if err != nil {
			panic(err)
		}
		txPr=append(txPr,pr)
	}
	responseData.Responses=txPr
	return &responseData,nil
}
