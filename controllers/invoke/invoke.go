package invoke

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk"
	redismodel "apiserver/models/redis"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type InvokeController struct {
	beego.Controller
}

var (
	Org1MSP = "CopyRightChain1MSP"
)

type Request struct {
	ConfigPath  string
	UserName    string
	OrdererName string
	ChannelID   string
	CCID        string
	Fcn         string //OrdID
	Args        []string
	TargetPeers []string
}

func (c *InvokeController) Query() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Query(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	if err != nil {
		//beego.Error("invoke query err", err)
		//tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		//payload,err:=RedisCLient.Get(req.Args[0])
		h:=sha256.Sum256([]byte(req.Args[0]))
		s:=hex.EncodeToString(h[:])
		if err != nil {
			n,_ :=strconv.Atoi(req.Args[0])
			if n-10000000>Txnum{
				tool.BackResError(c.Controller, http.StatusBadRequest, "key not exist")
			}

			tool.BackResData(c.Controller,s)
			return
		}
		tool.BackResData(c.Controller, s)
		return
	}
	beego.Info("query", string(res.Payload))
	tool.BackResData(c.Controller, string(res.Payload))

	return
}

var channelClientPool = make(map[string]*gosdk.ChannelClient)

func getClient(configPath, userName, channelID string) (*gosdk.ChannelClient, error) {

	if c, ok := channelClientPool[configPath+userName+channelID]; ok {
		return c, nil
	}
	c, err := gosdk.GetChannelClient(configPath, userName, channelID)
	channelClientPool[configPath+userName+channelID] = c
	return c, err
}

func (c *InvokeController) InvokeFunc() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := getClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	//defer channelClient.CloseSDK()

	_, err = channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	//go channelClient.Invoke(req.CCID,req.Fcn, tool.ChangeArgs(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	//beego.Info("%s, %v\n ",res.TransactionID,res.ChaincodeStatus)
	//tool.BackResData(c.Controller,res)
	tool.BackResSuccess(c.Controller)
	return
}

func (c *InvokeController) InvokeFunc1() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("%s, %v\n ", res.TransactionID, res.ChaincodeStatus,string(res.Payload))
	tool.BackResData(c.Controller, res)
	return
}

func (c *InvokeController) InvokeFunc2() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgsWithSm2(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("%s, %v\n ", res.TransactionID, res.ChaincodeStatus)
	tool.BackResData(c.Controller, res)
	return
}

func (c *InvokeController) InvokeFunc3() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgsWithSm3(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("%s, %v\n ", res.TransactionID, res.ChaincodeStatus)
	tool.BackResData(c.Controller, res)
	return
}

func (c *InvokeController) InvokeFunc4() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgsWithSm4(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("%s, %v\n ", res.TransactionID, res.ChaincodeStatus)
	tool.BackResData(c.Controller, res)
	return
}

func (c *InvokeController) InvokeFunc5() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	channelClient, err := gosdk.GetChannelClient(req.ConfigPath, req.UserName, req.ChannelID)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer channelClient.CloseSDK()
	res, err := channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgsWithSHA256(req.Args))
	if err != nil {
		beego.Error("invoke query err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("%s, %v\n ", res.TransactionID, res.ChaincodeStatus)
	tool.BackResData(c.Controller, res)
	return
}


func (c *InvokeController) InvokeEmpty() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	//channelClient, err := getClient(req.ConfigPath, req.UserName, req.ChannelID)
	//if err != nil {
	//	beego.Error("create resClient failed ", err)
	//	tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
	//	return
	//}
	////defer channelClient.CloseSDK()
	//
	//_, err = channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	////go channelClient.Invoke(req.CCID,req.Fcn, tool.ChangeArgs(req.Args))
	//if err != nil {
	//	beego.Error("invoke query err", err)
	//	tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
	//	return
	//}
	//beego.Info("%s, %v\n ",res.TransactionID,res.ChaincodeStatus)
	//tool.BackResData(c.Controller,res)
	tool.BackResSuccess(c.Controller)
	return
}
func (c *InvokeController) InvokeEmptyTime() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	time.Sleep(10*time.Millisecond)
	//channelClient, err := getClient(req.ConfigPath, req.UserName, req.ChannelID)
	//if err != nil {
	//	beego.Error("create resClient failed ", err)
	//	tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
	//	return
	//}
	////defer channelClient.CloseSDK()
	//
	//_, err = channelClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	////go channelClient.Invoke(req.CCID,req.Fcn, tool.ChangeArgs(req.Args))
	//if err != nil {
	//	beego.Error("invoke query err", err)
	//	tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
	//	return
	//}
	//beego.Info("%s, %v\n ",res.TransactionID,res.ChaincodeStatus)
	//tool.BackResData(c.Controller,res)
	tool.BackResSuccess(c.Controller)
	return
}

var InvokeClient  *gosdk.ChannelClient
var RedisCLient redismodel.RedisCLient
func init()  {
	var err error
	flag,err:=beego.AppConfig.Bool("redis")
	if err != nil {
		panic(err)
	}
	if !flag{
		return
	}
	InvokeClient,err=gosdk.GetChannelClient("/root/projectMod/apiserver/conf/config/multi-5host-ali/O1P0.yaml","Admin","mychannel")
	if err != nil {
		panic(err)
	}
	go Start()
	f,err:=os.Create("kv.txt")
	if err != nil {
		panic(err)
	}
	file=f
	RedisCLient=redismodel.Client
}


var file *os.File
var fileLock sync.Mutex


func write( kv []byte)  {
	fileLock.Lock()
	defer fileLock.Unlock()
	// 查找文件末尾的偏移量
	n, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		beego.Error("Get end of file",err)
	}
	// 从末尾的偏移量开始写入内容
	_, err = file.WriteAt(kv, n)
	if err != nil {
		beego.Error("Write kv error",err)
	}
}

var Top int =0
var TopLock  sync.Mutex
func getTop() int {
	return Top
}

var Txnum int
var TxNumLock  sync.Mutex
func getTxnum() int {
	TxNumLock.Lock()
	defer TxNumLock.Unlock()
	Txnum++
	return Txnum
}


var pool*redismodel.Pool

func Start()  {
	p:=redismodel.NewPool(2000,10000000)
	p.Start()
	pool=p

}



type Tx struct {
	Top int
}




func (s *Tx) Do() {
	data,err:=RedisCLient.Get(strconv.Itoa(100000000+s.Top))
	if err != nil {
		beego.Error("GET ERROR ", err)
	}
	req := &Request{}
	err = json.Unmarshal([]byte(data), req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
	}
	if len(req.Args)!=2{
		beego.Error("args error ", err)
		return
	}
	req.Args[0]=strconv.Itoa(100000000+s.Top)
	h:=sha256.Sum256([]byte(strconv.Itoa(100000000+s.Top)))
	req.Args[1]=hex.EncodeToString(h[:])
		_,err=InvokeClient.Invoke(req.CCID, req.Fcn, tool.ChangeArgs(req.Args))
	if err != nil {
		beego.Error("exec error ", err)
	}
	//time.Sleep(10 * time.Millisecond)
}


func (c *InvokeController) Invoke() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	num:=getTxnum()
	h:=sha256.Sum256([]byte(strconv.Itoa(100000000+num)))
	s:=hex.EncodeToString(h[:])
	write([]byte(strconv.Itoa(100000000+num)+"|  "+s+"\n"))
	_,err:=RedisCLient.Put(strconv.Itoa(100000000+num),string(data))
	if err != nil {
		beego.Error("PUT ERROR", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	go DO(num)
	tool.BackResSuccess(c.Controller)
	return
}

func DO(num int)  {
	pool.JobQueue<-&Tx{num}
}