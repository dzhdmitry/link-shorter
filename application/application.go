package application

import (
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	ProjectHost string `env:"PROJECT_HOST"`
	ProjectPort int    `env:"PROJECT_PORT"`
}

type Application struct {
	Config Config
	Logger Logger
}

func (app *Application) Serve() error {
	server := &http.Server{
		Addr:         app.Config.ProjectHost + ":" + strconv.Itoa(app.Config.ProjectPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
