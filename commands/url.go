package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
)

// URL command find and print the url for a given app name.
// Usage:
// `ergo url foo`
func URL(name string, config *proxy.Config) {
	for _, s := range config.Services {
		if name == s.Name {
			localUrl := `http://` + s.Name + config.Domain
			fmt.Println(localUrl)
		}
	}
}
