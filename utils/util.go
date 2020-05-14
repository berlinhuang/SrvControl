package util

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io"
	"log"
	//redis缓存引擎
	_ "github.com/astaxie/beego/cache/redis"
	//引入缓存模块
	"github.com/astaxie/beego/cache"
	"os"
	// 切记：导入驱动包
	_ "github.com/go-sql-driver/mysql"
)

const PAGELIMIT = 20


var db *sql.DB
/**
 * 获取redis连接实例
 */
func GetRedis() (adapter cache.Cache, err error) {

	redisKey := beego.AppConfig.String("rediskey")
	redisAddr := beego.AppConfig.String("redisaddr")
	redisPort := beego.AppConfig.String("redisport")
	redisdbNum := beego.AppConfig.String("redisdbnum")

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


func InitMysql(){
	LogInfo()
	driverName := beego.AppConfig.String( "mysql::driverName")
	//注册数据库驱动
	orm.RegisterDriver(driverName, orm.DRMySQL)

	//数据库连接
	user:=beego.AppConfig.String("mysql::user")
	pwd:=beego.AppConfig.String("mysql::pwd")
	host:=beego.AppConfig.String("mysql::host")
	port:=beego.AppConfig.String("mysql::port")
	dbname:=beego.AppConfig.String("mysql::dbname")
	//dbConn := "root:adf@tcp(127.0.0.1:3306)/cmsproject?charset=utf8
	dbConn := user+":"+pwd+"@tcp("+host+":"+port+")/"+dbname+"?charset=utf8"
	db1, err := sql.Open( driverName, dbConn)
	if err!=nil{
		LogError(err.Error())
		return
	}
	db=db1
	LogInfo("连接数据库成功")
	CreateTableWithUser()
}

//操作数据库
func ModifyDB( sql string, args ...interface{})(int64, error){
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil{
		log.Println(err)
		return 0, err
	}
	return count, nil
}

//查询
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}


/**
 * 判断当前path是否存在的工具方法
 */
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 *json格式数据转换为实体对象
 */
func JsonToEntity(data []byte, object interface{}) error {
	if len(data) <= 0 {
		return nil
	}
	return json.Unmarshal(data, object)
}

//向Map中存放数据
func PutParamToMap(mapp map[string]interface{}, key string, value interface{}) map[string]interface{} {
	mapp[key] = value
	return mapp
}


func MD5(str string) string{
	md5str :=fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

func MD5v2( str string) string{
	w := md5.New()
	io.WriteString(w, str) //将str写入到w中
	md5str := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
	return md5str
}


//创建用户表
func CreateTableWithUser() {
	sql := `CREATE TABLE IF NOT EXISTS t_user(
		id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		username VARCHAR(64),
		password VARCHAR(64),
		status INT(4),
		createtime INT(10)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	ModifyDB(sql)
}






/**
 * 根据开发模式进行判断是否输出日志
 */
func LogInfo(v ...interface{}) {

	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		beego.Info(v)
	}
}

/**
 * 错误
 */
func LogError(v ...interface{}) {
	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		beego.Error(v)
	}
}

func LogWarn(v ...interface{}) {
	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		beego.Warn(v)
	}
}

func LogDebug(v ...interface{}) {
	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		beego.Debug(v)
	}
}

func LogNotice(v ...interface{}) {
	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		beego.Notice(v)
	}
}