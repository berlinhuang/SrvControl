package controllers

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	//c.TplName = "index.html"
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}
