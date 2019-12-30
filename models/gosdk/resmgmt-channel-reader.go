package gosdk

import (
	"apiserver/models/gosdk/tool/utils"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/tools/configtxlator/update"
	"github.com/hyperledger/fabric/protos/common"
	"io"
)

func GetCreateChannelReader(channelID string, OrgNameList []string)(io.Reader,error){
	var err error
	configUpdate := &common.ConfigUpdate{}
	//value:=&common.ImplicitMetaPolicy{}
	channeljson:=`{"channel_id":"mychannel","read_set":{"groups":{"Application":{"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}}}},"values":{"Consortium":{}}},"write_set":{"groups":{"Application":{"version":1,"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}},"values":{"Capabilities":{"value":"CggKBFYxXzMSAA==","mod_policy":"Admins"}},"policies":{"Admins":{"policy":{"type":3,"value":"CgZBZG1pbnMQAg=="},"mod_policy":"Admins"},"Readers":{"policy":{"type":3,"value":"CgdSZWFkZXJz"},"mod_policy":"Admins"},"Writers":{"policy":{"type":3,"value":"CgdXcml0ZXJz"},"mod_policy":"Admins"}},"mod_policy":"Admins"}},"values":{"Consortium":{"value":"ChBTYW1wbGVDb25zb3J0aXVt"}}}}`
	if err=json.Unmarshal([]byte(channeljson),configUpdate);err!=nil{
		return nil ,err
	}
	if len(OrgNameList)>0{
		groups:= map[string]*common.ConfigGroup{}
		for _, v := range OrgNameList {
			groups[v]=&common.ConfigGroup{}
		}
		configUpdate.ReadSet.Groups["Application"].Groups=groups
		configUpdate.WriteSet.Groups["Application"].Groups=groups
	}else{
		return nil,errors.New("org name list is nil" )
	}
	return utils.CreateConfigUpdateEnvelope(channelID,configUpdate)
}


//通过configtxlator update模块生成
func GetAddOrgChannelConfigUpdate(request *ResmgmtRequest)(io.Reader,error){
	config,err:=GetOldChannelConfig(request.ConfigPath,request.OrgName,request.UserName,request.ChannelID,request.TargetPeers)
	if err != nil {
		return nil ,err
	}
	configByte,err:=proto.Marshal(config)
	if err != nil {
		return nil,err
	}
	newConfig:=&common.Config{}
	oldConfig:=&common.Config{}
	err=proto.Unmarshal(configByte,newConfig)
	if err != nil {
		return nil,err
	}
	err=proto.Unmarshal(configByte,oldConfig)
	if err != nil {
		return nil,err
	}
	for _, v := range request.AddOrgConfig {
		newOrg,err:= utils.NewOrgGroup(&v)
		if err != nil {
			return nil,err
		}
		newConfig.ChannelGroup.Groups["Application"].Groups[v.Name]=newOrg
	}
	configUpdate,err:=update.Compute(oldConfig,newConfig)
	if err != nil {
		return nil,err
	}
	configUpdate.ChannelId=request.ChannelID
	return  utils.CreateConfigUpdateEnvelope(configUpdate.ChannelId,configUpdate)

}