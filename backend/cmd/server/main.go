package main

import (
	"net"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"what-to-eat/config"
	"what-to-eat/internal/handler"
	"what-to-eat/internal/repository"
	"what-to-eat/internal/service"
	"what-to-eat/pkg/logger"
	"what-to-eat/pkg/middleware"
)

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	// 加载配置
	configPath := filepath.Join("config", "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// 配置加载失败时使用默认日志
		logger.InitDefault()
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 初始化日志
	if err := logger.Init(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		OutputPath: cfg.Log.OutputPath,
	}); err != nil {
		logger.InitDefault()
		logger.Error("Failed to init logger, using default", zap.Error(err))
	}
	defer logger.Sync()

	logger.Info("Config loaded",
		zap.String("port", cfg.Server.Port),
		zap.String("db_host", cfg.Database.Host),
		zap.String("db_name", cfg.Database.DBName),
	)

	// 初始化数据库（自动创建数据库和表）
	if err := repository.InitDB(&cfg.Database); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	db := repository.GetDB()

	// 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	decisionRepo := repository.NewDecisionRepository(db)

	// 初始化 Service
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	menuService := service.NewMenuService(menuRepo, restaurantRepo)
	decisionService := service.NewDecisionService(decisionRepo, menuRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authService)
	menuHandler := handler.NewMenuHandler(menuService)
	decisionHandler := handler.NewDecisionHandler(decisionService)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 实例
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.GinLogger()) // 使用自定义的 zap 日志中间件

	// 全局中间件
	r.Use(middleware.CORS())

	// 公开路由（无需认证）
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/guest", authHandler.GuestLogin) // 游客登录
		}
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		// 菜单管理
		menus := protected.Group("/menus")
		{
			menus.GET("", menuHandler.List)
			menus.POST("", menuHandler.Create)
			menus.DELETE("/:id", menuHandler.Delete)
		}

		// 餐厅列表（用于下拉选择）
		protected.GET("/restaurants", menuHandler.ListRestaurants)

		// 决策
		protected.POST("/decide", decisionHandler.Decide)
		protected.GET("/history", decisionHandler.History)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 获取本机IP
	localIP := getLocalIP()
	address := ":" + cfg.Server.Port

	logger.Info("Server starting",
		zap.String("port", cfg.Server.Port),
		zap.String("local_ip", localIP),
		zap.String("public_url", "http://"+localIP+address),
		zap.String("health_check_url", "http://"+localIP+address+"/health"),
	)

	// 自定义启动信息
	logger.Info("Listening and serving HTTP",
		zap.String("address", address),
		zap.String("local_ip", localIP),
		zap.String("public_url", "http://"+localIP+address),
		zap.String("health_check_url", "http://"+localIP+address+"/health"),
	)

	if err := r.Run(address); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
