package test

import (
	"github.com/dzhdmitry/link-shorter/internal/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Writer struct {
	Messages []string
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.Messages = append(w.Messages, string(p))

	return len(p), nil
}

type Clock struct {
	//
}

func (t *Clock) Now() time.Time {
	location, _ := time.LoadLocation("Europe/London")

	return time.Date(2024, 2, 7, 12, 0, 0, 0, location)
}

func resolveTestDSNs() (string, string) {
	defaultDsn := os.Getenv("TEST_DEFAULT_DSN")

	if defaultDsn == "" {
		defaultDsn = "postgres://go:pa55word@localhost:5432/postgres?sslmode=disable"
	}

	dsn := os.Getenv("TEST_DSN")

	if dsn == "" {
		dsn = "postgres://go:pa55word@localhost:5432/short_links_test?sslmode=disable"
	}

	return defaultDsn, dsn
}

func PrepareTestDB() string {
	defaultDsn, dsn := resolveTestDSNs()
	db, err := db.OpenPostgres(defaultDsn, 25, 25, "15m")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS short_links_test")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE short_links_test WITH OWNER go")

	if err != nil {
		panic(err)
	}

	db.Close()

	m, err := migrate.New("file://../../migrations", dsn)

	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		panic(err)
	}

	return dsn
}

func OpenTestRedis() (*redis.Client, error) {
	defaultDsn := os.Getenv("TEST_REDIS_DSN")

	if defaultDsn == "" {
		defaultDsn = "redis://localhost:6379/1"
	}

	return db.OpenRedis(defaultDsn)
}
