package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"link-shorter.dzhdmitry.net/db"
	"testing"
)

type RedisSuite struct {
	suite.Suite
	rdb *redis.Client
}

func (s *RedisSuite) SetupSuite() {
	rdb, err := db.OpenTestRedis()

	if err != nil {
		panic(err)
	}

	s.rdb = rdb
}

func (s *RedisSuite) SetupTest() {
	ctx := context.Background()
	_ = s.rdb.FlushAll(ctx)
}

func (s *RedisSuite) TearDownSuite() {
	if s.rdb != nil {
		s.rdb.Close()
	}
}

func (s *RedisSuite) TestRedisGetNonExisting() {
	c := NewRedisCache(s.rdb)
	result, exists, err := c.Get("n-ex")

	s.NoError(err)
	s.False(exists)
	s.Nil(result)
}

func (s *RedisSuite) TestRedisGetExisting() {
	c := NewRedisCache(s.rdb)
	ctx := context.Background()
	_ = s.rdb.Set(ctx, "ex", "value_ex", 0)
	result, exists, err := c.Get("ex")

	s.NoError(err)
	s.True(exists)
	s.Equal("value_ex", result)
}

func (s *RedisSuite) TestRedisPut() {
	c := NewRedisCache(s.rdb)
	err := c.Put("test-put", "put-value")

	s.NoError(err)

	ctx := context.Background()
	result, err := c.rdb.Get(ctx, "test-put").Result()

	s.NoError(err)
	s.Equal("put-value", result)
}

func TestSQLStorage(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}
