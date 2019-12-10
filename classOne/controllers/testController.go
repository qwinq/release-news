package controllers

import "github.com/astaxie/beego"

type MainController struct {
	beego.Controller
}
func (mc *MainController)Get(){
	mc.TplName="login.html"
}