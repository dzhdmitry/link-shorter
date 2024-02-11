package links

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type StorageInterface interface {
	StoreKeysURLs([][]string) error
	Restore() error
	GetURL(string) (string, error)
	GetLastKey() (string, error)
}

type FileStorage struct {
	filename string
	links    map[string]string
	lastKey  string
}

func NewFileStorage(filename string) (*FileStorage, error) {
	s := FileStorage{filename: filename, links: map[string]string{}}
	err := s.Restore()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (fs *FileStorage) StoreKeysURLs(keysURLs [][]string) error {
	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err = w.WriteAll(keysURLs); err != nil {
		return err
	}

	for _, ku := range keysURLs {
		key, URL := ku[0], ku[1]
		fs.links[key] = URL
		fs.lastKey = key
	}

	return nil
}

func (fs *FileStorage) Restore() error {
	file, err := os.Open(fs.filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lastKey := ""

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(record) != 2 {
			return errors.New("file has malformed data")
		}

		key, URL := record[0], record[1]
		fs.links[key] = URL
		lastKey = key
	}

	fs.lastKey = lastKey

	return nil
}

func (fs *FileStorage) GetURL(key string) (string, error) {
	return fs.links[key], nil
}

func (fs *FileStorage) GetLastKey() (string, error) {
	return fs.lastKey, nil
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
