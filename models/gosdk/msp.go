package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type MSPClient struct {
	ConfigPath string
	UserName   string
	OrgName    string

	SDK    *fabsdk.FabricSDK
	Client *msp.Client
}

// 创建账本客户端
func GetMSPClient(configPath, userName, orgName string) (*MSPClient, error) {
	SDK, err := InitializeSDK(configPath)
	if err != nil {
		return nil, err
	}

	mspClient, err := msp.New(SDK.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &MSPClient{
		configPath,
		userName,
		orgName,
		SDK,
		mspClient,
	}, nil
}

// 关闭SDK
func (MSPClient *MSPClient) CloseSDK() {
	MSPClient.SDK.Close()
}
