package routers

import (
	"apiserver/controllers/ledger"
	"github.com/astaxie/beego"
)

func init() {
	ns :=
		beego.NewNamespace("/api",
			beego.NSNamespace("/v1",
				beego.NSNamespace("/ledger",
					beego.NSRouter("/info", &ledger.LedgerController{}, "post:QueryInfo"),
					beego.NSRouter("/config", &ledger.LedgerController{}, "post:QueryConfig"),
					beego.NSRouter("/hash", &ledger.LedgerController{}, "post:QueryBlockByHash"),
					beego.NSRouter("/number", &ledger.LedgerController{}, "post:QueryBlockByNumber"),
					beego.NSRouter("/txid", &ledger.LedgerController{}, "post:QueryBlockByTxID"),
					beego.NSRouter("/range", &ledger.LedgerController{}, "post:QueryBlockByRange"),
					beego.NSRouter("/numbertest", &ledger.LedgerController{}, "post:QueryBlockByNumbertest"),
				),

			),
		)

	beego.AddNamespace(ns)
}
