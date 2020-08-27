package test

import (
	"SrvControl/models/db/mgo"
	"SrvControl/models/db/mysql"
	"SrvControl/models/db/redis"
	util "SrvControl/utils"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	util.InitLog()
	mgo.InitMongoDB()
	mysql.InitMysql()
	redis.InitRedisPool()
	tt := time.NewTicker(time.Second) //定时器设置为 1s
	for v := range tt.C {
		logs.Info("aaaaaaaaaaaaaa", v)
	}
}
