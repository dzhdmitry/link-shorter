package cache

import (
	"slices"
	"sync"
)

type LFUCache struct {
	cachedLinks map[string]string
	frequencies map[int][]string // key - frequency, value - slice of keys, from old to new
	length      int
	capacity    int
	mu          sync.Mutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		cachedLinks: make(map[string]string, capacity),
		frequencies: map[int][]string{},
		length:      0,
		capacity:    capacity,
	}
}

func (c *LFUCache) incrementFrequency(key string) {
	for frequency, keys := range c.frequencies {
		index := slices.Index(keys, key)

		if index == -1 {
			continue
		}

		c.frequencies[frequency] = slices.Delete(c.frequencies[frequency], index, index+1)
		c.frequencies[frequency+1] = append(c.frequencies[frequency+1], key)

		if len(c.frequencies[frequency]) == 0 {
			delete(c.frequencies, frequency)
		}

		break
	}
}

func (c *LFUCache) Get(key string) (string, bool) {
	URL, ok := c.cachedLinks[key]

	if !ok {
		return "", false
	}

	c.mu.Lock()

	defer c.mu.Unlock()

	c.incrementFrequency(key)

	return URL, ok
}

func (c *LFUCache) Remember(key, URL string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.cachedLinks[key]; ok == true {
		// what if prev has already set key?
		c.incrementFrequency(key)

		return
	}

	if c.length >= c.capacity {
		for frequency, keys := range c.frequencies {
			if len(keys) == 0 {
				continue
			}

			delete(c.cachedLinks, keys[0])
			c.frequencies[frequency] = slices.Delete(c.frequencies[frequency], 0, 1)

			break
		}
	} else {
		c.length++
	}

	c.frequencies[1] = append(c.frequencies[1], key)
	c.cachedLinks[key] = URL
}
