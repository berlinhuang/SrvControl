package main

import (
	_ "SrvControl/models"
	"SrvControl/models/db/mgo"
	"SrvControl/models/db/mysql"
	"SrvControl/models/db/redis"
	_ "SrvControl/routers" //_ 表示要执行这个包的init方法
	"SrvControl/utils"
	"github.com/astaxie/beego"
)

func main() {
	util.InitLog()
	mgo.InitMongoDB()
	mysql.InitMysql()
	redis.InitRedisPool()
	//beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run("127.0.0.1:8089")
}
