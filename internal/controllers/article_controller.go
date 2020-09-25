package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"net/http"
)


func Show_Article(context *gin.Context) {
	id :=context.Param("id")
	article := models.FindArticleById(id)
	article_user,_:= models.FindUserByArticle(&article)
	session := sessions.Default(context)
	fmt.Println(article.ArticleContent)
	current_user,_,_ :=  models.FindUserByUserName(fmt.Sprint(session.Get("sessionid")))
	context.HTML(200,"show_article.html",gin.H{
		"article":article,
		"article_user": article_user,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
	})
}

func NewArticle(context *gin.Context){
	session := sessions.Default(context)
	labels, _ := models.AllLabels()
	current_user,_,_ :=  models.FindUserByUserName(fmt.Sprint(session.Get("sessionid")))
	context.HTML(200, "new_article.html", gin.H{
		"user_session":session.Get("sessionid"),
		"current_user": current_user,
		"labels": labels,
	})
}

func SaveArticle(context *gin.Context){
	article_title := context.PostForm("article_title")
	article_content := context.PostForm("article_content")
	session := sessions.Default(context)
	article_label_id := context.PostForm("article_label")
	label := models.FindLabelById(article_label_id)
	is_public_str :=context.PostForm("is_public")
	is_public := true
	if is_public_str == "false"{
		is_public = false
	}
	//从session中获取当前登陆的用户名
	current_user_name := session.Get("sessionid")
	current_user,_,_ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	models.CreateArticle(article_title,article_content,current_user.ID,label,is_public)
	fmt.Println("保存文章成功")
	context.Redirect(http.StatusMovedPermanently, "/index")

}

func EditArticle(context *gin.Context){
	id :=context.Param("id")
	article := models.FindArticleById(id)
	article_user,_:= models.FindUserByArticle(&article)
	session := sessions.Default(context)
	labels, _ := models.AllLabels()
	label :=  article.FindLabelsByArticle()
	current_user,_,_ :=  models.FindUserByUserName(fmt.Sprint(session.Get("sessionid")))
	context.HTML(200,"edit_article.html",gin.H{
		"article":article,
		"article_user": article_user,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"labels": labels,
		"label": label,
	})
}

func UpdateArticle(context *gin.Context){
	id :=context.Param("id")
	current_article := models.FindArticleById(id)
	article_title := context.PostForm("article_title")
	article_content := context.PostForm("article_content")
	//article_label_id := context.PostForm("article_label")
	//label := models.FindLabelById(article_label_id)
	is_public_str :=context.PostForm("is_public")
	is_public := true
	if is_public_str == "false"{
		is_public = false
	}
	var update_map = map[string]interface{}{}
	if current_article.ArticleTitle != article_title{
		update_map["ArticleTitle"] = article_title
	}
	if current_article.ArticleContent != article_content{
		update_map["ArticleContent"] = article_content
	}
	if current_article.IsPublic != is_public{
		update_map["IsPublic"] = is_public
	}

	models.UpdateArticle(current_article,update_map)
	fmt.Println("更新文章成功")
	context.Redirect(http.StatusMovedPermanently, "/index")

}