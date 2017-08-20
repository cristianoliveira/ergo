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
	config := LoadConfig("../.ergo")
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
}
