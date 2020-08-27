package redis

import (
	"SrvControl/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"         //go缓存管理器  支持Memory File Redis Memcached
	_ "github.com/astaxie/beego/cache/redis" //redis缓存引擎
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPool *redis.Pool
)

func InitRedisPool() {
	maxIdle := beego.AppConfig.DefaultInt("redis::redis_maxidle", 20)
	maxActive := beego.AppConfig.DefaultInt("redis::redis_maxactive", 100)
	maxIdleTimeout := beego.AppConfig.DefaultInt64("redis::redis_maxIdleTimeout", 180)
	redisAddr := beego.AppConfig.String("redis::redisaddr")
	redisPort := beego.AppConfig.String("redis::redisport")
	//redisPassword := beego.AppConfig.String("redis::redis_password")
	//初始化连接池 信息
	redisPool = &redis.Pool{
		MaxIdle:     maxIdle,   //最大空闲连接
		MaxActive:   maxActive, // 最大激活连接
		IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr+":"+redisPort)//password
			//database
			//connection timeout
			//read timeout
			//write timeout

			if err != nil {
				util.LogError(err.Error())
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证密码
			//if _, err :=conn.Do("AUTH", redisPassword);err != nil {
			//	LogError(err)
			//	return nil, fmt.Errorf("redis auth password error: %s", err)
			//}
			return c, err
		},
		//TestOnBorrow 是一个测试链接可用性的方法
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	logs.Info("初始化redis成功")
}

func Set(k, v string) {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func GetStringValue(k string) string {
	c := redisPool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		fmt.Println("Get Error: ", err.Error())
		return ""
	}
	return username
}

func SetKeyExpire(k string, ex int) {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func CheckKey(k string) bool {
	c := redisPool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return exist
	}
}

func DelKey(k string) error {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SetJson(k string, data interface{}) error {
	c := redisPool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func getJsonByte(k string) ([]byte, error) {
	c := redisPool.Get()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonGet, nil
}

/**
 * 获取redis连接实例
 */
func GetRedis() (adapter cache.Cache, err error) {
	redisKey := beego.AppConfig.String("redis::rediskey")
	redisAddr := beego.AppConfig.String("redis::redisaddr")
	redisPort := beego.AppConfig.String("redis::redisport")
	redisdbNum := beego.AppConfig.String("redis::redisdbnum")

	redis_config_map := map[string]string{
		"key":   redisKey,
		"conn":  redisAddr + ":" + redisPort,
		"dbNum": redisdbNum,
	}
	redis_config, _ := json.Marshal(redis_config_map) //字符串

	cache_conn, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		return nil, err
	}
	return cache_conn, nil
}

func Redisget() {
	c1 := redisPool.Get()
	c2 := redisPool.Get()
	c3 := redisPool.Get()
	c4 := redisPool.Get()
	c5 := redisPool.Get()
	fmt.Println(c1, c2, c3, c4, c5)
}
