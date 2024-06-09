package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// NewErgoProxy returns the new reverse proxy.
func NewErgoProxy(config *Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		t := time.Now()
		fmt.Printf(
			"[%s] %v %v %v \n",
			t.Format("2006-01-02 15:04:05"),
			req.Method,
			req.URL,
			req.Proto,
		)

		if config.Verbose {
			fmt.Println(formatRequest(req))
		}

		service, err := config.GetService(req.Host)
		if err != nil {
			fmt.Printf("Error getting service: %v", err)
		}

		if service != nil {
			target := service.URL
			targetQuery := target.RawQuery

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

	return &httputil.ReverseProxy{
		Director: director,
	}
}

// ServeProxy listens & serves the HTTP proxy.
func ServeProxy(config *Config) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	go config.WatchConfigFile(ticker.C)

	http.HandleFunc("/proxy.pac", proxy(config))

	http.HandleFunc("/__ergo__/", list(config))

	http.Handle("/", NewErgoProxy(config))

	return http.ListenAndServe(":"+config.Port, nil)
}

func proxy(config *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content := `
		function FindProxyForURL (url, host) {
			if (dnsDomainIs(host, '` + config.Domain + `')) {
				return 'PROXY 127.0.0.1:` + config.Port + `';
			}

			return 'DIRECT';
		}
		`
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		w.Write([]byte(content))
	}
}

func list(config *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, s := range config.Services {
			localURL := `http://` + s.Name + config.Domain
			fmt.Fprintf(w, "- %s -> %s \n", localURL, s.URL.String())
		}
	}
}

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

func formatRequest(r *http.Request) string {
	var request []string
	request = append(request, fmt.Sprintf("Host: %v", r.Host))

	for name, headers := range r.Header {
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	return strings.Join(request, "\n")
}
