package generator

import (
	"errors"
	"strings"
)

var ErrLimitReached = errors.New("limit reached")

type alphabet struct {
	letters []string
}

func newAlphabet() *alphabet {
	letters := strings.Split("0123456789abcdefghijklmnopqrstuvwxyz", "")

	return &alphabet{letters: letters}
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

func Generate(lastKey string, keyMaxLength int) (string, error) {
	ab := newAlphabet()

	if lastKey == "" {
		return ab.firstLetter(), nil
	}

	if lastKey == strings.Repeat(ab.lastLetter(), keyMaxLength) {
		return "", ErrLimitReached
	}

	lastKeySplit := strings.Split(lastKey, "")

	for i := len(lastKeySplit) - 1; i >= 0; i-- {
		letter := lastKeySplit[i]
		letterIndex := ab.indexOfLetter(letter)

		// if letter is "z" (last)
		if letter == ab.lastLetter() {
			continue
		}

		// letter is not last
		newLetter := ab.letterOfIndex(letterIndex + 1)
		key := strings.Join(lastKeySplit[:i], "") + newLetter + strings.Join(lastKeySplit[i+1:], "")

		return key, nil
	}

	// all letters are "z" (last)
	key := strings.Repeat(ab.firstLetter(), len(lastKeySplit)+1)

	return key, nil
}
