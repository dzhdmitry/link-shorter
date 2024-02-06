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
		{
			name:     "Empty string",
			lastKey:  "",
			expected: "0",
		},
		{
			name:     "First non-empty",
			lastKey:  "0",
			expected: "1",
		},
		{
			name:     "First border between numbers and letters",
			lastKey:  "9",
			expected: "a",
		},
		{
			name:     "First border",
			lastKey:  "z",
			expected: "00",
		},
		{
			name:     "size=2 #1",
			lastKey:  "01",
			expected: "02",
		},
		{
			name:     "size=2 #2",
			lastKey:  "99",
			expected: "9a",
		},
		{
			name:     "Second border",
			lastKey:  "zz",
			expected: "000",
		},
		{
			name:     "Random #1",
			lastKey:  "ci2u1",
			expected: "ci2u2",
		},
		{
			name:     "Random #2",
			lastKey:  "2v15ia120eu9vaf5du8i",
			expected: "2v15ia120eu9vaf5du8j",
		},
		{
			name:     "Under limit",
			lastKey:  "yzzzzzzzzzzz",
			expected: "zzzzzzzzzzzz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := Generate(tt.lastKey, 12)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestGenerateLimit(t *testing.T) {
	key, err := Generate("zzzzzzzzzzzz", 12)

	assert.Error(t, err)
	assert.Empty(t, key)
}
