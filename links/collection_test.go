package links

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"link-shorter.dzhdmitry.net/generator"
	"testing"
)

type testStorage struct {
	//
}

func (t *testStorage) StoreKeysURLs(keysURLs [][]string) error {
	return nil
}

func (t *testStorage) Restore() error {
	return nil
}

func (t *testStorage) GetURL(key string) (string, error) {
	return "http://example.com", nil
}

func (t *testStorage) GetURLs(keys []string) (map[string]string, error) {
	return map[string]string{
		"key1": "http://example.com",
		"key2": "http://example.com",
	}, nil
}

func (t *testStorage) GetLastKey() (string, error) {
	return "", nil
}

func TestGenerateKey(t *testing.T) {
	collection := NewCollection(*generator.NewGenerator(5), &testStorage{})
	key, err := collection.GenerateKey("http://links.ru")

	require.NoError(t, err)
	assert.Equal(t, "0", key)
}

func TestGenerateKeys(t *testing.T) {
	tests := []struct {
		name     string
		key      []string
		expected map[string]string
	}{
		{"Empty", []string{}, map[string]string{}},
		{"Regular", []string{"https://example1", "https://example2"}, map[string]string{
			"https://example1": "0",
			"https://example2": "1",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := NewCollection(*generator.NewGenerator(5), &testStorage{})
			keys, err := collection.GenerateKeys(tt.key)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, keys)
		})
	}
}

func TestGetLink(t *testing.T) {
	collection := NewCollection(*generator.NewGenerator(5), &testStorage{})
	keys, err := collection.GetURL("2")

	require.NoError(t, err)
	assert.Equal(t, "http://example.com", keys)
}
