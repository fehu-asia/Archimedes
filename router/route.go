package router

import (
	"fehu/common/lib"
	"fehu/controller"
	"fehu/controller/base"
	"fehu/middleware"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	//programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	//router.Use(gin.LoggerWithFormatter(middleware.CustomLogger))
	//gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	//	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	//}

	// 错误中间件
	router.Use(middleware.RecoveryMiddleware())
	// 如果开启了限流，使用限流中间件
	if lib.GetIntConf("base.http.limitStatus") == 1 {
		router.Use(middleware.Limiter())
	}
	// 刷新token
	router.Use(middleware.RefreshToken())
	// 加解密中间件
	router.Use(middleware.Crpy())
	// 请求日志
	router.Use(middleware.RequestLog())
	// IP限制
	router.Use(middleware.IPAuthMiddleware())
	router.Use(middlewares...)
	//router.Use(func(context *gin.Context) {
	//})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//base 不需要加密的接口
	v1 := router.Group("/base")
	v1.Use(middleware.TranslationMiddleware())
	{
		base.BaseRegister(v1)
		controller.DemoRegister(v1)
	}
	//
	////非登陆接口
	//store := sessions.NewCookieStore([]byte("secret"))
	apiNormalGroup := router.Group("/login")
	apiNormalGroup.Use()
	{

		controller.NoLoginRequiredRegister(apiNormalGroup)
	}
	// 需要登录的接口
	apiNormalGroup.Use(middleware.JwtAuthMiddleware())
	{
		controller.ApiRegister(apiNormalGroup)
	}
	//
	////登陆接口
	//apiAuthGroup := router.Group("/api")
	//apiAuthGroup.Use(
	//	sessions.Sessions("mysession", store),
	//	middleware.RecoveryMiddleware(),
	//	middleware.RequestLog(),
	//	middleware.SessionAuthMiddleware(),
	//	middleware.TranslationMiddleware())
	//{
	//	controller.ApiLoginRegister(apiAuthGroup)
	//}
	return router
}
