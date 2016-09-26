package main

import (
	"beeblog/models"
	_ "beeblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	//注册数据库
	models.RegisterDB()
}
func main() {
	//开启ORM调式模式
	orm.Debug = true
	//time
	// orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Shanghai")
	//改变日志输出到自己的io.Writer
	// var w io.Writer
	// orm.DebugLog = orm.NewLog(w)
	//自动建表
	orm.RunSyncdb("default", false, true)
	//启动beego
	beego.Run()
}
