package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"strconv"
)
//分页的数量
const articleCount int = 8
//index的action
func ArticlePageIndex(context *gin.Context) {
	//获取当前页面的页码
	page := context.Param("page")
	//换算成int类型
	i, _ := strconv.Atoi(page)
	//获取到某一页的文章列表
	articles,_ := models.PublicArticleLimit(i,articleCount)
	//获取一共有多少文章
	count := models.ArticleCount()
	//通过文章的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count/articleCount
	if count%articleCount != 0{
		pageCount = (count/articleCount)+1
	}
	fmt.Println("******************************************************************",count)
	labels, _ := models.AllLabels()
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	//user,_,_:= models.FindUserByUserName(fmt.Sprint(session.Get("user_session")))
	//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++",user.Username)
	context.HTML(200, "index.html", gin.H{
		"articles":    articles,
		"labels": 	labels,
		"pageCount":pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": i,
	})
}


//默认首页
func ArticleIndex(context *gin.Context){
	//获取第一页的数据
	articles,_ := models.PublicArticleLimit(1,articleCount)
	//获取一共有多少文章
	count := models.ArticleCount()
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count/articleCount
	if count%articleCount != 0{
		pageCount = (count/articleCount)+1
	}
	labels, _ := models.AllLabels()
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	//user,_,_:= models.FindUserByUserName(fmt.Sprint(session.Get("user_session")))
	//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++",user.Username)
	context.HTML(200, "index.html", gin.H{
		"articles":    articles,
		"labels": 	labels,
		"pageCount":pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": 1,
	})

}