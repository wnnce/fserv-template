package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wnnce/fserv-template/config"
)

// redis.go
//
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/12 08:38

var (
	defaultClient *redis.Client
)

func RedisClient() *redis.Client {
	return defaultClient
}

func InitRedis(ctx context.Context) (func(), error) {
	host := config.ViperGet[string]("redis.host", "127.0.0.1")
	port := config.ViperGet[int]("redis.port", 6379)
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		DB:           config.ViperGet[int]("redis.index", 0),
		Username:     config.ViperGet[string]("redis.username"),
		Password:     config.ViperGet[string]("redis.password"),
		ReadTimeout:  config.ViperGet[time.Duration]("redis.timeout", 3*time.Second),
		WriteTimeout: config.ViperGet[time.Duration]("redis.timeout", 3*time.Second),
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("create redis defaultClient failed", slog.Group("data",
			slog.String("host", host),
			slog.Int("port", port),
		), slog.String("error", err.Error()))
		return nil, err
	}
	defaultClient = redisClient
	return func() {
		_ = defaultClient.Close()
	}, err
}
