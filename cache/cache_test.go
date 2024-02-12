package cache

import (
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
	cache.cachedLinks["a"] = "url"
	cache.frequencies[1] = []string{"a"}
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

		require.Equal(t, URL, cache.cachedLinks[key])
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

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cache.cachedLinks)

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

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cache.cachedLinks)

	require.Equal(t, map[int][]string{
		1: {"d"},
		3: {"b", "c"},
	}, cache.frequencies)
}
