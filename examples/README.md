# 中間件使用示例

## middleware_server.go

這是一個完整的 Web 服務器示例，展示了如何使用中間件注冊模組。

### 運行示例

1. 確保已經安裝依賴：

```bash
cd /Users/binbin/go/go_project/template
go mod tidy
```

2. 運行示例服務器：

```bash
go run examples/middleware_server.go
```

3. 測試 API：

```bash
# 健康檢查
curl http://localhost:8080/health

# Ping
curl http://localhost:8080/ping

# 登錄
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 註冊
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","email":"test@example.com"}'

# 使用 token 訪問受保護的 API
TOKEN="YOUR_TOKEN_HERE"
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"

# 測試 panic 恢復
curl http://localhost:8080/api/v1/test-panic

# 測試慢請求日誌
curl http://localhost:8080/api/v1/test-slow

# 測試可選認證
curl http://localhost:8080/api/v1/posts/public
curl http://localhost:8080/api/v1/posts/public \
  -H "Authorization: Bearer $TOKEN"
```

## 功能展示

### 1. 中間件註冊

示例展示了如何使用中間件注冊器：

```go
registry := cmn.NewMiddlewareRegistry(router)

registry.RegisterMultiple(
    cmn.RecoveryMiddleware(),
    cmn.LoggerMiddleware(),
    cmn.CORSMiddleware(),
    cmn.RateLimitByIP(100, 200),
    cmn.TimeoutMiddleware(30*time.Second),
)

registry.Apply()
```

### 2. JWT 認證

展示了如何使用認證中間件保護 API：

- 登錄獲取 token
- 使用 token 訪問受保護的路由
- 從上下文獲取用戶信息

### 3. 可選認證

展示了如何實現可選認證（訪客和登錄用戶都可以訪問）：

```go
optional := v1.Group("")
optional.Use(cmn.OptionalAuthMiddleware())
{
    optional.GET("/posts/public", publicPostsHandler)
}
```

### 4. 錯誤恢復

訪問 `/api/v1/test-panic` 測試 Recovery 中間件如何捕獲 panic。

### 5. 慢請求監控

訪問 `/api/v1/test-slow` 測試日誌中間件如何記錄慢請求。

### 6. 限流

快速連續發送多個請求測試限流中間件。

## API 端點

| 方法 | 路徑 | 認證 | 說明 |
|------|------|------|------|
| GET | /health | ❌ | 健康檢查 |
| GET | /ping | ❌ | Ping 測試 |
| POST | /api/v1/login | ❌ | 用戶登錄 |
| POST | /api/v1/register | ❌ | 用戶註冊 |
| GET | /api/v1/profile | ✅ | 獲取用戶信息 |
| PUT | /api/v1/profile | ✅ | 更新用戶信息 |
| GET | /api/v1/posts | ✅ | 獲取文章列表 |
| POST | /api/v1/posts | ✅ | 創建文章 |
| GET | /api/v1/posts/public | 🔶 | 獲取公開文章（可選認證） |
| GET | /api/v1/test-panic | ❌ | 測試 panic 恢復 |
| GET | /api/v1/test-slow | ❌ | 測試慢請求 |

圖例：
- ❌ 不需要認證
- ✅ 需要認證
- 🔶 可選認證

## 注意事項

1. 示例使用簡化的用戶驗證邏輯，生產環境需要實現真實的用戶驗證
2. 確保配置文件中設置了 `jwtSecret`
3. 根據實際需求調整限流參數
4. 生產環境建議關閉 Recovery 中間件的堆棧跟踪

