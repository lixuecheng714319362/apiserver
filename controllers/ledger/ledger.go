package ledger

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk"
	"apiserver/models/gosdk/tool/blockdata"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"net/http"
	"reflect"
)

type LedgerController struct {
	beego.Controller
}
type Request struct {
	Data  string
}


type ReqData struct {
	ConfigPath  string
	ChannelID   string
	UserName    string //组织用户名
	OrgName     string //组织在sdk配置文件中的标识
	BlockHash   []byte
	TxID        string
	BlockNumber uint64
	Start       uint64
	End         uint64
}

func getReq(c *LedgerController) (*ReqData,error)  {
	data := c.Ctx.Input.RequestBody
	//测试接口使用
	if beego.AppConfig.String("filter")=="false"{
		r:=&ReqData{}
		_=json.Unmarshal(data,r)
		return r,nil
	}
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil,err
	}
	reqData :=&ReqData{}
	err=json.Unmarshal([]byte(req.Data), reqData)
	return reqData,nil
}


func (c *LedgerController) QueryInfo() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	res, err := LedgerClient.QueryInfo()
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query channel info  status:",res.Status,"block:", res.BCI.Height, "current block hash:",base64.StdEncoding.EncodeToString(res.BCI.CurrentBlockHash))
	tool.BackResData(c.Controller, res.BCI)
	return
}

type ChannelConfig struct {
	ID          string
	BlockNumber uint64
	Version     *fab.Versions
	Orderers    []string
	MSPs        []*msp.MSPConfig
	AnchorPeers []*fab.OrgAnchorPeer
}

func (c *LedgerController) QueryConfig() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	res, err := LedgerClient.QueryConfig()
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	var channelcfg = &ChannelConfig{
		ID:          res.ID(),
		BlockNumber: res.BlockNumber(),
		Version:     res.Versions(),
		Orderers:    res.Orderers(),
		MSPs:        res.MSPs(),
		AnchorPeers: res.AnchorPeers(),
	}
	beego.Info("query channel config","channel id is ",res.ID())
	tool.BackResData(c.Controller, channelcfg)
	return
}
func (c *LedgerController) QueryBlockByHash() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	//hash,err:=base64.StdEncoding.DecodeString(req.BlockHash)
	//if err != nil {
	//	beego.Error("base64 decode failed ", err)
	//	tool.BackResError(c.Controller,http.StatusBadRequest,err.Error())
	//	return
	//}
	res, err := LedgerClient.QueryBlockByHash(req.BlockHash)
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b, err := Getinfo(res)
	if err != nil {
		beego.Error("get block info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b.CurrentBlockHash = req.BlockHash
	tool.BackResData(c.Controller, b)
	return
}

func (c *LedgerController) QueryBlockByTxID() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	res, err := LedgerClient.QueryBlockByTxID(req.TxID)
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b, err := Getinfo(res)
	if err != nil {
		beego.Error("get block info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := getBlockHashBynumber(LedgerClient, b.Number)
	if err != nil {
		beego.Error("get current hash failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b.CurrentBlockHash = hash
	tool.BackResData(c.Controller, b)
	return
}
func (c *LedgerController) QueryBlockByNumber() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	res, err := LedgerClient.QueryBlockByNumber(req.BlockNumber)
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b, err := Getinfo(res)
	if err != nil {
		beego.Error("get block info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query block info ,block number is",res.Header.Number)
	hash, err := getBlockHashBynumber(LedgerClient, req.BlockNumber)
	if err != nil {
		beego.Error("get current hash failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b.CurrentBlockHash = hash
	tool.BackResData(c.Controller, b)
	return
}
func (c *LedgerController) QueryBlockByRange() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	start := req.Start
	end := req.End
	var res []*Block
	for i := start; i <= end; i++ {
		b, err := LedgerClient.QueryBlockByNumber(i)
		if err != nil {
			beego.Error("query block  failed ", err)
			tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
			return
		}
		block, err := Getinfo(b)
		if err != nil {
			beego.Error("create resClient failed ", err)
			tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
			return
		}
		if i != start {
			res[i-start-1].CurrentBlockHash = block.PreviousHash
		}
		res = append(res, block)
	}
	hash, err := getBlockHashBynumber(LedgerClient, end)
	if err != nil {
		beego.Error("get blockhash by number failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	res[end-start].CurrentBlockHash = hash
	beego.Info("query block range  from",req.Start,"to",req.End)
	tool.BackResData(c.Controller, res)
	return
}

func getBlockHashBynumber(client *gosdk.LedgerClient, number uint64) ([]byte, error) {
	info, err := client.QueryInfo()
	if err != nil {
		return nil, err
	}
	if info.BCI.Height-1 == number {
		return info.BCI.CurrentBlockHash, nil
	} else if info.BCI.Height-2 == number {
		return info.BCI.PreviousBlockHash, nil
	}
	b, err := client.QueryBlockByNumber(number + 1)
	if err != nil {
		return nil, err
	}
	return b.Header.PreviousHash, nil
}
func (c *LedgerController) QueryBlockByNumbertest() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req.ConfigPath, req.ChannelID, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	res, err := LedgerClient.QueryBlockByNumber(req.BlockNumber)
	if err != nil {
		beego.Error("query ledger info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}


	// 将或取的区块信息全部解析成json字符串
	bt,err:=proto.Marshal(res)
	if err != nil {
		fmt.Println("marshal error " ,err)
	}
	msgType:=proto.MessageType("common.Block")
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	err =proto.Unmarshal(bt,msg)
	if err != nil {
		fmt.Println("Unmarshal  error",err)
	}
	err = protolator.DeepMarshalJSON(c.Controller.Ctx.ResponseWriter, msg)
	if err != nil {
		fmt.Println("Deep Marshal Json error",err)
	}





	b, err := blockdata.Getinfo(res)
	if err != nil {
		beego.Error("get block info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query block info ,block number is",res.Header.Number)
	hash, err := getBlockHashBynumber(LedgerClient, req.BlockNumber)
	if err != nil {
		beego.Error("get current hash failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	b.Header.CurrentBlockHash = hash
	tool.BackResData(c.Controller, b)
	return
}