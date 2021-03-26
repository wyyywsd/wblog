package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/db"
	"gorm_demo/internal/models"
	"net/http"
)

func Show_Article(context *gin.Context) {
	articleId := context.Param("id")
	article := models.FindArticleById(articleId)
	articleUser, _ := models.FindUserByArticle(&article)
	session := sessions.Default(context)
	//fmt.Println(article.ArticleContent)
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(session.Get("sessionId")))
	//该文章下第一页的评论
	comments := models.FindCommentByArticle(articleId, 1, articleCommentCount)
	commentCount := models.CommentCount(articleId)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := commentCount / articleCommentCount
	if commentCount%articleCommentCount != 0 {
		pageCount = (commentCount / articleCommentCount) + 1
	}

	context.HTML(200, "show_article.html", gin.H{
		"article":         article,
		"articleUser":     articleUser,
		"userSession":    session.Get("sessionId"),
		"currentUser":     currentUser,
		"comments":        comments,
		"commentCount":   commentCount,
		"pageCount":       pageCount,
		"currentPage":    1,
		"articleUserId": article.UserId,
	})
}

func NewArticle(context *gin.Context) {
	session := sessions.Default(context)
	labels, _ := models.AllLabels()
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(session.Get("sessionId")))
	context.HTML(200, "new_article.html", gin.H{
		"userSession": session.Get("sessionId"),
		"currentUser": currentUser,
		"labels":       labels,
	})
}

func SaveArticle(context *gin.Context) {
	articleTitle := context.PostForm("article_title")
	articleContent := context.PostForm("article_content")
	session := sessions.Default(context)
	articleLabelId := context.PostForm("article_label")
	label := models.FindLabelById(articleLabelId)
	isPublicStr := context.PostForm("is_public")
	isPublic := true
	if isPublicStr == "false" {
		isPublic = false
	}
	//从session中获取当前登陆的用户名
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	models.CreateArticle(articleTitle, articleContent, currentUser.ID, label, isPublic)
	fmt.Println("保存文章成功")
	context.Redirect(http.StatusMovedPermanently, "/index")

}

func EditArticle(context *gin.Context) {
	id := context.Param("id")
	article := models.FindArticleById(id)
	articleUser, _ := models.FindUserByArticle(&article)
	session := sessions.Default(context)
	labels, _ := models.AllLabels()
	label := article.FindLabelsByArticle()
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(session.Get("sessionId")))
	context.HTML(200, "edit_article.html", gin.H{
		"article":      article,
		"articleUser":  articleUser,
		"userSession": session.Get("sessionId"),
		"currentUser": currentUser,
		"labels":       labels,
		"label":        label,
	})
}

func UpdateArticle(context *gin.Context) {
	id := context.Param("id")
	currentArticle := models.FindArticleById(id)
	articleTitle := context.PostForm("article_title")
	articleContent := context.PostForm("article_content")
	isPublicStr := context.PostForm("is_public")
	isPublic := true
	if isPublicStr == "false" {
		isPublic = false
	}
	var updateMap = map[string]interface{}{}
	if currentArticle.ArticleTitle != articleTitle {
		updateMap["ArticleTitle"] = articleTitle
	}
	if currentArticle.ArticleContent != articleContent {
		updateMap["ArticleContent"] = articleContent
	}
	if currentArticle.IsPublic != isPublic {
		updateMap["IsPublic"] = isPublic
	}

	models.UpdateArticle(currentArticle, updateMap)
	fmt.Println("更新文章成功")
	context.Redirect(http.StatusMovedPermanently, "/index")

}

func DeleteArticle(context *gin.Context) {
	id := context.Param("id")
	currentArticle := models.FindArticleById(id)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	//删除
	db.W_Db.Delete(currentArticle)
	context.Redirect(http.StatusMovedPermanently, "/user/"+fmt.Sprint(currentUserName)+"")
}

func CollectArticle(context *gin.Context) {
	articleId := context.Param("id")
	isCollectStr := context.Param("is_collect")
	article := models.FindArticleById(articleId)
	articleUser, _ := models.FindUserByArticle(&article)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	collect, existsCollect, _ := models.FindCollectByUserIdAndArticleId(currentUser.ID, article.ID)
	isCollect := true
	if isCollectStr == "false" {
		isCollect = false
	}
	if existsCollect {
		//更新
		models.UpdateCollect(collect, isCollect)
	} else {
		//新增
		models.CreateCollect(currentUser.ID, article.ID, isCollect)
	}

	//context.Redirect(http.StatusMovedPermanently, "/article/"+articleId+"/")

	context.HTML(200, "_show_is_collect.html", gin.H{
		"article":      article,
		"currentUser": currentUser,
		"userSession": currentUserName,
		"articleUser":  articleUser,
	})

}
