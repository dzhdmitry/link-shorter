package main

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"link-shorter.dzhdmitry.net/application"
	"os"
	"strconv"
)

func main() {
	config := application.Config{}
	logger := application.NewLogger(os.Stdout)

	if err := env.Parse(&config); err != nil {
		logger.LogError(err)

		os.Exit(1)
	}

	flag.StringVar(&config.ProjectHost, "host", config.ProjectHost, "Project server host")
	flag.IntVar(&config.ProjectPort, "port", config.ProjectPort, "Project server port")
	flag.Parse()

	app := application.Application{
		Config: config,
		Logger: *logger,
	}

	logger.LogInfo("Start server on " + config.ProjectHost + ":" + strconv.Itoa(config.ProjectPort))

	err := app.Serve()

	if err != nil {
		logger.LogError(err)

		os.Exit(1)
	}
}
