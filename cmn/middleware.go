package cmn

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc = gin.HandlerFunc

// MiddlewareRegistry 中间件注册器
type MiddlewareRegistry struct {
	middlewares []MiddlewareFunc
	router      *gin.Engine
}

// NewMiddlewareRegistry 创建一个新的中间件注册器
func NewMiddlewareRegistry(router *gin.Engine) *MiddlewareRegistry {
	return &MiddlewareRegistry{
		middlewares: make([]MiddlewareFunc, 0),
		router:      router,
	}
}

// Register 注册单个中间件
func (r *MiddlewareRegistry) Register(middleware MiddlewareFunc) *MiddlewareRegistry {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// RegisterMultiple 批量注册中间件
func (r *MiddlewareRegistry) RegisterMultiple(middlewares ...MiddlewareFunc) *MiddlewareRegistry {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

// Apply 应用所有已注册的中间件到 gin 引擎
func (r *MiddlewareRegistry) Apply() {
	for _, middleware := range r.middlewares {
		r.router.Use(middleware)
	}
}

// ApplyDefault 应用默认中间件集合（Recovery + Logger + CORS）
func (r *MiddlewareRegistry) ApplyDefault() *MiddlewareRegistry {
	r.Register(RecoveryMiddleware())
	r.Register(LoggerMiddleware())
	r.Register(CORSMiddleware())
	return r
}

// GetMiddlewares 获取所有已注册的中间件
func (r *MiddlewareRegistry) GetMiddlewares() []MiddlewareFunc {
	return r.middlewares
}

// Clear 清空所有已注册的中间件
func (r *MiddlewareRegistry) Clear() {
	r.middlewares = make([]MiddlewareFunc, 0)
}

// Count 返回已注册中间件的数量
func (r *MiddlewareRegistry) Count() int {
	return len(r.middlewares)
}
