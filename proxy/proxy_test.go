package proxy

import (
	// "io/ioutil"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "net/http/httptest"
	// "strings"
	"net/url"
	"testing"
)

func TestWhenHasCollectionFile(t *testing.T) {
	config := NewConfig()
	config.ConfigFile = "../.ergo"
	err := config.LoadServices()
	if err != nil {
		t.Fatal("could not load requied configuration file for tests")
	}

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

func TestPollConfigChangeShouldNotFindFile(t *testing.T) {

	config := NewConfig()
	config.ConfigFile = ".notexistent"

	logbuf := new(bytes.Buffer)

	log.SetOutput(logbuf)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go config.WatchConfigFile(ticker.C)

	time.Sleep(1 * time.Second)

	if len(logbuf.String()) == 0 {
		t.Fatalf("Expected to get a read from the log. Got none")
	} else if !strings.Contains(logbuf.String(), "The system cannot find the file specified") &&
		!strings.Contains(logbuf.String(), "no such file or directory") {
		t.Fatalf("Expected the log to report a missing file. Got %s", logbuf.String())
	}
}

func TestPollConfigChangeWithInvalidConfigFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("No error expected while creating a temporary file. Got %s.", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("No error expected while writing initial config to a temporary file. Got %s.", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("No error expected while closing the temporary file. Got %s.", err.Error())
	}

	logbuf := new(bytes.Buffer)
	log.SetOutput(logbuf)

	config := NewConfig()
	config.ConfigFile = tmpfile.Name()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go config.WatchConfigFile(ticker.C)

	time.Sleep(2 * time.Second)

	err = ioutil.WriteFile(tmpfile.Name(), []byte("test.devlocalhost:9900"), 0644)

	if err != nil {
		t.Fatalf("Expected no error while rewriting the temporary config file. Got %s", err.Error())
	}

	time.Sleep(1 * time.Second)

	if len(logbuf.String()) == 0 {
		t.Fatalf("Expected to get a read from the log. Got none")
	} else if !strings.Contains(logbuf.String(), "Error reading the modified config file") {
		t.Fatalf("Expected the log to report an error reading an invalid config file. Got %s", logbuf.String())
	}
}

func TestPollConfigChangeWithValidConfigFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("No error expected while creating a temporary file. Got %s.", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("No error expected while writing initial config to a temporary file. Got %s.", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("No error expected while closing the temporary file. Got %s.", err.Error())
	}

	config := NewConfig()
	config.ConfigFile = tmpfile.Name()
	config.LoadServices()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go config.WatchConfigFile(ticker.C)

	time.Sleep(1 * time.Second)

	configFile, err := os.OpenFile(config.ConfigFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("Error while opening the config file %s", err.Error())
	}

	if _, err = configFile.WriteString("\ntest2.dev http://localhost:9900"); err != nil {
		t.Fatalf("Expected no error while rewriting the temporary config file. Got %s", err.Error())
	}

	configFile.Close()
	time.Sleep(1 * time.Second)

	if len(config.Services) != 2 {
		t.Fatalf("Expected to get 2 service Got %d", len(config.Services))
	}

	service := config.Services["test2.dev"]
	if service.URL.String() != "http://localhost:9900" {
		t.Fatalf("Expected to get 1 service with the URL http://localhost:9900 and the name test.dev. Got the URL: %s and the name: %s", service.URL, service.Name)
	}
}

// structure to mock a http ResponseWriter
type mockHTTPResponse struct {
	WrittenData   []byte
	WrittenHeader int
	MyHeader      http.Header
}

func (mr *mockHTTPResponse) Header() http.Header {
	return mr.MyHeader
}

func (mr *mockHTTPResponse) Write(data []byte) (int, error) {
	mr.WrittenData = data
	return len(data), nil
}

func (mr *mockHTTPResponse) WriteHeader(header int) {
	mr.WrittenHeader = header
}

func TestProxyFunction(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("No error expected while creating a temporary file. Got %s.", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("No error expected while writing initial config to a temporary file. Got %s.", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("No error expected while closing the temporary file. Got %s.", err.Error())
	}

	config := Config{}
	config.ConfigFile = tmpfile.Name()
	config.Domain = "dev"
	config.Port = PortDefault

	fncProxy := proxy(&config)
	m := &mockHTTPResponse{}
	m.MyHeader = make(http.Header)
	r := http.Request{}

	fncProxy(m, &r)

	if len(m.MyHeader) != 1 {
		t.Fatalf("Expected to get 1 header. Got %d", len(m.MyHeader))
	}
	if m.MyHeader["Content-Type"] == nil {
		t.Fatal("Expected to get a \"Content-Type\" header. Got none")
	}
	if m.MyHeader["Content-Type"][0] != "application/x-ns-proxy-autoconfig" {
		t.Fatalf("Expected to get a \"Content-Type\" header of \"application/x-ns-proxy-autoconfig\". Got %s", m.MyHeader["Content-Type"][0])
	}

	content := `
	function FindProxyForURL (url, host) {
		if (dnsDomainIs(host, '` + config.Domain + `')) {
			return 'PROXY 127.0.0.1:` + config.Port + `';
		}

		return 'DIRECT';
	}
	`

	content = strings.Replace(content, "\t", "", -1)

	response := strings.Replace(string(m.WrittenData), "\t", "", -1)

	if !strings.Contains(response, content) {
		t.Fatalf("Expected to get an response of %s. Got %s", content, response)
	}

}

func TestListFunction(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("No error expected while creating a temporary file. Got %s.", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("No error expected while writing initial config to a temporary file. Got %s.", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("No error expected while closing the temporary file. Got %s.", err.Error())
	}

	config := NewConfig()
	config.ConfigFile = tmpfile.Name()
	config.Domain = "dev"
	config.Port = PortDefault
	err = config.LoadServices()
	if err != nil {
		t.Fatalf("Expected no error while loading services from temp file. Got %s", err.Error())
	}

	fncList := list(config)
	m := &mockHTTPResponse{}
	m.MyHeader = make(http.Header)
	r := http.Request{}

	fncList(m, &r)

	content := fmt.Sprintf("- %s -> %s \n", "http:\\test.dev", "localhost:9999")

	if !strings.Contains(string(m.WrittenData), content) {
		t.Fatalf("Expected to get an response of %s. Got %s", content, m.WrittenData)
	}

}

func TestFormatRequest(t *testing.T) {

	r := http.Request{}
	r.Header = make(http.Header)
	r.Header.Add("test-header", "test-header-value")
	r.Host = "test-host"

	expected := `Host: test-host Test-Header: test-header-value`
	formated := formatRequest(&r)

	if strings.Replace(formated, "\n", " ", -1) != expected {
		t.Fatalf("Expected to get a formated string of \"%s\". Got \"%s\"", strings.Trim(expected, "\n"), strings.Trim(formated, "\n"))
	}
}

func TestSingleJoiningSlash(t *testing.T) {
	a1 := "path/"
	a2 := "path"
	b1 := "/anotherpath"
	b2 := "anotherpath"

	rez11 := singleJoiningSlash(a1, b1)
	rez12 := singleJoiningSlash(a1, b2)
	rez21 := singleJoiningSlash(a2, b1)
	rez22 := singleJoiningSlash(a2, b2)

	expected := "path/anotherpath"

	if rez11 != expected {
		t.Fatalf("Expected to get %s as the joining of %s and %s. Got %s", expected, a1, b1, rez11)
	}

	if rez12 != expected {
		t.Fatalf("Expected to get %s as the joining of %s and %s. Got %s", expected, a1, b2, rez12)
	}

	if rez21 != expected {
		t.Fatalf("Expected to get %s as the joining of %s and %s. Got %s", expected, a2, b1, rez21)
	}

	if rez22 != expected {
		t.Fatalf("Expected to get %s as the joining of %s and %s. Got %s", expected, a2, b2, rez22)
	}
}
