package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb2/models"
	"strconv"
	"math"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}

// 文章列表页显示
func (this *ArticleController)ShowArticleList()  {
	// 1.接受数据
	typeName := this.GetString("select")
	// 2.处理数据
	if typeName == ""{
		// 1.查询
		o := orm.NewOrm()
		qs:=o.QueryTable("Article")
		var articles []models.Article
		//qs.All(&articles)

		// 查询数据条目数
		count,err := qs.RelatedSel("ArticleType").Count()
		if err != nil {
			beego.Info("查询条目数错误")
			return
		}
		// 获取总页数
		pageSize := 2
		pageCount := float64(count)/float64(pageSize)
		pageCount2 := math.Ceil(pageCount)
		if err != nil {
			beego.Info("获取页码失败")
			return
		}
		// 每页显示内容设置
		pageIndex := this.GetString("pageIndex")
		pageIndex2,err := strconv.Atoi(pageIndex)
		if err != nil {
			pageIndex2 = 1
		}

		start := pageSize*(pageIndex2-1)
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)


		// 存储types
		var types []models.ArticleType
		// 从redis中获取文章类型数据
		conn,err := redis.Dial("tcp",":6379")
		if err != nil {
			beego.Info("redis数据库连接失败")
			return
		}
		rel,err := redis.Bytes(conn.Do("get","types"))
		if err != nil {
			beego.Info("获取redis数据错误")
			return
		}

		dec := gob.NewDecoder(bytes.NewReader(rel))
		dec.Decode(&types)
		beego.Info(types)

		if len(types) == 0 {
			// 从mysql中获取文章类型数据
			o.QueryTable("ArticleType").All(&types)
			// 把类型存入redis数据库
			var buffer bytes.Buffer
			enc := gob.NewEncoder(&buffer)
			enc.Encode(&types)

			_,err = conn.Do("set","types",buffer.Bytes())
			if err != nil {
				beego.Info("redis数据库操作失败")
				return
			}
			beego.Info("从mysql数据库中取数据")
		}


		// 获取用户名
		userName := this.GetSession("userName")

		// 2.把数据传递给视图
		this.Data["types"] = types
		this.Data["count"] = count
		this.Data["typeName"] = typeName
		this.Data["pageCount"] = pageCount2
		this.Data["pageIndex"] = pageIndex2
		this.Data["articles"] = articles
		this.Data["userName"] = userName.(string)

		// 返回视图
		this.Layout = "layout.html"
		this.TplName = "index.html"
	}else {
		// 3.查询数据
		o := orm.NewOrm()
		var articles []models.Article
		qs:=o.QueryTable("Article")
		//qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
		// 查询数据条目数
		count,err := qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		if err != nil {
			beego.Info("查询条目数错误")
			return
		}
		// 获取总页数
		pageSize := 2
		pageCount := float64(count)/float64(pageSize)
		pageCount2 := math.Ceil(pageCount)
		if err != nil {
			beego.Info("获取页码失败")
			return
		}
		// 每页显示内容设置
		pageIndex := this.GetString("pageIndex")
		pageIndex2,err := strconv.Atoi(pageIndex)
		if err != nil {
			pageIndex2 = 1
		}

		start := pageSize*(pageIndex2-1)
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

		//// 获取类型数据
		//var types []models.ArticleType
		//o.QueryTable("ArticleType").All(&types)

		// 存储types
		var types []models.ArticleType
		// 从redis中获取文章类型数据
		conn,err := redis.Dial("tcp",":6379")
		if err != nil {
			beego.Info("redis数据库连接失败")
			return
		}
		rel,err := redis.Bytes(conn.Do("get","types"))
		if err != nil {
			beego.Info("获取redis数据错误")
			return
		}

		dec := gob.NewDecoder(bytes.NewReader(rel))
		dec.Decode(&types)
		beego.Info(types)


		// 获取用户名
		userName := this.GetSession("userName")

		// 4.把数据传递给视图
		this.Data["userName"] = userName.(string)
		this.Data["types"] = types
		this.Data["count"] = count
		this.Data["typeName"] = typeName
		this.Data["pageCount"] = pageCount2
		this.Data["pageIndex"] = pageIndex2
		this.Data["articles"] = articles

		// 返回视图
		this.Layout = "layout.html"
		this.TplName = "index.html"
	}



}

// 添加文章页面显示
func (this *ArticleController)ShowAddArticle()  {
	// 1.数据库获取数据
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	// 2.传递数据
	this.Data["types"] = types
	this.TplName = "add.html"
}

// 添加文章页面处理
func (this *ArticleController)HandleAddArticle()  {
	//1.接受数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	f,h,err := this.GetFile("uploadname")
	defer f.Close()
	if err != nil {
		beego.Info("上传文件失败")
		return
	}
	//2.判断数据
	if typeName == ""{
		beego.Info("下拉框数据不能为空")
		return
	}
	//2.1判断文件格式
	ext := path.Ext(h.Filename)
	if ext!=".jpg" && ext!=".jpeg" && ext!=".png" {
		beego.Info("上传文件格式错误")
		return
	}
	//2.2判断文件大小
	if h.Size > 5000000 {
		beego.Info("上传文件太大")
		return
	}
	//2.3不能重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")

	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//3.插入数据库
	o := orm.NewOrm()
	article := models.Article{}
	article.Title = articleName
	article.Content = content
	article.Img = "./static/img/"+fileName+ext

	// 插入文章类型
	// 获取type对象
	articleType := models.ArticleType{TypeName:typeName}
	err = o.Read(&articleType,"TypeName")
	if err != nil {
		beego.Info("获取文章类型错误")
		return
	}
	article.ArticleType = &articleType

	_,err = o.Insert(&article)
	if err != nil {
		beego.Info("插入失败")
		return
	}

	//4.返回视图
	this.Redirect("/Article/ShowArticle",302)
}

// 显示文章详情页
func (this *ArticleController)ShowArticleContent()  {
	// 1.获取Id
	id := this.GetString("id")
	beego.Info(id)
	id2,err := strconv.Atoi(id)
	if err != nil {
		beego.Info("获取Id错误")
		return
	}

	// 2.查询数据
	o := orm.NewOrm()
	article := models.Article{Id:id2}
	err = o.Read(&article)
	if err != nil {
		beego.Info("查询数据错误")
		return
	}
	// 把count加1
	article.Count+=1
	o.Update(&article)

	// 添加阅读人，多对多插入
	// 2.1获取操作对象（已获取）

	// 2.2获取多对多操作对象
	m2m := o.QueryM2M(&article,"Users")
	// 2.3获取插入对象
	userName := this.GetSession("userName")
	user := models.User{}
	user.UserName = userName.(string)
	o.Read(&user,"UserName")
	// 2.4多对多插入
	_,err = m2m.Add(&user)
	if err != nil {
		beego.Info("插入阅读人错误")
		return
	}
	o.Update(&article)
	 // 第一种多对多查询方法，会有重复
	//o.LoadRelated(&article,"Users")
	// 第一种多对多查询方法
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id2).Distinct().All(&users)
	// 3.传递数据
	this.Data["article"] = article
	this.Data["users"] = users
	// 4.返回视图
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "head.html"
	this.TplName = "content.html"
}

// 处理删除业务
func (this *ArticleController)HandleDelete()  {
	// 1.获取id
	id := this.GetString("id")
	id2,err := strconv.Atoi(id)
	if err != nil {
		beego.Info("获取Id失败")
	}

	// 2.删除操作
	o := orm.NewOrm()
	article := models.Article{Id:id2}
	o.Delete(&article)

	// 3.重定向到列表页
	this.Redirect("/Article/ShowArticle",302)
}

// 编辑更新页面显示
func (this *ArticleController)ShowUpdata()  {
	// 1.获取id
	id := this.GetString("id")
	id2,err := strconv.Atoi(id)
	if err != nil {
		beego.Info("id错误")
		return
	}
	// 2.查询数据
	o := orm.NewOrm()
	article := models.Article{Id:id2}
	err = o.Read(&article)
	if err != nil {
		beego.Info("查询错误")
		return
	}
	// 3.传递数据
	this.Data["article"] = article

	// 4.展示数据
	this.TplName = "update.html"
}

// 处理编辑更新页面
func (this *ArticleController)HandleUptata()  {
	// 1.获取数据
	id := this.GetString("id")
	id2,err := strconv.Atoi(id)
	if err != nil {
		beego.Info("id错误")
		return
	}
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	f,h,err := this.GetFile("uploadname")
	defer f.Close()
	if err != nil {
		beego.Info("获取文件错误")
		return
	}
	// 2.判断数据
	if articleName == "" || content == "" {
		beego.Info("文章标题或内容不能为空")
		return
	}
	// 判断文件类型
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Info("文件类型错误")
		return
	}
	// 判断文件大小
	if h.Size > 5000000 {
		beego.Info("文件过大")
		return
	}
	// 文件名不重复
	fileName := time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	// 3.查询数据
	o := orm.NewOrm()
	article := models.Article{Id:id2}
	err = o.Read(&article)
	if err != nil {
		beego.Info("查询数据错误")
		return
	}
	// 4.更新数据
	article.Title = articleName
	article.Content = content
	article.Img = "./static/img/"+fileName+ext
	_,err =o.Update(&article)
	if err != nil {
		beego.Info("更新数据错误")
		return
	}
	// 5.重定向到列表页
	this.Redirect("/Article/ShowArticle",302)
}

// 添加文章类型显示
func (this *ArticleController)ShowAddType()  {
	// 1.数据库查询数据
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		beego.Info("查询类型错误")
	}
	// 2.传递数据
	this.Data["articleTypes"] = articleTypes
	// 3.返回视图
	this.TplName = "addType.html"

}

// 添加文章类型处理
func (this *ArticleController)HandleAddType()  {
	// 1.获取数据
	typeName := this.GetString("typeName")
	// 2.判断数据
	if typeName == ""{
		 beego.Info("添加数据类型为空")
		return
	}
	// 3.把数据插入数据库
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_,err := o.Insert(&articleType)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	// 4.重定向返回视图
	this.Redirect("/Article/AddArticleType",302)
}

// 删除文章类型
func (this *ArticleController)HandleDeleteArticleType()  {
	// 1.获取数据
	id := this.GetString("id")
	id2,err := strconv.Atoi(id)
	if err != nil {
		beego.Info("获取id错误")
		return
	}
	// 2.从数据库删除数据
	o := orm.NewOrm()
	articleType := models.ArticleType{Id:id2}
	o.Delete(&articleType)
	// 3.重定向视图
	this.Redirect("/Article/AddArticleType",302)
}

// 退出登陆
func (this *ArticleController)Logout()  {
	// 1.删除登陆状态
	this.DelSession("userName")
	// 2.跳转到登陆页面
	this.Redirect("/",302)
}

