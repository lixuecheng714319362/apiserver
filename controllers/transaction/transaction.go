package transaction

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	"apiserver/models/gosdk"
	"apiserver/models/gosdk/tool/blockdata"
	redismodel "apiserver/models/redis"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type InvokeController struct {
	beego.Controller
}

func (c *InvokeController) Invoke() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("invoke chainCode ccId:",req.CCID,", txId:",res.TransactionID,", statusCode:",res.ChaincodeStatus,
		", args:",tool.ChangeArgs(req.Args),", payload:", string(res.Payload))
	tool.BackResData(c.Controller, res)
	return
}

func (c *InvokeController) Query() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Query(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query chainCode ", req.CCID, tool.ChangeArgs(req.Args), string(res.Payload))
	beego.Debug("query result", res)
	tool.BackResData(c.Controller, string(res.Payload))
	return
}
func (c *InvokeController) QueryTx() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Query(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query chainCode ", req.CCID, tool.ChangeArgs(req.Args), string(res.Payload))
	beego.Debug("query result", res)
	cr,err:=blockdata.GetChannelResposeInfo(&res)
	if err != nil {
		beego.Error("decode channel response info ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	tool.BackResData(c.Controller, cr)
	return
}

func (c *InvokeController) QueryTest() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Query(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query chainCode ", req.CCID, tool.ChangeArgs(req.Args), string(res.Payload))
	beego.Debug("query result", res)
	cr,err:=blockdata.GetChannelResposeInfo(&res)
	if err != nil {
		beego.Error("decode channel response info ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	tool.BackResData(c.Controller, cr)
	return
}


func (c *InvokeController) InvokeEmpty() {
	defer tool.HanddlerError(c.Controller)
	_, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	beego.Info("query chainCode empty")
	tool.BackResSuccess(c.Controller)
	return
}

func (c *InvokeController) InvokeFunc() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := getClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	//defer channelClient.CloseSDK() //该sdk复用
	res, err := channelClient.Invoke(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("invoke chainCode ccId:",req.CCID,", txId:",res.TransactionID,", statusCode:",res.ChaincodeStatus,
		", args:",tool.ChangeArgs(req.Args),", payload:", string(res.Payload))
	tool.BackResData(c.Controller, res)
	return
}

func getReq(c *InvokeController) (*gosdk.ChannelRequest, error) {
	data := c.Ctx.Input.RequestBody
	beego.Debug("request data is ", string(data))
	//测试接口使用
	if filter.IsFilterVerify == "false" {
		r := &gosdk.ChannelRequest{}
		if err := json.Unmarshal(data, r); err != nil {
			return nil, err
		}
		gosdk.ChangeChannelRequestSingleConfig(r)
		return r, nil
	}
	req := &filter.ValidateRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	reqData := &gosdk.ChannelRequest{}
	err = json.Unmarshal([]byte(req.Data), reqData)
	gosdk.ChangeChannelRequestSingleConfig(reqData)
	return reqData, nil
}


//下面为缓存redis数据，和缓存client
var channelClientPool = make(map[string]*gosdk.ChannelClient)

func getClient(req *gosdk.ChannelRequest) (*gosdk.ChannelClient, error) {

	if c, ok := channelClientPool[req.ConfigPath+req.UserName+req.ChannelID]; ok {
		return c, nil
	}
	c, err := gosdk.GetChannelClient(req)
	channelClientPool[req.ConfigPath+req.UserName+req.ChannelID] = c
	return c, err
}

var InvokeClient *gosdk.ChannelClient
var RedisCLient redismodel.RedisCLient

func init() {
	var err error
	flag := beego.AppConfig.String("redis")
	if flag != "true" {
		return
	}
	req := &gosdk.ChannelRequest{
		ConfigPath: "/root/projectMod/apiserver/conf/config/multi-5host-ali/O1P0.yaml",
		UserName:   "Admin",
		ChannelID:  "mychannel",
	}
	InvokeClient, err = gosdk.GetChannelClient(req)
	if err != nil {
		panic(err)
	}
	go Start()
	f, err := os.Create("kv.txt")
	if err != nil {
		panic(err)
	}
	file = f
	RedisCLient = redismodel.Client
}

var file *os.File
var fileLock sync.Mutex

func write(kv []byte) {
	fileLock.Lock()
	defer fileLock.Unlock()
	// 查找文件末尾的偏移量
	n, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		beego.Error("Get end of file", err)
	}
	// 从末尾的偏移量开始写入内容
	_, err = file.WriteAt(kv, n)
	if err != nil {
		beego.Error("Write kv error", err)
	}
}

var Top = 0

//var TopLock  sync.Mutex
//func getTop() int {
//	return Top
//}

var Txnum int
var TxNumLock sync.Mutex

func getTxnum() int {
	TxNumLock.Lock()
	defer TxNumLock.Unlock()
	Txnum++
	return Txnum
}

var pool *redismodel.Pool

func Start() {
	p := redismodel.NewPool(2000, 10000000)
	p.Start()
	pool = p

}

type Tx struct {
	Top int
}

func (s *Tx) Do() {
	data, err := RedisCLient.Get(strconv.Itoa(100000000 + s.Top))
	if err != nil {
		beego.Error("GET ERROR ", err)
	}
	req := &gosdk.ChannelRequest{}
	err = json.Unmarshal([]byte(data), req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
	}
	if len(req.Args) != 2 {
		beego.Error("args error ", err)
		return
	}
	req.Args[0] = strconv.Itoa(100000000 + s.Top)
	h := sha256.Sum256([]byte(strconv.Itoa(100000000 + s.Top)))
	req.Args[1] = hex.EncodeToString(h[:])
	_, err = InvokeClient.Invoke(req, tool.ChangeArgs(req.Args),tool.ChangeTransientMap(req.TransientMap))
	if err != nil {
		beego.Error("exec error ", err)
	}
}

func (c *InvokeController) InvokeFuncTmp() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	num := getTxnum()
	h := sha256.Sum256([]byte(strconv.Itoa(100000000 + num)))
	s := hex.EncodeToString(h[:])
	write([]byte(strconv.Itoa(100000000+num) + "|  " + s + "\n"))
	_, err := RedisCLient.Put(strconv.Itoa(100000000+num), string(data))
	if err != nil {
		beego.Error("PUT ERROR", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	go DO(num)
	tool.BackResSuccess(c.Controller)
	return
}

func DO(num int) {
	pool.JobQueue <- &Tx{num}
}
