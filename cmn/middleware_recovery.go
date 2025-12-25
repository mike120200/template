package cmn

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware Recovery 中间件，用于捕获 panic
func RecoveryMiddleware() MiddlewareFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := string(debug.Stack())

				// 记录错误日志
				Logger().Error("发生panic",
					zap.Any("error", err),
					zap.String("stack", stack),
					zap.String("method", c.Request.Method),
					zap.String("uri", c.Request.RequestURI),
					zap.String("client_ip", c.ClientIP()),
				)

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "服务器内部错误",
					"error":   fmt.Sprintf("%v", err),
				})

				// 终止请求
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RecoveryConfig Recovery 中间件配置
type RecoveryConfig struct {
	EnableStackTrace bool   // 是否在响应中包含堆栈信息（仅用于调试）
	ErrorMessage     string // 自定义错误消息
}

// DefaultRecoveryConfig 默认 Recovery 配置
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		EnableStackTrace: false,
		ErrorMessage:     "服务器内部错误",
	}
}

// RecoveryMiddlewareWithConfig 带配置的 Recovery 中间件
func RecoveryMiddlewareWithConfig(config *RecoveryConfig) MiddlewareFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := string(debug.Stack())

				// 记录错误日志
				Logger().Error("发生panic",
					zap.Any("error", err),
					zap.String("stack", stack),
					zap.String("method", c.Request.Method),
					zap.String("uri", c.Request.RequestURI),
					zap.String("client_ip", c.ClientIP()),
				)

				// 构建响应
				response := gin.H{
					"code":    http.StatusInternalServerError,
					"message": config.ErrorMessage,
				}

				// 如果启用堆栈跟踪，添加到响应中（仅用于调试环境）
				if config.EnableStackTrace {
					response["error"] = fmt.Sprintf("%v", err)
					response["stack"] = stack
				}

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, response)

				// 终止请求
				c.Abort()
			}
		}()

		c.Next()
	}
}
