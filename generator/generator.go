package generator

import (
	"errors"
	"strings"
)

var ErrLimitReached = errors.New("limit reached")
var Letters = "0123456789abcdefghijklmnopqrstuvwxyz"

type alphabet struct {
	letters []string
}

func (a *alphabet) firstLetter() string {
	return a.letters[0]
}

func (a *alphabet) lastLetter() string {
	return a.letters[len(a.letters)-1]
}

func (a *alphabet) indexOfLetter(letter string) int {
	for k, v := range a.letters {
		if letter == v {
			return k
		}
	}

	return -1
}

func (a *alphabet) letterOfIndex(index int) string {
	if index >= len(a.letters) {
		return ""
	}

	return a.letters[index]
}

type Generator struct {
	alphabet     *alphabet
	keyMaxLength int
}

func NewGenerator(keyMaxLength int) *Generator {
	return &Generator{
		alphabet: &alphabet{
			letters: strings.Split(Letters, ""),
		},
		keyMaxLength: keyMaxLength,
	}
}

func (g *Generator) Generate(lastKey string) (string, error) {
	if lastKey == "" {
		return g.alphabet.firstLetter(), nil
	}

	if lastKey == strings.Repeat(g.alphabet.lastLetter(), g.keyMaxLength) {
		return "", ErrLimitReached
	}

	lastKeySplit := strings.Split(lastKey, "")

	for i := len(lastKeySplit) - 1; i >= 0; i-- {
		if lastKeySplit[i] == g.alphabet.lastLetter() {
			lastKeySplit[i] = g.alphabet.firstLetter()
		} else {
			lastKeySplit[i] = g.alphabet.letterOfIndex(g.alphabet.indexOfLetter(lastKeySplit[i]) + 1)

			return strings.Join(lastKeySplit, ""), nil
		}
	}

	lastKeySplit = append([]string{"0"}, lastKeySplit...)

	return strings.Join(lastKeySplit, ""), nil
}

// GenerateMany returns slice of slices [["key1", "URL1"], ["key2", "URL2"], ...]
func (g *Generator) GenerateMany(lastKey string, URLs []string) ([][]string, error) {
	generatedKeysSorted := [][]string{}
	currentLastKey := lastKey

	for _, URL := range URLs {
		key, err := g.Generate(currentLastKey)

		if err != nil {
			return nil, err
		}

		generatedKeysSorted = append(generatedKeysSorted, []string{key, URL})
		currentLastKey = key
	}

	return generatedKeysSorted, nil
}
