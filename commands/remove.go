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
	Service proxy.Service
}

// Execute apply the RemoveServiceCommand
func (c RemoveServiceCommand) Execute(config *proxy.Config) (string, error) {
	var oldService *proxy.Service
	for _, srv := range config.Services {
		if srv.URL == c.Service.URL || srv.Name == c.Service.Name {
			oldService = &srv
		}
	}

	if oldService == nil {
		return "", fmt.Errorf("Service %s not found", c.Service.Name)
	}

	services := []proxy.Service{}
	for _, srv := range config.Services {
		if srv.Name != c.Service.Name || srv.URL == c.Service.URL {
			services = append(services, srv)
		}
	}

	if len(services) > len(config.Services) {
		return "", fmt.Errorf("Failed to remove service %s", c.Service.Name)
	}

	config.Services = services
	err := proxy.RemoveService(config.ConfigFile, *oldService)
	if err != nil {
		return "", fmt.Errorf("Failed to remove service cause %s", err)
	}

	return "Service Removed", nil
}
