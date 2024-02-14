package cache

import (
	"fmt"
	"link-shorter.dzhdmitry.net/application"
)

type MemoryCacheInterface interface {
	Get(string) (interface{}, bool)
	Put(string, interface{})
}

type CachedCollection struct {
	collection application.LinksCollectionInterface
	cache      MemoryCacheInterface
}

func NewCachedCollection(collection application.LinksCollectionInterface, cache MemoryCacheInterface) *CachedCollection {
	return &CachedCollection{
		collection: collection,
		cache:      cache,
	}
}

func (c *CachedCollection) GenerateKey(URL string) (string, error) {
	return c.collection.GenerateKey(URL)
}

func (c *CachedCollection) GenerateKeys(URLs []string) (map[string]string, error) {
	return c.collection.GenerateKeys(URLs)
}

func (c *CachedCollection) GetURL(key string) (string, error) {
	cachedURL, ok := c.cache.Get(key)

	if ok {
		return fmt.Sprintf("%s", cachedURL), nil
	}

	URL, err := c.collection.GetURL(key)

	if err != nil {
		return "", err
	}

	c.cache.Put(key, URL)

	return URL, nil
}

func (c *CachedCollection) GetURLs(keys []string) (map[string]string, error) {
	URLs := make(map[string]string, len(keys))

	for _, key := range keys {
		URL, err := c.GetURL(key)

		if err != nil {
			return nil, err
		}

		URLs[key] = URL
	}

	return URLs, nil
}
