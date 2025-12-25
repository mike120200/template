# Gin 中間件注冊模組 - 完整指南

## 📦 已創建的文件

### 核心文件

```
cmn/
├── middleware.go              # 中間件注冊器（核心）
├── middleware_cors.go         # CORS 中間件
├── middleware_logger.go       # 日誌中間件  
├── middleware_recovery.go     # 恢復中間件（panic 捕獲）
├── middleware_auth.go         # JWT 認證中間件
├── middleware_rate_limit.go   # 限流中間件（令牌桶）
├── middleware_timeout.go      # 超時中間件
├── middleware_example.go      # 代碼示例
├── middleware_test.go         # 單元測試
├── token.go                   # JWT token 處理（已擴展）
└── 文檔/
    ├── README.md              # 概述文檔
    ├── MIDDLEWARE_README.md   # 詳細使用文檔
    └── INSTALL.md             # 安裝說明

examples/
├── middleware_server.go       # 完整的服務器示例
└── README.md                  # 示例說明

go.mod                         # 已添加 gin 依賴
```

## 🚀 快速開始

### 第 1 步：安裝依賴

```bash
cd /Users/binbin/go/go_project/template
go mod tidy
```

這會安裝 `github.com/gin-gonic/gin` 和其他必要的依賴。

### 第 2 步：準備配置文件

確保配置文件存在（如 `.conf_linux.json`）：

```json
{
  "log": {
    "level": "debug",
    "dir": "./cmn/log/"
  },
  "safe": {
    "jwtSecret": "your-secret-key-here-change-in-production"
  }
}
```

### 第 3 步：使用中間件

**最簡單的方式**（使用默認中間件）：

```go
package main

import (
    "my_template/cmn"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化
    cmn.ViperInit(".conf_linux.json")
    cmn.LoggerInit()
    
    // 創建路由器和注冊器
    router := gin.New()
    registry := cmn.NewMiddlewareRegistry(router)
    
    // 應用默認中間件（Recovery + Logger + CORS）
    registry.ApplyDefault().Apply()
    
    // 設置路由
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    // 啟動
    router.Run(":8080")
}
```

**自定義中間件配置**：

```go
// 創建注冊器
registry := cmn.NewMiddlewareRegistry(router)

// 自定義 CORS
corsConfig := &cmn.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowCredentials: true,
}

// 註冊中間件
registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),                    // 捕獲 panic
    cmn.LoggerMiddleware(),                      // 日誌記錄
    cmn.CORSMiddlewareWithConfig(corsConfig),    // 跨域配置
    cmn.RateLimitByIP(100, 200),                 // 限流：每秒100請求
    cmn.TimeoutMiddleware(30*time.Second),       // 30秒超時
)

// 應用中間件
registry.Apply()
```

### 第 4 步：添加認證保護

```go
// API 路由組
api := router.Group("/api/v1")
{
    // 公開路由
    api.POST("/login", loginHandler)
    api.POST("/register", registerHandler)
    
    // 需要認證的路由
    protected := api.Group("")
    protected.Use(cmn.AuthMiddleware())
    {
        protected.GET("/profile", profileHandler)
        protected.PUT("/profile", updateProfileHandler)
    }
}

// 在處理器中獲取用戶信息
func profileHandler(c *gin.Context) {
    userId, _ := cmn.GetUserId(c)
    username, _ := cmn.GetUsername(c)
    
    c.JSON(200, gin.H{
        "user_id": userId,
        "username": username,
    })
}
```

## 📚 完整示例

運行完整的示例服務器：

```bash
go run examples/middleware_server.go
```

測試 API：

```bash
# 健康檢查
curl http://localhost:8080/health

# 登錄獲取 token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 使用 token 訪問受保護的 API
TOKEN="your-token-here"
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"

# 測試 panic 恢復
curl http://localhost:8080/api/v1/test-panic

# 測試慢請求
curl http://localhost:8080/api/v1/test-slow
```

## 🧪 運行測試

```bash
cd cmn

# 運行所有測試
go test -v

# 運行特定測試
go test -v -run TestMiddlewareRegistry
go test -v -run TestCORSMiddleware
go test -v -run TestAuthMiddleware
go test -v -run TestRateLimitMiddleware
go test -v -run TestTimeoutMiddleware
```

## 📖 可用的中間件

| 中間件 | 功能 | 使用場景 |
|--------|------|----------|
| **RecoveryMiddleware** | 捕獲 panic，防止服務崩潰 | 必需，應該最先註冊 |
| **LoggerMiddleware** | 記錄所有 HTTP 請求 | 必需，便於調試和監控 |
| **CORSMiddleware** | 處理跨域請求 | 前後端分離項目 |
| **AuthMiddleware** | JWT 認證 | 保護需要登錄的 API |
| **OptionalAuthMiddleware** | 可選認證 | 訪客和用戶都能訪問的 API |
| **RateLimitByIP** | 按 IP 限流 | 防止濫用 |
| **RateLimitByUser** | 按用戶限流 | 精細化限流控制 |
| **TimeoutMiddleware** | 請求超時控制 | 防止長時間請求 |

## 🔧 配置選項

### CORS 配置

```go
corsConfig := &cmn.CORSConfig{
    AllowOrigins:     []string{"*"},                    // 允許的來源
    AllowMethods:     []string{"GET", "POST", "PUT"},   // 允許的方法
    AllowHeaders:     []string{"Content-Type", "Authorization"}, // 允許的頭
    AllowCredentials: true,                             // 是否允許憑證
    MaxAge:           86400,                            // 預檢緩存時間（秒）
}
```

### 日誌配置

```go
loggerConfig := cmn.DefaultLoggerConfig()
loggerConfig.SkipPaths = []string{"/health", "/metrics"}  // 跳過的路徑
loggerConfig.SlowThreshold = 200 * time.Millisecond       // 慢請求閾值
```

### 認證配置

```go
authConfig := cmn.DefaultAuthConfig()
authConfig.SkipPaths = []string{"/login", "/register"}    // 跳過認證的路徑
authConfig.TokenHeader = "Authorization"                  // Token 所在的頭
authConfig.TokenPrefix = "Bearer"                         // Token 前綴
```

### 限流配置

```go
rateLimitConfig := &cmn.RateLimitConfig{
    Rate:     100,                      // 每秒請求數
    Capacity: 200,                      // 桶容量
    KeyFunc: func(c *gin.Context) string {
        return c.ClientIP()             // 限流鍵（IP/用戶ID等）
    },
}
```

### 超時配置

```go
timeoutConfig := &cmn.TimeoutConfig{
    Timeout: 30 * time.Second,          // 超時時間
    ErrorHandler: func(c *gin.Context) {
        c.JSON(408, gin.H{"error": "請求超時"})
    },
}
```

## 💡 最佳實踐

### 1. 中間件順序

推薦的中間件註冊順序：

```go
registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),      // 1. 最先：捕獲 panic
    cmn.LoggerMiddleware(),        // 2. 日誌記錄
    cmn.CORSMiddleware(),          // 3. CORS 處理
    cmn.RateLimitByIP(100, 200),   // 4. 限流
    cmn.TimeoutMiddleware(30*s),   // 5. 超時控制
    // Auth 中間件通常在路由組級別使用
)
```

### 2. 認證中間件使用

全局認證（不推薦）：
```go
router.Use(cmn.AuthMiddleware())  // 所有路由都需要認證
```

路由組認證（推薦）：
```go
public := router.Group("/public")
{
    public.GET("/posts", listPublicPosts)
}

protected := router.Group("/api")
protected.Use(cmn.AuthMiddleware())
{
    protected.GET("/profile", getProfile)
}
```

### 3. 生產環境配置

```go
// Recovery 中間件 - 不顯示堆棧跟踪
recoveryConfig := &cmn.RecoveryConfig{
    EnableStackTrace: false,
    ErrorMessage:     "服務器內部錯誤",
}

// 日誌中間件 - 跳過健康檢查
loggerConfig := cmn.DefaultLoggerConfig()
loggerConfig.SkipPaths = []string{"/health", "/metrics"}

// 限流 - 根據業務調整
registry.Register(cmn.RateLimitByIP(1000, 2000))
```

### 4. JWT 密鑰管理

開發環境：
```json
{
  "safe": {
    "jwtSecret": "dev-secret-key"
  }
}
```

生產環境：
```bash
# 使用環境變量
export JWT_SECRET="your-strong-production-secret"
```

## 📝 文檔資源

- **[cmn/README.md](cmn/README.md)** - 概述和快速入門
- **[cmn/MIDDLEWARE_README.md](cmn/MIDDLEWARE_README.md)** - 詳細的 API 文檔
- **[cmn/INSTALL.md](cmn/INSTALL.md)** - 安裝說明
- **[examples/README.md](examples/README.md)** - 示例說明
- **[cmn/middleware_example.go](cmn/middleware_example.go)** - 代碼示例
- **[examples/middleware_server.go](examples/middleware_server.go)** - 完整服務器

## 🐛 故障排除

### 問題 1：找不到 gin 包

```bash
# 解決方案
go mod tidy
```

### 問題 2：Logger 相關錯誤

```go
// 確保先初始化日誌
cmn.LoggerInit()
```

### 問題 3：JWT 認證失敗

```json
// 確保配置文件中有 jwtSecret
{
  "safe": {
    "jwtSecret": "your-secret-key"
  }
}
```

### 問題 4：CORS 不生效

```go
// 確保 CORS 中間件在路由之前註冊
registry.Register(cmn.CORSMiddleware())
registry.Apply()  // 必須調用 Apply()
```

## ✅ 功能清單

- ✅ 中間件注冊器（支持鏈式調用）
- ✅ CORS 中間件（完全可配置）
- ✅ 日誌中間件（支持慢請求監控）
- ✅ 恢復中間件（panic 捕獲）
- ✅ JWT 認證中間件（支持可選認證）
- ✅ 限流中間件（令牌桶算法）
- ✅ 超時中間件（請求超時控制）
- ✅ 完整的單元測試
- ✅ 詳細的文檔
- ✅ 可運行的示例

## 📄 許可證

MIT License

## 🎉 開始使用

1. 運行 `go mod tidy` 安裝依賴
2. 查看 `examples/middleware_server.go` 了解完整示例
3. 閱讀 `cmn/MIDDLEWARE_README.md` 了解詳細 API
4. 在你的項目中使用中間件！

祝你使用愉快！🚀

