package controllers

import (
	"bytes"
	"classOne/models"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	"time"
	"github.com/gomodule/redigo/redis"
)

type ArticleController struct {
	beego.Controller
}

func (ac * ArticleController)ShowArticleList(){

	o:=orm.NewOrm() //获取orm对象
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		beego.Info("redis数据库连接失败",err)
		return
	}
	defer conn.Close()

	FirstPage:=false //是否为首页
	EndPage:=false //是否为末页
	var count int64 //数据总数
	pageCount:=0.0 //总页数
	pageSize:=2 //单页显示数
	start:=0 //起始页

	st:=ac.GetString("select")  //获取下拉选框数据
	pageIndex,err:=strconv.Atoi(ac.GetString("pageIndex")) //获取当前页码
	if err!=nil{
		pageIndex=1 //处理默认页码
	}
	un:=ac.GetSession("userName") //获取登陆用户名

	qta:=o.QueryTable("Article") //查询文章表所有数据
	qtat:=o.QueryTable("ArticleType") //查询文章类型表所有数据

	awts:=[]models.Article{} //初始化文章带类型的对象
	ats:=[]models.ArticleType{} //初始化文字类型对象

	rel,err:=redis.Bytes(conn.Do("get","types"))
	dec:=gob.NewDecoder(bytes.NewReader(rel))
	dec.Decode(&ats)
	//beego.Info(ats,"-----------")

	//ats没有数据则从mysql数据库提取
	if len(ats)==0{
		qtat.All(&ats)//获取全部文章类型,传到html
		//序列化存入redis
		buffer:=bytes.Buffer{}
		enc:=gob.NewEncoder(&buffer)
		err=enc.Encode(ats)
		_,err=conn.Do("set","types",buffer.Bytes())
		if err!=nil{
			beego.Info("redis set数据错误",err)
			return
		}
		beego.Info("从mysql数据库中取到数据")
	}

	beego.Info("从mysql跳过来了")
	
	if st==""{
		//下拉框无类型,获取全部带类型文章数目
		count,_=qta.RelatedSel("ArticleType").Count()
		//计算总页数
		pageCount=float64(count)/float64(pageSize)
		//向上取整
		pageCount=math.Ceil(pageCount)
		//计算起始页
		start=pageSize*(pageIndex-1)
		//将含文章类型文章分页显示,数据存于带文章类型对象中,传到html
		qta.Limit(pageSize,start).RelatedSel("ArticleType").All(&awts)
		//若为首页,html判断{{if compare .FirstPage true}} 上一页与首页不可点击
		if pageIndex==1{
			FirstPage=true
		}
		//若为末页,html判断{{if compare .EndPage true}} 下一页与末页不可点击
		if pageIndex==int(pageCount){
			EndPage=true
		}
	}else{
		//下拉框有类型选择,获取对应类型文章数目
		count,_=qta.RelatedSel("ArticleType").Filter("ArticleType__TypeName",st).Count()
		pageCount=float64(count)/float64(pageSize)
		pageCount=math.Ceil(pageCount)
		start=pageSize*(pageIndex-1)
		qta.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",st).All(&awts)
		if pageIndex==1{
			FirstPage=true
		}
		if pageIndex==int(pageCount){
			EndPage=true
		}
	}

	ac.Data["un"]=un //展示用户名
	ac.Data["st"]=st // $.st判断选框是否选中 &select={{.st}}传递已选中数据
	ac.Data["ats"]=ats //文章类型总数据
	ac.Data["arts"]=awts //带类型文章总数据
	ac.Data["FirstPage"]=FirstPage // {{if compare .FirstPage true}}
	ac.Data["EndPage"]=EndPage // {{if compare .EndPage true}}
	ac.Data["pageCount"]=pageCount //总页数 末页:pageIndex={{.pageCount}}
	ac.Data["count"]=count //文章记录总数
	ac.Data["pageIndex"]=pageIndex //首页=1,上一页={{.pageIndex | ShowPrePage}},下一页,末页

	ac.Layout="layout.html"
	ac.TplName="index.html"
}
/*
1.查询
2.传给视图显示
*/
//展示文章列表页 index.html
func (ac * ArticleController)ShowArticleList00(){
	//index.html
	//判断session
	//un:=ac.GetSession("userName")
	//if un==nil{//无session回到登陆页
	//	ac.Redirect("/",302)
	//	return
	//}

	//1.查询
	o:=orm.NewOrm()
	qta:=o.QueryTable("Article")
	//1.接收数据 name="select"
	st:=ac.GetString("select")
	var as[] models.Article
	FirstPage:=false
	EndPage:=false

	//总数据数
	count,err:=qta.RelatedSel("ArticleType").Filter("ArticleType__TypeName",st).Count() //返回数据条目数
	if err!=nil{
		beego.Info("查询错误")
		return
	}
	beego.Info(count,"-----------------------")
	//获取总页数 总数/单页显示数
	//pageIndex:=1//起始位置
	pageIndex,err:=strconv.Atoi(ac.GetString("pageIndex"))
	if err!=nil{
		pageIndex=1//处理默认页码
	}

	pageSize:=2 //单页显示数
	pageCount:=float64(count)/float64(pageSize)//总页数
	pageCount=math.Ceil(pageCount)
	start:=2*(pageIndex-1)
	//qs.All(&arts)//select * from Article
	//beego.Info(arts[0])
	//设置显示个数
	qta.Limit(pageSize,start).RelatedSel("ArticleType").All(&as)//1. pageSize 2. start开始位置
	//2.传给视图显示 {{range .arts }}  {{.Id}}   {{end}}
	//首页末页数据处理
	if pageIndex==1{
		FirstPage=true
	}
	if pageIndex==int(pageCount){
		EndPage=true
	}

	//获取文章类型
	qtat:=o.QueryTable("ArticleType")
	ats:=[]models.ArticleType{}
	_,err=qtat.All(&ats)
	if err!=nil{
		beego.Info("文章类型读取错误",err)
		ac.TplName="index.html"
		return
	}
	//根据类型获取数据
	//1.接收数据 name="select"
	//st:=ac.GetString("select")
	//beego.Info(tName)
	//处理数据
	awts:=[]models.Article{}
	if st==""{
		beego.Info("下拉框读取失败")
		//return
		//获取全部数据
		qta.Limit(pageSize,start).RelatedSel("ArticleType").All(&awts)
	}else{//获取对应数据
		//根据下拉框内容,获取相应内容数据
		//o:=orm.NewOrm()
		//arts:=[]models.Article{}
		beego.Info("查询结果----",st)
		_,err:=qta.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",st).All(&awts)
		if err!=nil{
			beego.Info("article查询失败")
		}
		beego.Info(awts,"--------------------")
	}
	un:=ac.GetSession("userName")



	//返回视图
	//ac.Redirect()

	ac.Data["un"]=un
	ac.Data["st"]=st
	ac.Data["ats"]=ats
	ac.Data["FirstPage"]=FirstPage
	ac.Data["EndPage"]=EndPage
	ac.Data["pageCount"]=pageCount
	ac.Data["count"]=count
	ac.Data["pageIndex"]=pageIndex
	ac.Data["arts"]=awts  // {{range $index,$val:=.arts }} {{$val.Id}} {{end}}
	ac.Layout="layout.html"
	ac.TplName="index.html"
}
//处理下拉框
func (ac *ArticleController)HandleSelect(){

}


//添加文章页 add.html
func (ac *ArticleController)ShowAddArticle(){

	//类型选择
	o:=orm.NewOrm()
	ats:=[]models.ArticleType{}
	_,err:=o.QueryTable("ArticleType").All(&ats)
	if err!=nil{
		beego.Info("文章类型读取错误",err)
	}
	ac.Data["ats"]=ats
	ac.TplName="add.html"
}
/*
1.拿数据
	Id int `orm:"pk;atuo"`
	Title string `orm:"size(20)"'`
	Content string `orm:size(500)`
	Img string `orm:size(50);null`
	//Type string
	//orm:"auto_now_add;type(datatime);  auto_now_add 每次设置
	// orm:"auto_now;type(data); auto_now 首次设置
	DateTime time.Time `orm:"auto_now,type(datatime);"`
	Count int `orm:"default(0)"'`
2.判断数据
3.插入数据
4.返回视图
*/
//添加文章 add.html
func (ac *ArticleController)HandleAddArticle(){
	// add.html
	//1.拿数据
	title:=ac.GetString("articleName") //name="articleName"
	content:=ac.GetString("content") //name="content"
	imgFile, head, err :=ac.GetFile("uploadname") //name="uploadname"
	defer imgFile.Close()
	//存储图片
	//1. 文件格式
	ext:=path.Ext(head.Filename)
	beego.Info(ext)
	if ext!=".jpg"&&ext!=".png"&&ext!=".jpeg"{
		beego.Info("文件格式错误")
		ac.TplName="add.html"
		return
	}
	//2. 文件大小
	if head.Size>5000000{
		beego.Info("文件过大,上传失败")
		ac.TplName="add.html"
		return
	}
	//3. 不能重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")+head.Filename
	ac.SaveToFile("uploadname","./static/img/"+fileName)
	if err!=nil{
		beego.Info("上传文件失败",err)
		ac.TplName="add.html"
		return
	}
	beego.Info(title,content,head.Filename,fileName)
	//2.判断数据
	if title==""||content==""{
		beego.Info("标题或内容不能为空")
		ac.TplName="add.html"
		return
	}
	//3.插入数据
	o:=orm.NewOrm()

	var art models.Article
	art.Title=title
	art.Content=content
	art.Img="/static/img/"+fileName
	art.DateTime=time.Now()
	//art.Id=3

	//为article对象赋值
	//获取下拉框传递过来的类型数据
	tName:=ac.GetString("select")
	if tName==""{
		beego.Info("下拉框数据获取失败")
		return
	}
	//从数据库获取ArticleType对象
	at:=models.ArticleType{}
	at.TypeName=tName
	err=o.Read(&at,"TypeName")
	if err!=nil{
		beego.Info("获取类型失败",err)
		return
	}
	art.ArticleType=&at

	n,err:=o.Insert(&art)
	if err!=nil{
		beego.Info("orm insert err",err)
		return
	}
	beego.Info("改变数量",n)
	//ac.Ctx.WriteString("insert成功")

	//返回视图
	ac.Redirect("/Article/showArticle",302)
}
//展示文章详情页
func (ac *ArticleController)ShowArticleDetail(){
	//获取数据 从index.html传数据过来
	//<a href="/showArticleDetail?articleId={{$val.Id}}">查看详情</a>
	id,err:=ac.GetInt("articleId")
	//数据校验
	if err!=nil{
		beego.Info("传递链接错误")
	}
	//操作数据
	o:=orm.NewOrm()
	var art models.Article
	art.Id=id
	o.Read(&art)
	//修改阅读量
	art.Count+=1
	//增加浏览者信息
	//1.获取对象
	//a2:=models.Article{Id:id}
	//2. 获取多对多操作对象
	m2m:=o.QueryM2M(&art,"Users")
	//获取插入对象
	un:=ac.GetSession("userName")
	u:=models.User{}
	u.UserName=un.(string)
	o.Read(&u,"UserName")

	//多对多插入
	_,err=m2m.Add(&u)
	if err!=nil {
		beego.Info("多对多插入失败")
	}
	o.Update(&art)

	//多对的查询
	//o.LoadRelated(&art,"Users")
	us:=[]models.User{}
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&us)
	//beego.Info(art)
	//返回视图页面
	ac.Data["us"]=us
	ac.Data["art"]=art //传数据到content.html {{.art.Title}}
	ac.Layout="layout.html"
	ac.LayoutSections=make(map[string]string)
	ac.LayoutSections["ContentHead"]="head.html"
	ac.TplName="content.html"  //转到content.html <img src={{.art.Img}}> {{.art.Count}}
	//{{.art.DateTime.Format "2006-01-02 15:04:05"}}
}
/*
1. URl传值
2. 执行delete
*/
//从index.html获取id删除数据
func (ac *ArticleController)HandleDelete(){
	//index.html
	//获取数据 <a href="/deleteArticle?id={{$val.Id}}"
	id,err:=ac.GetInt("id")
	if err!=nil{
		beego.Info("传递链接错误")
	}
	//操作数据
	o:=orm.NewOrm()
	var art models.Article
	art.Id=id

	n,err:=o.Delete(&art)
	if err!=nil{
		beego.Info("delete err",err)
		return
	}
	beego.Info("改变数量",n)
	//ac.Ctx.WriteString("insert成功")
	//返回视图
	//ac.TplName="/index.html"
	ac.Redirect("/Article/showArticle",302)
}
//从index.html获取🆔id编辑数据
func (ac * ArticleController)ShowUpdate(){
	//index.html
	//获取数据 <a href="/updateArticle?id={{$val.Id}}"
	id:=ac.GetString("id")
	//判断
	if id==""{
		beego.Info("传递链接错误")
	}
	//查询数据库,获取文章信息,传递数据库信息到update.html
	o:=orm.NewOrm()
	art:=models.Article{}
	id2Int,err:=strconv.Atoi(id)
	if err!=nil{
		beego.Info("转换错误")
		return
	}
	art.Id=id2Int

	err=o.Read(&art)
	if err!=nil{
		beego.Info("读取数据库信息错误")
		return
	}
	ac.Data["art"]=art
	ac.TplName="update.html"
}
func (ac *ArticleController)HandleUpdate(){
	//获取update.html post数据更新至数据库
	// 方式1: action="/updateArticle?id={{.art.Id}}
	//id,err:=ac.GetInt("id")
	//方式2: 隐藏域<input name="id" value="{{.art.Id}}" hidden="hidden">
	id,_:=ac.GetInt("id")
	title:=ac.GetString("articleName")//name = "articleName"
	content:=ac.GetString("content")//name="content"
	if title==""||content==""{
		beego.Info("标题与内容不能为空")
		ac.TplName="update.html"
		return
	}
	imgFile,head,err:=ac.GetFile("uploadname")//name="uploadname"
	defer imgFile.Close()
	if err!=nil{
		beego.Info("上传文件失败 err",err)
		return
	}
	//判断文件格式,大小,防止重名
	ext:=path.Ext(head.Filename)
	if ext!=".jpg"&&ext!=".jpeg"&&ext!="png"{
		beego.Info("文件格式错误")
		ac.TplName="update.html"
		return
	}
	if head.Size>5000000{
		beego.Info("文件太大")
		ac.TplName="update.html"
		return
	}
	fileName:=time.Now().Format("2006-01-02 15:04:05")+head.Filename
	//存储新图片到本地
	ac.SaveToFile("uploadname","./static/img/"+fileName)
	//更新数据至数据库
	o:=orm.NewOrm()
	art:=models.Article{}

	art.Id=id
	art.Img="/static/img/"+fileName
	art.Title=title
	art.Content=content
	art.DateTime=time.Now()
	n,err:=o.Update(&art)
	if err!=nil{
		beego.Info("update err",err)
		return
	}
	beego.Info("改变数量",n)
	//返回视图
	ac.Redirect("/Article/showArticle",302)
}
func (ac *ArticleController)ShowAddType(){
	//读取数据库,显示已有类型
	o:=orm.NewOrm()
	ats:=[]models.ArticleType{}
	qt:=o.QueryTable("ArticleType")
	_,err:=qt.All(&ats)
	if err!=nil{
		beego.Info("查询类型错误",err)
	}
	ac.Data["ats"]=ats
	ac.TplName="addType.html"
}

func (ac * ArticleController)HandleAddType(){
	//1. 获取数据
	//	name="typeName"
	tName:=ac.GetString("typeName")
	beego.Info(tName)
	//2.判断数据
	if tName==""{
		beego.Info("类型不能为空")
		ac.Redirect("/Article/addArticleType",302)
		return
	}
	//3.执行插入数据
	o:=orm.NewOrm()
	at:=models.ArticleType{}
	at.TypeName=tName
	n,err:=o.Insert(&at)
	if err!=nil{
		beego.Info("插入失败",err)
		return
	}
	beego.Info("更新数量----",n)
	//4.返回视图
	ac.Redirect("/Article/addArticleType",302)
}

//删除
func (ac *ArticleController)HandleDeleteType()  {
	//获取数据
	tId,_:=ac.GetInt("typeId")
	//组织数据
	o:=orm.NewOrm()
	at:=models.ArticleType{}
	at.Id=tId
	//执行删除
	o.Delete(&at)
	//返回视图
	ac.Redirect("/Article/addArticleType",302)
}

//退出登陆
func (ac *ArticleController)Logout()  {
	//1. 删除登陆状态
	ac.DelSession("userName")
	//2.跳转到登陆页
	ac.Redirect("/",302)
}