package utils

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/protos/common"
	"io"
)

//signer
func GetCreateChannelSinger(signerMap map[string]string,sdk *fabsdk.FabricSDK) (signerList []msp.SigningIdentity,err error) {
	for k, v := range signerMap {
		singerClient,err:=mspclient.New(sdk.Context(),mspclient.WithOrg(k))
		if err != nil {
			return nil ,err
		}
		si,err:=singerClient.GetSigningIdentity(v)
		if err != nil {
			return nil,err
		}
		signerList=append(signerList,si)
	}
	return signerList,err
}


func CreateConfigUpdateEnvelope(channelID string ,configUpdate *common.ConfigUpdate)(io.Reader, error){

	configUpdate.ChannelId=channelID
	configUpdateData,err:=proto.Marshal(configUpdate)
	if err!=nil{
		return nil ,err
	}

	configUpdateEnvelope := &common.ConfigUpdateEnvelope{}
	configUpdateEnvelope.ConfigUpdate= configUpdateData

	configUpdateEnvelopeData,err:=proto.Marshal(configUpdateEnvelope)
	if err != nil {
		return nil ,err
	}

	channelHeader:=&common.ChannelHeader{
		Type: 2,
		Timestamp:util.CreateUtcTimestamp(),
		ChannelId:channelID,
	}
	channelHeaderData,err:=proto.Marshal(channelHeader)
	if err != nil {
		return nil ,err
	}

	payload := &common.Payload{
		Header:&common.Header{},
	}
	payload.Header.ChannelHeader=channelHeaderData
	payload.Data= configUpdateEnvelopeData

	payloadData,err:=proto.Marshal(payload)
	if err != nil {
		return nil ,err
	}

	envelope := &common.Envelope{}
	envelope.Payload= payloadData

	envelopeData,err:=proto.Marshal(envelope)
	if err != nil {
		return nil ,err
	}

	return   bytes.NewReader(envelopeData),nil
}