package redis

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

// helper.go
//
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/12 11:56

func RedisGetStruct[T any](ctx context.Context, key string, client *redis.Client) (T, error) {
	var empty T
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return empty, err
	}
	err = sonic.UnmarshalString(result, &empty)
	return empty, err
}

func RedisGetAddrStruct[T any](ctx context.Context, key string, client *redis.Client) (*T, error) {
	var empty T
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	err = sonic.UnmarshalString(result, &empty)
	return &empty, err
}

func RedisGetSlice[T any](ctx context.Context, key string, client *redis.Client) ([]T, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	value := make([]T, 0)
	err = sonic.UnmarshalString(result, &value)
	return value, err
}
