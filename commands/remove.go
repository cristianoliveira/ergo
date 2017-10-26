package commands

import (
	"fmt"

	"github.com/cristianoliveira/ergo/proxy"
)

// RemoveService removes a service from the configuration and tells the proxy to
// remove it from the configFile
func RemoveService(config *proxy.Config, service proxy.Service, configFile string) {

	var oldService *proxy.Service
	for _, srv := range config.Services {
		if srv.URL == service.URL || srv.Name == service.Name {
			oldService = &srv
		}
	}

	if oldService == nil {
		fmt.Printf("Service %s not found\n", service.Name)
		return
	}

	services := []proxy.Service{}
	for _, srv := range config.Services {
		if srv.Name != service.Name || srv.URL == service.URL {
			services = append(services, srv)
		}
	}

	if len(services) > len(config.Services) {
		fmt.Println("Failed to remove service %s", service.Name)
		return
	}

	config.Services = services
	err := proxy.RemoveService(configFile, oldService)
	if err != nil {
		fmt.Println("Failed to remove service %s", service.Name)
		return
	}

	fmt.Println("Service Removed")
}
