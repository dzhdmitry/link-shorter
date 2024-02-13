package cache

import (
	"fmt"
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

	cache.Put("a", "url")

	URL, ok := cache.Get("a")

	require.True(t, ok)
	require.Equal(t, "url", URL)
	require.Equal(t, []string{"a"}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)
	require.Equal(t, 1, cache.frequencies.Len())
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

		cache.Put(key, URL)

		require.Equal(t, URL, cache.cachedEntries[key].value)
		require.Equal(t, []string{key}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)

		cachedURL, ok := cache.Get(key)

		require.True(t, ok)
		require.Equal(t, URL, cachedURL)
	}

	require.Equal(t, []string{"a", "b", "c"}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)
	require.Equal(t, 1, cache.frequencies.Len())
}

func TestRememberReplace(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Put("a", "url1")
	cache.Put("b", "url2")
	cache.Put("c", "url3")
	cache.Get("b")
	cache.Get("c")
	cache.Put("d", "url4")

	cachedEntries := map[string]string{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = fmt.Sprintf("%s", v.value)
	}

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cachedEntries)

	assert.Equal(t, []string{"d"}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)
	assert.Equal(t, []string{"b", "c"}, cache.frequencies.Back().Value.(*FrequencyEntry).keys)
	require.Equal(t, 2, cache.frequencies.Len())
}

func TestRememberReplaceEmptyFirstFrequency(t *testing.T) {
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

	cachedEntries := map[string]string{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = fmt.Sprintf("%s", v.value)
	}

	require.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
	}, cachedEntries)

	assert.Equal(t, []string{"d"}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)
	assert.Equal(t, []string{"b", "c"}, cache.frequencies.Back().Value.(*FrequencyEntry).keys)
	require.Equal(t, 2, cache.frequencies.Len())
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
			cache.Put(links[i][0], links[i][1])
		}
	}

	cache.Put("g", "url7")

	cachedEntries := map[string]string{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = fmt.Sprintf("%s", v.value)
	}

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"d": "url4",
		"e": "url5",
		"f": "url6",
		"g": "url7",
	}, cachedEntries)

	keys := []string{"g", "b", "c", "d", "e", "f"}
	e := cache.frequencies.Front()

	for i := 0; i < len(keys); i++ {
		assert.Equal(t, []string{keys[i]}, e.Value.(*FrequencyEntry).keys)

		e = e.Next()
	}

	require.Equal(t, 6, cache.frequencies.Len())
}

func TestRememberReplaceMultiple(t *testing.T) {
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

	cachedEntries := map[string]string{}

	for k, v := range cache.cachedEntries {
		cachedEntries[k] = fmt.Sprintf("%s", v.value)
	}

	assert.Equal(t, map[string]string{
		"b": "url2",
		"c": "url3",
		"f": "url6",
	}, cachedEntries)

	assert.Equal(t, []string{"f"}, cache.frequencies.Front().Value.(*FrequencyEntry).keys)
	assert.Equal(t, []string{"b", "c"}, cache.frequencies.Back().Value.(*FrequencyEntry).keys)
	require.Equal(t, 2, cache.frequencies.Len())
}
