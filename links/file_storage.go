package links

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"link-shorter.dzhdmitry.net/application"
	"os"
	"strconv"
	"sync"
)

type FileStorage struct {
	filename   string
	links      map[int64]string
	lastNumber int64
	mu         sync.Mutex
}

func NewFileStorage(filename string) (*FileStorage, error) {
	s := FileStorage{filename: filename, links: map[int64]string{}}
	err := s.Restore()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (fs *FileStorage) persist(idsURLs [][]string) error {
	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err = w.WriteAll(idsURLs); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) generate(URLs []string) ([][]string, map[string]string) {
	fs.mu.Lock()

	defer fs.mu.Unlock()

	var idsURLs [][]string
	keysByURLs := make(map[string]string, len(URLs))

	for _, URL := range URLs {
		fs.lastNumber++
		fs.links[fs.lastNumber] = URL
		idsURLs = append(idsURLs, []string{fmt.Sprintf("%d", fs.lastNumber), URL})
		keysByURLs[URL] = convertNumberToKey(fs.lastNumber)
	}

	return idsURLs, keysByURLs
}

// StoreURLs Returns map with key=URL, value=key
func (fs *FileStorage) StoreURLs(URLs []string) (map[string]string, error) {
	idsURLs, keysByURLs := fs.generate(URLs)

	if err := fs.persist(idsURLs); err != nil {
		return nil, err
	}

	return keysByURLs, nil
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

		idRaw, URL := record[0], record[1]
		id, err := strconv.ParseInt(idRaw, 10, 64)

		if err != nil {
			return err
		}

		fs.links[id] = URL
		fs.lastNumber = id
	}

	return nil
}

func (fs *FileStorage) GetURL(key string) (string, error) {
	return fs.links[convertKeyToNumber(key)], nil
}

func (fs *FileStorage) GetURLs(keys []string) (map[string]string, error) {
	URLs := make(map[string]string, len(keys))

	for _, key := range keys {
		if URL, ok := fs.links[convertKeyToNumber(key)]; ok {
			URLs[key] = URL
		}
	}

	return URLs, nil
}

type FileStorageAsync struct {
	logger     *application.Logger
	background *application.Background
	fs         *FileStorage
	mu         sync.Mutex
}

func NewFileStorageAsync(logger *application.Logger, background *application.Background, filename string) (*FileStorageAsync, error) {
	fs, err := NewFileStorage(filename)

	if err != nil {
		return nil, err
	}

	return &FileStorageAsync{
		logger:     logger,
		background: background,
		fs:         fs,
	}, nil
}

func (fsa *FileStorageAsync) StoreURLs(URLs []string) (map[string]string, error) {
	idsURLs, keysByURLs := fsa.fs.generate(URLs)

	fsa.background.Run(func() {
		fsa.mu.Lock()

		defer fsa.mu.Unlock()

		if err := fsa.fs.persist(idsURLs); err != nil {
			fsa.logger.LogError(err)
		}
	})

	return keysByURLs, nil
}

func (fsa *FileStorageAsync) Restore() error {
	return nil
}

func (fsa *FileStorageAsync) GetURL(key string) (string, error) {
	return fsa.fs.GetURL(key)
}

func (fsa *FileStorageAsync) GetURLs(keys []string) (map[string]string, error) {
	return fsa.fs.GetURLs(keys)
}
