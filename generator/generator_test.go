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
		{"Under limit", "zzzzzzzzzzzy", "zzzzzzzzzzzz"},
	}
	g := NewGenerator(12)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := g.Generate(tt.lastKey)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestGenerateMultipleTimes(t *testing.T) {
	g := NewGenerator(12)
	key := ""
	keys := []string{}

	for i := 0; i < 500; i++ {
		key, _ = g.Generate(key)
		keys = append(keys, key)
	}

	require.Equal(t, "0", keys[0])
	require.Equal(t, "1", keys[1])
	require.Equal(t, "z", keys[35])
	require.Equal(t, "00", keys[36])
	require.Equal(t, "01", keys[37])
	require.Equal(t, "02", keys[38])
	require.Equal(t, "0z", keys[71])
	require.Equal(t, "10", keys[72])
	require.Equal(t, "11", keys[73])
	require.Equal(t, "12", keys[74])
	require.Equal(t, "1z", keys[107])
	require.Equal(t, "20", keys[108])
	require.Equal(t, "21", keys[109])
	require.Equal(t, "22", keys[110])
	require.Equal(t, "2z", keys[143])
	require.Equal(t, "30", keys[144])
}

func TestGenerateLimit(t *testing.T) {
	g := NewGenerator(12)
	key, err := g.Generate("zzzzzzzzzzzz")

	assert.Error(t, err)
	assert.Empty(t, key)
}
