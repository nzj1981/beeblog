package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	name := c.Input().Get("nickname")
	content := c.Input().Get("content")
	err := models.ReplyAdd(tid, name, content)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 302)
}

func (c *ReplyController) Delete() {
	if !checkAccount(c.Ctx) {
		return
	}
	tid := c.Input().Get("tid")
	rid := c.Input().Get("rid")
	err := models.ReplyDelete(tid, rid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 302)
}
