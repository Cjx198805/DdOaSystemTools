package database

import (
	"context"
	"fmt"
	"github.com/ddoalistdownload/backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	RedisClient *redis.Client
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		// 连接池配置
		PoolSize:        100,  // 连接池最大连接数
		MinIdleConns:    10,   // 最小空闲连接数
		MaxConnAge:      1 * time.Hour, // 连接最大生命周期
		PoolTimeout:     30 * time.Second, // 获取连接超时时间
		IdleTimeout:     5 * time.Minute,  // 空闲连接超时时间
		IdleCheckFrequency: 1 * time.Minute, // 空闲连接检查频率
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logrus.Errorf("连接Redis失败: %v", err)
		return err
	}

	logrus.Infof("连接Redis成功，响应: %s", pong)
	RedisClient = client
	return nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
