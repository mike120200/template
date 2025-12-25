# Gin 中間件注冊模塊

這是一個基於 Gin 框架的中間件注冊和管理模塊，提供了靈活的中間件註冊機制和多個常用中間件。

## 功能特性

- ✅ **中間件注冊器**: 靈活的中間件註冊和管理
- ✅ **CORS中間件**: 跨域資源共享支持
- ✅ **日誌中間件**: 請求日誌記錄和慢請求監控
- ✅ **恢復中間件**: Panic 捕獲和恢復
- ✅ **認證中間件**: JWT 認證支持
- ✅ **限流中間件**: 基於令牌桶算法的限流
- ✅ **超時中間件**: 請求超時控制

## 安裝依賴

首先需要安裝 Gin 框架：

```bash
go get -u github.com/gin-gonic/gin
```

## 快速開始

### 1. 使用默認中間件

```go
package main

import (
    "my_template/cmn"
    "github.com/gin-gonic/gin"
)

func main() {
    // 創建 gin 引擎
    router := gin.New()
    
    // 創建中間件注冊器並使用默認中間件
    registry := cmn.NewMiddlewareRegistry(router)
    registry.ApplyDefault().Apply()
    
    // 設置路由
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    // 啟動服務
    router.Run(":8080")
}
```

### 2. 手動註冊中間件

```go
router := gin.New()
registry := cmn.NewMiddlewareRegistry(router)

// 單個註冊
registry.Register(cmn.RecoveryMiddleware())
registry.Register(cmn.LoggerMiddleware())
registry.Register(cmn.CORSMiddleware())

// 批量註冊
registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),
    cmn.LoggerMiddleware(),
    cmn.CORSMiddleware(),
    cmn.RateLimitByIP(10, 20),
    cmn.TimeoutMiddleware(30 * time.Second),
)

registry.Apply()
```

### 3. 使用自定義配置

```go
router := gin.New()
registry := cmn.NewMiddlewareRegistry(router)

// 自定義 CORS 配置
corsConfig := &cmn.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           86400,
}

// 自定義日誌配置
loggerConfig := cmn.DefaultLoggerConfig()
loggerConfig.SkipPaths = []string{"/health", "/metrics"}
loggerConfig.SlowThreshold = 500 * time.Millisecond

// 自定義認證配置
authConfig := cmn.DefaultAuthConfig()
authConfig.SkipPaths = []string{"/api/v1/login", "/api/v1/register"}

registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),
    cmn.LoggerMiddlewareWithConfigCustom(loggerConfig),
    cmn.CORSMiddlewareWithConfig(corsConfig),
    cmn.RateLimitByIP(100, 200),
)

registry.Apply()
```

## 中間件詳解

### 1. CORS 中間件

處理跨域請求：

```go
// 使用默認配置
router.Use(cmn.CORSMiddleware())

// 使用自定義配置
corsConfig := &cmn.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           86400,
}
router.Use(cmn.CORSMiddlewareWithConfig(corsConfig))
```

### 2. 日誌中間件

記錄所有 HTTP 請求：

```go
// 使用默認配置
router.Use(cmn.LoggerMiddleware())

// 使用自定義配置
loggerConfig := cmn.DefaultLoggerConfig()
loggerConfig.SkipPaths = []string{"/health", "/ping"}
loggerConfig.SlowThreshold = 200 * time.Millisecond
router.Use(cmn.LoggerMiddlewareWithConfigCustom(loggerConfig))
```

### 3. 恢復中間件

捕獲和恢復 panic：

```go
// 使用默認配置
router.Use(cmn.RecoveryMiddleware())

// 使用自定義配置（調試模式）
recoveryConfig := &cmn.RecoveryConfig{
    EnableStackTrace: true, // 在響應中包含堆棧信息
    ErrorMessage:     "服務器內部錯誤",
}
router.Use(cmn.RecoveryMiddlewareWithConfig(recoveryConfig))
```

### 4. 認證中間件

JWT 認證：

```go
// 全局認證
router.Use(cmn.AuthMiddleware())

// 路由組認證
api := router.Group("/api/v1")
{
    // 公開路由
    api.POST("/login", loginHandler)
    
    // 需要認證的路由
    auth := api.Group("")
    auth.Use(cmn.AuthMiddleware())
    {
        auth.GET("/profile", profileHandler)
        auth.PUT("/profile", updateProfileHandler)
    }
}

// 在處理器中獲取用戶信息
func profileHandler(c *gin.Context) {
    userId, _ := cmn.GetUserId(c)
    username, _ := cmn.GetUsername(c)
    claims, _ := cmn.GetClaims(c)
    
    c.JSON(200, gin.H{
        "user_id":  userId,
        "username": username,
    })
}
```

### 5. 限流中間件

基於令牌桶算法的限流：

```go
// 按 IP 限流（每秒10個請求，桶容量20）
router.Use(cmn.RateLimitByIP(10, 20))

// 按用戶限流
router.Use(cmn.RateLimitByUser(100, 200))

// 自定義限流配置
rateLimitConfig := cmn.DefaultRateLimitConfig()
rateLimitConfig.Rate = 100     // 每秒100個請求
rateLimitConfig.Capacity = 200 // 桶容量200
rateLimitConfig.KeyFunc = func(c *gin.Context) string {
    // 自定義限流鍵生成邏輯
    return c.ClientIP()
}
router.Use(cmn.RateLimitMiddleware(rateLimitConfig))
```

### 6. 超時中間件

控制請求超時：

```go
// 設置30秒超時
router.Use(cmn.TimeoutMiddleware(30 * time.Second))

// 使用自定義配置
timeoutConfig := &cmn.TimeoutConfig{
    Timeout: 30 * time.Second,
    ErrorHandler: func(c *gin.Context) {
        c.JSON(http.StatusRequestTimeout, gin.H{
            "code":    408,
            "message": "請求超時",
        })
    },
}
router.Use(cmn.TimeoutMiddlewareWithConfig(timeoutConfig))
```

### 7. 可選認證中間件

不強制要求認證，但如果提供了 token 會驗證：

```go
router.Use(cmn.OptionalAuthMiddleware())

func handler(c *gin.Context) {
    userId, exists := cmn.GetUserId(c)
    if exists {
        // 用戶已登錄
        c.JSON(200, gin.H{"user_id": userId})
    } else {
        // 訪客
        c.JSON(200, gin.H{"message": "guest"})
    }
}
```

## 完整示例

```go
package main

import (
    "my_template/cmn"
    "time"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化配置和日誌
    cmn.ViperInit(".conf_linux.json")
    cmn.LoggerInit()
    
    // 創建 gin 引擎
    router := gin.New()
    
    // 創建中間件注冊器
    registry := cmn.NewMiddlewareRegistry(router)
    
    // 註冊中間件
    registry.RegisterMultiple(
        cmn.RecoveryMiddleware(),
        cmn.LoggerMiddleware(),
        cmn.CORSMiddleware(),
        cmn.RateLimitByIP(100, 200),
        cmn.TimeoutMiddleware(30 * time.Second),
    )
    
    // 應用中間件
    registry.Apply()
    
    // 設置路由
    setupRoutes(router)
    
    // 啟動服務
    router.Run(":8080")
}

func setupRoutes(router *gin.Engine) {
    // 公開路由
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    // API 路由組
    api := router.Group("/api/v1")
    {
        // 認證相關
        api.POST("/login", loginHandler)
        api.POST("/register", registerHandler)
        
        // 需要認證的路由
        auth := api.Group("")
        auth.Use(cmn.AuthMiddleware())
        {
            auth.GET("/profile", profileHandler)
            auth.PUT("/profile", updateProfileHandler)
            auth.GET("/posts", listPostsHandler)
        }
    }
}

func loginHandler(c *gin.Context) {
    // 驗證用戶名密碼...
    token, _ := cmn.GenerateToken("user123", "testuser")
    c.JSON(200, gin.H{"token": token})
}

func registerHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "registered"})
}

func profileHandler(c *gin.Context) {
    userId, _ := cmn.GetUserId(c)
    username, _ := cmn.GetUsername(c)
    c.JSON(200, gin.H{
        "user_id":  userId,
        "username": username,
    })
}

func updateProfileHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "updated"})
}

func listPostsHandler(c *gin.Context) {
    c.JSON(200, gin.H{"posts": []string{"post1", "post2"}})
}
```

## API 參考

### MiddlewareRegistry

- `NewMiddlewareRegistry(router *gin.Engine) *MiddlewareRegistry` - 創建中間件注冊器
- `Register(middleware MiddlewareFunc) *MiddlewareRegistry` - 註冊單個中間件
- `RegisterMultiple(middlewares ...MiddlewareFunc) *MiddlewareRegistry` - 批量註冊中間件
- `ApplyDefault() *MiddlewareRegistry` - 應用默認中間件集合
- `Apply()` - 應用所有已註冊的中間件
- `GetMiddlewares() []MiddlewareFunc` - 獲取所有已註冊的中間件
- `Clear()` - 清空所有已註冊的中間件
- `Count() int` - 返回已註冊中間件的數量

### 輔助函數

- `GetUserId(c *gin.Context) (string, bool)` - 從上下文獲取用戶ID
- `GetUsername(c *gin.Context) (string, bool)` - 從上下文獲取用戶名
- `GetClaims(c *gin.Context) (*Claims, bool)` - 從上下文獲取完整的 claims

## 測試

運行測試：

```bash
cd cmn
go test -v
```

運行特定測試：

```bash
go test -v -run TestMiddlewareRegistry
go test -v -run TestCORSMiddleware
go test -v -run TestAuthMiddleware
```

## 注意事項

1. **中間件順序**: 中間件的執行順序很重要，建議順序：Recovery -> Logger -> CORS -> RateLimit -> Timeout -> Auth
2. **日誌依賴**: Logger 和 Recovery 中間件依賴於 `cmn.Logger`，需要先初始化日誌
3. **認證配置**: 使用認證中間件前需要配置 JWT 密鑰
4. **限流策略**: 根據實際業務需求調整限流參數
5. **超時設置**: 超時時間應該大於業務處理時間

## 許可證

MIT License

