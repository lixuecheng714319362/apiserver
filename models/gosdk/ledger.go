package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

type LedgerClient struct {
	ConfigPath string
	ChannelID  string
	UserName   string
	OrgName    string

	SDK    *fabsdk.FabricSDK
	Client *ledger.Client
}

// 创建账本客户端
func GetLedgerClient(configPath, channelID, userName, orgName string) (*LedgerClient, error) {
	SDK, err := InitializeSDK(configPath)
	if err != nil {
		return nil, err
	}
	ledgerClient, err := ledger.New(SDK.ChannelContext(channelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &LedgerClient{
		configPath,
		channelID,
		userName,
		orgName,
		SDK,
		ledgerClient,
	}, nil
}

//查询最新账本信息
func (LedgerClient *LedgerClient) QueryInfo() (*fab.BlockchainInfoResponse, error) {
	return LedgerClient.Client.QueryInfo()
}

//查询通道配置信息
func (LedgerClient *LedgerClient) QueryConfig() (fab.ChannelCfg, error) {
	return LedgerClient.Client.QueryConfig()
}

//交易ID为16进制字符串
func (LedgerClient *LedgerClient) QueryBlockByTxID(txID string) (*common.Block, error) {
	return LedgerClient.Client.QueryBlockByTxID(fab.TransactionID(txID))
}

//交易哈希使用base64
func (LedgerClient *LedgerClient) QueryBlockByHash(blockHash []byte) (*common.Block, error) {
	return LedgerClient.Client.QueryBlockByHash(blockHash)
}

//交易哈希使用base64
func (LedgerClient *LedgerClient) QueryBlockByNumber(number uint64) (*common.Block, error) {
	return LedgerClient.Client.QueryBlock(number)
}

// 关闭SDK
func (LedgerClient *LedgerClient) CloseSDK() {
	LedgerClient.SDK.Close()
}
