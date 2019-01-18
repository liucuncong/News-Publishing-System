package routers

import (
	"newsWeb2/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
        )

func init() {
    beego.InsertFilter("/Article/*",beego.BeforeRouter,filterFunc)
    beego.Router("/", &controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/register", &controllers.RegController{},"get:ShowReg;post:HandleReg")
    beego.Router("/Article/ShowArticle", &controllers.ArticleController{},"get:ShowArticleList")
    beego.Router("/Article/AddArticle", &controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    beego.Router("/ArticleContent", &controllers.ArticleController{},"get:ShowArticleContent")
    beego.Router("/Article/DeleteContent", &controllers.ArticleController{},"get:HandleDelete")
    beego.Router("/Article/UpdataArticle", &controllers.ArticleController{},"get:ShowUpdata;post:HandleUptata")
    beego.Router("/Article/AddArticleType", &controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/Article/DeleteArticleType", &controllers.ArticleController{},"get:HandleDeleteArticleType")
    beego.Router("/Article/Logout", &controllers.ArticleController{},"get:Logout")
}


var filterFunc = func(ctx *context.Context) {
   userName := ctx.Input.Session("userName")
   if userName == nil{
       ctx.Redirect(302,"/",)  // 如果没有session，回到登陆界面
   }
}








