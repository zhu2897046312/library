package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	
	"library/config"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.GlobalConfig.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			fmt.Printf("Error closing redis connection: %v\n", err)
		}
	}
}
