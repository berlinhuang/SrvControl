package mysql

import (
	"SrvControl/utils"
	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 切记：导入驱动包
	"log"
)

var db *sql.DB

func InitMysql() {
	driverName := beego.AppConfig.String("mysql::driverName")
	//注册数据库驱动
	orm.RegisterDriver(driverName, orm.DRMySQL)

	//数据库连接
	user := beego.AppConfig.String("mysql::user")
	pwd := beego.AppConfig.String("mysql::pwd")
	host := beego.AppConfig.String("mysql::host")
	port := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::dbname")
	//dbConn := "root:adf@tcp(127.0.0.1:3306)/cmsproject?charset=utf8
	dbConn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8"
	db1, err := sql.Open(driverName, dbConn)
	if err != nil {
		util.LogError(err.Error())
		return
	}
	db = db1
	logs.Info("MySQL Connected OK")
	CreateTableWithUser()
}

//操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

//查询
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
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
