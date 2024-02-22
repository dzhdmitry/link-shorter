package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfo(t *testing.T) {
	config := Config{
		ProjectPort:        80,
		ProjectStorageType: StorageTypeFile,
		CacheType:          CacheTypeDisabled,
		LimiterEnabled:     true,
		LimiterRPS:         2,
		LimiterBurst:       4,
	}

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"    Async:                false\n"+
		"  Cache:                  disabled\n"+
		"  Rate limiter enabled:   true\n"+
		"    RPS per IP:           2\n"+
		"    Maximum burst:        4", config.Info())
}

func TestInfoPostgres(t *testing.T) {
	config := Config{
		ProjectPort:        80,
		ProjectStorageType: StorageTypePostgres,
		DbDSN:              "dsn",
		DbMaxOpenConns:     25,
		DbMaxIdleConns:     25,
		DbMaxIdleTime:      "15m",
		DbTimeout:          1,
		CacheType:          CacheTypeDisabled,
		LimiterEnabled:     false,
	}

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                postgres\n"+
		"    DSN:                  dsn\n"+
		"    Max open connections: 25\n"+
		"    Max idle connections: 25\n"+
		"    Max idle time:        15m\n"+
		"    Timeout (seconds):    1\n"+
		"  Cache:                  disabled\n"+
		"  Rate limiter enabled:   false", config.Info())
}

func TestInfoCacheInMemory(t *testing.T) {
	config := Config{
		ProjectPort:        80,
		ProjectStorageType: StorageTypeFile,
		CacheType:          CacheTypeInMemory,
		CacheCapacity:      10,
		LimiterEnabled:     false,
	}

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"    Async:                false\n"+
		"  Cache:                  in-memory\n"+
		"    Capacity of cache:    10\n"+
		"  Rate limiter enabled:   false", config.Info())
}

func TestInfoCacheRedis(t *testing.T) {
	config := Config{
		ProjectPort:        80,
		ProjectStorageType: StorageTypeFile,
		CacheType:          CacheTypeRedis,
		CacheRedisDSN:      "redis://redis:6379/0",
		LimiterEnabled:     false,
	}

	assert.Equal(t, "Using config:\n"+
		"  Start server on:        \":80\"\n"+
		"  Storage:                file\n"+
		"    Async:                false\n"+
		"  Cache:                  redis\n"+
		"    Redis DSN:            redis://redis:6379/0\n"+
		"  Rate limiter enabled:   false", config.Info())
}
