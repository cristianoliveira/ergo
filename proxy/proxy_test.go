package proxy

import (
	// "io/ioutil"
	"net/http"
	// "net/http/httptest"
	// "strings"
	"net/url"
	"testing"
)

func TestWhenHasCollectionFile(t *testing.T) {
	config := NewConfig()
	services, err := LoadServices("../.ergo")
	if err != nil {
		t.Fatal("could not load requied configuration file for tests")
	}
	config.Services = services
	proxy := NewErgoProxy(config)

	t.Run("it redirects foo.dev to localhost 3000", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://foo.dev/", nil)
		if err != nil {
			t.Fatal(err)
		}

		proxy.Director(req)

		expected, _ := url.Parse("http://localhost:3000/")
		result := req.URL

		if expected.Host != result.Host {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Scheme != result.Scheme {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Path != result.Path {
			t.Errorf("Expected %s got %s", expected, result)
		}
	})

	t.Run("it redirects bla.dev to localhost 5000", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://bla.dev/", nil)
		if err != nil {
			t.Fatal(err)
		}

		proxy.Director(req)

		expected, _ := url.Parse("http://localhost:5000/")
		result := req.URL

		if expected.Host != result.Host {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Scheme != result.Scheme {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Path != result.Path {
			t.Errorf("Expected %s got %s", expected, result)
		}
	})

	t.Run("it doens't redirects others", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://google.com/", nil)
		if err != nil {
			t.Fatal(err)
		}

		proxy.Director(req)

		expected, _ := url.Parse("http://google.com/")
		result := req.URL

		if expected.Host != result.Host {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Scheme != result.Scheme {
			t.Errorf("Expected %s got %s", expected, result)
		}
		if expected.Path != result.Path {
			t.Errorf("Expected %s got %s", expected, result)
		}
	})

	t.Run("when subdomains", func(tt *testing.T) {
		tt.Run("it redirects one.domain to localhost 8081", func(_ *testing.T) {
			req, err := http.NewRequest("GET", "http://one.domain.dev/", nil)
			if err != nil {
				t.Fatal(err)
			}

			proxy.Director(req)

			expected, _ := url.Parse("http://localhost:8081/")
			result := req.URL

			if expected.Host != result.Host {
				tt.Errorf("Expected %s got %s", expected, result)
			}
			if expected.Scheme != result.Scheme {
				tt.Errorf("Expected %s got %s", expected, result)
			}
			if expected.Path != result.Path {
				tt.Errorf("Expected %s got %s", expected, result)
			}
		})

		tt.Run("it redirects two.domain to localhost 8082", func(_ *testing.T) {
			req, err := http.NewRequest("GET", "http://two.domain.dev/", nil)
			if err != nil {
				t.Fatal(err)
			}

			proxy.Director(req)

			expected, _ := url.Parse("http://localhost:8082/")
			result := req.URL

			if expected.Host != result.Host {
				tt.Errorf("Expected %s got %s", expected, result)
			}
			if expected.Scheme != result.Scheme {
				tt.Errorf("Expected %s got %s", expected, result)
			}
			if expected.Path != result.Path {
				tt.Errorf("Expected %s got %s", expected, result)
			}
		})
	})
}
