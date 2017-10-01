package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
)

// ListNames command lists all configured apps names and its urls.
// Usage:
// `ergo list-names`
func ListNames(config *proxy.Config) {
	fmt.Println("Ergo Proxy current list: ")
	for _, s := range config.Services {
		fmt.Printf(" - %s -> %s \n", s.Name, s.Url)
	}
}
