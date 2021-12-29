package utils

import (
	"net/http"
	"net/url"
)

func UnShortenURL(URL string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(URL)
	if err != nil && err != http.ErrUseLastResponse {
		return "", err
	}

	url, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
