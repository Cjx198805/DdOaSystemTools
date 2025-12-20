package main

import (
	"fmt"
	"github.com/ddoalistdownload/backend/config"
	"github.com/ddoalistdownload/backend/controller"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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

	// API分组
	api := router.Group("/api/v1")
	{
		// 集团公司管理
		company := api.Group("/company")
		company.GET("", companyController.List)
		company.POST("", companyController.Create)
		company.GET("/:id", companyController.Get)
		company.PUT("/:id", companyController.Update)
		company.DELETE("/:id", companyController.Delete)
		company.GET("/tree", companyController.GetTree)

		// 身份验证（免登）
		sso := api.Group("/sso")
		sso.GET("/config", ssoController.GetConfig)
		sso.POST("/config", ssoController.UpdateConfig)
		sso.GET("/test", ssoController.TestSSO)

		// AccessToken管理
		accessToken := api.Group("/access-token")
		accessToken.GET("", accessTokenController.GetAccessToken)
		accessToken.POST("", accessTokenController.CreateAccessToken)
		accessToken.PUT("/:id", accessTokenController.UpdateAccessToken)
		accessToken.DELETE("/:id", accessTokenController.DeleteAccessToken)
		accessToken.GET("/list", accessTokenController.GetAccessTokenList)
		accessToken.POST("/refresh", accessTokenController.RefreshAccessToken)
		accessToken.POST("/test", accessTokenController.TestAccessToken)

		// API配置管理
		apiConfig := api.Group("/api-config")
		apiConfig.GET("", apiConfigController.List)
		apiConfig.POST("", apiConfigController.Create)
		apiConfig.GET("/:id", apiConfigController.Get)
		apiConfig.PUT("/:id", apiConfigController.Update)
		apiConfig.DELETE("/:id", apiConfigController.Delete)
		apiConfig.POST("/test", apiConfigController.Test)

		// 用户管理
		user := api.Group("/user")
		user.GET("", userController.List)
		user.POST("", userController.Create)
		user.GET("/:id", userController.Get)
		user.PUT("/:id", userController.Update)
		user.DELETE("/:id", userController.Delete)
		user.PUT("/:id/reset-password", userController.ResetPassword)
		user.PUT("/update-password", userController.UpdatePassword)
		user.GET("/:id/roles", userController.GetRoles)
		user.PUT("/:id/assign-roles", userController.AssignRoles)
		user.POST("/login", userController.Login)

		// 角色管理
		role := api.Group("/role")
		role.GET("", roleController.List)
		role.POST("", roleController.Create)
		role.GET("/:id", roleController.Get)
		role.PUT("/:id", roleController.Update)
		role.DELETE("/:id", roleController.Delete)
		role.GET("/:id/menus", roleController.GetMenus)
		role.PUT("/:id/assign-menus", roleController.AssignMenus)

		// 菜单管理
		menu := api.Group("/menu")
		menu.GET("", menuController.List)
		menu.POST("", menuController.Create)
		menu.GET("/:id", menuController.Get)
		menu.PUT("/:id", menuController.Update)
		menu.DELETE("/:id", menuController.Delete)
		menu.GET("/tree", menuController.GetTree)
		menu.GET("/parent/:parent_id", menuController.GetByParentID)
		menu.GET("/all", menuController.GetAll)

		// 字段权限管理
		fieldPermission := api.Group("/field-permission")
		fieldPermission.GET("", fieldPermissionController.List)
		fieldPermission.POST("", fieldPermissionController.Create)
		fieldPermission.GET("/:id", fieldPermissionController.Get)
		fieldPermission.PUT("/:id", fieldPermissionController.Update)
		fieldPermission.DELETE("/:id", fieldPermissionController.Delete)
		fieldPermission.GET("/role/:role_id/module/:module", fieldPermissionController.GetByRoleAndModule)

		// 数据字典管理
		dataDictionary := api.Group("/data-dictionary")
		dataDictionary.GET("", dataDictionaryController.List)
		dataDictionary.POST("", dataDictionaryController.Create)
		dataDictionary.GET("/:id", dataDictionaryController.Get)
		dataDictionary.PUT("/:id", dataDictionaryController.Update)
		dataDictionary.DELETE("/:id", dataDictionaryController.Delete)
		dataDictionary.GET("/module/:module", dataDictionaryController.GetByModule)
		dataDictionary.GET("/module/:module/field/:field", dataDictionaryController.GetByModuleAndField)
	}
}
