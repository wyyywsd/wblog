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

func NewComment(context *gin.Context) {
	fmt.Println("==============================================")
	comment := context.PostForm("new_article_comment")
	articleId := context.Param("id")
	article := models.FindArticleById(articleId)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	//获取该文章所有的评论
	commentCount := models.CommentCount(articleId)
	models.CreateComment(article.ID, currentUser.ID, comment, commentCount+1)
	context.Redirect(http.StatusMovedPermanently, "/article/"+articleId+"")

}

func ShowCommentByArticle(context *gin.Context) {
	articleId := context.Param("id")
	article := models.FindArticleById(articleId)
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == "" {
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	session := sessions.Default(context)
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(session.Get("sessionId")))
	//获取某页的评论
	comments := models.FindCommentByArticle(articleId, i, articleCommentCount)
	//获取该文章所有的评论
	commentCount := models.CommentCount(articleId)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := commentCount / articleCommentCount
	if commentCount%articleCommentCount != 0 {
		pageCount = (commentCount / articleCommentCount) + 1
	}
	comments[1].FindZansByComment()

	context.HTML(200, "_show_comment.html", gin.H{
		"comments":        comments,
		"article":         article,
		"pageCount":       pageCount,
		"currentUser":    currentUser,
		"currentPage":    i,
		"articleUserId": article.UserId,
	})
}

//func LikeComment(context *gin.Context){
//	fmt.Println("jinlaile")
//	comment_id := context.Param("id")
//	comment := models.FindCommentById(comment_id)
//	var update_map = map[string]interface{}{}
//	comment_like_count := comment.CommentLikeCount
//	fmt.Println("该评论的点赞数为：",comment_like_count)
//	update_map["CommentLikeCount"] = comment_like_count+1
//	models.UpdateComment(comment,update_map)
//	currentPage := context.Param("page")
//	context.Redirect(http.StatusMovedPermanently, "/show_comment_by_article/"+fmt.Sprint(comment.ArticleId)+"/"+currentPage+"")
//}

func LikeComment(context *gin.Context) {
	commentId := context.Param("id")
	comment := models.FindCommentById(commentId)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	currentPage := context.Param("page")
	isZanStr := context.Param("is_zan")
	zan, existsZan, _ := models.FindZanByUserIDAndCommentID(currentUser.ID, comment.ID)
	//在这里判断 是否有zan 表 如果有  就更新  如何没有 就新增
	fmt.Println(zan)
	isZan := true
	if isZanStr == "false" {
		fmt.Println("000ddijfijifjijdladjilfilhfijlasdiwdhiuhfgiljdlasdijn")
		isZan = false
	}
	if existsZan {
		models.UpdateZan(zan, isZan)
	} else {
		models.CreateZan(currentUser.ID, comment.ID, isZan)
	}

	context.Redirect(http.StatusMovedPermanently, "/show_comment_by_article/"+fmt.Sprint(comment.ArticleId)+"/"+currentPage+"")
}
