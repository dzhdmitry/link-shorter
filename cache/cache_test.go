package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetNonExisting(t *testing.T) {
	cache := NewLFUCache(5)
	URL, ok := cache.Get("a")

	require.False(t, ok)
	require.Equal(t, "", URL)
}

func TestGetExisting(t *testing.T) {
	cache := NewLFUCache(5)

	cache.Remember("a", "url")

	URL, ok := cache.Get("a")

	require.True(t, ok)
	require.Equal(t, "url", URL)
	require.Equal(t, []string{"a"}, cache.frequencies[2])
}

func TestRemember(t *testing.T) {
	cache := NewLFUCache(5)
	links := [][]string{
		{"a", "url1"},
		{"b", "url2"},
		{"c", "url3"},
	}

	for _, keyURL := range links {
		key, URL := keyURL[0], keyURL[1]

		cache.Remember(key, URL)

		require.Equal(t, URL, cache.cachedEntries[key].URL)
		require.Equal(t, []string{key}, cache.frequencies[1])

		cachedURL, ok := cache.Get(key)

		require.True(t, ok)
		require.Equal(t, URL, cachedURL)
	}

	require.Equal(t, map[int][]string{2: {"a", "b", "c"}}, cache.frequencies)
}

func TestRememberReplace(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Remember("a", "url1")
	cache.Remember("b", "url2")
	cache.Remember("c", "url3")
	cache.Get("b")
	cache.Get("c")
	cache.Remember("d", "url4")

	cachedEntries := map[string]string{}
	cachedFrequencies := map[string]int{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = v.URL
		cachedFrequencies[k] = v.frequency
	}

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cachedEntries)

	assert.Equal(t, map[string]int{
		"b": 2,
		"c": 2,
		"d": 1,
	}, cachedFrequencies)

	require.Equal(t, map[int][]string{
		1: {"d"},
		2: {"b", "c"},
	}, cache.frequencies)
}

func TestRememberReplaceEmptyFirstFrequency(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Remember("a", "url1")
	cache.Remember("b", "url2")
	cache.Remember("c", "url3")
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Remember("d", "url4")

	cachedEntries := map[string]string{}
	cachedFrequencies := map[string]int{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = v.URL
		cachedFrequencies[k] = v.frequency
	}

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cachedEntries)

	assert.Equal(t, map[string]int{
		"b": 3,
		"c": 3,
		"d": 1,
	}, cachedFrequencies)

	require.Equal(t, map[int][]string{
		1: {"d"},
		3: {"b", "c"},
	}, cache.frequencies)
}

func TestRememberReplaceFirstFrequency(t *testing.T) {
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
			cache.Remember(links[i][0], links[i][1])
		}
	}

	cache.Remember("g", "url7")

	cachedEntries := map[string]string{}
	cachedFrequencies := map[string]int{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = v.URL
		cachedFrequencies[k] = v.frequency
	}

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
		"e": "url5",
		"f": "url6",
		"g": "url7",
	}, cachedEntries)

	assert.Equal(t, map[string]int{
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 1,
	}, cachedFrequencies)

	assert.Equal(t, map[int][]string{
		1: {"g"},
		2: {"b"},
		3: {"c"},
		4: {"d"},
		5: {"e"},
		6: {"f"},
	}, cache.frequencies)
}

func TestRememberReplaceMultiple(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Remember("a", "url1")
	cache.Remember("b", "url2")
	cache.Remember("c", "url3")

	cache.Get("a")
	cache.Get("b")
	cache.Get("c")

	cache.Remember("d", "url4")
	cache.Remember("e", "url5")
	cache.Remember("f", "url6")

	cachedEntries := map[string]string{}
	cachedFrequencies := map[string]int{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = v.URL
		cachedFrequencies[k] = v.frequency
	}

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"f": "url6",
	}, cachedEntries)

	assert.Equal(t, map[string]int{
		"b": 2,
		"c": 2,
		"f": 1,
	}, cachedFrequencies)

	assert.Equal(t, map[int][]string{
		1: {"f"},
		2: {"b", "c"},
	}, cache.frequencies)
}
