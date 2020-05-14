package routers

import (
	"SrvControl/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

    beego.Router("/chatroom", &controllers.ServerController{})
	beego.Router("/chatroom/WS", &controllers.ServerController{}, "get:WS")


    //主页
	beego.Router("/blog", &controllers.HomeController{})
	//注册
    beego.Router("/register", &controllers.RegisterController{})
    //登录
    beego.Router("/login", &controllers.LoginController{})
    //退出
    beego.Router("/exit", &controllers.ExitController{})
}
