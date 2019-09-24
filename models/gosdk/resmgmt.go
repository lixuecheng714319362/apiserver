package gosdk

import "C"
import (
	"apiserver/models/gosdk/tool/utils"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"strings"
)

type ResmgmtClient struct {
	ConfigPath string
	UserName   string
	OrgName    string

	SDK    *fabsdk.FabricSDK
	Client *resmgmt.Client
}

//参数 configPath sdk 配置文件路径  OrgAdmin组织管理账户名  orgName 组织名称
// 创建资源管理
func GetResMgmtClient(configPath, userName, orgName string) (ResClient *ResmgmtClient, err error) {
	// 创建资源管理客户端上下文
	SDK, err := InitializeSDK(configPath)
	if err != nil {
		return nil, err
	}
	AdminClient, err := resmgmt.New(SDK.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)))
	if err != nil {
		SDK.Close()
		return nil, err
	}
	return &ResmgmtClient{
		configPath,
		userName,
		orgName,
		SDK,
		AdminClient,
	}, nil
}

// 判断是否创建通道和加入通道
func (ResmgmtClient *ResmgmtClient) CreateAndJoinChannel(channelID, channelConfig, targetOrderer string, PeerNameList []string) error {
	channelInstall := false
	for _, v := range PeerNameList {
		channelRes, err := ResmgmtClient.QueryChannel(v)
		if err != nil {
			return err
		}
		if channelRes != nil {
			for _, v := range channelRes.Channels {
				if strings.EqualFold(channelID, v.ChannelId) {
					channelInstall = true
				}
			}
		}
	}
	if !channelInstall {

		if _, err := ResmgmtClient.CreateChannel(channelID, channelConfig, targetOrderer); err != nil {
			return err
		}
		if err := ResmgmtClient.JoinChannel(targetOrderer, channelID, PeerNameList); err != nil {
			return err
		}
	}
	return nil
}

//参数ChannelID  管道ID channelTxPath 管道配置文件路径
// 创建通道
func (ResmgmtClient *ResmgmtClient) CreateChannel(channelID, channelTxPath, targetOrderer string) (resmgmt.SaveChannelResponse, error) {
	saveChannelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelTxPath,
	}
	return ResmgmtClient.Client.SaveChannel(saveChannelReq, resmgmt.WithOrdererEndpoint(targetOrderer))
}

//参数 ordererID 指定发送ordererID(sdk 配置文件中的orderer配置名称或orderer的全域名)可以不指定orderer但不能指定错误或""否则报错
// 加入通道
func (ResmgmtClient *ResmgmtClient) JoinChannel(targetOrderer, channelID string, targetPeers []string) error {
	return ResmgmtClient.Client.JoinChannel(
		channelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(targetOrderer),
		resmgmt.WithTargetEndpoints(targetPeers...),
	)
}

//指定peer的名称，可以是sdk config文件中的制定名称也可以是peer的全域名
// 加入通道
func (ResmgmtClient *ResmgmtClient) QueryChannel(targetPeerName string) (*peer.ChannelQueryResponse, error) {
	return ResmgmtClient.Client.QueryChannels(
		resmgmt.WithTargetEndpoints(targetPeerName),
	)
}

//TODO 传peer列表可能会出现错误
// 判断是否安装链码和初始化链码
func (ResmgmtClient *ResmgmtClient) InstallAndInitChainCode(policy , chainCodePath, chainCodeGoPath,
	chainCodeID, chainCodeVersion, channelID, targetOrderer string, PeerNameList []string, args [][]byte,CollectionConfig []utils.CollectionConfig) error {

	var chainCodeInstalled, chainCodeInstantiated bool
	for _, v := range PeerNameList {
		chainCodeInstalled = false
		endpoints := resmgmt.WithTargetEndpoints(v)
		ccRes, err := ResmgmtClient.Client.QueryInstalledChaincodes(endpoints)
		if err != nil {
			return err
		}
		if ccRes != nil {
			for _, v := range ccRes.Chaincodes {
				if strings.EqualFold(v.Name, chainCodeID) {
					chainCodeInstalled = true
				}
			}
		}
		if !chainCodeInstalled {
			if _, err := ResmgmtClient.ChainCodeInstall(chainCodePath, chainCodeGoPath, chainCodeID, chainCodeVersion, PeerNameList); err != nil {
				return err
			}
		}

		//TODO
		//未验证初始化peer
		chainCodeInstantiated = false
		ccRes, err = ResmgmtClient.QueryInstantiatedChaincodes(channelID, v)
		if err != nil {
			return err
		}
		if ccRes.Chaincodes != nil && len(ccRes.Chaincodes) > 0 {
			for _, v := range ccRes.Chaincodes {
				if strings.EqualFold(v.Name, chainCodeID) {
					chainCodeInstantiated = true
				}
			}
		}
		if !chainCodeInstantiated {
			if _, err := ResmgmtClient.ChainCodeInit(policy, chainCodeID, chainCodeGoPath, chainCodeVersion,
				channelID, targetOrderer, args, PeerNameList,CollectionConfig); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ResmgmtClient *ResmgmtClient) QueryInstalledChaincodes(TargetPeerName string) (*peer.ChaincodeQueryResponse, error) {
	return ResmgmtClient.Client.QueryInstalledChaincodes(
		resmgmt.WithTargetEndpoints(TargetPeerName),
	)
}

func (ResmgmtClient *ResmgmtClient) QueryInstantiatedChaincodes(ChannelID, TargetPeerName string) (*peer.ChaincodeQueryResponse, error) {
	return ResmgmtClient.Client.QueryInstantiatedChaincodes(
		ChannelID,
		resmgmt.WithTargetEndpoints(TargetPeerName),
	)
}

// 安装链码
func (ResmgmtClient *ResmgmtClient) ChainCodeInstall(chainCodePath, chainCodeGoPath, chainCodeID,
	chainCodeVersion string, targetPeers []string) ([]resmgmt.InstallCCResponse, error) {
	// 打包链码
	ccPkg, err := packager.NewCCPackage(chainCodePath, chainCodeGoPath)
	if err != nil {
		return nil, errors.WithMessage(err, "NewCCPackage err")
	}
	installCCReq := resmgmt.InstallCCRequest{
		Name:    chainCodeID,
		Path:    chainCodePath, //不能使用绝对路径
		Version: chainCodeVersion,
		Package: ccPkg,
	}
	return ResmgmtClient.Client.InstallCC(installCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(targetPeers...),
	)
}

// 初始链码
func (ResmgmtClient *ResmgmtClient) ChainCodeInit(poicy ,chainCodeID, chainCodeGoPath,
	chainCodeVersion, channelID, targetOrderer string, args [][]byte, targetPeers []string ,CollectionConfig []utils.CollectionConfig) (resmgmt.InstantiateCCResponse, error) {
	// 设置背书策略,该参数依赖configtx.yaml文件中Organizations->ID
	//TODO
	// 需要区分背书策略
	//policy:cauthdsl.SignedByNOutOfGivenRole(int32(len(orgMSPID)), msp.MSPRole_MEMBER,orgMSPID)//and背书策略
	//ccPolicy := cauthdsl.SignedByAnyMember(orgMSPID)
	ccPolicy,err:=cauthdsl.FromString(poicy)
	if err != nil {
		return resmgmt.InstantiateCCResponse{},err
	}
	var collconfig=[]*common.CollectionConfig{}
	if CollectionConfig!=nil{
		for _, v := range CollectionConfig {
			cf,err:=utils.NewCollectionConfig(v.Name,v.MemberOrgsPolicy, v.RequiredPeerCount, v.MaximumPeerCount, v.BlockToLive)
			if err != nil {
				return resmgmt.InstantiateCCResponse{},err
			}
			collconfig=append(collconfig, cf)
		}
	}else {
		collconfig=nil
	}

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    chainCodeID,
		Path:    chainCodeGoPath,
		Version: chainCodeVersion,
		Args:    args,
		Policy:  ccPolicy,
		CollConfig:collconfig,
	}
	return ResmgmtClient.Client.InstantiateCC(
		channelID, instantiateCCReq,
		resmgmt.WithTargetEndpoints(targetPeers...),
		resmgmt.WithOrdererEndpoint(targetOrderer),
	)
}

// 升级链码
func (ResmgmtClient *ResmgmtClient) ChainCodeUpgrade(policy string, targetPeers []string,  chainCodeID,
	chainCodeGoPath, chainCodeVersion, channelID, targetOrderer string, args [][]byte,CollectionConfig []utils.CollectionConfig) (rep resmgmt.UpgradeCCResponse, err error) {
	// 安装链码
	// 需要区分背书策略
	ccPolicy ,err:= cauthdsl.FromString(policy)
	if err != nil {
		return resmgmt.UpgradeCCResponse{},err
	}
	var collconfig=[]*common.CollectionConfig{}
	if CollectionConfig!=nil{
		for _, v := range CollectionConfig {
			cf,err:=utils.NewCollectionConfig(v.Name,v.MemberOrgsPolicy, v.RequiredPeerCount, v.MaximumPeerCount, v.BlockToLive)
			if err != nil {
				return resmgmt.UpgradeCCResponse{},err
			}
			collconfig=append(collconfig, cf)
		}
	}else {
		collconfig=nil
	}
	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    chainCodeID,
		Path:    chainCodeGoPath,
		Version: chainCodeVersion,
		Args:    args,
		Policy:  ccPolicy,
	}
	return  ResmgmtClient.Client.UpgradeCC(
		channelID,
		upgradeCCReq,
		resmgmt.WithTargetEndpoints(targetPeers...),
		resmgmt.WithOrdererEndpoint(targetOrderer),
	)
}

// 关闭SDK
func (ResmgmtClient *ResmgmtClient) CloseSDK() {
	ResmgmtClient.SDK.Close()
}
