package cmn

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() MiddlewareFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "缺少认证令牌",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 验证 token
		token := parts[1]
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "无效的认证令牌",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// AuthConfig 认证中间件配置
type AuthConfig struct {
	TokenHeader  string                    // Token 所在的请求头字段，默认 "Authorization"
	TokenPrefix  string                    // Token 前缀，默认 "Bearer"
	SkipPaths    []string                  // 跳过认证的路径
	ErrorHandler func(*gin.Context, error) // 自定义错误处理
}

// DefaultAuthConfig 默认认证配置
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		TokenHeader: "Authorization",
		TokenPrefix: "Bearer",
		SkipPaths:   []string{"/login", "/register", "/health", "/ping"},
		ErrorHandler: func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "认证失败",
				"error":   err.Error(),
			})
		},
	}
}

// AuthMiddlewareWithConfig 带配置的认证中间件
func AuthMiddlewareWithConfig(config *AuthConfig) MiddlewareFunc {
	return func(c *gin.Context) {
		// 检查是否需要跳过此路径
		path := c.Request.URL.Path
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// 从请求头获取 token
		authHeader := c.GetHeader(config.TokenHeader)
		if authHeader == "" {
			config.ErrorHandler(c, ErrInvalidToken)
			c.Abort()
			return
		}

		// 检查前缀
		var token string
		if config.TokenPrefix != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == config.TokenPrefix) {
				config.ErrorHandler(c, ErrInvalidToken)
				c.Abort()
				return
			}
			token = parts[1]
		} else {
			token = authHeader
		}

		// 验证 token
		claims, err := ParseToken(token)
		if err != nil {
			config.ErrorHandler(c, err)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证，但如果提供了 token 会验证）
func OptionalAuthMiddleware() MiddlewareFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有 token，继续处理
			c.Next()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			// 验证 token
			token := parts[1]
			claims, err := ParseToken(token)
			if err == nil {
				// token 有效，将用户信息存入上下文
				c.Set("user_id", claims.UserId)
				c.Set("username", claims.Username)
				c.Set("claims", claims)
			}
		}

		c.Next()
	}
}

// GetUserId 从上下文中获取用户ID
func GetUserId(c *gin.Context) (string, bool) {
	userId, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	userIdStr, ok := userId.(string)
	return userIdStr, ok
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	usernameStr, ok := username.(string)
	return usernameStr, ok
}

// GetClaims 从上下文中获取完整的 claims
func GetClaims(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	claimsObj, ok := claims.(*Claims)
	return claimsObj, ok
}
