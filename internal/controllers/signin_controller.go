package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"log"
	"net/http"
)

func Signin(context *gin.Context){
	log.Println("进入登录controller")
	//value,err := context.Get("message")
	//fmt.Println("dwifqwjifioojfiwqjijoidjqid",value,err)
	context.HTML(200, "signin.html", gin.H{"message":"12321"})
}


//登陆
func Login(context *gin.Context){
	log.Println("进入登陆action")
	//查询用户信息
	//user,existsUser,err:=models.FindUserByUserName(req.Username)
	log.Println("22222")
	username := context.PostForm("username")
	password := context.PostForm("password")
	user,existsUser,err:=models.FindUserByUserName(username)
	log.Println(username)

	log.Println("查到用户了")

	if err!=nil {
		log.Println("报错")
		log.Fatalf("UserLogin:查询用户信出错，err:%v",err)
		//context.HTML(200, "signin.html", gin.H{"message":"查询用户信出错，请重新登陆"})
		context.Redirect(http.StatusMovedPermanently, "/index")
		return
	}
	if !existsUser{
		log.Println("找不到用户")
		//context.HTML(200, "signin.html", gin.H{"message":"找不到用户，请检查用户名是否正确"})
		context.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}
	if user.PassWord != password {
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
	session.Set("sessionid", username)
	session_err := session.Save()

	if session_err != nil {
		fmt.Println("保存session失败！错误代码：",session_err)
	} else {
		fmt.Println("保存session成功")
	}
	v := session.Get("sessionid")
	fmt.Println("sessionid的值是:", v)

	context.Redirect(http.StatusMovedPermanently, "/index")

}

//判断是否登陆的中间件
func AuthRequiredSession() gin.HandlerFunc{
	return func(context *gin.Context) {
		session := sessions.Default(context)
		session_id := session.Get("sessionid")
		fmt.Println("sessionID+++++",session_id)
		if session_id == nil{
			context.HTML(200, "signin.html", gin.H{"message":"请先登录"})
			//context.Next()
			context.Abort()
		}else{
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
	fmt.Println("sessionid===================",session.Get("sessionid"))
	session.Delete("sessionid")
	session.Save()
	c.HTML(200, "signin.html", gin.H{"message":"session已清除"})
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


