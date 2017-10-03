package commands

import (
	"fmt"
	"regexp"

	"github.com/cristianoliveira/ergo/proxy"
)

// URL command find and print the url for a given app name.
// Usage:
// `ergo url foo`
func URL(name string, config *proxy.Config) {
	for _, s := range config.Services {

		if name == s.Name {

			localURL := s.Name + config.Domain

			//a protocol legth between 3 and 8 should fit all protocols
			valid := regexp.MustCompile("^\\w{3,8}\\:\\/\\/.*$")

			if !valid.MatchString(s.Name) {
				localURL = `http://` + localURL
			}

			fmt.Println(localURL)
		}
	}
}
