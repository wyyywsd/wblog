package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"net/http"
	"strconv"
)
const articleCommentCount int = 5
func NewComment(context *gin.Context){
	fmt.Println("==============================================")
	comment := context.PostForm("new_article_comment")
	article_id := context.Param("id")
	article := models.FindArticleById(article_id)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取该文章所有的评论
	comment_count := models.CommentCount(article_id)
	models.CreateComment(article.ID,current_user.ID,comment,comment_count+1)
	context.Redirect(http.StatusMovedPermanently, "/article/"+article_id+"")

}

func ShowCommentByArticle(context *gin.Context){
	article_id := context.Param("id")
	article := models.FindArticleById(article_id)
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == ""{
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	//session := sessions.Default(context)
	//current_user_name := session.Get("sessionid")
	//获取某页的评论
	comments := models.FindCommentByArticle(article_id,i,articleCommentCount)
	//获取该文章所有的评论
	comment_count := models.CommentCount(article_id)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := comment_count/articleCommentCount
	if comment_count%articleCommentCount != 0{
		pageCount = (comment_count/articleCommentCount)+1
	}
	context.HTML(200, "_show_comment.html", gin.H{
		"comments":    comments,
		"article": article,
		"pageCount":pageCount,
		//"user_session": session.Get("sessionid"),
		"current_page": i,
		"article_user_id": article.UserId,
	})





}
