package cache

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func collectCachedEntries(cachedEntries map[string]*CachedEntry) map[string]string {
	collectedEntries := map[string]string{}

	for k, v := range cachedEntries {
		collectedEntries[k] = fmt.Sprintf("%s", v.value)
	}

	return collectedEntries
}

func collectFrequencies(frequencies *list.List) map[int][]string {
	e := frequencies.Front()
	keysByFrequency := map[int][]string{}

	for {
		if e == nil {
			break
		}

		keysByFrequency[e.Value.(*FrequencyEntry).frequency] = e.Value.(*FrequencyEntry).keys
		e = e.Next()
	}

	return keysByFrequency
}

func TestGetNonExisting(t *testing.T) {
	cache := NewLFUCache(5)
	URL, ok := cache.Get("a")

	require.False(t, ok)
	require.Equal(t, "", URL)
}

func TestGetExisting(t *testing.T) {
	cache := NewLFUCache(5)

	cache.Put("a", "url")

	URL, ok := cache.Get("a")

	require.True(t, ok)
	require.Equal(t, "url", URL)
	assert.Equal(t, map[int][]string{
		2: {"a"},
	}, collectFrequencies(cache.frequencies))
}

func TestPut(t *testing.T) {
	cache := NewLFUCache(5)
	links := [][]string{
		{"a", "url1"},
		{"b", "url2"},
		{"c", "url3"},
	}

	for _, keyURL := range links {
		key, URL := keyURL[0], keyURL[1]

		cache.Put(key, URL)

		require.Equal(t, URL, cache.cachedEntries[key].value)

		cachedURL, ok := cache.Get(key)

		require.True(t, ok)
		require.Equal(t, URL, cachedURL)
	}

	assert.Equal(t, map[int][]string{
		2: {"a", "b", "c"},
	}, collectFrequencies(cache.frequencies))
}

func TestPutEvict(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Put("a", "url1")
	cache.Put("b", "url2")
	cache.Put("c", "url3")
	cache.Get("b")
	cache.Get("c")
	cache.Put("d", "url4")

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, collectCachedEntries(cache.cachedEntries))

	assert.Equal(t, map[int][]string{
		1: {"d"},
		2: {"b", "c"},
	}, collectFrequencies(cache.frequencies))
}

func TestPutEvictEmptyFirstFrequency(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Put("a", "url1")
	cache.Put("b", "url2")
	cache.Put("c", "url3")
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Put("d", "url4")

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, collectCachedEntries(cache.cachedEntries))

	assert.Equal(t, map[int][]string{
		1: {"d"},
		3: {"b", "c"},
	}, collectFrequencies(cache.frequencies))
}

func TestPutEvictFirstFrequency(t *testing.T) {
	cache := NewLFUCache(6)
	links := [][]string{
		{"a", "url1"},
		{"b", "url2"},
		{"c", "url3"},
		{"d", "url4"},
		{"e", "url5"},
		{"f", "url6"},
	}

	for i := 0; i < len(links); i++ {
		for j := 0; j < i+1; j++ {
			cache.Put(links[i][0], links[i][1])
		}
	}

	cache.Put("g", "url7")

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
		"e": "url5",
		"f": "url6",
		"g": "url7",
	}, collectCachedEntries(cache.cachedEntries))

	assert.Equal(t, map[int][]string{
		1: {"g"},
		2: {"b"},
		3: {"c"},
		4: {"d"},
		5: {"e"},
		6: {"f"},
	}, collectFrequencies(cache.frequencies))
}

func TestPutEvictMultiple(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Put("a", "url1")
	cache.Put("b", "url2")
	cache.Put("c", "url3")
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Put("d", "url4")
	cache.Put("e", "url5")
	cache.Put("f", "url6")

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"f": "url6",
	}, collectCachedEntries(cache.cachedEntries))

	assert.Equal(t, map[int][]string{
		1: {"f"},
		2: {"b", "c"},
	}, collectFrequencies(cache.frequencies))
}

func TestPutRangeBetweenFrequencies(t *testing.T) {
	cache := NewLFUCache(5)

	cache.Put("a", "url1")
	cache.Put("b", "url2")
	cache.Put("c", "url3")
	cache.Put("d", "url4")

	cache.Get("a")
	cache.Get("a")
	cache.Get("b")

	assert.Equal(t, map[string]string{
		"a": "url1",
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, collectCachedEntries(cache.cachedEntries))

	assert.Equal(t, map[int][]string{
		1: {"c", "d"},
		2: {"b"},
		3: {"a"},
	}, collectFrequencies(cache.frequencies))
}
