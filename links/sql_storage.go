package links

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type StorageInterface interface {
	StoreKeysURLs([][]string) error
	Restore() error
	GetURL(string) (string, error)
	GetLastKey() (string, error)
}

type SQLStorage struct {
	db      *sql.DB
	timeout time.Duration
	lastKey string
}

func NewSQLStorage(db *sql.DB, timeout int) (*SQLStorage, error) {
	s := SQLStorage{
		db:      db,
		timeout: time.Second * time.Duration(timeout),
	}
	err := s.Restore()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *SQLStorage) StoreKeysURLs(keysURLs [][]string) error {
	if len(keysURLs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var placeholders []string
	var values []interface{}
	n := 1

	for _, keyURL := range keysURLs {
		key, URL := keyURL[0], keyURL[1]
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", n, n+1))
		values = append(values, key, URL)
		n += 2
	}

	query := "INSERT INTO links(key, url) VALUES " + strings.Join(placeholders, ", ")
	_, err := s.db.ExecContext(ctx, query, values...)

	if err != nil {
		return err
	}

	s.lastKey = keysURLs[len(keysURLs)-1][0]

	return nil
}

func (s *SQLStorage) Restore() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var lastKey string
	query := "SELECT key FROM links ORDER BY id DESC LIMIT 1"
	err := s.db.QueryRowContext(ctx, query).Scan(&lastKey)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err

	}

	s.lastKey = lastKey

	return nil
}

func (s *SQLStorage) GetURL(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var URL string
	query := "SELECT url FROM links WHERE key = $1 LIMIT 1"
	err := s.db.QueryRowContext(ctx, query, key).Scan(&URL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return URL, nil
}

func (s *SQLStorage) GetLastKey() (string, error) {
	return s.lastKey, nil
}
