package gosdk

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var SinglePeerModel string
var PeerConfig SinglePeerConfig

type SinglePeerConfig struct {
	ConfigPath string
	UserName   string
	OrgName    string
	TargetPeers string
	TargetOrderer string
}

func init() {
	SinglePeerModel = beego.AppConfig.String("singlepeer")
	if SinglePeerModel != "true" {
		return
	}
	PeerConfig = SinglePeerConfig{
		ConfigPath: beego.AppConfig.String("configpath"),
		UserName:   beego.AppConfig.String("username"),
		OrgName:    beego.AppConfig.String("orgname"),
		TargetPeers: beego.AppConfig.String("targetpeers"),
		TargetOrderer: beego.AppConfig.String("targetorderer"),
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

func ChangeResmgmtRequetSingleConfig(req *ResmgmtRequest)  {
	if SinglePeerModel=="true"{
		if PeerConfig.TargetPeers != "" {
			targetPeers := make([]string, 0)
			err := json.Unmarshal([]byte(PeerConfig.TargetPeers), &targetPeers)
			if err == nil && len(targetPeers) > 0 {
				req.TargetPeers = targetPeers
			}
		}
		req.TargetOrderer=PeerConfig.TargetOrderer
		req.OrgName=PeerConfig.OrgName
		req.ConfigPath=PeerConfig.ConfigPath
		req.UserName=PeerConfig.UserName

	}
	beego.Debug("peerConfig :",PeerConfig)
	return
}

func ChangeLedgerRequestSingleConfig(req *LedgerRequest) {
	if SinglePeerModel == "true" {
		req.OrgName = PeerConfig.OrgName
		req.ConfigPath = PeerConfig.ConfigPath
		req.UserName = PeerConfig.UserName
	}
	return
}

func ChangeChannelRequestSingleConfig(req *ChannelRequest) {
	if SinglePeerModel == "true" {
		if PeerConfig.TargetPeers != "" {
			targetPeers := make([]string, 0)
			err := json.Unmarshal([]byte(PeerConfig.TargetPeers), &targetPeers)
			if err == nil && len(targetPeers) > 0 {
				req.TargetPeers = targetPeers
			} else {
				req.TargetPeers = nil
			}
		}
		req.ConfigPath = PeerConfig.ConfigPath
		req.UserName = PeerConfig.UserName
	}
	return
}
