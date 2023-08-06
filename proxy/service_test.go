package proxy;

import (
	"testing"
)

func TestNewService(t *testing.T) {
	testCases := []struct {
		name    string;

		serviceName string;
		serviceURL  string;

		expectError bool
	}{
		{
			name:    "empty name",
			serviceName: "",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},
		{
			name:    "empty url",
			serviceName: "test",
			serviceURL:  "",
			expectError: true,
		},
		{
			name:    "a valid service with name and an valid URL",
			serviceName: "test",
			serviceURL:  "http://localhost:8080",
			expectError: false,
		},
		{
			name:    "invalid url - double port",
			serviceName: "test",
			serviceURL:  "http://localhost:8080:8888",
			expectError: true,
		},
		{
			name:    "invalid url - missing scheme",
			serviceName: "test",
			serviceURL:  "localhost:8080",
			expectError: true,
		},
		{
			name:    "invalid name - name contains an URL",
			serviceName: "http://localhost:3333",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},
		{
			name:    "invalid name - name contains an URL port",
			serviceName: "localhost:3333",
			serviceURL:  "http://localhost:8080",
			expectError: true,
		},

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewService(tc.serviceName, tc.serviceURL)
			if (err != nil) != tc.expectError {
				t.Errorf("NewService() error = %v, expectError %v", err, tc.expectError)
				return
			}
		})
	}
}
