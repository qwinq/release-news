package controllers

import (
	"classOne/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)


type RegController struct {
	beego.Controller
}
//注册页 register.html
func (rc *RegController)ShowGet(){
	rc.TplName="register.html"
}
/*
1. 拿到浏览器传递数据
2. 处理数据
3. 插入数据库(数据库User)
4. 返回视图
*/
//注册 register.html
func (rc *RegController)HandelPost(){
	//method="post" action="register.html"
	//1. 拿到浏览器传递数据
	name:=rc.GetString("userName") //name="userName"
	pwd:=rc.GetString("password") //name = "password"
	//2. 处理数据
	//beego.Info(name,pwd)
	if name==""||pwd==""{
		beego.Info("用户名或密码不能为空!")
		rc.TplName="register.html"
		return
	}
	//3. 插入数据库(数据库User)
	o:=orm.NewOrm()
	var user models.User
	user.UserName=name
	user.Passwd=pwd
	n,err:=o.Insert(&user)
	if err!=nil{
		beego.Info("orm insert err",err)
		return
	}
	beego.Info("改变数量",n)
	//rc.TplName="login.html"
	rc.Redirect("/",302) //http://192.168.146.132:8080/
	//rc.Ctx.WriteString("insert成功")
}

type LoginController struct {
	beego.Controller
}
//登陆页 login.html
func (lc *LoginController)ShowLogin(){
	uName:=lc.Ctx.GetCookie("userName")
	if uName!=""{
		lc.Data["un"]=uName
		lc.Data["ck"]="checked"
	}
	lc.TplName="login.html"
}
/*
登录业务流程
1. 拿到浏览器数据
2. 判断数据
3. 查找数据库
4. 返回视图
*/
//登陆 login.html
func (lc *LoginController)HandleLogin(){
	//action="/" method="post"
	//1. 拿到浏览器数据
	name:=lc.GetString("userName") //name = "userName"
	pwd:=lc.GetString("password") //name = "password"
	rm:=lc.GetString("remember")
	beego.Info(name,pwd)
	//2. 判断数据
	if name==""||pwd==""{
		beego.Info("用户名或密码不能为空!")
		lc.TplName="register.html"
		return
	}
	//3. 按用户名查找数据库
	o:=orm.NewOrm()
	user:=models.User{}
	user.UserName=name
	err:=o.Read(&user,"UserName") //UserName string `orm:"size(20)"`
	if err!=nil{
		beego.Info("用户名错误",err)
		lc.TplName="login.html"
		return
	}
	//4. 判断密码
	if user.Passwd!=pwd{
		beego.Info("密码错误",err)
		lc.TplName="login.html"
		return
	}
	beego.Info(rm)
	if rm=="on" {
		beego.Info("记录密码")
		lc.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		//删除cookie
		lc.Ctx.SetCookie("userName",name,-1)
	}
	lc.SetSession("userName",name)

	//lc.Ctx.WriteString("登陆成功")
	//beego.Info(user)
	//返回视图
	lc.Redirect("/Article/showArticle",302)
}
