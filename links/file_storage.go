package links

import (
	"encoding/csv"
	"errors"
	"io"
	"link-shorter.dzhdmitry.net/application"
	"os"
	"strconv"
	"sync"
)

type FileStorage struct {
	filename   string
	links      map[string]string
	lastKey    string
	lastNumber int
}

func NewFileStorage(filename string) (*FileStorage, error) {
	s := FileStorage{filename: filename, links: map[string]string{}}
	err := s.Restore()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (fs *FileStorage) persist(keysURLs [][]string) error {
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

	return nil
}

func (fs *FileStorage) generate(URLs []string) ([][]string, map[string]string) {
	var keysURLs [][]string
	keysByURLs := make(map[string]string, len(URLs))

	for _, URL := range URLs {
		fs.lastNumber++
		key := numberToKey(fs.lastNumber)
		keysURLs = append(keysURLs, []string{key, URL})
		keysByURLs[URL] = key
		// todo mutex
	}

	return keysURLs, keysByURLs
}

func (fs *FileStorage) remember(keysURLs [][]string) { // todo map?
	for _, keysURL := range keysURLs {
		key, URL := keysURL[0], keysURL[1]
		fs.links[key] = URL
		fs.lastKey = key
	}
}

// StoreURLs Returns map with key=URL, value=key
func (fs *FileStorage) StoreURLs(URLs []string) (map[string]string, error) {
	keysURLs, keysByURLs := fs.generate(URLs)

	if err := fs.persist(keysURLs); err != nil {
		return nil, err
	}

	fs.remember(keysURLs)

	return keysByURLs, nil
}

// todo remove
func (fs *FileStorage) StoreKeysURLs(keysURLs [][]string) error {
	if err := fs.persist(keysURLs); err != nil {
		return err
	}

	fs.remember(keysURLs)

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

		id, URL := record[0], record[1]
		idx, _ := strconv.Atoi(id) // todo err
		fs.links[numberToKey(idx)] = URL
		fs.lastNumber = idx
	}

	return nil
}

func (fs *FileStorage) GetURL(key string) (string, error) {
	return fs.links[key], nil
}

func (fs *FileStorage) GetURLs(keys []string) (map[string]string, error) {
	URLs := make(map[string]string, len(keys))

	for _, k := range keys {
		if URL, ok := fs.links[k]; ok {
			URLs[k] = URL
		}
	}

	return URLs, nil
}

// todo delete
func (fs *FileStorage) GetLastKey() (string, error) {
	return fs.lastKey, nil
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
	keysURLs, keysByURLs := fsa.fs.generate(URLs)

	fsa.background.Run(func() {
		fsa.mu.Lock()

		defer fsa.mu.Unlock()

		if err := fsa.fs.persist(keysURLs); err != nil {
			fsa.logger.LogError(err)
		}
	})

	fsa.fs.remember(keysURLs)

	return keysByURLs, nil
}

// todo remove
func (fsa *FileStorageAsync) StoreKeysURLs(keysURLs [][]string) error {
	fsa.background.Run(func() {
		fsa.mu.Lock()

		defer fsa.mu.Unlock()

		if err := fsa.fs.persist(keysURLs); err != nil {
			fsa.logger.LogError(err)
		}
	})

	fsa.fs.remember(keysURLs)

	return nil
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

// todo remove
func (fsa *FileStorageAsync) GetLastKey() (string, error) {
	return fsa.fs.GetLastKey()
}
