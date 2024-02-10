package links_in_memory

import (
	"link-shorter.dzhdmitry.net/generator"
	"sync"
)

type LinksCollection struct {
	generator    generator.Generator
	links        map[string]string
	lastKey      string
	keyMaxLength int
	storage      StorageInterface
	mu           sync.Mutex
}

func NewLinksCollection(gen generator.Generator, storage StorageInterface, keyMaxLength int) (*LinksCollection, error) {
	lc := &LinksCollection{
		generator:    gen,
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
	keys, err := lc.GenerateKeys([]string{URL})

	if err != nil {
		return "", err
	}

	return keys[URL], nil
}

func (lc *LinksCollection) GenerateKeys(URLs []string) (map[string]string, error) {
	lc.mu.Lock()

	defer lc.mu.Unlock()

	urlsByGeneratedKeys := map[string]string{} // -> lc.storage
	generatedKeysSorted := [][]string{}        // -> lc.links
	keysByURLs := map[string]string{}          // -> output
	currentLastKey := lc.lastKey

	for _, URL := range URLs {
		key, err := lc.generator.Generate(currentLastKey, lc.keyMaxLength)

		if err != nil {
			return nil, err
		}

		urlsByGeneratedKeys[key] = URL
		generatedKeysSorted = append(generatedKeysSorted, []string{key, URL})
		currentLastKey = key
	}

	err := lc.storage.StoreURLs(urlsByGeneratedKeys)

	if err != nil {
		return nil, err
	}

	for _, keyURL := range generatedKeysSorted {
		key, URL := keyURL[0], keyURL[1]
		lc.links[key] = URL
		lc.lastKey = key
		keysByURLs[URL] = key
	}

	return keysByURLs, nil
}

func (lc *LinksCollection) GetLink(key string) string {
	return lc.links[key]
}
