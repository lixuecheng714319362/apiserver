package gosdk

import (
	"github.com/astaxie/beego"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var SinglePeerModel string
var PeerConfig SinglePeerConfig
type SinglePeerConfig struct {
	ConfigPath string
	UserName string
	OrgName string
}

func init(){
	SinglePeerModel=beego.AppConfig.String("singlepeer")
	if SinglePeerModel!="true"{
		return
	}
	PeerConfig=SinglePeerConfig{
		ConfigPath:beego.AppConfig.String("configpath"),
		UserName:beego.AppConfig.String("username"),
		OrgName:beego.AppConfig.String("orgname"),
	}

}

// 初始化SDK
func InitializeSDK(configFile string) (*fabsdk.FabricSDK, error) {
	// 读取配置文件
	if len(configFile) < 500 {
		file := config.FromFile(configFile)
		return fabsdk.New(file)
	}
	return fabsdk.New(config.FromRaw([]byte(configFile), "yaml"))
}
