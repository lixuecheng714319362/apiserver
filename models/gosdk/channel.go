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

type ChannelRequest struct {
	ConfigPath  string
	UserName    string

	ChannelID   string
	CCID        string
	Fcn         string //OrdID
	Args        []string
	TargetPeers []string
	TransientMap map[string]string

}

// 创建通道客户端
func GetChannelClient(channelRequest *ChannelRequest) (*ChannelClient, error) {
	SDK, err := InitializeSDK(channelRequest.ConfigPath)
	if err != nil {
		return nil, err
	}
	client, err := channel.New(SDK.ChannelContext(channelRequest.ChannelID, fabsdk.WithUser(channelRequest.UserName),))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &ChannelClient{
		channelRequest.ConfigPath,
		channelRequest.UserName,
		channelRequest.ChannelID,
		SDK,
		client,
	}, nil
}

func (ChannelClient *ChannelClient) Query(channelRequest *ChannelRequest, args [][]byte,transientMap map[string][]byte) (channel.Response, error) {
	//TODO
	//TransientMap、InvocationChain功能
	return ChannelClient.Client.Query(
		channel.Request{
			ChaincodeID: channelRequest.CCID,
			Fcn:         channelRequest.Fcn,
			Args:        args,
			TransientMap: transientMap,
		},
		channel.WithTargetEndpoints(channelRequest.TargetPeers...),
	)
}

func (ChannelClient *ChannelClient) Invoke(channelRequest *ChannelRequest, args [][]byte,transientMap map[string][]byte) (channel.Response, error) {
	//TODO
	//TransientMap、InvocationChain功能
	return ChannelClient.Client.Execute(
		channel.Request{
			ChaincodeID: channelRequest.CCID,
			Fcn:         channelRequest.Fcn,
			Args:        args,
			TransientMap: transientMap,
		},
		channel.WithTargetEndpoints(channelRequest.TargetPeers...),
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
