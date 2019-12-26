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
//根据tx文件创建channel
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
		beego.Error("create channel failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("join channel successed", res.TransactionID)
	tool.BackResData(c.Controller, res)
	return
}
//使用orgName创建依据系统通道的channel
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
		beego.Error("create channel failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("join channel successed", res.TransactionID)
	tool.BackResData(c.Controller, res)
	return
}
//根据组织信息将新组织添加到channel中
func (c *ChanController) AddOrgUpdateChannel() {
	defer tool.HanddlerError(c.Controller)
	req, err := getReq(c)
	if err != nil {
		beego.Error("request json unmarshal failed ", err)
		tool.BackResError(c.Controller, http.StatusForbidden, err.Error())
		return
	}

	reader ,err:= gosdk.GetAddOrgChannelConfigUpdate(req)
	if err != nil {
		beego.Error("get  add org channel config update reader failed", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	ResClient, err := gosdk.GetResMgmtClient(req)
	if err != nil {
		beego.Error("create resClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer ResClient.CloseSDK()
	res, err := ResClient.UpdateChannel(req,reader)
	if err != nil || res.TransactionID == "" {
		beego.Error("update channel failed  ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info("add org update channel success", res.TransactionID)
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
	beego.Debug("request data is ",string(data))
	//测试接口使用
	if filter.IsFilterVerify == "false" {
		r := &gosdk.ResmgmtRequest{}
		if err:=json.Unmarshal(data,r);err!=nil{
			return nil,err
		}
		gosdk.ChangeResmgmtRequetSingleConfig(r)
		return r, nil
	}
	req := &filter.ValidateRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	reqData := &gosdk.ResmgmtRequest{}
	err = json.Unmarshal([]byte(req.Data), reqData)
	gosdk.ChangeResmgmtRequetSingleConfig(reqData)
	return reqData, nil
}
