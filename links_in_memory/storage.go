package links_in_memory

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type StorageInterface interface {
	StoreURLs(URLs map[string]string) error
	Restore(map[string]string) (string, error)
}

type FileStorage struct {
	Filename string
}

func (fs *FileStorage) StoreURLs(URLs map[string]string) error {
	file, err := os.OpenFile(fs.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string

	for key, URL := range URLs {
		data = append(data, []string{key, URL})
	}

	if err = w.WriteAll(data); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) Restore(links map[string]string) (string, error) {
	file, err := os.Open(fs.Filename)

	if err != nil {
		return "", err
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
			return "", err
		}

		if len(record) != 2 {
			return "", errors.New("file has malformed data")
		}

		key, URL := record[0], record[1]
		links[key] = URL
		lastKey = key
	}

	return lastKey, nil
}
