package gosdk

import (
	"apiserver/models/gosdk/tool/utils"
	"encoding/json"
	"errors"
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
