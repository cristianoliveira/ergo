package proxy

import (
	"testing"
)

func TestNewService(t *testing.T) {
	testCases := []struct {
		title       string
		serviceName string
		serviceURL  string

		expectError bool
	}{
		{
			title:       "a service with name and url is valid",
			serviceName: "test",
			serviceURL:  "http://localhost:8080",
			expectError: false,
		},
		{
			title:       "a service with empty name is invalid",
			serviceName: "",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},
		{
			title:       "a service with empty url is invalid",
			serviceName: "test",
			serviceURL:  "",
			expectError: true,
		},
		{
			title:       "a service with doubled port url is invalid",
			serviceName: "test",
			serviceURL:  "http://localhost:8080:8888",
			expectError: true,
		},
		{
			title:       "a service with missing scheme in URL is invalid",
			serviceName: "test",
			serviceURL:  "localhost:8080",
			expectError: true,
		},
		{
			title:       "a service with name containing an URL is invalid",
			serviceName: "http://localhost:3333",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},
		{
			title:       "a service with name containing a port is invalid",
			serviceName: "localhost:3333",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},
		{
			title:       "a service containing an invalid URL is invalid",
			serviceName: "ergoproxy",
			serviceURL:  "http:///localhost:3000\n",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := NewService(tc.serviceName, tc.serviceURL)
			if (err != nil) != tc.expectError {
				t.Errorf("NewService() error = %v, expectError %v", err, tc.expectError)
				return
			}
		})
	}
}
