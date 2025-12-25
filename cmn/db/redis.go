package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 客户端连接
func InitRedis() error {
	// 从配置中读取 Redis 连接信息
	host := viper.GetString("Redis.host")
	port := viper.GetString("Redis.port")
	password := viper.GetString("Redis.password")
	db := viper.GetInt("Redis.DB")
	maxIdleConns := viper.GetInt("Redis.MaxIdleConns")
	maxActiveConns := viper.GetInt("Redis.MaxActiveConns")

	// 创建 Redis 客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           db,
		MaxIdleConns: maxIdleConns,
		PoolSize:     maxActiveConns,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接测试失败: %w", err)
	}

	zap.L().Info("Redis 客户端初始化成功",
		zap.String("host", host),
		zap.Int("db", db))

	return nil
}

func GetRedis() *redis.Client {
	return RedisClient
}
