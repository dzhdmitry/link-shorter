package application

import (
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	ProjectHost         string `env:"PROJECT_HOST"`
	ProjectPort         int    `env:"PROJECT_PORT"`
	ProjectKeyMaxLength int    `env:"PROJECT_KEY_MAX_LENGTH"`
}

type Application struct {
	Config Config
	Logger Logger
	Links  LinksStorage
}

type LinksStorage interface {
	GenerateKey(URL string) (string, error)
	GetLink(key string) string
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

func (app *Application) composeShortLink(key string) string {
	host := app.Config.ProjectHost

	if host == "" {
		host = "localhost"
	}

	return "http://" + host + "/go/" + key
}
