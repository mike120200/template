package cmn

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 設置路由和中間件的示例
func SetupRouter() *gin.Engine {
	// 創建 gin 引擎
	router := gin.New()

	// 創建中間件注冊器
	registry := NewMiddlewareRegistry(router)

	// 方式1: 使用默認中間件集合
	registry.ApplyDefault()

	// 方式2: 手動注冊中間件
	// registry.Register(RecoveryMiddleware())
	// registry.Register(LoggerMiddleware())
	// registry.Register(CORSMiddleware())

	// 方式3: 批量注冊中間件
	// registry.RegisterMultiple(
	// 	RecoveryMiddleware(),
	// 	LoggerMiddleware(),
	// 	CORSMiddleware(),
	// 	RateLimitByIP(10, 20),
	// 	TimeoutMiddleware(30 * time.Second),
	// )

	// 應用所有已注冊的中間件
	registry.Apply()

	return router
}

// SetupRouterWithCustomConfig 使用自定義配置設置路由
func SetupRouterWithCustomConfig() *gin.Engine {
	router := gin.New()
	registry := NewMiddlewareRegistry(router)

	// 自定義 CORS 配置
	corsConfig := &CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	// 自定義日志配置
	loggerConfig := DefaultLoggerConfig()
	loggerConfig.SkipPaths = []string{"/health", "/metrics"}

	// 自定義認證配置
	authConfig := DefaultAuthConfig()
	authConfig.SkipPaths = []string{"/api/v1/login", "/api/v1/register", "/health"}

	// 自定義限流配置
	rateLimitConfig := DefaultRateLimitConfig()
	rateLimitConfig.Rate = 100     // 每秒100個請求
	rateLimitConfig.Capacity = 200 // 桶容量200

	// 註冊中間件
	registry.RegisterMultiple(
		RecoveryMiddleware(),
		LoggerMiddlewareWithConfigCustom(loggerConfig),
		CORSMiddlewareWithConfig(corsConfig),
		RateLimitMiddleware(rateLimitConfig),
		TimeoutMiddleware(30000000000), // 30秒
	)

	registry.Apply()

	// 設置路由組，添加認證中間件
	api := router.Group("/api/v1")
	{
		// 公開路由
		api.POST("/login", loginHandler)
		api.POST("/register", registerHandler)

		// 需要認證的路由
		auth := api.Group("")
		auth.Use(AuthMiddlewareWithConfig(authConfig))
		{
			auth.GET("/profile", profileHandler)
			auth.PUT("/profile", updateProfileHandler)
		}
	}

	return router
}

// 示例處理器函數
func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "login"})
}

func registerHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "register"})
}

func profileHandler(c *gin.Context) {
	// 從上下文獲取用戶信息
	userId, _ := GetUserId(c)
	username, _ := GetUsername(c)

	c.JSON(200, gin.H{
		"user_id":  userId,
		"username": username,
	})
}

func updateProfileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "profile updated"})
}
