package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
)

func Show_Article_By_Label(context *gin.Context) {
	labelId := context.Param("id")
	labels, _ := models.AllLabels()
	articles := models.FindArticlesByLabel(labelId)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
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
		"userSession":  session.Get("sessionId"),
		"currentUser":  currentUser,
		"pageCount":     pageCount,
		"currentPage":  1,
		"nextPage":     2,
		"PreviousPage": 0,
	})
}
