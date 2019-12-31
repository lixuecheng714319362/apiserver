package gosdk

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk/tool/utils"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/pkg/errors"
	"io"
)

type ResmgmtClient struct {
	ConfigPath string
	UserName   string
	OrgName    string

	SDK    *fabsdk.FabricSDK
	Client *resmgmt.Client
}

type ResmgmtRequest struct {
	//资源管理配置
	ConfigPath    string
	UserName      string
	OrgName       string

	ChannelID     string
	TargetOrderer string
	//TargetPeer    string   //查询连码是否安装
	TargetPeers   []string //default 安装全部peer

	//Channel 请求内容
	//AddOrgList 	[]string			//动态添加组织，待添加列表
	CRL 			[]string
	DeleteOrgList 	[]string
	AddOrgConfig []localconfig.Organization
	OrgNameList   []string			//orgName 在系统通道中组织的字段
	//MSP 签名组织对应制定用户  key为签名组织在sdk配置文件中的org标示，value为组织用户的User名
	SignerMap     map[string]string
	ChannelTxPath string

	//Chain Code请求配置
	CCPath           string
	CCGoPath         string
	CCID             string
	CCVersion        string
	CCPackage        string
	Policy           string //背书策略
	CollectionConfig []utils.CollectionConfig
	Args             []string
}

//参数 configPath sdk 配置文件路径  OrgAdmin组织管理账户名  orgName 组织名称
// 创建资源管理
func GetResMgmtClient(request *ResmgmtRequest) (ResClient *ResmgmtClient, err error) {
	// 创建资源管理客户端上下文
	SDK, err := InitializeSDK(request.ConfigPath)
	if err != nil {
		return nil, err
	}
	AdminClient, err := resmgmt.New(SDK.Context(fabsdk.WithUser(request.UserName), fabsdk.WithOrg(request.OrgName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &ResmgmtClient{
		request.ConfigPath,
		request.UserName,
		request.OrgName,
		SDK,
		AdminClient,
	}, nil
}
// 创建通道,不通过channel TxPath
func (ResmgmtClient *ResmgmtClient) UpdateChannel(request *ResmgmtRequest ,reader io.Reader) ( resmgmt.SaveChannelResponse, error) {
	var err error
	var signer []msp.SigningIdentity

	if signer, err = utils.GetCreateChannelSinger(request.SignerMap, ResmgmtClient.SDK); err != nil {
		return resmgmt.SaveChannelResponse{}, err
	}
	saveChannelReq := resmgmt.SaveChannelRequest{
		ChannelID:         request.ChannelID,
		ChannelConfig:     reader,
		SigningIdentities: signer,
	}
	return ResmgmtClient.Client.SaveChannel(saveChannelReq, resmgmt.WithOrdererEndpoint(request.TargetOrderer))
}
// 创建通道,不通过channel TxPath
func (ResmgmtClient *ResmgmtClient) CreateNewChannel(request *ResmgmtRequest) ( resmgmt.SaveChannelResponse, error) {

	var err error
	var reader io.Reader
	var signer []msp.SigningIdentity
	var saveChannelReq resmgmt.SaveChannelRequest
	if request.ChannelTxPath == "" {

		if reader, err = GetCreateChannelReader(request.ChannelID, request.OrgNameList); err != nil {
			return resmgmt.SaveChannelResponse{}, err
		}
		if signer, err = utils.GetCreateChannelSinger(request.SignerMap, ResmgmtClient.SDK); err != nil {
			return resmgmt.SaveChannelResponse{}, err
		}
		saveChannelReq = resmgmt.SaveChannelRequest{
			ChannelID:         request.ChannelID,
			ChannelConfig:     reader,
			SigningIdentities: signer,
		}
	}else {
			saveChannelReq = resmgmt.SaveChannelRequest{
				ChannelID:         request.ChannelID,
				ChannelConfigPath: request.ChannelTxPath,
			}
	}

	return ResmgmtClient.Client.SaveChannel(saveChannelReq, resmgmt.WithOrdererEndpoint(request.TargetOrderer))
}

//参数ChannelID  管道ID channelTxPath 管道配置文件路径
// 创建通道
func (ResmgmtClient *ResmgmtClient) CreateChannel(request *ResmgmtRequest) (resmgmt.SaveChannelResponse, error) {
	saveChannelReq := resmgmt.SaveChannelRequest{
		ChannelID:         request.ChannelID,
		ChannelConfigPath: request.ChannelTxPath,
	}
	return ResmgmtClient.Client.SaveChannel(saveChannelReq, resmgmt.WithOrdererEndpoint(request.TargetOrderer))
}

//参数 ordererID 指定发送ordererID(sdk 配置文件中的orderer配置名称或orderer的全域名)可以不指定orderer但不能指定错误或""否则报错
// 加入通道
func (ResmgmtClient *ResmgmtClient) JoinChannel(request *ResmgmtRequest) error {
	return ResmgmtClient.Client.JoinChannel(
		request.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(request.TargetOrderer),
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
	)
}

//指定peer的名称，可以是sdk config文件中的制定名称也可以是peer的全域名
// 加入通道
func (ResmgmtClient *ResmgmtClient) QueryChannel(request *ResmgmtRequest) (*peer.ChannelQueryResponse, error) {
	return ResmgmtClient.Client.QueryChannels(
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
	)
}

func (ResmgmtClient *ResmgmtClient) QueryInstalledChaincodes(request *ResmgmtRequest) (*peer.ChaincodeQueryResponse, error) {
	return ResmgmtClient.Client.QueryInstalledChaincodes(
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
	)
}

func (ResmgmtClient *ResmgmtClient) QueryInstantiatedChaincodes(request *ResmgmtRequest) (*peer.ChaincodeQueryResponse, error) {
	return ResmgmtClient.Client.QueryInstantiatedChaincodes(
		request.ChannelID,
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
	)
}

// 安装链码
func (ResmgmtClient *ResmgmtClient) ChainCodeInstall(request *ResmgmtRequest) ([]resmgmt.InstallCCResponse, error) {
	// 打包链码
	var ccPkg = &resource.CCPackage{}
	var err error
	if request.CCPackage == "" {
		ccPkg, err = packager.NewCCPackage(request.CCPath, request.CCGoPath)
	} else {
		err = json.Unmarshal([]byte(request.CCPackage), ccPkg)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "NewCCPackage err")
	}
	installCCReq := resmgmt.InstallCCRequest{
		Name:    request.CCID,
		Path:    request.CCPath, //不能使用绝对路径
		Version: request.CCVersion,
		Package: ccPkg,
	}
	return ResmgmtClient.Client.InstallCC(installCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
	)
}

// 初始链码
func (ResmgmtClient *ResmgmtClient) ChainCodeInit(request *ResmgmtRequest) (resmgmt.InstantiateCCResponse, error) {
	// 设置背书策略,该参数依赖configtx.yaml文件中Organizations->MSPID
	// 需要区分背书策略
	ccPolicy, err := cauthdsl.FromString(request.Policy)
	if err != nil {
		return resmgmt.InstantiateCCResponse{}, err
	}
	colConfigs := make([]*common.CollectionConfig, 0)
	if request.CollectionConfig != nil {
		for _, v := range request.CollectionConfig {
			cf, err := utils.NewCollectionConfig(v.Name, v.MemberOrgsPolicy, v.RequiredPeerCount, v.MaximumPeerCount, v.BlockToLive,v.MemberOnlyRead)
			if err != nil {
				return resmgmt.InstantiateCCResponse{}, err
			}
			colConfigs = append(colConfigs, cf)
		}
	} else {
		colConfigs = nil
	}

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:       request.CCID,
		Path:       request.CCGoPath,
		Version:    request.CCVersion,
		Args:       tool.ChangeArgs(request.Args),
		Policy:     ccPolicy,
		CollConfig: colConfigs,
	}
	return ResmgmtClient.Client.InstantiateCC(
		request.ChannelID, instantiateCCReq,
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
		resmgmt.WithOrdererEndpoint(request.TargetOrderer),
	)
}

// 升级链码
func (ResmgmtClient *ResmgmtClient) ChainCodeUpgrade(request *ResmgmtRequest) (rep resmgmt.UpgradeCCResponse, err error) {
	// 安装链码
	// 需要区分背书策略
	ccPolicy, err := cauthdsl.FromString(request.Policy)
	if err != nil {
		return resmgmt.UpgradeCCResponse{}, err
	}
	colConfig := make([]*common.CollectionConfig, 0)
	if request.CollectionConfig != nil {
		for _, v := range request.CollectionConfig {
			cf, err := utils.NewCollectionConfig(v.Name, v.MemberOrgsPolicy, v.RequiredPeerCount, v.MaximumPeerCount, v.BlockToLive,v.MemberOnlyRead)
			if err != nil {
				return resmgmt.UpgradeCCResponse{}, err
			}
			colConfig = append(colConfig, cf)
		}
	} else {
		colConfig = nil
	}
	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    request.CCID,
		Path:    request.CCGoPath,
		Version: request.CCVersion,
		Args:    tool.ChangeArgs(request.Args),
		Policy:  ccPolicy,
		CollConfig:colConfig,
	}
	return ResmgmtClient.Client.UpgradeCC(
		request.ChannelID,
		upgradeCCReq,
		resmgmt.WithTargetEndpoints(request.TargetPeers...),
		resmgmt.WithOrdererEndpoint(request.TargetOrderer),
	)
}

// 关闭SDK
func (ResmgmtClient *ResmgmtClient) CloseSDK() {
	ResmgmtClient.SDK.Close()
}
