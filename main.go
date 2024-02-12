package main

import (
	"github.com/caarlos0/env/v10"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/container"
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

	config.Parse()
	config.Print(logger)

	Container := container.Container{
		Logger:     logger,
		Background: background,
	}

	linksCollection, dbConn, err := Container.CreateLinksCollection(config)

	if dbConn != nil {
		defer dbConn.Close()
	}

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	app := application.Application{
		Config: config,
		Logger: logger,
		Validator: application.Validator{
			KeyMaxLength: config.ProjectKeyMaxLength,
		},
		Links:      linksCollection,
		Background: background,
	}

	logger.LogInfo("Start server on " + config.ProjectHost + ":" + strconv.Itoa(config.ProjectPort))

	err = app.Serve()

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}
}
