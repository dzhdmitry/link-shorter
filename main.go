package main

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v10"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/db"
	"link-shorter.dzhdmitry.net/generator"
	"link-shorter.dzhdmitry.net/links"
	"os"
	"strconv"
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
	flag.IntVar(&config.ProjectKeyMaxLength, "key_max_length", config.ProjectKeyMaxLength, "Max length of the key")
	flag.StringVar(&config.ProjectStorageType, "storage", config.ProjectStorageType, "Storage type (file|postgres)")
	flag.BoolVar(&config.FileAsync, "file-async", config.FileAsync, "File storage is asynchronous|synchronous (true|false)")
	flag.StringVar(&config.DatabaseDSN, "db-dsn", config.DatabaseDSN, "PostgreSQL DSN")
	flag.IntVar(&config.DatabaseMaxOpenConns, "db-max-open-conns", config.DatabaseMaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&config.DatabaseMaxIdleConns, "db-max-idle-conns", config.DatabaseMaxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&config.DatabaseMaxIdleTime, "db-max-idle-time", config.DatabaseMaxIdleTime, "PostgreSQL max connection idle time")
	flag.IntVar(&config.DatabaseTimeout, "db-timeout", config.DatabaseTimeout, "PostgreSQL queries execution timeout")
	flag.Parse()

	config.Print(logger)

	var storage links.StorageInterface
	var err error

	if config.ProjectStorageType == "file" {
		if config.FileAsync {
			storage, err = links.NewFileStorageAsync(*logger, background, "tmp/storage.csv")
		} else {
			storage, err = links.NewFileStorage("tmp/storage.csv")
		}
	} else if config.ProjectStorageType == "postgres" {
		dbx, errDB := db.Open(config.DatabaseDSN, config.DatabaseMaxOpenConns, config.DatabaseMaxIdleConns, config.DatabaseMaxIdleTime)

		if errDB != nil {
			logger.LogError(errDB)
			os.Exit(1)
		}

		defer dbx.Close()

		storage, err = links.NewSQLStorage(dbx, config.DatabaseTimeout)
	} else {
		logger.LogError(errors.New("unknown type " + config.ProjectStorageType))
		os.Exit(1)
	}

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	app := application.Application{
		Config: config,
		Logger: *logger,
		Validator: application.Validator{
			KeyMaxLength: config.ProjectKeyMaxLength,
		},
		Links:      links.NewCollection(*generator.NewGenerator(config.ProjectKeyMaxLength), storage),
		Background: background,
	}

	logger.LogInfo("Start server on " + config.ProjectHost + ":" + strconv.Itoa(config.ProjectPort))

	err = app.Serve()

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}
}
