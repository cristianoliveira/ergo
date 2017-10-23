package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

//this should be called ( quit.Stop() )
//when the configuration watcher should stop
var quit chan struct{}

const pollIntervall = 500

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

func pollConfigChange(config *Config) {
	ticker := time.NewTicker(pollIntervall * time.Millisecond)

	quit = make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				info, err := os.Stat(config.ConfigFile)
				if err != nil {
					log.Printf("Error reading config file: %s\r\n", err.Error())
					continue
				}

				if info.ModTime().Before(modTime) || info.Size() != size {
					services, err := LoadServices(config.ConfigFile)
					if err != nil {
						log.Printf("Error reading the modified config file: %s\r\n", err.Error())
						continue
					}
					//clear the data if there is any
					select {
					case <-configChan:
					default:
					}
					configChan <- services
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

//NewErgoProxy returns the new reverse proxy.
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

		service := config.GetService(req.URL.Host)
		if service != nil {
			target, _ := url.Parse(service.URL)
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
			fmt.Fprintf(w, "- %s -> %s \n", localURL, s.URL)
		}
	}
}

//ServeProxy listens & serves the HTTP proxy.
func ServeProxy(config *Config) error {

	pollConfigChange(config)

	http.HandleFunc("/proxy.pac", proxy(config))

	http.HandleFunc("/_ergo/list", list(config))

	http.Handle("/", NewErgoProxy(config))

	return http.ListenAndServe(":"+config.Port, nil)
}
