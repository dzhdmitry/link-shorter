package links

type StorageInterface interface {
	StoreURLs(URLs []string) (map[string]string, error)
	GetURL(string) (string, error)
	GetURLs([]string) (map[string]string, error)
}

type Collection struct {
	storage StorageInterface
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
	return c.storage.StoreURLs(URLs)
}

func (c *Collection) GetURL(key string) (string, error) {
	return c.storage.GetURL(key)
}

func (c *Collection) GetURLs(keys []string) (map[string]string, error) {
	return c.storage.GetURLs(keys)
}
