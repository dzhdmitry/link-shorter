package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfo(t *testing.T) {
	conf := NewConfig()
	confInfo := conf.Info()

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"  Async:                  false\n"+
		"  Cache:                  disabled\n"+
		"  Rate limiter enabled:   true\n"+
		"    RPS per IP:           2\n"+
		"    Maximum burst:        4", confInfo)
}

func TestInfoPostgres(t *testing.T) {
	conf := NewConfig()
	conf.ProjectStorageType = StorageTypePostgres
	conf.LimiterEnabled = false
	confInfo := conf.Info()

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                postgres\n"+
		"    DSN:                  postgres://go:pa55word@postgres:5432/short_links?sslmode=disable\n"+
		"    Max open connections: 25\n"+
		"    Max idle connections: 25\n"+
		"    Max idle time:        15m\n"+
		"    Timeout (seconds):    1\n"+
		"  Cache:                  disabled\n"+
		"  Rate limiter enabled:   false", confInfo)
}

func TestInfoCacheInMemory(t *testing.T) {
	conf := NewConfig()
	conf.CacheType = CacheTypeInMemory
	conf.LimiterEnabled = false
	confInfo := conf.Info()

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"  Async:                  false\n"+
		"  Cache:                  in-memory\n"+
		"    Capacity of cache:    0\n"+
		"  Rate limiter enabled:   false", confInfo)
}

func TestInfoCacheRedis(t *testing.T) {
	conf := NewConfig()
	conf.CacheType = CacheTypeRedis
	conf.LimiterEnabled = false
	confInfo := conf.Info()

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"  Async:                  false\n"+
		"  Cache:                  redis\n"+
		"    Redis DSN:            redis://redis:6379/0\n"+
		"  Rate limiter enabled:   false", confInfo)
}
