package main

import (
	"my_template/cmn"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 這是一個完整的使用示例，展示如何使用中間件注冊模組
// 運行方式: go run examples/middleware_server.go

func main() {
	// 1. 初始化配置（必須）
	if err := cmn.ViperInit(".conf_linux.json"); err != nil {
		panic("初始化配置失敗: " + err.Error())
	}

	// 2. 初始化日誌（必須）
	if err := cmn.LoggerInit(); err != nil {
		panic("初始化日誌失敗: " + err.Error())
	}

	// 3. 創建 gin 引擎
	router := gin.New()

	// 4. 創建中間件注冊器
	registry := cmn.NewMiddlewareRegistry(router)

	// 5. 註冊中間件
	// 方式1: 使用默認中間件（Recovery + Logger + CORS）
	// registry.ApplyDefault()

	// 方式2: 手動註冊中間件（推薦，可以自定義順序和配置）
	registry.RegisterMultiple(
		cmn.RecoveryMiddleware(),              // 捕獲 panic
		cmn.LoggerMiddleware(),                // 記錄請求日誌
		cmn.CORSMiddleware(),                  // 處理跨域
		cmn.RateLimitByIP(100, 200),           // IP 限流（每秒100個請求）
		cmn.TimeoutMiddleware(30*time.Second), // 30秒超時
	)

	// 6. 應用中間件
	registry.Apply()

	// 7. 設置路由
	setupRoutes(router)

	// 8. 啟動服務器
	router.Run(":8080")
}

func setupRoutes(router *gin.Engine) {
	// 健康檢查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// API v1 路由組
	v1 := router.Group("/api/v1")
	{
		// 公開路由
		v1.POST("/login", loginHandler)
		v1.POST("/register", registerHandler)

		// 測試路由
		v1.GET("/test-panic", testPanicHandler)
		v1.GET("/test-slow", testSlowHandler)

		// 需要認證的路由組
		auth := v1.Group("")
		auth.Use(cmn.AuthMiddleware())
		{
			auth.GET("/profile", profileHandler)
			auth.PUT("/profile", updateProfileHandler)
			auth.GET("/posts", listPostsHandler)
			auth.POST("/posts", createPostHandler)
		}

		// 可選認證的路由組
		optional := v1.Group("")
		optional.Use(cmn.OptionalAuthMiddleware())
		{
			optional.GET("/posts/public", publicPostsHandler)
		}
	}
}

// ============ 處理器函數 ============

func loginHandler(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "請求參數錯誤",
			"error":   err.Error(),
		})
		return
	}

	// 這裡應該驗證用戶名和密碼
	// 為了演示，我們假設驗證成功
	token, err := cmn.GenerateToken("user_"+req.Username, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成 token 失敗",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登錄成功",
		"data": gin.H{
			"token":    token,
			"username": req.Username,
		},
	})
}

func registerHandler(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "請求參數錯誤",
			"error":   err.Error(),
		})
		return
	}

	// 這裡應該保存用戶到數據庫
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "註冊成功",
		"data": gin.H{
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

func profileHandler(c *gin.Context) {
	// 從上下文獲取用戶信息
	userId, _ := cmn.GetUserId(c)
	username, _ := cmn.GetUsername(c)
	claims, _ := cmn.GetClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "獲取成功",
		"data": gin.H{
			"user_id":  userId,
			"username": username,
			"claims":   claims,
		},
	})
}

func updateProfileHandler(c *gin.Context) {
	userId, _ := cmn.GetUserId(c)

	type UpdateRequest struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "請求參數錯誤",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "更新成功",
		"data": gin.H{
			"user_id":  userId,
			"nickname": req.Nickname,
			"email":    req.Email,
		},
	})
}

func listPostsHandler(c *gin.Context) {
	userId, _ := cmn.GetUserId(c)

	// 模擬文章列表
	posts := []gin.H{
		{"id": 1, "title": "文章1", "author_id": userId},
		{"id": 2, "title": "文章2", "author_id": userId},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "獲取成功",
		"data":    posts,
	})
}

func createPostHandler(c *gin.Context) {
	userId, _ := cmn.GetUserId(c)

	type CreatePostRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "請求參數錯誤",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "創建成功",
		"data": gin.H{
			"id":        123,
			"title":     req.Title,
			"content":   req.Content,
			"author_id": userId,
		},
	})
}

func publicPostsHandler(c *gin.Context) {
	// 檢查是否有用戶登錄
	userId, hasUser := cmn.GetUserId(c)

	posts := []gin.H{
		{"id": 1, "title": "公開文章1", "is_public": true},
		{"id": 2, "title": "公開文章2", "is_public": true},
	}

	response := gin.H{
		"code":    http.StatusOK,
		"message": "獲取成功",
		"data":    posts,
	}

	if hasUser {
		response["user_id"] = userId
	}

	c.JSON(http.StatusOK, response)
}

func testPanicHandler(c *gin.Context) {
	// 測試 Recovery 中間件
	panic("這是一個測試 panic")
}

func testSlowHandler(c *gin.Context) {
	// 測試慢請求日誌
	time.Sleep(500 * time.Millisecond)
	c.JSON(http.StatusOK, gin.H{
		"message": "這是一個慢請求",
		"delay":   "500ms",
	})
}
