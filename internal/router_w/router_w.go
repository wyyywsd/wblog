package router_w

//var (
//	tomlFilePath,mode string
//)
//
////初始化engin
//func InitEngine() (engine *gin.Engine,tomlConfig *config.MyApiConfig,err error){
//	flag.StringVar(&tomlFilePath, "config", "/Users/mac/work_space/wblog/internal/config/database.toml", "服务配置文件")
//	flag.StringVar(&mode, "mode", "release", "模型-debug还是release还是test")
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	flag.Parse()
//	gin.SetMode(mode)
//
//	//解析配置文件
//	tomlConfig,err = config.UnmarshalConfig(tomlFilePath)
//	if err != nil {
//		return nil ,tomlConfig,fmt.Errorf("解析配置文件出错")
//	}
//	//初始化路由
//	engine = gin.New()
//	//初始化中间件
//	initMiddleware(engine,tomlConfig)
//	//初始化路由
//	//addHandler(engine)
//	//返回服务器对象实例，配置文件和错误信息
//	return engine,tomlConfig,nil
//}
//
//
////初始化路由
////func addHandler(engine *gin.Engine) {
////	//子路由
////	userGroup(engine)
////}
//
////initMiddleware 初始化中间件
//func initMiddleware(router *gin.Engine, apiConfig *config.MyApiConfig) {
//	//捕获错误，
//	router.Use(gin.Recovery())
//	//将配置文件引入中间件
//	router.Use(middleware.SetMiddleware(middleware.MiddlewareConfig, apiConfig))
//	//返回一个DBServer的结构体
//	testdbCfg, ok := apiConfig.DBServerConf(middleware.MiddlewareTestDB)
//	if !ok {
//		panic(fmt.Sprintf("initMiddleware: %v配置不存在\n", middleware.MiddlewareTestDB))
//	}
//
//	testdborm, err := testdbCfg.NewGormDB(20)
//	if err != nil {
//		panic(fmt.Sprintf("initMiddleware: 连接数据库%v出错, err:%v\n", middleware.MiddlewareTestDBORM, err))
//	}
//	router.Use(middleware.SetMiddleware(middleware.MiddlewareTestDBORM, testdborm))
//
//	usercacheCfg, ok := apiConfig.RedisServerConf(middleware.MiddlewareUserCache)
//	if !ok {
//		panic(fmt.Sprintf("initMiddleware: %v配置不存在\n", middleware.MiddlewareUserCache))
//	}
//	usercache := usercacheCfg.NewRedisPool(15)
//	router.Use(middleware.SetMiddleware(middleware.MiddlewareUserCache, usercache))
//}
