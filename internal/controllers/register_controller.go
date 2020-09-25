package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"log"
)

func Register(context *gin.Context){
	context.HTML(200, "register.html", gin.H{"keywords":"macbook pro"})
}

func Signup(context *gin.Context){
	log.Println("进入注册action")
	username := context.PostForm("username")
	password := context.PostForm("password")
	email := context.PostForm("email")
	//验证用户名和邮箱是否存在
	_,existsUser,_:=models.FindUserByUserName(username)
	_,existEmail :=models.FindUserByEmail(email)
	if !existsUser{
		log.Println("找不到用户 继续注册")
		//继续注册
		//判断邮箱 暂时在后端判断 邮箱唯一性
		if !existEmail{
			//保存并跳转到登陆页面
			user,_ := models.CreateUser(username,password,email)
			fmt.Println(user)
			//context.Redirect(http.StatusMovedPermanently, "/index")
			context.HTML(200, "signin.html", gin.H{"message":"注册成功，请登录！"})

		}else{
			context.HTML(200, "register.html", gin.H{"message":"邮箱已存在"})
		}

	}else{
		context.HTML(200, "register.html", gin.H{"message":"用户名已存在"})
	}
}