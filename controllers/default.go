package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beeblog.com"
	c.Data["Email"] = "18322313385@163.com"
	c.Data["IsHome"] = true
	c.TplName = "index.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	topics, err := models.TopicGetAll(c.Input().Get("cate"), c.Input().Get("lable"), true)
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Topics"] = topics
	}
	//文章分类 begin
	cates, err := models.GetAllCategories(true)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Cates"] = cates
	//文章分类 end
}
