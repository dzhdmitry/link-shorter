package application

import (
	"context"
	"errors"
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
	ProjectHost          string `env:"PROJECT_HOST"`
	ProjectPort          int    `env:"PROJECT_PORT"`
	ProjectKeyMaxLength  int    `env:"PROJECT_KEY_MAX_LENGTH"`
	ProjectStorageType   string `env:"PROJECT_STORAGE_TYPE"`
	FileAsync            bool   `env:"FILE_ASYNC"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
	DatabaseMaxOpenConns int    `env:"DATABASE_MAX_OPEN_CONNS"`
	DatabaseMaxIdleConns int    `env:"DATABASE_MAX_IDLE_CONNS"`
	DatabaseMaxIdleTime  string `env:"DATABASE_MAX_OPEN_TIME"`
	DatabaseTimeout      int    `env:"DATABASE_TIMEOUT"`
}

func NewConfig() Config {
	return Config{
		ProjectHost:          "",
		ProjectPort:          80,
		ProjectKeyMaxLength:  12,
		ProjectStorageType:   "file",
		FileAsync:            false,
		DatabaseDSN:          "postgres://go:pa55word@postgres:5432/short_links?sslmode=disable",
		DatabaseMaxOpenConns: 25,
		DatabaseMaxIdleConns: 25,
		DatabaseMaxIdleTime:  "15m",
		DatabaseTimeout:      1,
	}
}

func (c Config) Print(l *Logger) {
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
		fmt.Sprintf("      DSN:                  %s", c.DatabaseDSN),
		fmt.Sprintf("      Max open connections: %d", c.DatabaseMaxOpenConns),
		fmt.Sprintf("      Max idle connections: %d", c.DatabaseMaxIdleConns),
		fmt.Sprintf("      Max idle time:        %s", c.DatabaseMaxIdleTime),
		fmt.Sprintf("      Timeout (seconds):    %d", c.DatabaseTimeout),
	}

	l.LogInfo(strings.Join(lines, "\n"))
}

type Application struct {
	Config     Config
	Logger     Logger
	Validator  Validator
	Links      LinksCollectionInterface
	Background *Background
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
