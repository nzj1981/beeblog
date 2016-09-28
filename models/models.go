package models

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// _DB_NAME        = "data/beeblog.db?charset=utf8&loc=url.QueryEscape(" + "Asia/Shanghai" + ")"
	_DB_NAME        = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

//分类
type Category struct {
	Id              int64
	Title           string    `orm:"null"`
	Created         time.Time `orm:"index;null"`
	Views           int64     `orm:"null"`
	TopicTime       time.Time `orm:"index;null"`
	TopicCount      int64     `orm:"null"`
	TopicLastUserId int64     `orm:"null"`
}

//文章
type Topic struct {
	Id              int64
	Uid             int64     `orm:"null"`
	Title           string    `orm:"null"`
	Content         string    `orm:"size(5000);null"`
	Attachment      string    `orm:"null"`
	Created         time.Time `orm:"index;null"`
	Updated         time.Time `orm:"index;null"`
	Views           int64     `orm:"null"`
	Author          string    `orm:"null"`
	ReplyTime       time.Time `orm:"index;null"`
	ReplyCount      int64     `orm:"null"`
	ReplyLastUserId int64     `orm:"null"`
}

//评论
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"auto_now_add;type(datetime);index"`
}

//`orm:"auto_now;type(datetime)"`
func RegisterDB() {
	//检查数据文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	//注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	//注册驱动
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	//注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)

}

//获取数据库时间值
func GetDate() time.Time {
	var t time.Time
	orm.DefaultTimeLoc = time.UTC
	orm.NewOrm().Raw("select datetime('now','localtime');").QueryRow(&t)
	return t
	//获取数据库时间值
}

//string to int64
func GetInt64(id string) (int64, error) {
	tidNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return tidNum, err
}

//评论
func ReplyAdd(tid, nickname, content string) error {
	tidNum, err := GetInt64(tid)
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: GetDate(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	// 向文章表增加最后回复时间和回复次数 begin
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return err
	}
	topic.ReplyTime = GetDate()
	topic.ReplyCount++
	_, err = o.Update(topic)

	// 向文章表增加最后回复时间和回复次数 end
	return err

}
func RepliesGetAll(tid string, isDesc bool) (replies []*Comment, err error) {
	tidNum, err := GetInt64(tid)
	replies = make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	qs = qs.Filter("tid", tidNum)
	if isDesc {
		_, err = qs.OrderBy("-created").All(&replies)
	} else {
		_, err = qs.All(&replies)
	}
	return replies, err
}
func ReplyDelete(tid, rid string) error {
	tidNum, err := GetInt64(tid)
	ridNum, err := GetInt64(rid)
	o := orm.NewOrm()
	reply := &Comment{Id: ridNum}
	var oldtid int64
	if o.Read(reply) == nil {
		oldtid = reply.Tid
		reply.Id = ridNum
		o.Delete(reply)
	}
	topic := new(Topic)
	if oldtid == tidNum {
		qs := o.QueryTable("topic")
		err = qs.Filter("id", oldtid).One(topic)
		if err != nil {
			return err
		}
		topic.Created = GetDate()
		topic.ReplyCount--
		_, err = o.Update(topic)
	}
	// 更新文章表中最后回复时间和回复总数begin

	// 更新文章表中最后回复时间和回复总数end
	return err
}

// 文章
func TopicAdd(title, content, uid string) error {
	timeNow := GetDate()
	cid, _ := GetInt64(uid)
	o := orm.NewOrm()
	topic := &Topic{
		Uid:     cid,
		Title:   title,
		Content: content,
		Created: timeNow,
		Updated: timeNow,
	}
	//判断重复提交 begin
	qs := o.QueryTable("topic")
	err := qs.Filter("title", title).One(topic)
	if err == nil {
		return err
	}
	//判断重复提交 end
	// fmt.Println(topic)
	_, err = o.Insert(topic)
	//update topic_time,topic_count views begin
	category := new(Category)
	qs = o.QueryTable("category")
	err = qs.Filter("id", cid).One(category)
	if err != nil {
		return err
	}
	category.TopicTime = GetDate()
	category.TopicCount++
	_, err = o.Update(category)
	//update end
	return err
}
func TopicGetAll(cateId string, isDesc bool) ([]*Topic, error) {
	//add topic uid query
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(cateId) > 0 {
			uid, _ := GetInt64(cateId)
			qs = qs.Filter("uid", uid)
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}

func TopicGet(tid string) (*Topic, error) {
	tidNum, _ := GetInt64(tid)
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err := qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_, err = o.Update(topic)
	return topic, err
}
func TopicModify(tid, title, content, uid string) error {
	tidNum, _ := GetInt64(tid)
	cidNum, _ := GetInt64(uid)
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	var oldUid int64
	if o.Read(topic) == nil {
		oldUid = topic.Uid
		topic.Uid = cidNum
		topic.Title = title
		topic.Content = content
		topic.Updated = GetDate()
		o.Update(topic)
	}
	//修改统计文章分类次数、文章分类最后变更时间 begin
	if oldUid != cidNum {
		category := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("id", oldUid).One(category)
		if err != nil {
			return err
		}
		category.TopicTime = GetDate()
		category.TopicCount--
		_, err = o.Update(category)

		err = qs.Filter("id", cidNum).One(category)
		if err != nil {
			return err
		}
		category.TopicTime = GetDate()
		category.TopicCount++
		_, err = o.Update(category)
	}
	//修改统计文章分类次数、文章分类最后变更时间 end
	return nil
}
func TopicDelete(tid string) error {
	tidNum, err := GetInt64(tid)
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	var uid int64
	if o.Read(topic) == nil {
		uid = topic.Uid
		topic.Id = tidNum
		o.Delete(topic)
	}
	// _, err = o.Delete(topic)
	//更新分类中文章统计数 begin
	cate := new(Category)
	if uid > 0 {
		qs := o.QueryTable("category")
		err = qs.Filter("id", uid).One(cate)
		if err != nil {
			return err
		}
		cate.TopicTime = GetDate()
		cate.TopicCount--
		_, err = o.Update(cate)
	}
	//更新分类中文章统计数 end
	return err
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{
		Title:   name,
		Created: GetDate(),
	}
	qs := o.QueryTable("Category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}
func GetAllCategories(isDesc bool) ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("Category")
	var err error
	if isDesc {
		_, err = qs.OrderBy("-Id").All(&cates)
	} else {
		_, err = qs.All(&cates)
	}
	return cates, err
}

//得到一条文章分类信息 begin
func GetCategory(cid string) (*Category, error) {
	cidNum, _ := GetInt64(cid)
	o := orm.NewOrm()
	cate := new(Category)
	qs := o.QueryTable("Category")
	err := qs.Filter("id", cidNum).One(cate)
	if err != nil {
		return nil, err
	}
	return cate, err
}

//得到一条文章分类信息 end
//修改文章分类 begin
func UpdateCategory(tid, title string) error {
	tidNum, _ := GetInt64(tid)
	o := orm.NewOrm()
	cate := &Category{Id: tidNum}
	if o.Read(cate) == nil {
		cate.Id = tidNum
		cate.Title = title
		cate.Created = GetDate()
		o.Update(cate)
	}
	return nil
}

//修改文章分类 end
func DelCategory(id string) error {
	cid, err := GetInt64(id)
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	//文章分类下有文章不充许删除文章分类 begin
	var tCount int64
	if o.Read(cate) == nil {
		tCount = cate.TopicCount
	}
	if tCount == 0 {
		_, err = o.Delete(cate)
		return err
	}
	return err
	//文章分类下有文章不充许删除文章分类 beginr

}
