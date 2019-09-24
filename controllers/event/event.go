package event

import (
	"apiserver/controllers/tool"
	"apiserver/models/gosdk"
	"github.com/astaxie/beego"

	"net/http"
	"time"
)

type EvController struct {
	beego.Controller
}

var (
	Org1MSP = "CopyRightChain1MSP"
)

var (
	configPath  = "conf/config/config.yaml"
	OrgAdmin    = "Admin"
	OrgName     = "copyrightchain1"
	ChannelID   = "copyrightchainchannel"
	UserName    = "Admin"
	Org1        = "copyrightchain1"
	ChainCodeID = "copyrightchain"
	ConfigPath  = "../conf/config/config.yaml"

	PreviousBlockHash string
	TxID              string
)

func (c *EvController) RegisterChaincodeEvent() {
	EventClient, err := gosdk.GetEventClient(configPath, OrgAdmin, ChannelID)
	if err != nil {
		beego.Error("create EventClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.CloseSDK()
	reg, ccEvent, err := EventClient.RegisterChaincodeEvent(ChainCodeID)
	if err != nil {
		beego.Error("RegisterChaincodeEvent info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.Unregister(reg)
	for {
		select {
		case event := <-ccEvent:
			//TODO
			beego.Info(event)
		case <-time.After(50 * time.Second):
			tool.BackResTimeOut(c.Controller)
			return
		}
	}
}

func (c *EvController) RegisterBlockEvent() {
	EventClient, err := gosdk.GetEventClient(configPath, OrgAdmin, ChannelID)
	if err != nil {
		beego.Error("create GetEventClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.CloseSDK()
	//TODO
	//调用会报错
	reg, ccEvent, err := EventClient.RegisterBlockEvent(ChainCodeID)
	if err != nil {
		beego.Error("RegisterBlockEvent info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.Unregister(reg)
	for {
		select {
		case event := <-ccEvent:
			beego.Info(event)
		case <-time.After(50 * time.Second):
			tool.BackResTimeOut(c.Controller)
			return
		}
	}
}

func (c *EvController) RegisterFilteredBlockEvent() {
	EventClient, err := gosdk.GetEventClient(configPath, OrgAdmin, ChannelID)
	if err != nil {
		beego.Error("create GetEventClient failed ", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.CloseSDK()
	reg, ccEvent, err := EventClient.RegisterFilteredBlockEvent()
	if err != nil {
		beego.Error("RegisterFilteredBlockEvent  info err", err)
		tool.BackResError(c.Controller, http.StatusBadRequest, err.Error())
		return
	}
	defer EventClient.Unregister(reg)
	for {
		select {
		case event := <-ccEvent:
			beego.Info(event)
		case <-time.After(50 * time.Second):
			tool.BackResTimeOut(c.Controller)
			return
		}
	}
}

//
//func (c *EvController )RegisterTxStatusEvent () {
//	EventClient, err := gosdk.GetEventClient(configPath, UserName, ChannelID)
//	if err != nil {
//		beego.Error("create GetEventClient failed ", err)
//		tool.BackResError(c.Controller,http.StatusBadRequest,err.Error())
//		return
//	}
//	defer EventClient.CloseSDK()
//
//	for  {
//
//		select {
//		case txid:=<-Txchan:
//			go func(txid string) {
//				beego.Info("start listen txid",txid)
//				reg,ccEvent,err:=EventClient.RegisterTxStatusEvent(txid)
//				if err != nil {
//					beego.Error("RegisterTxStatusEvent info err",err)
//					tool.BackResError(c.Controller,http.StatusBadRequest,err.Error())
//					return
//				}
//				defer EventClient.Unregister(reg)
//				for  {
//					select {
//					case event:=<-ccEvent:
//						beego.Info("get event" ,txid)
//						beego.Info(event)
//					case <-time.After(50 * time.Second):
//						beego.Info("do not listen txid",txid)
//						return
//					}
//				}
//			}(txid)
//		case <-time.After(200 * time.Second):
//			return
//		}
//
//	}
//}
