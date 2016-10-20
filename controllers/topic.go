package controllers

import (
	"beeblog/models"
	"beeblog/pager"
	"fmt"
	"github.com/astaxie/beego"
	"path"
	"strconv"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	/*topics, err := models.TopicGetAll("", "", false)
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Topics"] = topics

	}*/
	//增加分页功能 begin
	pno, _ := c.GetInt("pno") //获取当前请求页
	fmt.Println(pno, "$$$")
	var tTopic []models.Topic
	var conditions string = " order by id desc" //定义文章查询条件
	var po pager.PageOptions
	//定义一个分页对象
	po.TableName = "topic"
	//指定分页表名
	po.EnableFirstLastLink = true
	//是否显示首页尾页默认false
	po.EnablePreNexLink = true
	//是否显示上一页下一页默认false
	po.Conditions = conditions
	//传递分页条件默认值全表
	po.Currentpage = int(pno)
	//传递当前页数，默认为1
	po.PageSize = 10
	//页面大小 默认为20
	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalpages, rs, pagerhtml := pager.GetPagerLinks(&po, c.Ctx)
	rs.QueryRows(&tTopic)           //把当前页面的数据序列化进一个切片内
	c.Data["List"] = tTopic         //把当前页面的数据传递到前台
	c.Data["PagerHtml"] = pagerhtml //分页组件
	c.Data["TotalItem"] = totalItem //总记录数
	// c.Data["PageSize"] = po.PageSize                       //每页记录数
	preItem := (po.Currentpage-1)*po.PageSize + 1 //当前页开始记录数
	nextItem := po.Currentpage * po.PageSize      //当前页结束记录数
	if nextItem > totalItem {
		nextItem = totalItem
	}
	c.Data["PreItem"] = preItem
	c.Data["NextItem"] = nextItem
	c.Data["TotalPages"] = totalpages //总页数
	c.Data["FuncNum"] = funcNum       //每页记录数序列
	c.Data["StrSlice"] = strSlice     //文章标签的转换
	//增加分页功能 end
	c.TplName = "topic.html"
}

//通过当前页开始记录数与当前记录索引数，得出每页的序列数
func funcNum(m, n int) int {
	num := m + n
	return num
}

//文章标签转换slice topic中lables标签
func strSlice(str string) []string {
	str = strings.Replace(strings.Replace(str, "#", " ", -1), "$", "", -1)
	strs := strings.Split(str, " ")
	return strs
}

//获取序列号
/*func getList(pageSize, Currentpage, totalItem int, list []models.Topic) []int {
	nums := make([]int, 0)
	for i := 0; i < pageSize; i++ {
		seq := Currentpage*pageSize - pageSize + i + 1
		if seq > totalItem {

		}
		nums = append(nums, seq, list[i])
	}
	return nums
}*/
func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	//解析表单
	tid := c.Input().Get("tid")
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	uid := c.Input().Get("areaSelect")
	lable := c.Input().Get("lable")
	var err error
	//获取附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	if len(tid) == 0 {
		err = models.TopicAdd(title, content, uid, lable, attachment)
	} else {
		err = models.TopicModify(tid, title, content, uid, lable, attachment)
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
	//获取文章标签分类的数组 begin
	c.Data["Lables"] = strings.Split(topic.Lables, " ")
	//获取文章标签分类的数组 end
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
