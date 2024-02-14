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

type SQLStorage struct {
	db      *sql.DB
	timeout time.Duration
	lastKey string
}

func NewSQLStorage(db *sql.DB, timeout int) *SQLStorage {
	s := SQLStorage{
		db:      db,
		timeout: time.Second * time.Duration(timeout),
	}

	return &s
}

func (s *SQLStorage) StoreURLs(URLs []string) (map[string]string, error) {
	if len(URLs) == 0 {
		return map[string]string{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var placeholders []string
	var values []interface{}
	n := 1

	for _, URL := range URLs {
		placeholders = append(placeholders, fmt.Sprintf("($%d)", n))
		values = append(values, URL)
		n++
	}

	query := "INSERT INTO links(url) VALUES " + strings.Join(placeholders, ", ") + " RETURNING id"
	rows, err := s.db.QueryContext(ctx, query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keysByURLs := make(map[string]string, len(URLs))
	i := 0

	for rows.Next() {
		var id int

		err := rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		key := convertNumberToKey(id)
		keysByURLs[URLs[i]] = key
		i++
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keysByURLs, nil
}

func (s *SQLStorage) GetURL(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	var URL string
	query := "SELECT url FROM links WHERE id = $1 LIMIT 1"
	err := s.db.QueryRowContext(ctx, query, convertKeyToNumber(key)).Scan(&URL)

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
		values = append(values, convertKeyToNumber(key))
	}

	query := "SELECT id, url FROM links WHERE id IN (" + strings.Join(placeholders, ", ") + ") LIMIT " + strconv.Itoa(len(keys))
	rows, err := s.db.QueryContext(ctx, query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	URLs := make(map[string]string, len(keys))

	for rows.Next() {
		var id int
		var URL string

		err := rows.Scan(&id, &URL)

		if err != nil {
			return nil, err
		}

		URLs[convertNumberToKey(id)] = URL
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return URLs, nil
}
