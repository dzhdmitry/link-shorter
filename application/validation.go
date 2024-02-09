package application

import (
	"errors"
	"net/url"
)

func validateKey(key string, maxLength int) error {
	if key == "" {
		return errors.New("key must be at least 1 letter long")
	}

	if len(key) > maxLength {
		return errors.New("key is invalid")
	}

	return nil
}

func validateKeys(keys []string, maxLength int) error {
	for _, key := range keys {
		if err := validateKey(key, maxLength); err != nil {
			return err
		}
	}

	return nil
}

func validateURL(URL string) error {
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

func validateURLs(URLs []string) error {
	uniqueURLs := map[string]bool{}

	for _, URL := range URLs {
		err := validateURL(URL)

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
