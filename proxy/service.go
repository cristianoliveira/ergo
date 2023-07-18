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

// NewService creates a new service from a name and a URL
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

// UnsafeNewService creates a new service from a name and a URL
// without checking if the URL is valid. Must only be used in tests
func UnsafeNewService(name string, rawURL string) Service {
	url, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Invalid URL %s", rawURL)
	}
	return Service{Name: name, URL: url}
}
