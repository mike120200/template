package cmn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareRegistry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("創建中間件注冊器", func(t *testing.T) {
		router := gin.New()
		registry := NewMiddlewareRegistry(router)

		assert.NotNil(t, registry)
		assert.Equal(t, 0, registry.Count())
	})

	t.Run("註冊單個中間件", func(t *testing.T) {
		router := gin.New()
		registry := NewMiddlewareRegistry(router)

		registry.Register(func(c *gin.Context) {
			c.Next()
		})

		assert.Equal(t, 1, registry.Count())
	})

	t.Run("批量註冊中間件", func(t *testing.T) {
		router := gin.New()
		registry := NewMiddlewareRegistry(router)

		registry.RegisterMultiple(
			func(c *gin.Context) { c.Next() },
			func(c *gin.Context) { c.Next() },
			func(c *gin.Context) { c.Next() },
		)

		assert.Equal(t, 3, registry.Count())
	})

	t.Run("應用默認中間件", func(t *testing.T) {
		router := gin.New()
		registry := NewMiddlewareRegistry(router)

		registry.ApplyDefault()

		// 默認中間件包括: Recovery + Logger + CORS
		assert.Equal(t, 3, registry.Count())
	})

	t.Run("清空中間件", func(t *testing.T) {
		router := gin.New()
		registry := NewMiddlewareRegistry(router)

		registry.RegisterMultiple(
			func(c *gin.Context) { c.Next() },
			func(c *gin.Context) { c.Next() },
		)
		assert.Equal(t, 2, registry.Count())

		registry.Clear()
		assert.Equal(t, 0, registry.Count())
	})
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("處理CORS預檢請求", func(t *testing.T) {
		router := gin.New()
		router.Use(CORSMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("處理正常CORS請求", func(t *testing.T) {
		router := gin.New()
		router.Use(CORSMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("捕獲panic", func(t *testing.T) {
		// 初始化日誌（Recovery 中間件需要）
		LoggerInit()

		router := gin.New()
		router.Use(RecoveryMiddleware())

		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/panic", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("缺少token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("token格式錯誤", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())

		router.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidToken")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("有效token", func(t *testing.T) {
		// 生成測試token
		token, _ := GenerateToken("user123", "testuser")

		router := gin.New()
		router.Use(AuthMiddleware())

		router.GET("/protected", func(c *gin.Context) {
			userId, _ := GetUserId(c)
			username, _ := GetUsername(c)
			c.JSON(200, gin.H{
				"user_id":  userId,
				"username": username,
			})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("限流測試", func(t *testing.T) {
		router := gin.New()
		// 設置較小的限流參數便於測試
		router.Use(RateLimitByIP(2, 2)) // 每秒2個請求，桶容量2

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// 前兩個請求應該成功
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// 第三個請求應該被限流
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		// 等待一段時間後應該可以再次請求
		time.Sleep(600 * time.Millisecond)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("請求超時", func(t *testing.T) {
		router := gin.New()
		router.Use(TimeoutMiddleware(100 * time.Millisecond))

		router.GET("/slow", func(c *gin.Context) {
			time.Sleep(200 * time.Millisecond)
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/slow", nil)

		router.ServeHTTP(w, req)

		// 注意：由於gin的特性，超時後的狀態碼可能不是標準的408
		// 這裡我們只檢查請求是否被處理
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("請求正常完成", func(t *testing.T) {
		router := gin.New()
		router.Use(TimeoutMiddleware(200 * time.Millisecond))

		router.GET("/fast", func(c *gin.Context) {
			time.Sleep(50 * time.Millisecond)
			c.JSON(200, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/fast", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("無token - 允許訪問", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			userId, exists := GetUserId(c)
			c.JSON(200, gin.H{
				"has_user": exists,
				"user_id":  userId,
			})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("有效token - 提取用戶信息", func(t *testing.T) {
		token, _ := GenerateToken("user123", "testuser")

		router := gin.New()
		router.Use(OptionalAuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			userId, exists := GetUserId(c)
			c.JSON(200, gin.H{
				"has_user": exists,
				"user_id":  userId,
			})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
