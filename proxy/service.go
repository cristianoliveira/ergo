package proxy

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
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

	if strings.Contains(name, "://") || strings.Contains(name, ":") {
		return Service{}, fmt.Errorf("Name '%v' is invalid, it can't be an URL", name)
	}

	url, err := url.ParseRequestURI(rawURL)
	isInvalidHostname := len(url.Hostname()) == 0 || strings.Contains(url.Hostname(), ":")
	if err != nil || isInvalidHostname {
		return Service{}, fmt.Errorf("URL '%v' is invalid, example of valid URL 'http://example.com:8080'", rawURL)
	}

	return Service{Name: name, URL: url}, nil
}

// Empty service means no name or no url
func (s Service) Empty() bool {
	return s.Name == "" || s.URL == nil
}

// String returns a string representation of the service
// Example: 
//  Service{Name: "test", URL: "http://localhost:8080"}
//  produces the string "test http://localhost:8080"
func (s Service) String() string {
	return s.Name + " " + s.URL.String()
}

// GetOriginalURL returns the original URL of the service
func (s Service) GetOriginalURL() *url.URL {
	return s.URL
}

// GetServiceURL returns the local URL to be used
// by the proxy server to redirect all request to the service
func (s Service) GetServiceURL(localTLD string) string {
	return fmt.Sprintf("%s://%s%s", s.URL.Scheme, s.Name, localTLD)
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
