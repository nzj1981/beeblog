package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
	"strconv"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	topics, err := models.TopicGetAll("", false)
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Topics"] = topics
	}
	c.TplName = "topic.html"
}

func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	tid := c.Input().Get("tid")
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	uid := c.Input().Get("areaSelect")
	var err error
	if len(tid) == 0 {
		err = models.TopicAdd(title, content, uid)
	} else {
		err = models.TopicModify(tid, title, content, uid)
	}
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
	return
}
func (c *TopicController) Add() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	var err error
	c.Data["Categories"], err = models.GetAllCategories(true)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "topic_add.html"
}
func (c *TopicController) View() {
	c.TplName = "topic_view.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	// fmt.Println("*************************", c.Ctx.Input.Param("0"), "++++++++", c.Ctx.Input.Param("1"))
	tid := c.Ctx.Input.Param("0")
	topic, err := models.TopicGet(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	//根据文章分类id获取文章分类标题 begin
	uid := strconv.FormatInt(topic.Uid, 10)
	c.Data["Cate"], err = models.GetCategory(uid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	//根据文章分类id获取文章分类标题 end
	c.Data["Topic"] = topic
	// c.Data["Tid"] = c.Ctx.Input.Param("0")
	// c.Data["Uid"] = c.Ctx.Input.Param("1")

	// 获取评论回复begin
	replies, err := models.RepliesGetAll(tid, true)
	if err != nil {
		beego.Error(err)
		return
	}
	c.Data["Replies"] = replies
	// 获取评论回复end
}
func (c *TopicController) Modify() {
	c.TplName = "topic_modify.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	tid := c.Input().Get("tid")
	topic, err := models.TopicGet(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Categories"], err = models.GetAllCategories(true)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}
func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	err := models.TopicDelete(c.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
	return
}
