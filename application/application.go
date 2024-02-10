package application

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	ProjectHost          string `env:"PROJECT_HOST"`
	ProjectPort          int    `env:"PROJECT_PORT"`
	ProjectKeyMaxLength  int    `env:"PROJECT_KEY_MAX_LENGTH"`
	ProjectStorageType   string `env:"PROJECT_STORAGE_TYPE"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
	DatabaseMaxOpenConns int    `env:"DATABASE_MAX_OPEN_CONNS"`
	DatabaseMaxIdleConns int    `env:"DATABASE_MAX_IDLE_CONNS"`
	DatabaseMaxIdleTime  string `env:"DATABASE_MAX_OPEN_TIME"`
}

type Application struct {
	Config    Config
	Logger    Logger
	Validator Validator
	Links     LinksCollectionInterface
}

type LinksCollectionInterface interface {
	GenerateKey(URL string) (string, error)
	GenerateKeys(URLs []string) (map[string]string, error)
	GetURL(key string) (string, error)
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

		app.Logger.LogInfo("completing work...")

		// todo complete work...

		app.Logger.LogInfo("work completed")

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
