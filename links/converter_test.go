package links

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumberToKey(t *testing.T) {
	tests := []struct {
		name        string
		number      int
		expectedKey string
	}{
		{"Zero", 0, ""},
		{"#1", 1, "1"},
		{"#2", 9, "9"},
		{"#3", 10, "a"},
		{"#4", 11, "b"},
		{"#5", 35, "z"},
		{"#6", 36, "10"},
		{"#7", 37, "11"},
		{"#8", 46, "1a"},
		{"#9", 71, "1z"},
		{"#10", 72, "20"},
		{"#11", 73, "21"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := numberToKey(tt.number)

			assert.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestKeyToNumber(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedNumber int
	}{
		{"Zero", "", 0},
		{"#1", "1", 1},
		{"#2", "9", 9},
		{"#3", "a", 10},
		{"#4", "b", 11},
		{"#5", "z", 35},
		{"#6", "10", 36},
		{"#7", "11", 37},
		{"#8", "1a", 46},
		{"#9", "1z", 71},
		{"#10", "20", 72},
		{"#11", "21", 73},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number := keyToNumber(tt.key)

			assert.Equal(t, tt.expectedNumber, number)
		})
	}
}
