package chaincode

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	"apiserver/models/gosdk"
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
)

type CcController struct {
	beego.Controller
}


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
	beego.Info("cc install :",req.CCID,res)
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
	beego.Info("query  installed cc success")
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
	beego.Info("query instantiated cc success")
	tool.BackResData(c.Controller, res)
	return
}

func getReq(c *CcController) (*gosdk.ResmgmtRequest,error)  {
	data := c.Ctx.Input.RequestBody
	beego.Debug("request data is ",string(data))
	//测试接口使用
	if filter.IsFilterVerify =="false"{
		r:=&gosdk.ResmgmtRequest{}
		if err:=json.Unmarshal(data,r);err!=nil{
			return nil,err
		}
		return r,nil
	}
	req := &filter.ValidateRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil,err
	}
	reqData :=&gosdk.ResmgmtRequest{}
	err=json.Unmarshal([]byte(req.Data), reqData)
	return reqData,nil
}

