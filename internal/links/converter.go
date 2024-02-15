package links

import (
	"math"
	"slices"
	"strings"
)

var Letters = "0123456789abcdefghijklmnopqrstuvwxyz"
var alphabet = strings.Split(Letters, "")

func convertNumberToKey(number int64) string {
	var digits []int64

	for {
		if number == 0 {
			break
		}

		reminder := number % int64(len(alphabet))
		number = number / int64(len(alphabet))
		digits = append(digits, reminder)
	}

	key := make([]string, len(digits))

	for i, digit := range digits {
		key[len(digits)-i-1] = alphabet[digit]
	}

	return strings.Join(key, "")
}

func pow(x float64, y int) int64 {
	return int64(math.Pow(x, float64(y)))
}

func convertKeyToNumber(key string) int64 {
	number := int64(0)
	alphabetLen := float64(len(alphabet))

	for i, letter := range strings.Split(key, "") {
		index := int64(slices.Index(alphabet, letter))
		number += index * pow(alphabetLen, len(key)-i-1)
	}

	return number
}
