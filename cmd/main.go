// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"llmapisrv/config"
	"llmapisrv/internal/api"
	"llmapisrv/internal/api/admin"
	"llmapisrv/internal/api/chat"
	"llmapisrv/internal/api/dashboard"
	"llmapisrv/internal/middleware"
	"llmapisrv/internal/service"
	"llmapisrv/pkg/cache"
	"llmapisrv/pkg/cron"
	"llmapisrv/pkg/logger"
	"llmapisrv/pkg/oss"
	"llmapisrv/pkg/queue"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Setup(config.AppConfig.Logger)

	gormConf := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info), // 打印所有 SQL,
	}

	// 初始化网关数据库
	gatewayDB, err := gorm.Open(mysql.Open(config.AppConfig.Database.GatewayDSN), gormConf)
	if err != nil {
		log.Fatalf("Failed to connect to gateway database: %v", err)
	}

	// 初始化New API数据库
	newAPIDB, err := gorm.Open(mysql.Open(config.AppConfig.Database.NewAPIDSN), gormConf)
	if err != nil {
		log.Fatalf("Failed to connect to New API database: %v", err)
	}

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	// 初始化缓存
	redisCache := cache.NewRedisCache(redisClient)

	// 初始化队列
	redisQueue := queue.NewRedisQueue(redisClient)

	// 初始化OSS客户端
	ossClient, err := oss.NewOSSClient(&config.AppConfig.OSS.Aliyun)
	if err != nil {
		log.Fatalf("Failed to initialize OSS client: %v", err)
	}

	// 初始化同步服务
	syncService := service.NewSyncService(gatewayDB, newAPIDB, &config.AppConfig, redisCache)

	// 初始化服务
	userService := service.NewUserService(gatewayDB, newAPIDB, redisCache, syncService)
	logService := service.NewLogService(gatewayDB, newAPIDB, &config.AppConfig)
	newAPIService := service.NewNewAPIService(&config.AppConfig, redisCache)
	modelService := service.NewModelService(gatewayDB, newAPIDB, &config.AppConfig)
	redemptionService := service.NewRedemptionService(gatewayDB)

	// 初始化处理器
	statusHandler := api.NewStatusHandler(newAPIService)
	billingHandler := dashboard.NewBillingHandler(newAPIService, userService)
	pricingHandler := api.NewPricingHandler(newAPIService, modelService)
	chatHandler := chat.NewChatHandler(newAPIService, logService, redisQueue)
	redemptionHandler := api.NewRedemptionHandler(newAPIService, redemptionService, userService)
	adminRedemptionHandler := admin.NewRedemptionAdminHandler(redemptionService, userService)
	adminUploadHandler := admin.NewUploadHandler(ossClient)
	logHandler := api.NewLogHandler(logService)
	proxyHandler := api.NewProxyHandler(ossClient)

	// 启动定时任务
	cronManager := cron.NewCronManager(logService, modelService, syncService)
	// cronManager.Start()
	defer cronManager.Stop()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.TraceIDMiddleware())
	r.Use(middleware.ClientInfoMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.MetricsMiddleware())

	// 添加静态文件支持
	r.Static("/static", "./web")
	r.LoadHTMLGlob("web/*.html")

	// 添加前端页面路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "主页",
		})
	})

	r.GET("/pricing", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pricing.html", gin.H{
			"title": "价格",
		})
	})

	r.GET("/redeem", func(c *gin.Context) {
		c.HTML(http.StatusOK, "redeem.html", gin.H{
			"title": "兑换码",
		})
	})

	r.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.html", gin.H{
			"title": "状态",
		})
	})

	// 注册路由
	// 健康检查
	r.GET("/api/about", statusHandler.HealthCheck)

	// Prometheus 指标
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 价格查询
	r.GET("/api/pricing", pricingHandler.GetPricing)

	// 图片代理接口（公开接口）
	r.GET("/api/image-proxy", proxyHandler.ImageProxy)

	// 认证路由
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware(userService, redisCache))

	// 账单相关
	authGroup.GET("/v1/dashboard/billing/subscription", billingHandler.GetSubscription)
	authGroup.GET("/v1/dashboard/billing/usage", billingHandler.GetUsage)

	// 聊天完成 openai兼容的接口调用方式
	authGroup.POST("/v1/chat/completions", chatHandler.ChatCompletions)

	// 兑换码
	authGroup.GET("/api/redeem", redemptionHandler.RedeemCodeInfo)
	authGroup.POST("/api/redeem", redemptionHandler.RedeemCode)

	// 日志查询
	authGroup.GET("/api/logs", logHandler.GetLogs)

	// 管理员路由
	adminGroup := r.Group("api/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware(&config.AppConfig))
	{
		// 管理员兑换码管理
		adminGroup.POST("/redemption/generate", adminRedemptionHandler.GenerateCodes)
		adminGroup.POST("/quota/add", adminRedemptionHandler.AddQuota)

		// 管理员手动同步
		adminGroup.POST("/sync/user", admin.NewSyncHandler(syncService).SyncUser)
		adminGroup.POST("/sync/logs", admin.NewSyncHandler(syncService).SyncLogs)
		adminGroup.POST("/sync/all", admin.NewSyncHandler(syncService).SyncAll)

		// 管理员手动删除旧日志
		adminGroup.POST("/cleanup/logs", logHandler.CleanupOldLogs)

		// 管理员图片上传
		adminGroup.POST("/upload/image", adminUploadHandler.UploadImage)

	}

	// 启动服务
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
