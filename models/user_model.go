package models

import (
	"SrvControl/models/db/mysql"
	"fmt"
)

type User struct {
	Id         int
	Username   string
	Password   string
	Status     int // 0 正常状态， 1删除
	Createtime int64
}

//

//插入用户
func InsertUser(user User) (int64, error) {
	//sql:=fmt.Sprintf("insert into t_user(username, password, status, createtime) values('%s','%s',%d,%d)",user.Username, user.Password, user.Status, user.Createtime)
	//return util.ModifyDB(sql)
	return mysql.ModifyDB("insert into t_user(username, password, status, createtime) values(?,?,?,?)",
		user.Username, user.Password, user.Status, user.Createtime)
}

//按条件查询
func QueryUserWightConn(con string) int {
	sql := fmt.Sprintf("select id from t_user %s", con)
	fmt.Println(sql)
	row := mysql.QueryRowDB(sql)
	id := 0
	row.Scan(&id)
	return id
}

//根据用户名查询id
func QueryUserWithUsername(username string) int {
	sql := fmt.Sprintf("where username = '%s'", username)
	return QueryUserWightConn(sql)
}

//根据用户名和
func QueryUserWithParam(username, password string) int {
	sql := fmt.Sprintf("where username = '%s' and password = '%s'", username, password)
	return QueryUserWightConn(sql)
}
