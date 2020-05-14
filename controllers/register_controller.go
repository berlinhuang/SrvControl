package controllers

import (
	"SrvControl/models"
	"SrvControl/utils"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

type RegisterController struct {
	beego.Controller
}


func ( this *RegisterController) Get(){
	this.TplName = "register.html"
}


func( this *RegisterController) Post(){
	// 获取表单信息
	username := this.GetString("username")
	password := this.GetString("password")
	repassword := this.GetString("repassword")
	fmt.Println( username, password, repassword )
	util.LogInfo( username, password, repassword)

	//注册之前先判断该用户是否已被注册，如果已注册，则返回错误
	id := models.QueryUserWithUsername(username)
	fmt.Println("id:",id)
	if id > 0 {
		this.Data["json"] = map[string]interface{}{"code":0,"message":"用户名已经存在"}
		this.ServeJSON()
		return
	}

	// 注册用户名和密码
	password = util.MD5v2(password)
	fmt.Println("md5后：", password)

	user:= models.User{0,username,password,0,time.Now().Unix()}
	_,err:=models.InsertUser(user)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code":0,"message":"注册失败"}
	}else{
		this.Data["json"]=map[string]interface{}{"code":1,"message":"注册成功"}
	}
	//ServeJSONP（）进行渲染，会设置内容类型为application / javascript，然后同时把数据进行JSON序列化，然后根据请求的回调参数设置jsonp输出。
	this.ServeJSONP()
	//重定向到登录页面
	this.Redirect("/login",302)
}