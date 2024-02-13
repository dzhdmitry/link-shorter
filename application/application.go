package application

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	ProjectHost         string `env:"PROJECT_HOST"`
	ProjectPort         int    `env:"PROJECT_PORT"`
	ProjectKeyMaxLength int    `env:"PROJECT_KEY_MAX_LENGTH"`
	ProjectStorageType  string `env:"PROJECT_STORAGE_TYPE"`
	FileAsync           bool   `env:"FILE_ASYNC"`
	DbDSN               string `env:"DB_DSN"`
	DbMaxOpenConns      int    `env:"DB_MAX_OPEN_CONNS"`
	DbMaxIdleConns      int    `env:"DB_MAX_IDLE_CONNS"`
	DbMaxIdleTime       string `env:"DB_MAX_OPEN_TIME"`
	DbTimeout           int    `env:"DATABASE_TIMEOUT"`
	CacheEnabled        bool   `env:"CACHE_ENABLED"`
	CacheCapacity       int    `env:"CACHE_CAPACITY"`
	LimiterEnabled      bool   `env:"LIMITER_ENABLED"`
	LimiterRPS          int    `env:"LIMITER_RPS"`
	LimiterBurst        int    `env:"LIMITER_BURST"`
}

func NewConfig() Config {
	return Config{
		ProjectHost:         "",
		ProjectPort:         80,
		ProjectKeyMaxLength: 12,
		ProjectStorageType:  "file",
		FileAsync:           false,
		DbDSN:               "postgres://go:pa55word@postgres:5432/short_links?sslmode=disable",
		DbMaxOpenConns:      25,
		DbMaxIdleConns:      25,
		DbMaxIdleTime:       "15m",
		DbTimeout:           1,
		CacheEnabled:        false,
		CacheCapacity:       0,
		LimiterEnabled:      true,
		LimiterRPS:          2,
		LimiterBurst:        4,
	}
}

func (c *Config) Parse() {
	flag.StringVar(&c.ProjectHost, "host", c.ProjectHost, "Project server host")
	flag.IntVar(&c.ProjectPort, "port", c.ProjectPort, "Project server port")
	flag.IntVar(&c.ProjectKeyMaxLength, "key-max", c.ProjectKeyMaxLength, "Max length of the key")
	flag.StringVar(&c.ProjectStorageType, "storage", c.ProjectStorageType, "Storage type (file|postgres)")
	flag.BoolVar(&c.FileAsync, "file-async", c.FileAsync, "File storage is asynchronous|synchronous (true|false)")
	flag.StringVar(&c.DbDSN, "db-dsn", c.DbDSN, "PostgreSQL DSN")
	flag.IntVar(&c.DbMaxOpenConns, "db-max-open-conns", c.DbMaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&c.DbMaxIdleConns, "db-max-idle-conns", c.DbMaxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&c.DbMaxIdleTime, "db-max-idle-time", c.DbMaxIdleTime, "PostgreSQL max connection idle time")
	flag.IntVar(&c.DbTimeout, "db-timeout", c.DbTimeout, "PostgreSQL queries execution timeout")
	flag.BoolVar(&c.CacheEnabled, "cache", c.CacheEnabled, "Caching of short links enabled")
	flag.IntVar(&c.CacheCapacity, "cache-cap", c.CacheCapacity, "Capacity of cache")
	flag.BoolVar(&c.LimiterEnabled, "limiter", c.LimiterEnabled, "Rate limiter is enabled")
	flag.IntVar(&c.LimiterRPS, "limiter-rps", c.LimiterRPS, "Rate limiter maximum RPS per IP")
	flag.IntVar(&c.LimiterBurst, "limiter-burst", c.LimiterBurst, "Rate limiter maximum burst")
	flag.Parse()
}

func (c *Config) Info() string {
	lines := []string{
		"Using config:",
		"   Project:",
		fmt.Sprintf("      Host:                 %s", c.ProjectHost),
		fmt.Sprintf("      Port:                 %d", c.ProjectPort),
		fmt.Sprintf("      Key MaxLength:        %d", c.ProjectKeyMaxLength),
		fmt.Sprintf("      Storage Type:         %s", c.ProjectStorageType),
		"   File storage:",
		fmt.Sprintf("      Is acync:             %t", c.FileAsync),
		"   Database:",
		fmt.Sprintf("      DSN:                  %s", c.DbDSN),
		fmt.Sprintf("      Max open connections: %d", c.DbMaxOpenConns),
		fmt.Sprintf("      Max idle connections: %d", c.DbMaxIdleConns),
		fmt.Sprintf("      Max idle time:        %s", c.DbMaxIdleTime),
		fmt.Sprintf("      Timeout (seconds):    %d", c.DbTimeout),
		"   Cache:",
		fmt.Sprintf("      Caching enabled:      %t", c.CacheEnabled),
		fmt.Sprintf("      Capacity of cache:    %d", c.CacheCapacity),
		"   Rate limiter:",
		fmt.Sprintf("      Rate limiter enabled: %t", c.LimiterEnabled),
		fmt.Sprintf("      RPS per IP:           %d", c.LimiterRPS),
		fmt.Sprintf("      Maximum burst:        %d", c.LimiterBurst),
	}

	return strings.Join(lines, "\n")
}

type Application struct {
	Config     Config
	Logger     *Logger
	Validator  Validator
	Links      LinksCollectionInterface
	Background *Background
}

type LinksCollectionInterface interface {
	GenerateKey(URL string) (string, error)
	GenerateKeys(URLs []string) (map[string]string, error)
	GetURL(key string) (string, error)
	GetURLs(keys []string) (map[string]string, error)
}

func (app *Application) Serve() error {
	server := &http.Server{
		Addr:         app.Config.ProjectHost + ":" + strconv.Itoa(app.Config.ProjectPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal)

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
