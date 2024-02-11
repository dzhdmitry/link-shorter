package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func Open(dsn string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func ResolveTestDSNs() (string, string) {
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
	defaultDsn, dsn := ResolveTestDSNs()
	db, err := Open(defaultDsn, 25, 25, "15m")

	if err != nil {
		fmt.Println("failed:")
		fmt.Println(defaultDsn)

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

	m, err := migrate.New("file://../migrations", dsn)

	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		panic(err)
	}

	return dsn
}
