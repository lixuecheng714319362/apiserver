package channel

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk"
	"encoding/json"
	"github.com/astaxie/beego"

	"net/http"
)

var (
	Org1MSP = "CopyRightChain1MSP"
)

type ChanController struct {
	beego.Controller
}

type Request struct {
	ConfigPath string
	UserName   string
	OrgName    string

	ChannelID string

	ChannelTxPath string

	TargetOrderer string //OrdID

	TargetPeer string //查询单个peer

	TargetPeers []string //defatut安装过程中查询所有peer是否安装过

}

func (c *ChanController) CreateChannel() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req.ConfigPath, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
	}
	defer ResClient.CloseSDK()

	res, err := ResClient.CreateChannel(req.ChannelID, req.ChannelTxPath, req.TargetOrderer)
	if err != nil || res.TransactionID == "" {
		beego.Error("create channel failed %s ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("join channel successed", res.TransactionID)
	tool.BackResData(c.Controller, res)
	return
}

func (c *ChanController) JoinChannel() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req.ConfigPath, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	err = ResClient.JoinChannel(req.TargetOrderer, req.ChannelID, req.TargetPeers)
	if err != nil {
		beego.Error("join channel failed %s ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("Join channel successed")
	tool.BackResSuccess(c.Controller)
	return
}

func (c *ChanController) QueryChannel() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req.ConfigPath, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	res, err := ResClient.QueryChannel(req.TargetPeer)
	if err != nil {
		beego.Error("query channel err ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query channel successed", res.Channels)
	tool.BackResData(c.Controller, res.Channels)
	return

}

func (c *ChanController) CreateAndJoinChannel() {
	defer tool.HanddlerError(c.Controller)
	data := c.Ctx.Input.RequestBody
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req.ConfigPath, req.UserName, req.OrgName)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	err = ResClient.CreateAndJoinChannel(req.ChannelID, req.ChannelTxPath, req.TargetOrderer, req.TargetPeers)
	if err != nil {
		beego.Error("create and join channel err ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("create and join channel successed ")
	tool.BackResSuccess(c.Controller)
	return
}
