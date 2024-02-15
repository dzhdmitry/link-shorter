package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func OpenRedis(DSN string) (*redis.Client, error) {
	opts, err := redis.ParseURL(DSN)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rdb := redis.NewClient(opts)
	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func OpenTestRedis() (*redis.Client, error) {
	defaultDsn := os.Getenv("TEST_REDIS_DSN")

	if defaultDsn == "" {
		defaultDsn = "redis://localhost:6379/1"
	}

	return OpenRedis(defaultDsn)
}
