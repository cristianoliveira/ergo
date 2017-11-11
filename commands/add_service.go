package commands

import (
	"errors"

	"github.com/cristianoliveira/ergo/proxy"
)

// AddServiceCommand Allows to add new services to existing configuration file
// Usage:
// `ergo add service servicehost:port`
type AddServiceCommand struct {
	Service proxy.Service
}

// Execute apply the AddServiceCommand
func (c AddServiceCommand) Execute(config *proxy.Config) (string, error) {
	oldService := config.GetService(c.Service.Name + config.Domain)
	if oldService != nil {
		return "", errors.New("Service already present")
	}

	config.Services = append(config.Services, c.Service)

	err := proxy.AddService(config.ConfigFile, c.Service)
	if err != nil {
		return "", err
	}

	return "Service added successfully", nil
}
