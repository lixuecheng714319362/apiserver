package channel

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	"apiserver/models/gosdk"
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
)

type ChanController struct {
	beego.Controller
}

func (c *ChanController) CreateChannel() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()

	res, err := ResClient.CreateChannel(req)
	if err != nil || res.TransactionID == "" {
		beego.Error("create channel failed %s ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("join channel successed", res.TransactionID)
	tool.BackResData(c.Controller, res)
	return
}
func (c *ChanController) CreateNewChannel() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()

	res, err := ResClient.CreateNewChannel(req)
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
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	err = ResClient.JoinChannel(req)
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
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	res, err := ResClient.QueryChannel(req)
	if err != nil {
		beego.Error("query channel err ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("query channel successed", res.Channels)
	tool.BackResData(c.Controller, res.Channels)
	return
}
func getReq(c *ChanController) (*gosdk.ResmgmtRequest, error) {
	data := c.Ctx.Input.RequestBody
	//测试接口使用
	if filter.IsFilterVerify == "false" {
		r := &gosdk.ResmgmtRequest{}
		_ = json.Unmarshal(data, r)
		return r, nil
	}
	req := &filter.ValidateRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	reqData := &gosdk.ResmgmtRequest{}
	err = json.Unmarshal([]byte(req.Data), reqData)
	return reqData, nil
}
