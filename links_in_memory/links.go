package links_in_memory

import (
	"link-shorter.dzhdmitry.net/generator"
	"sync"
)

type LinksCollection struct {
	links        map[string]string
	lastKey      string
	keyMaxLength int
	storage      StorageInterface
	mu           sync.Mutex
}

func NewLinksCollection(storage StorageInterface, keyMaxLength int) (*LinksCollection, error) {
	lc := &LinksCollection{
		links:        map[string]string{},
		keyMaxLength: keyMaxLength,
		storage:      storage,
	}

	lastKey, err := storage.Restore(lc.links)

	if err != nil {
		return nil, err
	}

	lc.lastKey = lastKey

	return lc, nil
}

func (lc *LinksCollection) GenerateKey(URL string) (string, error) {
	lc.mu.Lock()

	defer lc.mu.Unlock()

	key, err := generator.Generate(lc.lastKey, lc.keyMaxLength)

	if err != nil {
		return "", err
	}

	lc.links[key] = URL
	lc.lastKey = key
	err = lc.storage.Store(key, URL)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (lc *LinksCollection) GetLink(key string) string {
	return lc.links[key]
}
