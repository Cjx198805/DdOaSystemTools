package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ddoalistdownload/backend/config"
	"github.com/ddoalistdownload/backend/controller"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("启动 DdOaListDownload 后端服务")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		logrus.Fatalf("初始化MySQL失败: %v", err)
	}

	// 初始化Redis连接
	if err := database.InitRedis(&cfg.Redis); err != nil {
		logrus.Fatalf("初始化Redis失败: %v", err)
	}

	// 执行数据库迁移
	if err := database.MigrateDB(); err != nil {
		logrus.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建Gin引擎
	router := gin.Default()

	// 添加中间件
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.RecoverMiddleware())

	// 注册路由
	registerRoutes(router)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logrus.Infof("服务器正在启动，监听地址: %s", serverAddr)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			logrus.Fatalf("启动服务器失败: %v", err)
		}
	}()

	<-quit
	logrus.Info("正在关闭服务器...")

	// 关闭数据库连接
	sqlDB, _ := database.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// 关闭Redis连接
	database.CloseRedis()

	logrus.Info("服务器已关闭")
}

// registerRoutes 注册路由
func registerRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "DdOaListDownload 服务运行正常",
		})
	})

	// 初始化服务和控制器
	companyController := controller.NewCompanyController()
	accessTokenController := controller.NewAccessTokenController()
	ssoController := controller.NewSSOController()
	apiConfigController := controller.NewAPIConfigController()
	userController := controller.NewUserController()
	roleController := controller.NewRoleController()
	menuController := controller.NewMenuController()
	fieldPermissionController := controller.NewFieldPermissionController()
	dataDictionaryController := controller.NewDataDictionaryController()
	downloadTaskController := controller.NewDownloadTaskController()
	apiTestController := controller.NewAPITestController()

	// API分组
	api := router.Group("/api/v1")
	{
		// 登录路由（不需要认证）
		api.POST("/user/login", userController.Login)

		// 需要认证的路由分组
		authAPI := api.Group("")
		authAPI.Use(middleware.AuthMiddleware())
		{
			// 集团公司管理
			company := authAPI.Group("/company")
			company.Use(middleware.PermissionMiddleware("company:manage"))
			company.GET("", companyController.List)
			company.POST("", companyController.Create)
			company.GET("/:id", companyController.Get)
			company.PUT("/:id", companyController.Update)
			company.DELETE("/:id", companyController.Delete)
			company.GET("/tree", companyController.GetTree)

			// 身份验证（免登）
			sso := authAPI.Group("/sso")
			sso.Use(middleware.PermissionMiddleware("sso:manage"))
			sso.GET("/config", ssoController.GetConfig)
			sso.POST("/config", ssoController.UpdateConfig)
			sso.GET("/test", ssoController.TestSSO)

			// AccessToken管理
			accessToken := authAPI.Group("/access-token")
			accessToken.Use(middleware.PermissionMiddleware("access_token:manage"))
			accessToken.GET("", accessTokenController.GetAccessToken)
			accessToken.POST("", accessTokenController.CreateAccessToken)
			accessToken.PUT("/:id", accessTokenController.UpdateAccessToken)
			accessToken.DELETE("/:id", accessTokenController.DeleteAccessToken)
			accessToken.GET("/list", accessTokenController.GetAccessTokenList)
			accessToken.POST("/refresh", accessTokenController.RefreshAccessToken)
			accessToken.POST("/test", accessTokenController.TestAccessToken)

			// API配置管理
			apiConfig := authAPI.Group("/api-config")
			apiConfig.Use(middleware.PermissionMiddleware("api_config:manage"))
			apiConfig.GET("", apiConfigController.List)
			apiConfig.POST("", apiConfigController.Create)
			apiConfig.GET("/:id", apiConfigController.Get)
			apiConfig.PUT("/:id", apiConfigController.Update)
			apiConfig.DELETE("/:id", apiConfigController.Delete)
			apiConfig.POST("/test", apiConfigController.Test)

			// 用户管理
			user := authAPI.Group("/user")
			user.Use(middleware.PermissionMiddleware("user:manage"))
			user.GET("", userController.List)
			user.POST("", userController.Create)
			user.GET("/:id", userController.Get)
			user.PUT("/:id", userController.Update)
			user.DELETE("/:id", userController.Delete)
			user.PUT("/:id/reset-password", userController.ResetPassword)
			user.PUT("/update-password", userController.UpdatePassword)
			user.GET("/:id/roles", userController.GetRoles)
			user.PUT("/:id/assign-roles", userController.AssignRoles)

			// 角色管理
			role := authAPI.Group("/role")
			role.Use(middleware.PermissionMiddleware("role:manage"))
			role.GET("", roleController.List)
			role.POST("", roleController.Create)
			role.GET("/:id", roleController.Get)
			role.PUT("/:id", roleController.Update)
			role.DELETE("/:id", roleController.Delete)
			role.GET("/:id/menus", roleController.GetMenus)
			role.PUT("/:id/assign-menus", roleController.AssignMenus)

			// 菜单管理
			menu := authAPI.Group("/menu")
			menu.Use(middleware.PermissionMiddleware("menu:manage"))
			menu.GET("", menuController.List)
			menu.POST("", menuController.Create)
			menu.GET("/:id", menuController.Get)
			menu.PUT("/:id", menuController.Update)
			menu.DELETE("/:id", menuController.Delete)
			menu.GET("/tree", menuController.GetTree)
			menu.GET("/parent/:parent_id", menuController.GetByParentID)
			menu.GET("/all", menuController.GetAll)

			// 字段权限管理
			fieldPermission := authAPI.Group("/field-permission")
			fieldPermission.Use(middleware.PermissionMiddleware("field_permission:manage"))
			fieldPermission.GET("", fieldPermissionController.List)
			fieldPermission.POST("", fieldPermissionController.Create)
			fieldPermission.GET("/:id", fieldPermissionController.Get)
			fieldPermission.PUT("/:id", fieldPermissionController.Update)
			fieldPermission.DELETE("/:id", fieldPermissionController.Delete)
			fieldPermission.GET("/role/:role_id/module/:module", fieldPermissionController.GetByRoleAndModule)

			// 数据字典管理
			dataDictionary := authAPI.Group("/data-dictionary")
			dataDictionary.Use(middleware.PermissionMiddleware("data_dictionary:manage"))
			dataDictionary.GET("", dataDictionaryController.List)
			dataDictionary.POST("", dataDictionaryController.Create)
			dataDictionary.GET("/:id", dataDictionaryController.Get)
			dataDictionary.PUT("/:id", dataDictionaryController.Update)
			dataDictionary.DELETE("/:id", dataDictionaryController.Delete)
			dataDictionary.GET("/module/:module", dataDictionaryController.GetByModule)
			dataDictionary.GET("/module/:module/field/:field", dataDictionaryController.GetByModuleAndField)

			// 下载任务管理
			downloadTask := authAPI.Group("/download-task")
			downloadTask.Use(middleware.PermissionMiddleware("download_task:manage"))
			downloadTask.GET("", downloadTaskController.List)
			downloadTask.POST("", downloadTaskController.Create)
			downloadTask.GET("/user/:user_id", downloadTaskController.GetTaskByUserID)
			downloadTask.GET("/result/:task_id", downloadTaskController.GetResult)
			downloadTask.GET("/:id", downloadTaskController.Get)
			downloadTask.DELETE("/:id", downloadTaskController.Delete)

			// API测试管理
			apiTest := authAPI.Group("/api-test")
			apiTest.Use(middleware.PermissionMiddleware("api_test:manage"))

			// 测试用例相关路由
			testCase := apiTest.Group("/case")
			testCase.GET("", apiTestController.ListTestCases)
			testCase.POST("", apiTestController.CreateTestCase)
			testCase.GET("/:id", apiTestController.GetTestCase)
			testCase.PUT("/:id", apiTestController.UpdateTestCase)
			testCase.DELETE("/:id", apiTestController.DeleteTestCase)
			testCase.POST("/:id/run", apiTestController.RunTestCase)

			// 测试历史记录相关路由
			testHistory := apiTest.Group("/history")
			testHistory.GET("", apiTestController.ListTestHistory)
			testHistory.GET("/:id", apiTestController.GetTestHistory)
			testHistory.DELETE("/:id", apiTestController.DeleteTestHistory)
			testHistory.POST("/clear", apiTestController.ClearTestHistory)
		}
	}
}
