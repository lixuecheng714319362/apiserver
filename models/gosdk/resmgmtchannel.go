package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/common/util"
	"io"
)

func GetCreateChannelReader(channelID string,MSPList []string)(io.Reader,error){
	var err error
	envelope := &common.Envelope{}
	payload := &common.Payload{
		Header:&common.Header{},
	}
	channelHeader := &common.ChannelHeader{}
	configEnvelope := &common.ConfigUpdateEnvelope{}
	configUpdate := &common.ConfigUpdate{}
	if err=proto.Unmarshal(configEnvelope.ConfigUpdate, configUpdate);err!=nil{
		return nil ,err
	}
	//value:=&common.ImplicitMetaPolicy{}

	channeljson:=`{"channel_id":"mychannel","read_set":{"groups":{"Application":{"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}}}},"values":{"Consortium":{}}},"write_set":{"groups":{"Application":{"version":1,"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}},"values":{"Capabilities":{"value":"CggKBFYxXzMSAA==","mod_policy":"Admins"}},"policies":{"Admins":{"policy":{"type":3,"value":"CgZBZG1pbnMQAg=="},"mod_policy":"Admins"},"Readers":{"policy":{"type":3,"value":"CgdSZWFkZXJz"},"mod_policy":"Admins"},"Writers":{"policy":{"type":3,"value":"CgdXcml0ZXJz"},"mod_policy":"Admins"}},"mod_policy":"Admins"}},"values":{"Consortium":{"value":"ChBTYW1wbGVDb25zb3J0aXVt"}}}}`
	channelj:=&common.ConfigUpdate{}
	if err=json.Unmarshal([]byte(channeljson),channelj);err!=nil{
		return nil ,err
	}

	configUpdate=channelj


	if len(MSPList)>0{
		groups:= map[string]*common.ConfigGroup{}
		for _, v := range MSPList {
			groups[v]=&common.ConfigGroup{}
		}
		configUpdate.ReadSet.Groups["Application"].Groups=groups
		configUpdate.WriteSet.Groups["Application"].Groups=groups
	}

	channelheadj:=&common.ChannelHeader{
		Type: 2,
		Timestamp:util.CreateUtcTimestamp(),
		ChannelId:channelID,
	}

	channelHeader=channelheadj

	configUpdate.ChannelId=channelID

	channelHeader.ChannelId=channelID
	channelHeader.Timestamp=util.CreateUtcTimestamp()

	b,_:=json.Marshal(configUpdate)
	fmt.Println(string(b))

	cu2,err:=proto.Marshal(configUpdate)
	if err!=nil{
		return nil ,err
	}
	h2,err:=proto.Marshal(channelHeader)
	if err != nil {
		return nil ,err
	}
	payload.Header.ChannelHeader=h2

	configEnvelope.ConfigUpdate=cu2

	da2,err:=proto.Marshal(configEnvelope)
	if err != nil {
		return nil ,err
	}
	payload.Data=da2


	pa2,err:=proto.Marshal(payload)
	if err != nil {
		return nil ,err
	}

	envelope.Payload=pa2

	en2,err:=proto.Marshal(envelope)

	return  bytes.NewReader(en2),nil


}
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