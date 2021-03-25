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
	articles, _ := models.PublicArticleLimit(i, articleCount)
	//获取一共有多少文章
	count := models.ArticleCount()
	//通过文章的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count / articleCount
	if count%articleCount != 0 {
		pageCount = (count / articleCount) + 1
	}
	fmt.Println("******************************************************************", count)
	labels, _ := models.AllLabels()
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//user,_,_:= models.FindUserByUserName(fmt.Sprint(session.Get("user_session")))
	//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++",user.Username)
	context.HTML(200, "index.html", gin.H{
		"articles":     articles,
		"labels":       labels,
		"pageCount":    pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": i,
	})
}

//默认首页
func ArticleIndex(context *gin.Context) {
	//获取第一页的数据
	articles, _ := models.PublicArticleLimit(1, articleCount)
	//获取一共有多少文章
	count := models.ArticleCount()
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count / articleCount
	if count%articleCount != 0 {
		pageCount = (count / articleCount) + 1
	}
	labels, _ := models.AllLabels()
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//user,_,_:= models.FindUserByUserName(fmt.Sprint(session.Get("user_session")))
	//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++",user.Username)
	context.HTML(200, "index.html", gin.H{
		"articles":     articles,
		"labels":       labels,
		"pageCount":    pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": 1,
	})

}

func SearchArticle(context *gin.Context) {

	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == "" {
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	key_word := context.Query("key_word")
	fmt.Println("808080895i6u5634853247812342vn93vn4293c4n02938402394vb782b7u842cn7", key_word)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//此处返回的error 暂时不处理
	articles, _ := models.FindArticleByKeyWord(key_word, i, articleCount, true, current_user.ID)
	for _, v := range articles {
		fmt.Println(v.ArticleTitle)
		fmt.Println(key_word)
	}
	labels, _ := models.AllLabels()
	//获取一共搜索到多少文章
	count := models.KeyWordArticleCount(key_word, true, current_user.ID)
	////通过文章的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count / articleCount
	if count%articleCount != 0 {
		pageCount = (count / articleCount) + 1
	}
	context.HTML(200, "index.html", gin.H{
		"articles":     articles,
		"labels":       labels,
		"pageCount":    pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": i,
		"page_type":    "search",
		"key_word":     key_word,
	})

}
