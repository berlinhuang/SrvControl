package main

import (
	_ "SrvControl/models"
	_ "SrvControl/routers" //_ 表示要执行这个包的init方法
	util "SrvControl/utils"
	"github.com/astaxie/beego"
)


func main() {
	util.InitMysql()
	//beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run("127.0.0.1:8089")
}

