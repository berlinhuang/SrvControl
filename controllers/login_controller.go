package controllers

import (
	"SrvControl/models"
	util "SrvControl/utils"
	"fmt"
	"strings"
)

type LoginController struct {
	BaseController
}

//func( this *LoginController ) NestPrepare() {
//
//	if app, ok := this.AppController.(ModelPreparer); ok {
//		app.ModelPrepare()
//		return
//	}
//}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")
	fmt.Println("username:", username, ",password:", password)

	id := models.QueryUserWithParam(username, util.MD5v2(password))
	//fmt.Println("id:", id)
	if id > 0 {
		/*
			设置了session后悔将数据处理设置到cookie，然后再浏览器进行网络请求的时候回自动带上cookie
			因为我们可以通过获取这个cookie来判断用户是谁，这里我们使用的是session的方式进行设置
		*/
		this.SetSession("loginuser", username)
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败"}
	}
	this.ServeJSON()
}

func (this *LoginController) Index() {
	if this.Ctx.Request.Method == "POST" {
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		if len(username) > 0 && len(password) > 0 {
			fmt.Println("username:", username, ",password:", password)
			id := models.QueryUserWithParam(username, util.MD5v2(password))
			fmt.Println("id:", id)
			if id > 0 {
				/*
					设置了session后悔将数据处理设置到cookie，然后再浏览器进行网络请求的时候回自动带上cookie
					因为我们可以通过获取这个cookie来判断用户是谁，这里我们使用的是session的方式进行设置
				*/
				this.SetSession("loginuser", username)
				this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功"}
			} else {
				this.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败"}
			}
			this.ServeJSON()
		}
	}
	this.TplName = "login/index.html"
}
