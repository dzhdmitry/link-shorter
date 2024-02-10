package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		lastKey  string
		expected string
	}{
		{"Empty string", "", "0"},
		{"First non-empty", "0", "1"},
		{"First border between numbers and letters", "9", "a"},
		{"First border", "z", "00"},
		{"size=2 #1", "01", "02"},
		{"size=2 #2", "99", "9a"},
		{"Second border", "zz", "000"},
		{"Random #1", "ci2u1", "ci2u2"},
		{"Random #2", "2v15ia120eu9vaf5du8i", "2v15ia120eu9vaf5du8j"},
		{"Under limit", "yzzzzzzzzzzz", "zzzzzzzzzzzz"},
	}
	g := NewGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := g.Generate(tt.lastKey, 12)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestGenerateLimit(t *testing.T) {
	g := NewGenerator()
	key, err := g.Generate("zzzzzzzzzzzz", 12)

	assert.Error(t, err)
	assert.Empty(t, key)
}
