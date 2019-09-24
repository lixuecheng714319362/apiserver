package main

import (
	"apiserver/filter"
	_ "apiserver/routers"
	"github.com/astaxie/beego"
)

func main() {
	level:=beego.AppConfig.String("loglevel")
	switch level {
	case "error":
		beego.SetLevel(beego.LevelError)
	case "info":
		beego.SetLevel(beego.LevelInformational)
	case "debug":
		beego.SetLevel(beego.LevelDebug)
	default:
		beego.SetLevel(beego.LevelInformational)
	}

	beego.InsertFilter("*",beego.BeforeExec,filter.UserFilter)
	beego.Run()
}
