package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)
import _"github.com/go-sql-driver/mysql"

type User struct {
	Id int `orm:"pk;auto"`
	UserName string `orm:"size(20)"`
	Passwd string `orm:"size(20)"`
	Articles []*Article `orm:"rel(m2m)"`
}
//自动创建 user_articles
//文章表与文章类型表为 1对多
type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(20)"`
	Content string `orm:"size(500)"`
	Img string `orm:"size(50);null"`
	DateTime time.Time `orm:"auto_now,type(datatime);"`
	Count int `orm:"default(0);null"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"reverse(many)"`
}
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles[] *Article `orm:"reverse(many)"`
}
func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/qwindb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}