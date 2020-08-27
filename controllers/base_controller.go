package controllers

import (
	"SrvControl/models"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
)

// 约定：如果子controller 存在NestPrepare()方法，就实现了该接口，
type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	beego.Controller
	User      models.User
	IsMobile  bool
	IsLogin   bool
	Loginuser interface{}
}

//判断是否登录
/*
	这个函数主要是为了用户扩展用的，这个函数会在下面定义的这些 Method 方法之前执行，
    用户可以重写这个函数实现类似用户验证之类。
*/
func (this *BaseController) Prepare() {
	loginuser := this.GetSession("loginuser")
	fmt.Println("loginuser---->", loginuser)
	if loginuser != nil {
		this.IsLogin = true
		this.Loginuser = loginuser
	} else {
		this.IsLogin = false
	}
	this.Data["IsLogin"] = this.IsLogin

	// c.AppController 代表接口，可以调用子类的方法 判断当前运行的 Controller 是否是 NestPreparer 实现
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

//
//// check if user not active then redirect
//func (this *BaseController) CheckActiveRedirect(args ...interface{}) bool {
//	var redirect_to string
//	code := 302
//	needActive := true
//	for _, arg := range args {
//		switch v := arg.(type) {
//		case bool:
//			needActive = v
//		case string:
//			// custom redirect url
//			redirect_to = v
//		case int:
//			code = v
//		}
//	}
//	if needActive {
//		// check login
//		if this.CheckLoginRedirect() {
//			return true
//		}
//
//		// redirect to active page
//		if !this.User.IsActive {
//			this.FlashRedirect("/settings/profile", code, "NeedActive")
//			return true
//		}
//	} else {
//		// no need active
//		if this.User.IsActive {
//			if redirect_to == "" {
//				redirect_to = "/"
//			}
//			this.Redirect(redirect_to, code)
//			return true
//		}
//	}
//	return false
//
//}

// check if not login then redirect
func (this *BaseController) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !this.IsLogin {
		if len(redirect_to) == 0 {
			req := this.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/login?to=" + url.QueryEscape(redirect_to)
		this.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && this.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		this.Redirect(redirect_to, code)
		return true
	}
	return false
}
