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
	cache.frequencies["a"] = 1
	URL, ok := cache.Get("a")

	require.True(t, ok)
	require.Equal(t, "url", URL)
	require.Equal(t, 2, cache.frequencies["a"])
}

func TestRemember(t *testing.T) {
	cache := NewLFUCache(5)

	links := map[string]string{
		"a": "url1",
		"b": "url2",
		"c": "url3",
	}

	for k, u := range links {
		cache.Remember(k, u)

		require.Equal(t, u, cache.cachedLinks[k])
		require.Equal(t, 1, cache.frequencies[k])

		URL, ok := cache.Get(k)

		require.True(t, ok)
		require.Equal(t, u, URL)
	}
	require.Equal(t, map[string]int{
		"a": 2,
		"b": 2,
		"c": 2,
	}, cache.frequencies)
}

func TestRememberDisplace(t *testing.T) {
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
	require.Equal(t, map[string]int{
		"b": 2,
		"c": 2,
		"d": 1,
	}, cache.frequencies)
}
