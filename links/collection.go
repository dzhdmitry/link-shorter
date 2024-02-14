package links

import (
	"sync"
)

type StorageInterface interface {
	StoreURLs(URLs []string) (map[string]string, error)
	Restore() error
	GetURL(string) (string, error)
	GetURLs([]string) (map[string]string, error)
}

type Collection struct {
	storage StorageInterface
	mu      sync.Mutex
}

func NewCollection(storage StorageInterface) *Collection {
	return &Collection{storage: storage}
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

	keysByURLs, err := c.storage.StoreURLs(URLs)

	if err != nil {
		return nil, err
	}

	return keysByURLs, nil
}

func (c *Collection) GetURL(key string) (string, error) {
	return c.storage.GetURL(key)
}

func (c *Collection) GetURLs(keys []string) (map[string]string, error) {
	return c.storage.GetURLs(keys)
}
