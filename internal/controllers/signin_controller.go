package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"log"
	"net/http"
)

func Signin(context *gin.Context) {
	log.Println("进入登录controller")
	context.HTML(200, "signin.html", gin.H{"message": "12321"})
}

//登陆
func Login(context *gin.Context) {
	log.Println("进入登陆action")
	//查询用户信息
	//user,existsUser,err:=models.FindUserByUserName(req.Username)
	log.Println("22222")
	userName := context.PostForm("username")
	passWord := context.PostForm("password")
	user, existsUser, err := models.FindUserByUserName(userName)
	log.Println(userName)

	log.Println("查到用户了")

	if err != nil {
		log.Println("报错")
		log.Fatalf("UserLogin:查询用户信出错，err:%v", err)
		//context.HTML(200, "signin.html", gin.H{"message":"查询用户信出错，请重新登陆"})
		context.Redirect(http.StatusMovedPermanently, "/index")
		return
	}
	if !existsUser {
		log.Println("找不到用户")
		//context.HTML(200, "signin.html", gin.H{"message":"找不到用户，请检查用户名是否正确"})
		context.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}
	if user.PassWord != passWord {
		log.Println("密码不对")
		context.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}
	log.Println("调转成功页面")

	//初始化session
	session := sessions.Default(context)
	//设置sessions的相关参数
	option := sessions.Options{MaxAge: 3600}
	session.Options(option)
	session.Set("sessionId", userName)
	sessionErr := session.Save()

	if sessionErr != nil {
		fmt.Println("保存session失败！错误代码：", sessionErr)
	} else {
		fmt.Println("保存session成功")
	}
	v := session.Get("sessionId")
	fmt.Println("sessionId的值是:", v)

	context.Redirect(http.StatusMovedPermanently, "/index")

}

//判断是否登陆的中间件
func AuthRequiredSession() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionId := session.Get("sessionId")
		fmt.Println("sessionID+++++", sessionId)
		if sessionId == nil {
			context.HTML(200, "signin.html", gin.H{"message": "请先登录"})
			//context.Next()
			context.Abort()
		} else {
			option := sessions.Options{MaxAge: 3600}
			session.Options(option)
			session.Save()
		}
		//context.Next()
	}
}
func Logout(c *gin.Context) {
	fmt.Println("进入注销动作")
	session := sessions.Default(c)
	fmt.Println("sessionId===================", session.Get("sessionId"))
	session.Delete("sessionId")
	session.Save()
	c.HTML(200, "signin.html", gin.H{"message": "session已清除"})
}

//是否登陆的中间件
//func AuthRequired() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		cookie, _ := c.Request.Cookie("username")
//		if cookie == nil {
//			c.HTML(200, "signin.html", gin.H{"message":"请先登录"})
//			c.Abort()
//		}
//		// 实际生产中应校验cookie是否合法
//		c.Next()
//	}
//}
