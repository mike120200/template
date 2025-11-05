package cmn

import (
	"github.com/spf13/viper"
	"runtime"
	"sync"

	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// 全局 ants goroutine 池
var (
	pool     *ants.Pool
	initOnce sync.Once
	initErr  error
)

// InitGoroutinePool 使用 sync.Once 初始化全局 goroutine 池；重复调用将直接返回相同结果。
// size <= 0 时默认使用 runtime.NumCPU()*4；preAlloc 控制是否预分配。
func InitGoroutinePool(preAlloc bool) error {
	initOnce.Do(func() {
		size := viper.GetInt("pool.size")
		if size <= 0 {
			size = runtime.NumCPU() * 4
		}
		p, err := ants.NewPool(size,
			ants.WithPreAlloc(preAlloc),
			ants.WithPanicHandler(func(panicVal interface{}) {
				zap.L().Error("panic in ants worker", zap.Any("panic", panicVal))
			}),
		)
		if err != nil {
			initErr = err
			return
		}
		pool = p
		Logger().Info("ants pool initialized", zap.Int("size", size), zap.Bool("preAlloc", preAlloc))
	})
	return initErr
}

// GetGoroutinePool 获取全局池（未初始化或初始化失败时返回 nil）。
func GetGoroutinePool() *ants.Pool {
	if pool == nil {
		panic("ants pool not initialized")
	}
	return pool
}
