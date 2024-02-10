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
	alphabet *alphabet
}

func NewGenerator() *Generator {
	return &Generator{alphabet: &alphabet{
		letters: strings.Split(Letters, ""),
	}}
}

func (g *Generator) Generate(lastKey string, keyMaxLength int) (string, error) {
	if lastKey == "" {
		return g.alphabet.firstLetter(), nil
	}

	if lastKey == strings.Repeat(g.alphabet.lastLetter(), keyMaxLength) {
		return "", ErrLimitReached
	}

	lastKeySplit := strings.Split(lastKey, "")

	for i := len(lastKeySplit) - 1; i >= 0; i-- {
		letter := lastKeySplit[i]
		letterIndex := g.alphabet.indexOfLetter(letter)

		// if letter is "z" (last)
		if letter == g.alphabet.lastLetter() {
			continue
		}

		// letter is not last
		newLetter := g.alphabet.letterOfIndex(letterIndex + 1)
		key := strings.Join(lastKeySplit[:i], "") + newLetter + strings.Join(lastKeySplit[i+1:], "")

		return key, nil
	}

	// all letters are "z" (last)
	key := strings.Repeat(g.alphabet.firstLetter(), len(lastKeySplit)+1)

	return key, nil
}
