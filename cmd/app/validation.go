package app

import (
	"errors"
	"net/url"
)

type Validator struct {
	KeyMaxLength   int
	allowedLetters string
}

func NewValidator(allowedLetters string) *Validator {
	return &Validator{
		allowedLetters: allowedLetters,
	}
}

func (v *Validator) validateKey(key string) error {
	if key == "" {
		return errors.New("key must be at least 1 letter long")
	}

	if len(key) > 10 {
		return errors.New("key is invalid")
	}

	for _, keyLetter := range key {
		found := false

		for _, generatorLetter := range v.allowedLetters {
			if keyLetter == generatorLetter {
				found = true

				break
			}
		}

		if !found {
			return errors.New("invalid letter")
		}
	}

	return nil
}

func (v *Validator) validateKeys(keys []string) error {
	for _, key := range keys {
		if err := v.validateKey(key); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateURL(URL string) error {
	if len(URL) > 2000 {
		return errors.New("URL must be maximum 2000 letters long")
	}

	if URL == "" {
		return errors.New("URL must be a valid URL string")
	}

	parsedURL, err := url.Parse(URL)

	if err != nil {
		return errors.New("URL must be a valid URL string")
	} else if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("URL must be an absolute URL")
	} else if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must begin with http or https")
	}

	return nil
}

func (v *Validator) validateURLs(URLs []string) error {
	uniqueURLs := make(map[string]bool, len(URLs))

	for _, URL := range URLs {
		err := v.validateURL(URL)

		if err != nil {
			return err
		}

		if uniqueURLs[URL] {
			return errors.New("URLs must be non-repeated")
		}

		uniqueURLs[URL] = true
	}

	return nil
}
