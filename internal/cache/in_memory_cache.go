package cache

import (
	"container/list"
	"slices"
	"sync"
)

type CachedEntry struct {
	value   interface{}
	freqRef *list.Element
}

type FrequencyEntry struct {
	frequency int
	keys      []string // keys of frequency - from old to new
}

func NewFrequencyEntry(frequency int) *FrequencyEntry {
	return &FrequencyEntry{
		frequency: frequency,
		keys:      []string{},
	}
}

type LFUCache struct {
	cachedEntries map[string]*CachedEntry
	frequencies   *list.List // list of *FrequencyEntry
	length        int
	capacity      int
	mu            sync.Mutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		cachedEntries: make(map[string]*CachedEntry, capacity),
		frequencies:   list.New(),
		length:        0,
		capacity:      capacity,
	}
}

func (c *LFUCache) incrementFrequency(key string) {
	freqRef := c.cachedEntries[key].freqRef
	keyPosition := slices.Index(freqRef.Value.(*FrequencyEntry).keys, key)

	if keyPosition == -1 {
		panic("not found key in cache frequencies")
	}

	freqRef.Value.(*FrequencyEntry).keys = slices.Delete(freqRef.Value.(*FrequencyEntry).keys, keyPosition, keyPosition+1)
	next := freqRef.Next()

	if next == nil {
		next = c.frequencies.PushBack(NewFrequencyEntry(freqRef.Value.(*FrequencyEntry).frequency + 1))
	} else {
		if next.Value.(*FrequencyEntry).frequency != freqRef.Value.(*FrequencyEntry).frequency+1 {
			next = c.frequencies.InsertAfter(NewFrequencyEntry(freqRef.Value.(*FrequencyEntry).frequency+1), freqRef)
		}
	}

	next.Value.(*FrequencyEntry).keys = append(next.Value.(*FrequencyEntry).keys, key)
	c.cachedEntries[key].freqRef = next

	if len(freqRef.Value.(*FrequencyEntry).keys) == 0 {
		c.frequencies.Remove(freqRef)
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool, error) {
	entry, ok := c.cachedEntries[key]

	if !ok {
		return "", false, nil
	}

	c.mu.Lock()

	defer c.mu.Unlock()

	c.incrementFrequency(key)

	return entry.value, ok, nil
}

func (c *LFUCache) Put(key string, value interface{}) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.cachedEntries[key]; ok {
		// what if parallel task has already set key?
		c.incrementFrequency(key)

		return nil
	}

	if c.length >= c.capacity {
		for e := c.frequencies.Front(); e != nil; e.Next() {
			if len(e.Value.(*FrequencyEntry).keys) == 0 {
				continue
			}

			delete(c.cachedEntries, e.Value.(*FrequencyEntry).keys[0])
			e.Value.(*FrequencyEntry).keys = slices.Delete(e.Value.(*FrequencyEntry).keys, 0, 1)

			break
		}
	} else {
		c.length++
	}

	if c.frequencies.Len() == 0 {
		c.frequencies.PushBack(NewFrequencyEntry(1))
	} else {
		if c.frequencies.Front().Value.(*FrequencyEntry).frequency != 1 {
			c.frequencies.PushFront(NewFrequencyEntry(1))
		}
	}

	c.frequencies.Front().Value.(*FrequencyEntry).keys = append(c.frequencies.Front().Value.(*FrequencyEntry).keys, key)
	c.cachedEntries[key] = &CachedEntry{value: value, freqRef: c.frequencies.Front()}

	return nil
}
