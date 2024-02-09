package links_in_memory

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type StorageInterface interface {
	Store(key, URL string) error
	Restore(map[string]string) error
}

type FileStorage struct {
	Filename string
}

func (fs *FileStorage) Store(key, URL string) error {
	file, err := os.OpenFile(fs.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err = w.Write([]string{key, URL}); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) Restore(links map[string]string) error {
	file, err := os.Open(fs.Filename)

	if err != nil {
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

		links[record[0]] = record[1]
	}

	return nil
}
