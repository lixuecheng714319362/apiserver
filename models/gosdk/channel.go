package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type ChannelClient struct {
	ConfigPath string
	UserName   string
	ChannelID  string

	SDK    *fabsdk.FabricSDK
	Client *channel.Client
}

// 创建通道客户端
func GetChannelClient(configPath, userName, channelID string) (*ChannelClient, error) {
	SDK, err := InitializeSDK(configPath)
	if err != nil {
		return nil, err
	}
	client, err := channel.New(SDK.ChannelContext(channelID, fabsdk.WithUser(userName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &ChannelClient{
		configPath,
		userName,
		channelID,
		SDK,
		client,
	}, nil
}

func (ChannelClient *ChannelClient) Query(chainCodeID, Fcn string, args [][]byte) (channel.Response, error) {
	//TODO
	//TransientMap、InvocationChain功能
	return ChannelClient.Client.Query(
		channel.Request{
			ChaincodeID: chainCodeID,
			Fcn:         Fcn,
			Args:        args,
		},
	)
}

func (ChannelClient *ChannelClient) Invoke(chainCodeID, Fcn string, args [][]byte) (channel.Response, error) {
	//TODO
	//TransientMap、InvocationChain功能
	return ChannelClient.Client.Execute(
		channel.Request{
			ChaincodeID: chainCodeID,
			Fcn:         Fcn,
			Args:        args,
		},
		//channel.WithTargetEndpoints(targetPeers...),
	)
}

func (ChannelClient *ChannelClient) RegisterChaincodeEven(chainCodeID string) (fab.Registration, <-chan *fab.CCEvent, error) {
	//TODO-
	//eventFilter正则表达式匹配所有事件
	return ChannelClient.Client.RegisterChaincodeEvent(chainCodeID, ".*")
}
func (ChannelClient *ChannelClient) UnregisterChaincodeEvent(Registration fab.Registration) {

	ChannelClient.Client.UnregisterChaincodeEvent(Registration)
}

// 关闭SDK
func (ChannelClient *ChannelClient) CloseSDK() {
	ChannelClient.SDK.Close()
}
