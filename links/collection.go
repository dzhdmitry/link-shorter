package links

import (
	"link-shorter.dzhdmitry.net/generator"
	"sync"
)

type Collection struct {
	generator generator.Generator
	storage   StorageInterface
	mu        sync.Mutex
}

func NewCollection(generator generator.Generator, storage StorageInterface) *Collection {
	return &Collection{generator: generator, storage: storage}
}

func (c *Collection) GenerateKey(URL string) (string, error) {
	keys, err := c.GenerateKeys([]string{URL})

	if err != nil {
		return "", err
	}

	return keys[URL], nil
}

func (c *Collection) GenerateKeys(URLs []string) (map[string]string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	lastKey, err := c.storage.GetLastKey()

	if err != nil {
		return nil, err
	}

	generatedKeysSorted, err := c.generator.GenerateMany(lastKey, URLs)

	if err != nil {
		return nil, err
	}

	err = c.storage.StoreKeysURLs(generatedKeysSorted)

	if err != nil {
		return nil, err
	}

	keysByURLs := map[string]string{}

	for _, keyURL := range generatedKeysSorted {
		key, URL := keyURL[0], keyURL[1]
		keysByURLs[URL] = key
	}

	return keysByURLs, nil
}

func (c *Collection) GetURL(key string) (string, error) {
	return c.storage.GetURL(key)
}
