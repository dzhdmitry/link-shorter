package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dzhdmitry/link-shorter/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const StorageTypeFile = "file"
const StorageTypePostgres = "postgres"
const CacheTypeDisabled = "disabled"
const CacheTypeInMemory = "in-memory"
const CacheTypeRedis = "redis"

type Config struct {
	ProjectHost        string `env:"PROJECT_HOST" env-default:""`
	ProjectPort        int    `env:"PROJECT_PORT" env-default:"80"`
	ProjectStorageType string `env:"PROJECT_STORAGE_TYPE" env-default:"file"`
	FileAsync          bool   `env:"FILE_ASYNC" env-default:"false"`
	DbDSN              string `env:"DB_DSN" env-default:"postgres://go:pa55word@localhost:5432/short_links?sslmode=disable"`
	DbMaxOpenConns     int    `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
	DbMaxIdleConns     int    `env:"DB_MAX_IDLE_CONNS" env-default:"25"`
	DbMaxIdleTime      string `env:"DB_MAX_OPEN_TIME" env-default:"15m"`
	DbTimeout          int    `env:"DATABASE_TIMEOUT" env-default:"1"`
	CacheType          string `env:"CACHE_TYPE" env-default:"disabled"`
	CacheCapacity      int    `env:"CACHE_CAPACITY" env-default:"10"`
	CacheRedisDSN      string `env:"CACHE_REDIS_DSN" env-default:"redis://localhost:6379/0"`
	LimiterEnabled     bool   `env:"LIMITER_ENABLED" env-default:"true"`
	LimiterRPS         int    `env:"LIMITER_RPS" env-default:"2"`
	LimiterBurst       int    `env:"LIMITER_BURST" env-default:"4"`
}

func (c *Config) Info() string {
	inf := info{basePadding: 26}

	inf.addString(2, "Start server on", fmt.Sprintf("\"%s:%d\"", c.ProjectHost, c.ProjectPort))
	inf.addString(2, "Storage", c.ProjectStorageType)

	if c.ProjectStorageType == StorageTypeFile {
		inf.addBool(4, "Async", c.FileAsync)
	} else if c.ProjectStorageType == StorageTypePostgres {
		inf.addString(4, "DSN", c.DbDSN)
		inf.addInt(4, "Max open connections", c.DbMaxOpenConns)
		inf.addInt(4, "Max idle connections", c.DbMaxIdleConns)
		inf.addString(4, "Max idle time", c.DbMaxIdleTime)
		inf.addInt(4, "Timeout (seconds)", c.DbTimeout)
	}

	inf.addString(2, "Cache", c.CacheType)

	if c.CacheType == CacheTypeInMemory {
		inf.addInt(4, "Capacity of cache", c.CacheCapacity)
	} else if c.CacheType == CacheTypeRedis {
		inf.addString(4, "Redis DSN", c.CacheRedisDSN)
	}

	inf.addBool(2, "Rate limiter enabled", c.LimiterEnabled)

	if c.LimiterEnabled {
		inf.addInt(4, "RPS per IP", c.LimiterRPS)
		inf.addInt(4, "Maximum burst", c.LimiterBurst)
	}

	return "Using config:\n" + inf.getLines()
}

type LinksCollectionInterface interface {
	GenerateKey(URL string) (string, error)
	GenerateKeys(URLs []string) (map[string]string, error)
	GetURL(key string) (string, error)
	GetURLs(keys []string) (map[string]string, error)
}

type Application struct {
	Config     Config
	Logger     *utils.Logger
	Validator  Validator
	Links      LinksCollectionInterface
	Background *utils.Background
}

func (app *Application) Serve() error {
	server := &http.Server{
		Addr:         app.Config.ProjectHost + ":" + strconv.Itoa(app.Config.ProjectPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		receivedSignal := <-quit

		app.Logger.LogInfo("received signal " + receivedSignal.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		app.Logger.LogInfo("wait for background tasks...")
		app.Background.Wait()
		app.Logger.LogInfo("background tasks completed")

		shutdownError <- server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError

	if err != nil {
		return err
	}

	app.Logger.LogInfo("server stopped")

	return nil
}

func (app *Application) composeShortLink(key string) string {
	host := app.Config.ProjectHost

	if host == "" {
		host = "localhost"
	}

	return "http://" + host + "/go/" + key
}
