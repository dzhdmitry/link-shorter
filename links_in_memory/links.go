package links_in_memory

import (
	"link-shorter.dzhdmitry.net/generator"
	"sync"
)

type LinksCollection struct {
	links        map[string]string
	lastKey      string
	keyMaxLength int
	mu           sync.Mutex
}

func NewLinksCollection(keyMaxLength int) *LinksCollection {
	return &LinksCollection{
		keyMaxLength: keyMaxLength,
		links:        map[string]string{},
	}
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

	return key, nil
}

func (lc *LinksCollection) GetLink(key string) string {
	return lc.links[key]
}
