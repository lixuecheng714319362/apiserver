package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type EventClient struct {
	ConfigPath string
	UserName   string
	ChannelID  string

	SDK    *fabsdk.FabricSDK
	Client *event.Client
}

// 创建通道客户端
func GetEventClient(configPath, userName, channelID string) (*EventClient, error) {
	SDK, err := InitializeSDK(configPath)
	if err != nil {
		return nil, err
	}
	client, err := event.New(SDK.ChannelContext(channelID, fabsdk.WithUser(userName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}

	return &EventClient{
		configPath,
		userName,
		channelID,
		SDK,
		client,
	}, nil
}

func (EventClient *EventClient) RegisterChaincodeEvent(chainCodeID string) (fab.Registration, <-chan *fab.CCEvent, error) {
	//TODO
	//eventFilter正则表达式匹配所有事件
	return EventClient.Client.RegisterChaincodeEvent(chainCodeID, ".*")
}

//Err: err block events are not permitted
func (EventClient *EventClient) RegisterBlockEvent(chainCodeID string) (fab.Registration, <-chan *fab.BlockEvent, error) {
	//TODO
	//eventFilter正则表达式匹配所有事件
	return EventClient.Client.RegisterBlockEvent()
}

func (EventClient *EventClient) RegisterFilteredBlockEvent() (fab.Registration, <-chan *fab.FilteredBlockEvent, error) {
	return EventClient.Client.RegisterFilteredBlockEvent()
}

//根据交易id注册时间，需要在提交proposal后获得id并进行注册，然后提交交易才可监听到时间。
func (EventClient *EventClient) RegisterTxStatusEvent(txID string) (fab.Registration, <-chan *fab.TxStatusEvent, error) {
	return EventClient.Client.RegisterTxStatusEvent(txID)
}

func (EventClient *EventClient) Unregister(reg fab.Registration) {
	EventClient.Client.Unregister(reg)
}

// 关闭SDK
func (EventClient *EventClient) CloseSDK() {
	EventClient.SDK.Close()
}
