package routers

import (
	"classOne/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

    //beego.Router("/", &controllers.MainController{})
    beego.InsertFilter("/Article/*",beego.BeforeRouter,filterFunc)
	beego.Router("/register",&controllers.RegController{},"get:ShowGet;post:HandelPost")
    beego.Router("/",&controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/logout",&controllers.ArticleController{},"get:Logout")
	beego.Router("/Article/showArticle",&controllers.ArticleController{},"get:ShowArticleList;post:HandleSelect")
	beego.Router("/Article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/showArticleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
	beego.Router("/Article/deleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/Article/updateArticle",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/Article/addArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/Article/deleteArticleType",&controllers.ArticleController{},"get:HandleDeleteType")
}

var filterFunc= func(ctx *context.Context) {
	un:=ctx.Input.Session("userName")
	if un==nil{
		ctx.Redirect(302,"/")
	}
}