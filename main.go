package main

import (
	"apiserver/controllers/tool"
	"apiserver/filter"
	_ "apiserver/routers"
	"github.com/astaxie/beego"
)



func main() {
	       logSet()
	       filter.Init()
	       beego.Run()
	}


func logSet()  {

	tool.LogLevel=beego.AppConfig.String("loglevel")
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
}
