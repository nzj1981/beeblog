package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["IsCategory"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	op := c.Input().Get("op")
	switch op {
	case "add":
		if !checkAccount(c.Ctx) {
			c.Redirect("/login", 302)
			return
		}
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return
	case "del":
		if !checkAccount(c.Ctx) {
			c.Redirect("/login", 302)
			return
		}
		id := c.Input().Get("id")
		// tcn := c.Input().Get("tCn")
		if len(id) == 0 {
			break
		}
		err := models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 301)
		return
	}
	var err error
	id := c.Input().Get("id")
	if len(id) > 0 {
		c.Data["Cate"], err = models.GetCategory(id)
		// fmt.Println("****************", c.Data["Cate"])
		if err != nil {
			beego.Error(err)
		}
		c.Data["Cid"] = id
	}
	c.Data["Categories"], err = models.GetAllCategories(false)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "category.html"
}

// func (c *CategoryController) Modify() {
// 	c.Data["IsCategory"] = true
// 	c.Data["IsLogin"] = checkAccount(c.Ctx)
// 	c.TplName = "category_modify.html"
// 	var err error
// 	id := c.Input().Get("id")
// 	c.Data["Cate"], err = models.GetCategory(id)
// 	// fmt.Println("****************", c.Data["Cate"])
// 	if err != nil {
// 		beego.Error(err)
// 	}
// 	c.Data["Cid"] = id
// 	c.Data["Categories"], err = models.GetAllCategories(false)
// 	if err != nil {
// 		beego.Error(err)
// 	}
// }
func (c *CategoryController) Post() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	id := c.Input().Get("id")
	title := c.Input().Get("name")
	err := models.UpdateCategory(id, title)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/category?id="+id, 301)
	return
}
