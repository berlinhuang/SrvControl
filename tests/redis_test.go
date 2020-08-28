package test

import (
	"SrvControl/models/db/redis"
	util "SrvControl/utils"
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	util.InitLog()
	redis.InitRedisPool()
	// set 数据
	v, err := redis.Set("test", 3, "EX", "180")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := string(v[:])
	fmt.Println(data)

	// mset 数据
	v1, err := redis.Mset("test", 2, "test1", "2")
	if err != nil {
		fmt.Println(err)
		return
	}
	data1 := string(v1[:])
	fmt.Println(data1)

	// get 数据
	v2, err := redis.Get("pool")
	data2 := string(v2[:])
	fmt.Println(data2, err)

	// MGET 数据
	v3, err := redis.MGet("test", "test1")
	fmt.Println(v3, err)

	//批量获取数据
	//key_data := []string{}
	//fmt.Println(key_data)
	//v4, err := redis.MGet()
	//data4 := string(v4[:])
	//fmt.Println(data4, err)

	//删除数据
	v5, err := redis.Del("test1")
	fmt.Println(v5, err)

	//修改过期时间（毫秒级）
	v6, err := redis.Pexp("test1", 100000)
	if v6 == 1 {
		fmt.Println("修改成功")
	}
	fmt.Println(v6, err)

}
