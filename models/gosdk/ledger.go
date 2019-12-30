package gosdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

type LedgerClient struct {
	ConfigPath string
	ChannelID  string
	UserName   string
	OrgName    string

	SDK    *fabsdk.FabricSDK
	Client *ledger.Client
}


type LedgerRequest struct {
	ConfigPath  string
	UserName    string //组织用户名
	OrgName     string //组织在sdk配置文件中的标识

	ChannelID   string
	BlockHash   string
	TxID        string
	BlockNumber uint64
	Start       uint64
	End         uint64
}

// 创建账本客户端
func GetLedgerClient(ledgerRequest *LedgerRequest) (*LedgerClient, error) {
	SDK, err := InitializeSDK(ledgerRequest.ConfigPath)
	if err != nil {
		return nil, err
	}
	ledgerClient, err := ledger.New(SDK.ChannelContext(ledgerRequest.ChannelID, fabsdk.WithUser(ledgerRequest.UserName), fabsdk.WithOrg(ledgerRequest.OrgName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &LedgerClient{
		ledgerRequest.ConfigPath,
		ledgerRequest.ChannelID,
		ledgerRequest.UserName,
		ledgerRequest.OrgName,
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

//
func (LedgerClient *LedgerClient) QueryTransactionByTxID(txId string) (*peer.ProcessedTransaction, error) {
	return LedgerClient.Client.QueryTransaction(fab.TransactionID(txId))
}
// 关闭SDK
func (LedgerClient *LedgerClient) CloseSDK() {
	LedgerClient.SDK.Close()
}


//获取指定channel的配置信息
func GetOldChannelConfig(configPath, orgName,userName ,channelID string,targetPeers []string)(*common.Config,error) {
	cl,err:= GetLedgerClient(&LedgerRequest{
		ConfigPath:configPath,
		OrgName:orgName,
		UserName:userName,
		ChannelID:channelID,
	})
	if err != nil {
		return nil,err
	}
	defer cl.CloseSDK()
	config,err:=cl.Client.QueryConfig(ledger.WithTargetEndpoints(targetPeers...))
	if err != nil {
		return nil,err
	}
	block,err:=cl.Client.QueryBlock(config.BlockNumber(),ledger.WithTargetEndpoints(targetPeers...))
	if err != nil {
		return nil,err
	}
	configEnv, err := resource.CreateConfigEnvelope(block.Data.Data[0])
	if err != nil {
		return nil,err
	}
	return configEnv.Config,nil
}