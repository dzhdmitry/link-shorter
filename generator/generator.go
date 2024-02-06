package generator

import (
	"errors"
	"strings"
)

var ErrLimitReached = errors.New("limit reached")

type Alphabet struct {
	letters []string
}

func NewAlphabet() *Alphabet {
	letters := strings.Split("0123456789abcdefghijklmnopqrstuvwxyz", "")

	return &Alphabet{letters: letters}
}

func (a *Alphabet) firstLetter() string {
	return a.letters[0]
}

func (a *Alphabet) lastLetter() string {
	return a.letters[len(a.letters)-1]
}

func (a *Alphabet) indexOfLetter(letter string) int {
	for k, v := range a.letters {
		if letter == v {
			return k
		}
	}

	return -1
}

func (a *Alphabet) letterOfIndex(index int) string {
	if index >= len(a.letters) {
		return ""
	}

	return a.letters[index]
}

func Generate(lastKey string, keyMaxLength int) (string, error) {
	alphabet := NewAlphabet()

	if lastKey == "" {
		return alphabet.firstLetter(), nil
	}

	if lastKey == strings.Repeat(alphabet.lastLetter(), keyMaxLength) {
		return "", ErrLimitReached
	}

	lastKeySplit := strings.Split(lastKey, "")

	for i := len(lastKeySplit) - 1; i >= 0; i-- {
		letter := lastKeySplit[i]
		letterIndex := alphabet.indexOfLetter(letter)

		// if letter is "z" (last)
		if letter == alphabet.lastLetter() {
			continue
		}

		// letter is not last
		newLetter := alphabet.letterOfIndex(letterIndex + 1)
		key := strings.Join(lastKeySplit[:i], "") + newLetter + strings.Join(lastKeySplit[i+1:], "")

		return key, nil
	}

	// all letters are "z" (last)
	key := strings.Repeat(alphabet.firstLetter(), len(lastKeySplit)+1)

	return key, nil
}
