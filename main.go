package main

import (
	"github.com/caarlos0/env/v10"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/container"
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

	config.Parse()
	logger.LogInfo(config.Info())

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

	logger.LogInfo("Start server on " + config.ProjectHost + ":" + strconv.Itoa(config.ProjectPort))

	err = app.Serve()

	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}
}
