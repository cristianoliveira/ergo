package commands

import (
	"fmt"

	"github.com/cristianoliveira/ergo/proxy"
)

// RemoveServiceCommand removes a service from the configuration
// and tells the proxy to remove it from the config file.
// USAGE:
// ergo remove myservicename
type RemoveServiceCommand struct {
	SearchTerm string
}

func findService(searchTerm string, services map[string]proxy.Service) (*proxy.Service, bool) {
	for _, srv := range services {
		if srv.URL.String() == searchTerm || srv.Name == searchTerm {
			return &srv, true
		}
	}

	return nil, false
}

// Execute apply the RemoveServiceCommand
func (c RemoveServiceCommand) Execute(config *proxy.Config) (string, error) {
	srv, isPresent := findService(c.SearchTerm, config.Services)

	if !isPresent {
		return "", fmt.Errorf("Service %s not found", c.SearchTerm)
	}

	err := proxy.RemoveService(config.ConfigFile, *srv)
	if err != nil {
		return "", fmt.Errorf("Failed to remove service cause %s", err)
	}

	return "Service Removed", nil
}
