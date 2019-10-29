package chaincode

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk"
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
)

type CcController struct {
	beego.Controller
}
type Request struct {
	Data  string
}
var filter =beego.AppConfig.String("filter")


func (c *CcController) InstallChainCode() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
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
	res, err := ResClient.ChainCodeInstall(req)
	if err != nil {
		beego.Error("install chaincode failed", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info(res)
	tool.BackResSuccess(c.Controller)
	return

}

func (c *CcController) InstantiateChainCode() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
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
	res, err := ResClient.ChainCodeInit(req)
	if err != nil {
		beego.Error("chaincode init failed", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info(res)
	tool.BackResData(c.Controller, res)
	return
}


func (c *CcController) UpgradeChainCode() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
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
	res, err := ResClient.ChainCodeUpgrade(req)
	if err != nil {
		beego.Error(err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	tool.BackResData(c.Controller, res)
	return

}

func (c *CcController) QueryInstallChainCode() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
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
	res, err := ResClient.QueryInstalledChaincodes(req)
	if err != nil {
		beego.Error("query installed ChainCode failed", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info(res)
	tool.BackResData(c.Controller, res)
	return
}

func (c *CcController) QueryInstantiateChainCode() {
	defer tool.HanddlerError(c.Controller)
	req,err:=getReq(c)
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
	res, err := ResClient.QueryInstantiatedChaincodes(req)
	if err != nil {
		beego.Error("query instantiated ChainCode failed", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	beego.Info(res)
	tool.BackResData(c.Controller, res)
	return
}

func getReq(c *CcController) (*gosdk.ResmgmtRequest,error)  {
	data := c.Ctx.Input.RequestBody
	//测试接口使用
	if filter=="false"{
		r:=&gosdk.ResmgmtRequest{}
		_=json.Unmarshal(data,r)
		return r,nil
	}
	req := &Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil,err
	}
	reqData :=&gosdk.ResmgmtRequest{}
	err=json.Unmarshal([]byte(req.Data), reqData)
	return reqData,nil
}

