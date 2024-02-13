package cache

import (
	"slices"
	"sort"
	"sync"
)

type CachedEntry struct {
	URL       string
	frequency int
}

type LFUCache struct {
	cachedEntries map[string]*CachedEntry
	frequencies   map[int][]string // key - frequency, value - slice of keys, from old to new
	length        int
	capacity      int
	mu            sync.Mutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		cachedEntries: make(map[string]*CachedEntry, capacity),
		frequencies:   map[int][]string{},
		length:        0,
		capacity:      capacity,
	}
}

func (c *LFUCache) incrementFrequency(key string) {
	keyFrequency := c.cachedEntries[key].frequency
	keyPosition := slices.Index(c.frequencies[keyFrequency], key)

	if keyPosition == -1 {
		panic("not found key in cache frequencies")
	}

	c.frequencies[keyFrequency] = slices.Delete(c.frequencies[keyFrequency], keyPosition, keyPosition+1)
	c.frequencies[keyFrequency+1] = append(c.frequencies[keyFrequency+1], key)
	c.cachedEntries[key].frequency++

	if len(c.frequencies[keyFrequency]) == 0 {
		delete(c.frequencies, keyFrequency)
	}
}

func (c *LFUCache) getFrequenciesOrdered() []int {
	frequencies := make([]int, len(c.frequencies))
	i := 0

	for frequency := range c.frequencies {
		frequencies[i] = frequency
		i++
	}

	sort.Ints(frequencies)

	return frequencies
}

func (c *LFUCache) Get(key string) (string, bool) {
	entry, ok := c.cachedEntries[key]

	if !ok {
		return "", false
	}

	c.mu.Lock()

	defer c.mu.Unlock()

	c.incrementFrequency(key)

	return entry.URL, ok
}

func (c *LFUCache) Remember(key, URL string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.cachedEntries[key]; ok == true {
		// what if parallel task has already set key?
		c.incrementFrequency(key)

		return
	}

	if c.length >= c.capacity {
		for _, frequency := range c.getFrequenciesOrdered() {
			if len(c.frequencies[frequency]) == 0 {
				continue
			}

			delete(c.cachedEntries, c.frequencies[frequency][0])
			c.frequencies[frequency] = slices.Delete(c.frequencies[frequency], 0, 1)

			break
		}
	} else {
		c.length++
	}

	c.cachedEntries[key] = &CachedEntry{URL: URL, frequency: 1}
	c.frequencies[1] = append(c.frequencies[1], key)
}
