package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewErgoProxy(config *Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		service := config.GetService(req.URL.Host)
		if service != nil {
			fmt.Println("request", req.URL)
			target, _ := url.Parse(service.Url)
			targetQuery := target.RawQuery
			fmt.Println(config)
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}

			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "")
			}
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

type Service struct {
	Name string
	Url  string
}

type Config struct {
	UrlPattern string
	Services   []Service
}

func (c *Config) GetService(host string) *Service {
	isDev := regexp.MustCompile(`.*\.dev$`)
	if !isDev.MatchString(host) {
		return nil
	}

	name := strings.Split(host, ".")[0]
	for _, s := range c.Services {
		if s.Name == name {
			return &s
		}
	}

	return nil
}

func loadConfig() *Config {
	file, e := os.Open("./.ergo")
	defer file.Close()
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	services := []Service{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		config := strings.Split(line, "=")
		name, url := config[0], config[1]
		services = append(services, Service{Name: name, Url: url})
	}

	return &Config{
		UrlPattern: `.*\.dev$`,
		Services:   services,
	}

}

func main() {
	config := loadConfig()
	proxy := NewErgoProxy(config)
	http.Handle("/", proxy)
	http.ListenAndServe(":8080", nil)
}
