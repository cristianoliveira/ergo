package commands

import (
	"fmt"

	"github.com/cristianoliveira/ergo/proxy"
)

// AddService Allows to add new services to existing configuration file
// Usage:
// `ergo add service servicehost:port`
func AddService(config *proxy.Config, service proxy.Service, configFile string) {
	oldService := config.GetService(service.Name + config.Domain)
	if oldService != nil {
		fmt.Println("Service already present!")
	} else {
		config.Services = append(config.Services, service)
		err := proxy.AddService(configFile, service)
		if err == nil {
			fmt.Println("Service added successfully!")
		} else {
			fmt.Println("Error while adding new service!")
		}
	}
}
