package cmn

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Timeout      time.Duration      // 超时时间
	ErrorHandler func(*gin.Context) // 超时时的错误处理
}

// DefaultTimeoutConfig 默认超时配置
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Timeout: 30 * time.Second,
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    http.StatusRequestTimeout,
				"message": "请求超时",
			})
		},
	}
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) MiddlewareFunc {
	return TimeoutMiddlewareWithConfig(&TimeoutConfig{
		Timeout: timeout,
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    http.StatusRequestTimeout,
				"message": "请求超时",
			})
		},
	})
}

// TimeoutMiddlewareWithConfig 带配置的超时中间件
func TimeoutMiddlewareWithConfig(config *TimeoutConfig) MiddlewareFunc {
	return func(c *gin.Context) {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// 替换请求的上下文
		c.Request = c.Request.WithContext(ctx)

		// 使用通道来等待请求完成或超时
		finished := make(chan struct{})

		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-ctx.Done():
			// 请求超时
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				config.ErrorHandler(c)
				c.Abort()
			}
		case <-finished:
			// 请求正常完成
		}
	}
}
