package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"path/filepath"
)

const PAGELIMIT = 20

func InitLog() {
	exepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exedir := filepath.Dir(exepath)
	//logLevel := beego.AppConfig.String("log_level")
	//logPath := beego.AppConfig.String("log_path")
	logConf := make(map[string]interface{})
	logConf["filename"] = fmt.Sprintf("%v\\logs\\srvctl.log", exedir)
	logConf["level"] = logs.LevelDebug
	logConfJson, err := json.Marshal(logConf)
	if err != nil {
		panic(err)
	}
	logs.Async()
	//输出文件名和位置
	logs.EnableFuncCallDepth(true)
	// log 的配置
	logs.SetLogger(logs.AdapterFile, string(logConfJson))
	// log打印文件名和行数
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)

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

func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

func MD5v2(str string) string {
	w := md5.New()
	io.WriteString(w, str)                  //将str写入到w中
	md5str := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	return md5str
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
