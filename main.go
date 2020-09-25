package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/controllers"
	"gorm_demo/internal/db"
	"gorm_demo/internal/helpers"
	"gorm_demo/internal/models"
	"html/template"
	"net/http"
)



func main() {
	db.InitDbConnection()
	//router_w.InitEngine()
	db.W_Db.AutoMigrate(&models.User{},&models.Article{},&models.Comment{},&models.Label{},&models.Sort{},&models.UserFriend{})
	//a := "12333B"
	//b := strings.Split(a,"")
	//for _,c := range a{
	//	fmt.Println("////////////////////////////////////////////////////////////",unicode.IsLetter(c))
	//}

	//article := models.FindArticleById("1")
	//label := models.FindLabelById("1")
	//db.W_Db.Model(&article).Association("Labels").Append(label)
	router := gin.Default()
	setTemplate(router)

	//在路由中使用中间件调用store
	var store = cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sessionid", store))

	//路由
	router.GET("/signin", controllers.Signin)
	router.POST("/login",controllers.Login)
	router.GET("/register", controllers.Register)
	router.POST("/signup",controllers.Signup)
	router.GET("/logout",controllers.Logout)
	router.Static("/static", "./internal/static")
	router.GET("/index",controllers.ArticleIndex)
	router.GET("/index/:page",controllers.ArticlePageIndex)
	router.GET("/label/:id",controllers.Show_Article_By_Label)

	router.StaticFS("/file", http.Dir("public"))
	auth := router.Group("")
	auth.Use(controllers.AuthRequiredSession())
	{
		auth.GET("/article/:id",controllers.Show_Article)
		auth.GET("/user/:user_session",controllers.ShowUser)
		auth.GET("/new_article",controllers.NewArticle)
		auth.POST("/save_article",controllers.SaveArticle)
		auth.GET("/user_setting",controllers.UserSetting)
		auth.GET("/_basic_setting",controllers.BasicSetting)
		auth.POST("/update_user",controllers.UpdateUser)
		//auth.GET("/show_user_articles",controllers.ShowUserArticles)
		auth.GET("/edit_article/:id",controllers.EditArticle)
		auth.POST("/update_article/:id",controllers.UpdateArticle)
		auth.GET("/show_user_articles/:page",controllers.ShowUserArticles)
		auth.GET("/picture_recognition",controllers.PictureRecognition)
		auth.POST("/submit_picture_recognition",controllers.SubmitPictureRecognition)


	}

	//auth := router.Group("")
	//auth.Use(controllers.AuthRequired())
	//{
	//	auth.GET("/index",controllers.ArticleIndex)
	//	auth.GET("/article/:id",controllers.Show_Article)
	//	auth.GET("/label/:id",controllers.Show_Article_By_Label)
	//}

	router.Run(":8080")


}


func setTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": helpers.DateFormat,
		"truncate":   helpers.Truncate,
		"replaceHtml": helpers.ReplaceHtml,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob("./internal/views/*/*")
}

