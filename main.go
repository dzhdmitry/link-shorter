package main

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/container"
	"link-shorter.dzhdmitry.net/links"
	"os"
)

func main() {
	config := application.NewConfig()
	logger := application.NewLogger(os.Stdout, &application.Clock{})
	background := &application.Background{}

	if err := env.Parse(&config); err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	flag.StringVar(&config.ProjectHost, "host", config.ProjectHost, "Project server host")
	flag.IntVar(&config.ProjectPort, "port", config.ProjectPort, "Project server port")
	flag.StringVar(&config.ProjectStorageType, "storage", config.ProjectStorageType, "Storage type (file|postgres)")
	flag.BoolVar(&config.FileAsync, "file-async", config.FileAsync, "File storage is asynchronous|synchronous (true|false)")
	flag.StringVar(&config.DbDSN, "db-dsn", config.DbDSN, "PostgreSQL DSN")
	flag.IntVar(&config.DbMaxOpenConns, "db-max-open-conns", config.DbMaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&config.DbMaxIdleConns, "db-max-idle-conns", config.DbMaxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&config.DbMaxIdleTime, "db-max-idle-time", config.DbMaxIdleTime, "PostgreSQL max connection idle time")
	flag.IntVar(&config.DbTimeout, "db-timeout", config.DbTimeout, "PostgreSQL queries execution timeout")
	flag.StringVar(&config.CacheType, "cache", config.CacheType, "Cache type (disabled|in-memory|redis)")
	flag.IntVar(&config.CacheCapacity, "cache-cap", config.CacheCapacity, "Capacity of in-memory cache")
	flag.StringVar(&config.CacheRedisDSN, "redis", config.CacheRedisDSN, "Redis DSN")
	flag.BoolVar(&config.LimiterEnabled, "limiter", config.LimiterEnabled, "Rate limiter is enabled")
	flag.IntVar(&config.LimiterRPS, "limiter-rps", config.LimiterRPS, "Rate limiter maximum RPS per IP")
	flag.IntVar(&config.LimiterBurst, "limiter-burst", config.LimiterBurst, "Rate limiter maximum burst")
	flag.Parse()

	Container := container.Container{
		Logger:     logger,
		Background: background,
	}

	linksCollection, dbConn, rdb, err := Container.CreateLinksCollection(config)

	if dbConn != nil {
		defer dbConn.Close()
	}

	if rdb != nil {
		defer rdb.Close()
	}

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	app := application.Application{
		Config:     config,
		Logger:     logger,
		Validator:  *application.NewValidator(links.Letters),
		Links:      linksCollection,
		Background: background,
	}

	logger.LogInfo(config.Info())

	err = app.Serve()

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}
}
