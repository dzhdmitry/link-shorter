package links

import (
	"math"
	"slices"
	"strings"
)

var Letters = "0123456789abcdefghijklmnopqrstuvwxyz"
var alphabet = strings.Split(Letters, "")

func numberToKey(number int) string {
	var digits []int

	for {
		if number == 0 {
			break
		}

		reminder := number % len(alphabet)
		number = number / len(alphabet)
		digits = append(digits, reminder)
	}

	key := make([]string, len(digits))

	for i, digit := range digits {
		key[len(digits)-i-1] = alphabet[digit]
	}

	return strings.Join(key, "")
}

func pow(x float64, y int) int {
	return int(math.Pow(x, float64(y)))
}

func keyToNumber(key string) int {
	number := 0
	alphabetLen := float64(len(alphabet))

	for i, letter := range strings.Split(key, "") {
		index := slices.Index(alphabet, letter)
		number += index * pow(alphabetLen, len(key)-i-1)
	}

	return number
}
