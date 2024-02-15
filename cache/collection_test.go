package cache

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCollection struct {
	//
}

func (c *testCollection) GenerateKey(URL string) (string, error) {
	return "key", nil
}

func (c *testCollection) GenerateKeys(URLs []string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (c *testCollection) GetURL(key string) (string, error) {
	return "url", nil
}

func (c *testCollection) GetURLs(keys []string) (map[string]string, error) {
	return map[string]string{}, nil
}

type testCache struct {
	data map[string]string
}

func (c *testCache) Get(key string) (interface{}, bool, error) {
	v, ok := c.data[key]

	return v, ok, nil
}

func (c *testCache) Put(key string, URL interface{}) error {
	c.data[key] = fmt.Sprintf("%s", URL)
	return nil
}

func TestGetURL(t *testing.T) {
	c := NewCachedCollection(
		&testCollection{},
		&testCache{
			data: map[string]string{"a": "url"},
		},
	)

	url, err := c.GetURL("a")

	require.NoError(t, err)
	require.Equal(t, "url", url)
}

func TestGetURLRemember(t *testing.T) {
	cache := &testCache{
		data: map[string]string{},
	}
	c := NewCachedCollection(
		&testCollection{},
		cache,
	)

	url, err := c.GetURL("a")

	require.NoError(t, err)
	require.Equal(t, "url", url)
	require.Equal(t, map[string]string{
		"a": "url",
	}, cache.data)
}
