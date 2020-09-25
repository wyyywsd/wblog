package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)
const userArticleCount int = 5
func ShowUser(context *gin.Context) {

	current_user_name :=context.Param("user_session")
	current_user,_,_:= models.FindUserByUserName(current_user_name)
	session := sessions.Default(context)
	//获取第一页的数据
	articles,_ := models.UserArticleLimit(1,userArticleCount,current_user.ID)
	//获取一共有多少文章
	count := models.UserArticleCount(current_user.ID)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count/userArticleCount
	if count%userArticleCount != 0{
		pageCount = (count/userArticleCount)+1
	}
	context.HTML(200,"show_user.html",gin.H{
		"current_user":current_user,
		"user_session": session.Get("sessionid"),
		"articles":    articles,
		"pageCount":pageCount,
		"current_page": 1,
	})
}

func UpdateUser(context *gin.Context) {
	fmt.Println("进入更新用户界面")
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取到图片  如果获取不到图片 目前暂时默认为没有上传
	file, header, err := context.Request.FormFile("file")
	profilephoto := ""
	if err != nil {
		//暂时什么都不做
	}else{
		filename := header.Filename
		//获取文件的后缀名   为数组的最后一位
		filename_slice := strings.Split(filename, ".")
		le := len(filename_slice) -1
		filename = fmt.Sprint(current_user.ID)+"."+filename_slice[le]
		out, err := os.Create("public/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		profilephoto = "/file/"+filename
	}
	//修改页面的数据
	username := context.PostForm("username")
	password := context.PostForm("password")
	email := context.PostForm("email")
	//判断修改的邮箱 用户名 是否存在
	existsUser := false
	if username != ""{
		_,existsUser,_=models.FindUserByUserName(username)
	}
	existEmail := false
	if email != ""{
		_,existEmail =models.FindUserByEmail(email)
	}

	if !existsUser{
		log.Println("修改的用户名不存在，可以修改")
		//继续修改
		//判断邮箱 暂时在后端判断 邮箱唯一性
		if !existEmail{
			var update_map = map[string]interface{}{}
			if username != "" {
				update_map["Username"] = username
			}
			if password != "" {
				update_map["PassWord"] = password
			}
			if email != "" {
				update_map["Email"] = email
			}
			if profilephoto != "" {
				update_map["ProfilePhoto"] = profilephoto
			}
			err := models.UpdateUser(*current_user,update_map)
			if err!= nil{

			}else{
				//修改用户信息之后  在这里重新设置一下session的值
				if username != "" && current_user_name != username{
					session.Set("sessionid", username)
					session.Save()
				}

				context.HTML(200,"index.html",gin.H{
					"current_user":current_user,
					"user_session": session.Get("sessionid"),
				})
			}

		}else{
			context.HTML(200, "user_setting.html", gin.H{"message":"邮箱已存在"})
		}

	}else{
		context.HTML(200, "user_setting.html", gin.H{"message":"保存失败"})
	}



}
func UserSetting(context *gin.Context) {
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	context.HTML(200,"user_setting.html",gin.H{
		"current_user":current_user,
		"user_session": session.Get("sessionid"),
	})

}

func BasicSetting(context *gin.Context) {
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	context.HTML(200,"_basic_setting.html",gin.H{
		"current_user":current_user,
		"user_session": session.Get("sessionid"),
	})
}


//func ImgTest(c *gin.Context){
//	file, header, err := c.Request.FormFile("file")
//	if err != nil {
//		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
//		return
//	}
//
//	filename := header.Filename
//
//	out, err := os.Create("public/" + filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, file)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	filepath := "http://localhost:8080/file/" + filename
//	fmt.Println("dahsiid******************************************",filepath)
//	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
//
//}

func ShowUserArticles(context *gin.Context){
	//如果获取不到page 默认就是1
	page := context.Param("page")
	if page == ""{
		page = "1"

	}
	fmt.Println("------------------------------------------------------------------",page)
	//将string类型的page 设置成int
	i, _ := strconv.Atoi(page)
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取某一页的数据
	articles,_ := models.UserArticleLimit(i,userArticleCount,current_user.ID)
	//获取一共有多少文章
	count := models.UserArticleCount(current_user.ID)
	//通过文章的数量 算出分页一共有多少页    如果有余数  就加一 目前先都加1  后面再改
	pageCount := count/userArticleCount
	if count%userArticleCount != 0{
		pageCount = (count/userArticleCount)+1
	}
	//labels, _ := models.AllLabels()
	context.HTML(200, "_user_articles.html", gin.H{
		"articles":    articles,
		//"labels": 	labels,
		"pageCount":pageCount,
		"user_session": session.Get("sessionid"),
		"current_user": current_user,
		"current_page": i,
	})
}