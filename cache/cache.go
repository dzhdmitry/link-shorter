package cache

import "sync"

type LFUCache struct {
	cachedLinks map[string]string
	frequencies map[string]int
	capacity    int
	mu          sync.Mutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		cachedLinks: map[string]string{},
		frequencies: map[string]int{},
		capacity:    capacity,
	}
}

func (c *LFUCache) Get(key string) (string, bool) {
	URL, ok := c.cachedLinks[key]

	if ok {
		c.mu.Lock()

		defer c.mu.Unlock()

		c.frequencies[key] += 1

		return URL, ok
	}

	return "", false
}

func (c *LFUCache) displace() {
	leastFrequentlyUsed := ""
	minFrequency := 1

	for k, f := range c.frequencies {
		if leastFrequentlyUsed == "" {
			// first element
			leastFrequentlyUsed = k
			minFrequency = f

			continue
		}

		if f < minFrequency {
			// replace
			leastFrequentlyUsed = k
			minFrequency = f
		}
	}

	if leastFrequentlyUsed == "" {
		panic("not found leastFrequentlyUsed in cache")
	}

	delete(c.cachedLinks, leastFrequentlyUsed)
	delete(c.frequencies, leastFrequentlyUsed)
}

func (c *LFUCache) Remember(key, URL string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.cachedLinks[key]; ok == true {
		// what if prev has already set key?
		c.frequencies[key] += 1

		return
	}

	if len(c.cachedLinks) >= c.capacity {
		c.displace()
	}

	c.cachedLinks[key] = URL
	c.frequencies[key] = 1
}
