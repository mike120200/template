package cmn

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	rate     float64            // 令牌生成速率（每秒）
	capacity int                // 桶容量
	buckets  map[string]*bucket // 每个客户端的令牌桶
	mu       sync.RWMutex
}

type bucket struct {
	tokens     float64
	lastUpdate time.Time
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rate float64, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:     rate,
		capacity: capacity,
		buckets:  make(map[string]*bucket),
	}
}

// Allow 检查是否允许请求
func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, exists := l.buckets[key]

	if !exists {
		// 新客户端，创建令牌桶
		l.buckets[key] = &bucket{
			tokens:     float64(l.capacity) - 1,
			lastUpdate: now,
		}
		return true
	}

	// 计算应该添加的令牌数
	elapsed := now.Sub(b.lastUpdate).Seconds()
	tokensToAdd := elapsed * l.rate

	// 更新令牌数量
	b.tokens = min(b.tokens+tokensToAdd, float64(l.capacity))
	b.lastUpdate = now

	// 检查是否有足够的令牌
	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Rate         float64                   // 每秒允许的请求数
	Capacity     int                       // 桶容量
	KeyFunc      func(*gin.Context) string // 生成限流键的函数（默认使用 ClientIP）
	ErrorHandler func(*gin.Context)        // 限流时的错误处理
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Rate:     10, // 每秒10个请求
		Capacity: 20, // 桶容量20
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
		},
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config *RateLimitConfig) MiddlewareFunc {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	limiter := NewTokenBucketLimiter(config.Rate, config.Capacity)

	return func(c *gin.Context) {
		key := config.KeyFunc(c)

		if !limiter.Allow(key) {
			config.ErrorHandler(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByIP IP限流中间件（简化版）
func RateLimitByIP(rate float64, capacity int) MiddlewareFunc {
	config := &RateLimitConfig{
		Rate:     rate,
		Capacity: capacity,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
		},
	}
	return RateLimitMiddleware(config)
}

// RateLimitByUser 用户限流中间件（基于用户ID）
func RateLimitByUser(rate float64, capacity int) MiddlewareFunc {
	config := &RateLimitConfig{
		Rate:     rate,
		Capacity: capacity,
		KeyFunc: func(c *gin.Context) string {
			userId, exists := GetUserId(c)
			if !exists {
				// 如果没有用户ID，使用IP
				return c.ClientIP()
			}
			return userId
		},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
		},
	}
	return RateLimitMiddleware(config)
}
