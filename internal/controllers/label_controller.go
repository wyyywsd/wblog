package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
)

func Show_Article_By_Label(context *gin.Context) {
	label_id := context.Param("id")
	labels, _ := models.AllLabels()
	articles := models.FindArticlesByLabel(label_id)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取到该标签下的所有文章数量
	count := len(articles)
	//通过文章的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count / articleCount
	if count%articleCount != 0 {
		pageCount = (count / articleCount) + 1
	}
	context.HTML(200, "index.html", gin.H{
		"articles":      articles,
		"labels":        labels,
		"user_session":  session.Get("sessionid"),
		"current_user":  current_user,
		"pageCount":     pageCount,
		"current_page":  1,
		"next_page":     2,
		"Previous_page": 0,
	})
}
