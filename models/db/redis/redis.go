package redis

import (
	"SrvControl/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"         //go缓存管理器  支持Memory File Redis Memcached
	_ "github.com/astaxie/beego/cache/redis" //redis缓存引擎
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"syscall"
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
			c, err := redis.Dial(
				"tcp",
				redisAddr+":"+redisPort)
			//password
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
			//	c.Close()
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

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		redisPool.Close()
		os.Exit(0)
	}()
}

func Set(key string, value interface{}, args ...interface{}) ([]byte, error) {
	args_len := len(args)
	new_args := make([]interface{}, args_len+2)
	copy(new_args[2:], args)
	new_args[0] = key
	new_args[1] = value

	conn := redisPool.Get()
	defer conn.Close()
	var data []byte
	data, err := redis.Bytes(conn.Do("SET", new_args...))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

func Mset(args ...interface{}) ([]byte, error) {
	conn := redisPool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("MSET", args...))
	if err != nil {
		return data, fmt.Errorf("error get key: %v", err)
	}
	return data, err
}

func Get(key string) ([]byte, error) {
	conn := redisPool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

// mget 批量读取 mget key1, key2, 返回map结构
func MGet(key ...interface{}) ([]string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	var data []string
	data, err := redis.Strings(conn.Do("MGET", key...))

	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

// DEL
func Del(key string) (int64, error) {
	conn := redisPool.Get()
	defer conn.Close()
	var data int64
	data, err := redis.Int64(conn.Do("DEL", key))

	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, nil
}

//设置过期时间(毫秒级)
func Pexp(key string, timeint int64) (int64, error) {
	conn := redisPool.Get()
	defer conn.Close()
	data, err := redis.Int64(conn.Do("PEXPIRE", key, timeint))

	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, nil
}

//同一个类型的操作只需要一个conn.Receive()接受
//不同的操作需要多个conn.Receive()
//只支持SET功能
func Pipe(mapdata map[string][]map[string]interface{}) (interface{}, error) {
	conn := redisPool.Get()
	defer conn.Close()

	senddata := []interface{}{"test", 1}
	for _, v := range mapdata {
		for _, sdata := range v {
			senddata = []interface{}{}
			for sk, sv := range sdata {
				senddata = append(senddata, sk)
				senddata = append(senddata, sv)
			}
			fmt.Println(senddata)
			conn.Send("SET", senddata...)
		}
	}
	conn.Flush()
	data, err := conn.Receive()
	if err != nil {
		return data, fmt.Errorf("error get key : %v", err)
	}
	return data, nil
}

// ExistsKey
func ExistsKey(key string) (bool, error) {
	rds := redisPool.Get()
	defer rds.Close()
	return redis.Bool(rds.Do("EXISTS", key))
}

// ttl 返回剩余时间
func TTLKey(key string) (int64, error) {
	rds := redisPool.Get()
	defer rds.Close()
	return redis.Int64(rds.Do("TTL", key))
}

// Incr 自增
func Incr(key string) (int64, error) {
	conn := redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("INCR", key))
}

// Decr 自减
func Decr(key string) (int64, error) {
	conn := redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("DECR", key))
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
