package links_in_memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	lc := NewLinksCollection(5)
	expectedKeys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a"}

	for _, expectedKey := range expectedKeys {
		key, err := lc.GenerateKey("http://links.ru")

		require.NoError(t, err)
		assert.Equal(t, expectedKey, key)
	}
}

func TestGetLink(t *testing.T) {
	lc := NewLinksCollection(5)
	key, _ := lc.GenerateKey("http://example.com")

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"empty", "", ""},
		{"non-existed", "1d32g", ""},
		{"existed", key, "http://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := lc.GetLink(tt.key)

			require.Equal(t, tt.expected, link)
		})
	}
}
