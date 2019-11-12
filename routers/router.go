package routers

import (
	"apiserver/controllers/chaincode"
	"apiserver/controllers/channel"
	"apiserver/controllers/ledger"
	"github.com/astaxie/beego"
)

func init() {
	ns :=
		beego.NewNamespace("/api",
			beego.NSNamespace("/v1",


				beego.NSNamespace("/ledger",
					beego.NSRouter("/info", &ledger.LedgerController{}, "post:QueryInfo"),
					//beego.NSRouter("/config", &ledger.LedgerController{}, "post:QueryConfig"),
					//beego.NSRouter("/hash", &ledger.LedgerController{}, "post:QueryBlockByHash"),
					//beego.NSRouter("/number", &ledger.LedgerController{}, "post:QueryBlockByNumber"),
					//beego.NSRouter("/txid", &ledger.LedgerController{}, "post:QueryBlockByTxID"),
					beego.NSRouter("/range", &ledger.LedgerController{}, "post:QueryBlockByRange"),
					//beego.NSRouter("/numbertest", &ledger.LedgerController{}, "post:QueryBlockByNumbertest"),
				),


				beego.NSNamespace("/chaincode",
					beego.NSRouter("/install", &chaincode.CcController{}, "post:InstallChainCode"),
					beego.NSRouter("/instantiate", &chaincode.CcController{}, "post:InstantiateChainCode"),
					beego.NSRouter("/upgrade", &chaincode.CcController{}, "post:UpgradeChainCode"),
					beego.NSNamespace("/query",
						beego.NSRouter("/install", &chaincode.CcController{}, "post:QueryInstallChainCode"),
						beego.NSRouter("/instantiate", &chaincode.CcController{}, "post:QueryInstantiateChainCode"),
					),
				),


				beego.NSNamespace("/channel",
					beego.NSRouter("/create", &channel.ChanController{}, "post:CreateChannel"),
					beego.NSRouter("/new", &channel.ChanController{}, "post:CreateNewChannel"),
					beego.NSRouter("/join", &channel.ChanController{}, "post:JoinChannel"),
					//beego.NSRouter("/update", &app.AppController{}, "post:UpdateChannel"),
					beego.NSRouter("/query", &channel.ChanController{}, "post:QueryChannel"),
				),


			),
		)

	beego.AddNamespace(ns)
}
