package proxy

import (
	"errors"
	"fmt"
	"net/url"
)

// Service holds the details of the service (Name and URL)
type Service struct {
	Name string
	URL  *url.URL
}

func NewService(name string, rawURL string) (Service, error) {
	if name == "" || rawURL == "" {
		return Service{}, errors.New("Name and URL are required")
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		return Service{}, fmt.Errorf("URL '%v' is invalid, example of valid URL 'http://example.com:8080'", rawURL)
	}

	return Service{Name: name, URL: url}, nil
}

// Empty service means no name or no url
func (s Service) Empty() bool {
	return s.Name == "" || s.URL == nil
}

func UnsafeNewService(name string, rawURL string) Service {
	url, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Invalid URL, example of a valid format is http://example.com:8080")
	}
	return Service{Name: name, URL: url}
}
