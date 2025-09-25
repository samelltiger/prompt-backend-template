// pkg/cache/redis.go
package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value string, expireSeconds int) error {
	return c.client.Set(
		context.Background(),
		key,
		value,
		time.Duration(expireSeconds)*time.Second,
	).Err()
}

// Delete 删除缓存
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(key string) (bool, error) {
	result, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// TTL 获取缓存剩余时间
func (c *RedisCache) TTL(key string) (time.Duration, error) {
	return c.client.TTL(context.Background(), key).Result()
}
