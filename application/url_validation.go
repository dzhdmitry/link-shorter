package application

import (
	"errors"
	"net/url"
)

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
