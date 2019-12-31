package gosdk

import (
	"apiserver/models/gosdk/tool/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/tools/configtxlator/update"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
	"io"
)

const nilChannel =`{"channel_id":"mychannel","read_set":{"groups":{"Application":{"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}}}},"values":{"Consortium":{}}},"write_set":{"groups":{"Application":{"version":1,"groups":{"ShuQinOrgOne":{},"ShuQinOrgTwo":{}},"values":{"Capabilities":{"value":"CggKBFYxXzMSAA==","mod_policy":"Admins"}},"policies":{"Admins":{"policy":{"type":3,"value":"CgZBZG1pbnMQAg=="},"mod_policy":"Admins"},"Readers":{"policy":{"type":3,"value":"CgdSZWFkZXJz"},"mod_policy":"Admins"},"Writers":{"policy":{"type":3,"value":"CgdXcml0ZXJz"},"mod_policy":"Admins"}},"mod_policy":"Admins"}},"values":{"Consortium":{"value":"ChBTYW1wbGVDb25zb3J0aXVt"}}}}`

func GetCreateChannelReader(channelID string, OrgNameList []string)(io.Reader,error){
	var err error
	configUpdate := &common.ConfigUpdate{}
	//value:=&common.ImplicitMetaPolicy{}
		if err=json.Unmarshal([]byte(nilChannel),configUpdate);err!=nil{
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
func GetAddOrgChannelConfigUpdateReader(request *ResmgmtRequest)(io.Reader,error){
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



func GetDeleteOrgChannelConfigUpdateReader(request *ResmgmtRequest)(io.Reader,error){
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
	for _, v := range request.DeleteOrgList {
		delete(newConfig.ChannelGroup.Groups["Application"].Groups,v)
	}
	configUpdate,err:=update.Compute(oldConfig,newConfig)
	if err != nil {
		return nil,err
	}
	configUpdate.ChannelId=request.ChannelID
	return  utils.CreateConfigUpdateEnvelope(configUpdate.ChannelId,configUpdate)
}




//通过configtxlator update模块生成
func GetRevokeReader(request *ResmgmtRequest)(io.Reader,error){
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
	var crl  [][]byte
	for _,e := range request.CRL {
		c,err:=base64.StdEncoding.DecodeString(e)
		if err != nil {
			return nil,err
		}
		crl=append(crl,c)
	}
	for _, v := range newConfig.ChannelGroup.Groups["Application"].Groups {
		val:=v.Values["MSP"].Value
		mspCfg := &msp.MSPConfig{}
		err = proto.Unmarshal(val, mspCfg)
		if err != nil {
			return nil,err
		}
		fabMspCfg := &msp.FabricMSPConfig{}
		err = proto.Unmarshal(mspCfg.Config, fabMspCfg)
		if err != nil {
			return nil,err
		}
		if len(request.CRL)<1{
			fabMspCfg.RevocationList=nil
		}else {
			fabMspCfg.RevocationList=append(fabMspCfg.RevocationList,crl...)
		}
		fabMspBytes, err := proto.Marshal(fabMspCfg)
		if err != nil {
			return nil,err
		}
		mspCfg.Config = fabMspBytes
		mspBytes, err := proto.Marshal(mspCfg)
		if err != nil {
			return nil,err
		}
		v.Values["MSP"].Value=mspBytes
	}

	configUpdate,err:=update.Compute(oldConfig,newConfig)
	if err != nil {
		return nil,err
	}
	//从通道配置中获取的channelID为空！需赋值。
	configUpdate.ChannelId=request.ChannelID
	return  utils.CreateConfigUpdateEnvelope(configUpdate.ChannelId,configUpdate)

}







