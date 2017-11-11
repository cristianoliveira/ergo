package commands

import (
	"errors"
	"regexp"

	"github.com/cristianoliveira/ergo/proxy"
)

// URLCommand find and print the url for a given app name.
// Usage:
// `ergo url foo`
type URLCommand struct {
	FilterName string
}

// Execute apply the URLCommand
func (c URLCommand) Execute(config *proxy.Config) (string, error) {
	for _, s := range config.Services {

		if c.FilterName == s.Name {

			localURL := s.Name + config.Domain

			//a protocol legth between 3 and 8 should fit all protocols
			valid := regexp.MustCompile("^\\w{3,8}\\:\\/\\/.*$")

			if !valid.MatchString(s.Name) {
				localURL = `http://` + localURL
			}

			return localURL, nil
		}
	}

	return "", errors.New("Url not found")
}
