package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"net/http"
)

func NewComment(context *gin.Context){
	fmt.Println("==============================================")
	comment := context.PostForm("new_article_comment")
	article_id := context.Param("id")
	article := models.FindArticleById(article_id)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	fmt.Println(article.ArticleTitle)
	fmt.Println(current_user.Username)
	fmt.Println(comment,article_id)
	models.CreateComment(article.ID,current_user.ID,comment)
	context.Redirect(http.StatusMovedPermanently, "/article/"+article_id+"")

}

