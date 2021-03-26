package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const userArticleCount int = 5

func ShowUser(context *gin.Context) {

	currentUserName := context.Param("user_session")
	currentUser, _, _ := models.FindUserByUserName(currentUserName)
	session := sessions.Default(context)
	//获取第一页的数据
	articles, _ := models.UserArticleLimit(1, userArticleCount, currentUser.ID)
	//获取一共有多少文章
	count := models.UserArticleCount(currentUser.ID)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count / userArticleCount
	if count%userArticleCount != 0 {
		pageCount = (count / userArticleCount) + 1
	}
	context.HTML(200, "show_user.html", gin.H{
		"currentUser": currentUser,
		"userSession": session.Get("sessionId"),
		"articles":     articles,
		"pageCount":    pageCount,
		"currentPage": 1,
	})
}

func UpdateUser(context *gin.Context) {
	fmt.Println("进入更新用户界面")
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	//获取到图片  如果获取不到图片 目前暂时默认为没有上传
	file, header, err := context.Request.FormFile("file")
	profilePhoto := ""
	if err != nil {
		//暂时什么都不做
	} else {
		fileName := header.Filename
		//获取文件的后缀名   为数组的最后一位
		fileNameSlice := strings.Split(fileName, ".")
		le := len(fileNameSlice) - 1
		fileName = fmt.Sprint(currentUser.ID) + "." + fileNameSlice[le]
		out, err := os.Create("file/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		profilePhoto = "/file/" + fileName
	}
	//修改页面的数据
	userName := context.PostForm("username")
	passWord := context.PostForm("password")
	email := context.PostForm("email")
	//判断修改的邮箱 用户名 是否存在
	existsUser := false
	if userName != "" {
		_, existsUser, _ = models.FindUserByUserName(userName)
	}
	existEmail := false
	if email != "" {
		_, existEmail = models.FindUserByEmail(email)
	}

	if !existsUser {
		log.Println("修改的用户名不存在，可以修改")
		//继续修改
		//判断邮箱 暂时在后端判断 邮箱唯一性
		if !existEmail {
			var update_map = map[string]interface{}{}
			if userName != "" {
				update_map["Username"] = userName
			}
			if passWord != "" {
				update_map["PassWord"] = passWord
			}
			if email != "" {
				update_map["Email"] = email
			}
			if profilePhoto != "" {
				update_map["ProfilePhoto"] = profilePhoto
			}
			err := models.UpdateUser(*currentUser, update_map)
			if err != nil {

			} else {
				//修改用户信息之后  在这里重新设置一下session的值
				if userName != "" && currentUserName != userName {
					session.Set("sessionId", userName)
					session.Save()
				}

				//context.HTML(200,"index.html",gin.H{
				//	"currentUser":currentUser,
				//	"user_session": session.Get("sessionId"),
				//})
				context.Redirect(http.StatusMovedPermanently, "/user/"+fmt.Sprint(currentUserName)+"")
			}

		} else {
			context.HTML(200, "user_setting.html", gin.H{"message": "邮箱已存在"})
		}

	} else {
		context.HTML(200, "user_setting.html", gin.H{"message": "保存失败"})
	}

}
func UserSetting(context *gin.Context) {
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	context.HTML(200, "user_setting.html", gin.H{
		"currentUser": currentUser,
		"userSession": session.Get("sessionId"),
	})

}

func BasicSetting(context *gin.Context) {
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	context.HTML(200, "_basic_setting.html", gin.H{
		"currentUser": currentUser,
		"userSession": session.Get("sessionId"),
	})
}

func ShowUserArticles(context *gin.Context) {
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == "" {
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	//获取某一页的数据
	articles, _ := models.UserArticleLimit(i, userArticleCount, currentUser.ID)
	//获取一共有多少文章
	count := models.UserArticleCount(currentUser.ID)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count / userArticleCount
	if count%userArticleCount != 0 {
		pageCount = (count / userArticleCount) + 1
	}
	//labels, _ := models.AllLabels()
	context.HTML(200, "_user_articles.html", gin.H{
		"articles": articles,
		//"labels": 	labels,
		"pageCount":    pageCount,
		"userSession": session.Get("sessionId"),
		"currentUser": currentUser,
		"currentPage":  i,
		"pageType":    "user_articles",
	})
}

func ShowUserCollects(context *gin.Context) {
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == "" {
		page = "1"
	}
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	myCollectArticles, _ := models.FindUserCollectArticles(currentUser.ID, userArticleCount, i)
	//获取一共有多少文章
	count := models.UserCollectArticlesCount(currentUser.ID)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count / userArticleCount
	if count%userArticleCount != 0 {
		pageCount = (count / userArticleCount) + 1
	}

	context.HTML(200, "_user_collect_articles.html", gin.H{
		"articles": myCollectArticles,
		//"labels": 	labels,
		"pageCount":    pageCount,
		"userSession": session.Get("sessionId"),
		"currentUser": currentUser,
		"currentPage":  i,
		"pageType":    "collect",
	})
}
