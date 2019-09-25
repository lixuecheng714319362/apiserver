package main

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	_ "apiserver/routers"
	"github.com/astaxie/beego"
)



func main() {
	tool.LogLevel=beego.AppConfig.String("loglevel")
	isFilter:=beego.AppConfig.String("filter")
	switch tool.LogLevel {
	case "error":
		beego.SetLevel(beego.LevelError)
	case "info":
		beego.SetLevel(beego.LevelInformational)
	case "debug":
		beego.SetLevel(beego.LevelDebug)
	default:
		beego.SetLevel(beego.LevelInformational)
	}
	if isFilter!="false"{
		beego.InsertFilter("*",beego.BeforeExec,filter.UserFilter)
	}

	beego.Run()
}
