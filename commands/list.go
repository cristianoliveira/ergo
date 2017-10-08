package commands

import (
	"fmt"

	"github.com/cristianoliveira/ergo/proxy"
)

// List command lists all configured apps local url and its original urls.
// Usage:
// `ergo list`
func List(config *proxy.Config) {
	fmt.Println("Ergo Proxy current list: ")
	for _, s := range config.Services {
		localURL := `http://` + s.Name + config.Domain
		fmt.Printf(" - %s -> %s \n", localURL, s.URL)
	}
}
