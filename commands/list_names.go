package commands

import (
	"fmt"

	"github.com/cristianoliveira/ergo/proxy"
)

// ListNameCommand lists all configured apps names and its urls.
// Usage:
// `ergo list-names`
type ListNameCommand struct{}

// Execute apply the ListNameCommand
func (c ListNameCommand) Execute(config *proxy.Config) (string, error) {
	output := "Ergo Proxy current list:\n"

	for _, s := range config.Services {
		output = fmt.Sprintf("%s - %s -> %s \n", output, s.Name, s.URL)
	}

	return output, nil
}
