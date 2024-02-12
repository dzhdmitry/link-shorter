package links

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StorageInterface interface {
	StoreKeysURLs([][]string) error
	Restore() error
	GetURL(string) (string, error)
	GetURLs([]string) (map[string]string, error)
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

func (s *SQLStorage) GetURLs(keys []string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var placeholders []string
	var values []interface{}

	for i, key := range keys {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		values = append(values, key)
	}

	query := "SELECT key, url FROM links WHERE key IN (" + strings.Join(placeholders, ", ") + ") LIMIT " + strconv.Itoa(len(keys))
	rows, err := s.db.QueryContext(ctx, query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	URLs := make(map[string]string, len(keys))

	for rows.Next() {
		var key, URL string

		err := rows.Scan(&key, &URL)

		if err != nil {
			return nil, err
		}

		URLs[key] = URL
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return URLs, nil
}

func (s *SQLStorage) GetLastKey() (string, error) {
	return s.lastKey, nil
}
