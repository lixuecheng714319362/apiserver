package ledger

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	"apiserver/models/gosdk"
	"apiserver/models/gosdk/tool/blockdata"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/common/tools/protolator"
	_ "github.com/hyperledger/fabric/protos/common"
	"net/http"
	"reflect"
	"sync"
	"time"
	//cb "github.com/hyperledger/fabric/protos/common" // Import these to register the proto types
	_ "github.com/hyperledger/fabric/protos/msp"
	_ "github.com/hyperledger/fabric/protos/orderer"
	_ "github.com/hyperledger/fabric/protos/orderer/etcdraft"
	_ "github.com/hyperledger/fabric/protos/peer"
)

type LedgerController struct {
	beego.Controller
}

func (c *LedgerController) QueryInfo() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
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
	beego.Info("query channel info  status:%d, block:%d, current block hash:%s",
		res.Status, res.BCI.Height, hex.EncodeToString(res.BCI.CurrentBlockHash))
	beego.Debug(res)
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
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
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
	beego.Info("query channel config", "channel id is ", res.ID())
	beego.Debug(res)
	tool.BackResData(c.Controller, channelcfg)
	return
}

func (c *LedgerController) QueryBlockByHash() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()
	blockHash, err := hex.DecodeString(req.BlockHash)
	if err != nil {
		beego.Error("block hash hex decode err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	res, err := LedgerClient.QueryBlockByHash(blockHash)
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
	beego.Info("query block by hash ,block hash  is", req.BlockHash)
	tool.BackResData(c.Controller, b)
	return
}

func (c *LedgerController) QueryBlockByTxID() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
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
	beego.Info("query block info by  txid  ,txid is", req.TxID)
	tool.BackResData(c.Controller, b)
	return
}

func (c *LedgerController) QueryBlockByNumber() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
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
	beego.Info("query block info ,block number is", res.Header.Number)
	tool.BackResData(c.Controller, b)
	return
}

func (c *LedgerController) QueryBlockByRange() {
	beego.Debug("start queryblockrange", time.Now().Unix())
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer LedgerClient.CloseSDK()

	res, err := getBlockbyRange(req.Start, req.End, LedgerClient)
	if err != nil {
		beego.Error("get block by range failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query block range  from", req.Start, "to", req.End)
	tool.BackResData(c.Controller, res)
	beego.Debug("end queryblockrange", time.Now().Unix())
	return
}

//func getBlockHashBynumber(client *gosdk.LedgerClient, number uint64) ([]byte, error) {
//	info, err := client.QueryInfo()
//	if err != nil {
//		return nil, err
//	}
//	if info.BCI.Height-1 == number {
//		return info.BCI.CurrentBlockHash, nil
//	} else if info.BCI.Height-2 == number {
//		return info.BCI.PreviousBlockHash, nil
//	}
//	b, err := client.QueryBlockByNumber(number + 1)
//	if err != nil {
//		return nil, err
//	}
//	return b.Header.PreviousHash, nil
//}

func (c *LedgerController) QueryBlockByNumberTest() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	LedgerClient, err := gosdk.GetLedgerClient(req)
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
	//TODO 必须注册所有引用的proto包参见configtxlator代码
	bt, err := proto.Marshal(res)
	if err != nil {
		fmt.Println("marshal error ", err)
	}
	msgType := proto.MessageType("common.Block")
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	err = proto.Unmarshal(bt, msg)
	if err != nil {
		fmt.Println("Unmarshal  error", err)
	}
	err = protolator.DeepMarshalJSON(c.Controller.Ctx.ResponseWriter, msg)
	if err != nil {
		fmt.Println("Deep Marshal Json error", err)
	}
	b, err := blockdata.Getinfo(res)
	if err != nil {
		beego.Error("get block info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query block test ,block number is", res.Header.Number)
	beego.Debug(b)
	beego.Debug()
	tool.BackResData(c.Controller, b)
	return
}

func getReq(c *LedgerController) (*gosdk.LedgerRequest, error) {
	data := c.Ctx.Input.RequestBody
	beego.Debug("request data is ", string(data))
	//测试接口使用
	if filter.IsFilterVerify == "false" {
		r := &gosdk.LedgerRequest{}
		if err := json.Unmarshal(data, r); err != nil {
			return nil, err
		}
		gosdk.ChangeLedgerRequestSingleConfig(r)
		return r, nil
	}
	req := &filter.ValidateRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	reqData := &gosdk.LedgerRequest{}
	err = json.Unmarshal([]byte(req.Data), reqData)
	gosdk.ChangeLedgerRequestSingleConfig(reqData)
	return reqData, nil
}



func getBlockbyRange(start uint64, end uint64, LedgerClient *gosdk.LedgerClient) ([]*Block, error) {
	beego.Debug("start getblockrange", time.Now().Unix())
	var res = make([]*Block, end-start+1)
	var wg sync.WaitGroup
	//异步获取所有区块
	flag := true
	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(i uint64) {
			for k := 0; ; k++ {
				b, err := LedgerClient.QueryBlockByNumber(i)
				if err != nil {
					if k > 3 {
						flag = false
						break
					} else {
						continue
					}
				}
				block, err := Getinfo(b)
				if err != nil {
					flag = false
					break
				}
				res[i-start] = block
				break
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if flag == false {
		return nil, errors.New("get block error")
	}
	beego.Debug("end getblockrange", time.Now().Unix())
	return res, nil
}
