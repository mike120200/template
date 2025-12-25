package cmn

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认 CORS 配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 小时
	}
}

// CORSMiddleware 创建 CORS 中间件（使用默认配置）
func CORSMiddleware() MiddlewareFunc {
	return CORSMiddlewareWithConfig(DefaultCORSConfig())
}

// CORSMiddlewareWithConfig 创建 CORS 中间件（使用自定义配置）
func CORSMiddlewareWithConfig(config *CORSConfig) MiddlewareFunc {
	return func(c *gin.Context) {
		// 设置允许的来源
		if len(config.AllowOrigins) > 0 {
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				// 检查来源是否在允许列表中
				allowed := false
				for _, allowOrigin := range config.AllowOrigins {
					if allowOrigin == "*" || allowOrigin == origin {
						allowed = true
						break
					}
				}
				if allowed {
					if config.AllowOrigins[0] == "*" {
						c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
					} else {
						c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					}
				}
			}
		}

		// 设置允许的方法
		if len(config.AllowMethods) > 0 {
			methods := ""
			for i, method := range config.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		}

		// 设置允许的请求头
		if len(config.AllowHeaders) > 0 {
			headers := ""
			for i, header := range config.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		}

		// 设置暴露的响应头
		if len(config.ExposeHeaders) > 0 {
			headers := ""
			for i, header := range config.ExposeHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Writer.Header().Set("Access-Control-Expose-Headers", headers)
		}

		// 设置是否允许携带凭证
		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 设置预检请求的缓存时间
		if config.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
