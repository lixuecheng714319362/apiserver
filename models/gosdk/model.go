package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 初始化SDK
func InitializeSDK(configFile string) (*fabsdk.FabricSDK, error) {
	// 读取配置文件
	if len(configFile) < 500 {
		file := config.FromFile(configFile)
		return fabsdk.New(file)
	}
	return fabsdk.New(config.FromRaw([]byte(configFile), "yaml"))
}
