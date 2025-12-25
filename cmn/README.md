# Gin 中間件注冊模組

## 概述

這是一個功能完整的 Gin 中間件注冊和管理模組，提供了：

- **中間件注冊器**: 靈活的中間件註冊和管理系統
- **常用中間件**: 包括 CORS、日誌、恢復、認證、限流、超時等
- **易於使用**: 鏈式調用、配置靈活
- **生產就緒**: 包含完整的測試和文檔

## 文件結構

```
cmn/
├── middleware.go              # 中間件注冊器核心
├── middleware_cors.go         # CORS 中間件
├── middleware_logger.go       # 日誌中間件
├── middleware_recovery.go     # 恢復中間件（捕獲 panic）
├── middleware_auth.go         # JWT 認證中間件
├── middleware_rate_limit.go   # 限流中間件（令牌桶算法）
├── middleware_timeout.go      # 超時中間件
├── middleware_example.go      # 使用示例
├── middleware_test.go         # 測試文件
├── MIDDLEWARE_README.md       # 詳細使用文檔
├── INSTALL.md                 # 安裝說明
└── README.md                  # 本文件
```

## 快速開始

### 1. 安裝依賴

```bash
cd /Users/binbin/go/go_project/template
go mod tidy
```

### 2. 初始化配置

確保配置文件存在（如 `.conf_linux.json`）：

```json
{
  "log": {
    "level": "debug",
    "dir": "./cmn/log/"
  },
  "safe": {
    "jwtSecret": "your-secret-key-here"
  }
}
```

### 3. 使用示例

```go
package main

import (
    "my_template/cmn"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化配置和日誌
    cmn.ViperInit(".conf_linux.json")
    cmn.LoggerInit()
    
    // 創建 gin 引擎和中間件注冊器
    router := gin.New()
    registry := cmn.NewMiddlewareRegistry(router)
    
    // 應用默認中間件並註冊
    registry.ApplyDefault().Apply()
    
    // 設置路由
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    // 啟動服務器
    router.Run(":8080")
}
```

### 4. 運行完整示例

查看 `example_server.go` 了解完整的使用示例：

```bash
go run example_server.go
```

測試 API：

```bash
# 健康檢查
curl http://localhost:8080/health

# Ping
curl http://localhost:8080/ping

# 登錄
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 使用 token 訪問受保護的 API
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN"

# 測試 panic 恢復
curl http://localhost:8080/api/v1/test-panic

# 測試慢請求
curl http://localhost:8080/api/v1/test-slow
```

## 核心功能

### 1. 中間件注冊器

```go
registry := cmn.NewMiddlewareRegistry(router)

// 單個註冊
registry.Register(cmn.RecoveryMiddleware())

// 批量註冊
registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),
    cmn.LoggerMiddleware(),
    cmn.CORSMiddleware(),
)

// 應用默認中間件
registry.ApplyDefault()

// 應用所有已註冊的中間件
registry.Apply()
```

### 2. 可用中間件

| 中間件 | 功能 | 配置函數 |
|--------|------|----------|
| `RecoveryMiddleware()` | 捕獲 panic | `RecoveryMiddlewareWithConfig()` |
| `LoggerMiddleware()` | 記錄請求日誌 | `LoggerMiddlewareWithConfigCustom()` |
| `CORSMiddleware()` | 處理跨域請求 | `CORSMiddlewareWithConfig()` |
| `AuthMiddleware()` | JWT 認證 | `AuthMiddlewareWithConfig()` |
| `OptionalAuthMiddleware()` | 可選認證 | - |
| `RateLimitByIP()` | IP 限流 | `RateLimitMiddleware()` |
| `RateLimitByUser()` | 用戶限流 | - |
| `TimeoutMiddleware()` | 請求超時控制 | `TimeoutMiddlewareWithConfig()` |

### 3. 自定義配置

每個中間件都提供了配置函數：

```go
// CORS 配置
corsConfig := &cmn.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowCredentials: true,
}

// 日誌配置
loggerConfig := cmn.DefaultLoggerConfig()
loggerConfig.SkipPaths = []string{"/health", "/metrics"}
loggerConfig.SlowThreshold = 500 * time.Millisecond

// 認證配置
authConfig := cmn.DefaultAuthConfig()
authConfig.SkipPaths = []string{"/login", "/register"}

// 註冊
registry.RegisterMultiple(
    cmn.CORSMiddlewareWithConfig(corsConfig),
    cmn.LoggerMiddlewareWithConfigCustom(loggerConfig),
    cmn.AuthMiddlewareWithConfig(authConfig),
)
```

## 運行測試

```bash
cd cmn

# 運行所有測試
go test -v

# 運行特定測試
go test -v -run TestMiddlewareRegistry
go test -v -run TestCORSMiddleware
go test -v -run TestAuthMiddleware
go test -v -run TestRateLimitMiddleware
```

## 文檔

- [詳細使用文檔](./MIDDLEWARE_README.md) - 完整的 API 文檔和使用指南
- [安裝說明](./INSTALL.md) - 安裝和配置說明
- [示例代碼](./middleware_example.go) - 代碼示例
- [完整示例服務器](../example_server.go) - 可運行的完整示例

## 特性

✅ **靈活的註冊系統**: 支持單個註冊、批量註冊、鏈式調用  
✅ **豐富的中間件**: 涵蓋常見需求（CORS、日誌、認證、限流等）  
✅ **高度可配置**: 每個中間件都支持自定義配置  
✅ **生產就緒**: 包含完整的錯誤處理和日誌記錄  
✅ **完整測試**: 所有中間件都有單元測試  
✅ **詳細文檔**: 提供完整的使用文檔和示例  
✅ **性能優化**: 使用令牌桶算法進行高效限流  
✅ **易於擴展**: 可以輕鬆添加自定義中間件  

## 注意事項

1. **初始化順序**: 必須先初始化日誌，再使用需要日誌的中間件
2. **中間件順序**: 建議順序為 Recovery -> Logger -> CORS -> RateLimit -> Timeout -> Auth
3. **JWT 密鑰**: 生產環境必須設置自定義的 `jwtSecret`
4. **限流參數**: 根據實際業務需求調整限流參數
5. **超時設置**: 超時時間應該大於業務處理時間

## 許可證

MIT License

## 作者

Created for my_template project

