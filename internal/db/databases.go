package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"time"
)

func OpenPostgres(dsn string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}

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
