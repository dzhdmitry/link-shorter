package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (c *RedisCache) Get(key string) (interface{}, bool, error) {
	ctx := context.Background()
	result, err := c.rdb.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return result, true, nil
}

func (c *RedisCache) Put(key string, value interface{}) error {
	ctx := context.Background()
	err := c.rdb.Set(ctx, key, value, 0).Err()

	return err
}
