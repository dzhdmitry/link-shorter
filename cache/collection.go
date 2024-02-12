package cache

import (
	"link-shorter.dzhdmitry.net/application"
)

type MemoryCacheInterface interface {
	Get(string) (string, bool)
	Remember(string, string)
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
		return cachedURL, nil
	}

	URL, err := c.collection.GetURL(key)

	if err != nil {
		return "", err
	}

	c.cache.Remember(key, URL)

	return URL, nil
}