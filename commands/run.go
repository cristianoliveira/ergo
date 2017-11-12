package commands

import (
	"fmt"
	"strings"

	"github.com/cristianoliveira/ergo/proxy"
)

// RunCommand starts the ergo proxy server.
//
// Usage:
// `ergo run`
type RunCommand struct{}

// Execute apply the RunCommand
func (c RunCommand) Execute(config *proxy.Config) (string, error) {
	if !strings.HasPrefix(config.Domain, ".") {
		return "", fmt.Errorf("Domain has a wrong format")
	}

	fmt.Println("Ergo Proxy listening on port " + config.Port + " for domains " + config.Domain)

	err := proxy.ServeProxy(config)
	if err != nil {
		return "", err
	}

	return "Finishied", nil
}
