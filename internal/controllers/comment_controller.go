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
	article_id := context.Param("id")
	article := models.FindArticleById(article_id)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取该文章所有的评论
	comment_count := models.CommentCount(article_id)
	models.CreateComment(article.ID, current_user.ID, comment, comment_count+1)
	context.Redirect(http.StatusMovedPermanently, "/article/"+article_id+"")

}

func ShowCommentByArticle(context *gin.Context) {
	article_id := context.Param("id")
	article := models.FindArticleById(article_id)
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == "" {
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	session := sessions.Default(context)
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(session.Get("sessionid")))
	//获取某页的评论
	comments := models.FindCommentByArticle(article_id, i, articleCommentCount)
	//获取该文章所有的评论
	comment_count := models.CommentCount(article_id)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := comment_count / articleCommentCount
	if comment_count%articleCommentCount != 0 {
		pageCount = (comment_count / articleCommentCount) + 1
	}
	comments[1].FindZansByComment()

	context.HTML(200, "_show_comment.html", gin.H{
		"comments":        comments,
		"article":         article,
		"pageCount":       pageCount,
		"current_user":    current_user,
		"current_page":    i,
		"article_user_id": article.UserId,
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
//	current_page := context.Param("page")
//	context.Redirect(http.StatusMovedPermanently, "/show_comment_by_article/"+fmt.Sprint(comment.ArticleId)+"/"+current_page+"")
//}

func LikeComment(context *gin.Context) {
	comment_id := context.Param("id")
	comment := models.FindCommentById(comment_id)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	current_page := context.Param("page")
	is_zan_s := context.Param("is_zan")
	zan, existsZan, _ := models.FindZanByUserIDAndCommentID(current_user.ID, comment.ID)
	//在这里判断 是否有zan 表 如果有  就更新  如何没有 就新增
	fmt.Println(zan)
	is_zan := true
	if is_zan_s == "false" {
		fmt.Println("000ddijfijifjijdladjilfilhfijlasdiwdhiuhfgiljdlasdijn")
		is_zan = false
	}
	if existsZan {
		models.UpdateZan(zan, is_zan)
	} else {
		models.CreateZan(current_user.ID, comment.ID, is_zan)
	}

	context.Redirect(http.StatusMovedPermanently, "/show_comment_by_article/"+fmt.Sprint(comment.ArticleId)+"/"+current_page+"")
}
