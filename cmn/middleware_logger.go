package cmn

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() MiddlewareFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志记录
		Logger().Info("HTTP请求",
			zap.Int("status_code", statusCode),
			zap.String("latency", latencyTime.String()),
			zap.String("client_ip", clientIP),
			zap.String("method", reqMethod),
			zap.String("uri", reqUri),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 如果有错误，记录错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				Logger().Error("请求错误",
					zap.String("error", e.Error()),
					zap.String("method", reqMethod),
					zap.String("uri", reqUri),
				)
			}
		}
	}
}

// LoggerConfig 带配置的日志中间件
type LoggerConfig struct {
	SkipPaths     []string      // 跳过记录的路径
	SkipMethods   []string      // 跳过记录的方法
	SlowThreshold time.Duration // 慢请求阈值
}

// DefaultLoggerConfig 默认日志配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:     []string{"/health", "/ping"},
		SkipMethods:   []string{},
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LoggerMiddlewareWithConfigCustom 自定义配置的日志中间件
func LoggerMiddlewareWithConfigCustom(config *LoggerConfig) MiddlewareFunc {
	return func(c *gin.Context) {
		// 检查是否需要跳过此路径
		path := c.Request.URL.Path
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// 检查是否需要跳过此方法
		method := c.Request.Method
		for _, skipMethod := range config.SkipMethods {
			if method == skipMethod {
				c.Next()
				return
			}
		}

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 根据执行时间选择日志级别
		if latencyTime > config.SlowThreshold {
			Logger().Warn("慢请求",
				zap.Int("status_code", statusCode),
				zap.Duration("latency", latencyTime),
				zap.String("client_ip", clientIP),
				zap.String("method", reqMethod),
				zap.String("uri", reqUri),
			)
		} else {
			Logger().Info("HTTP请求",
				zap.Int("status_code", statusCode),
				zap.Duration("latency", latencyTime),
				zap.String("client_ip", clientIP),
				zap.String("method", reqMethod),
				zap.String("uri", reqUri),
			)
		}

		// 如果有错误，记录错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				Logger().Error("请求错误",
					zap.String("error", e.Error()),
					zap.String("method", reqMethod),
					zap.String("uri", reqUri),
				)
			}
		}
	}
}
