package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

// 用户表和文章表是多对多
type User struct {
	Id int
	UserName string
	Password string
	Articles []*Article  `orm:"rel(m2m)"`
}

// 文章表和文章类型表是一对多
type Article struct {
	Id      int `orm:"pk;auto"`
	Title   string  `orm:"size(20)"`  //文章标题
	Content string  `orm:"size(500)"`  //内容
	Img     string  `orm:"size(50);null"`//图片路径
	Time time.Time  `orm:"type(datatime);auto_now_add"`  //发布时间
	Count int  `orm:"dafault(0)"`  //阅读量
	ArticleType *ArticleType `orm:"rel(fk)"`  //
	Users []*User  `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init()  {
	orm.RegisterDataBase("default","mysql","root:root@tcp(127.0.0.1:3306)/newsWeb2?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}



