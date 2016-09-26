package routers

import (
	"beeblog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//注册beego路由
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/category", &controllers.CategoryController{})
	// beego.Router("/category/modify", &controllers.CategoryController{}, "get:Modify")
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
}