package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb2/models"
	"time"
)

// 注册页面显示
type RegController struct {
	beego.Controller
}

// 注册页面显示
func (this *RegController)ShowReg()  {
	this.TplName = "register.html"
}

// 注册页面处理
func (this *RegController)HandleReg()  {
	// 1.接收数据
	name := this.GetString("userName")
	password := this.GetString("password")
	// 2.判断数据
	if name == "" || password == ""{
		beego.Info("用户名或密码不能为空")
		this.TplName = "register.html"
		return
	}
	// 3.插入数据库
	// 3.1获取orm对象
	o := orm.NewOrm()
	// 3.2获取插入对象
	user := models.User{}
	user.UserName = name
	user.Password = password
	_,err := o.Insert(&user)
	if err != nil {
		beego.Info("插入数据错误",err)
		this.TplName = "register.html"
		return
	}
	// 4.返回视图
	this.Redirect("/",302)
}


// 登陆页面控制器
type LoginController struct {
	beego.Controller
}

// 登陆页面显示
func (this *LoginController)ShowLogin()  {
	// 获取cookie
	name := this.Ctx.GetCookie("userName")
	if name != ""{
		this.Data["name"] = name
		this.Data["check"] = "checked"
	}
	this.Data["name"] = name
	this.TplName = "login.html"
}

// 登陆页面处理
func (this *LoginController)HandleLogin()  {
	// 1.获取数据
	name := this.GetString("userName")
	password := this.GetString("password")
	// 2.判断数据
	if name == "" || password == ""{
		beego.Info("用户名或密码不能为空")
		this.TplName = "login.html"
		return
	}
	// 3.查找数据库
	o :=orm.NewOrm()
	user := models.User{}
	user.UserName = name
	err := o.Read(&user,"UserName")
	if err != nil {
		beego.Info("登录失败")
		this.TplName = "login.html"
		return
	}
	if user.Password != password {
		beego.Info("登录失败")
		this.TplName = "login.html"
		return
	}

	//设置cookie，记住用户名
	check := this.GetString("remember")
	if check == "on" {
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else {
		this.Ctx.SetCookie("userName","ss",-1)
	}
	this.SetSession("userName",name)

	// 4.返回视图
	this.Redirect("/Article/ShowArticle",302)
}