package commands

import (
	"fmt"
	"runtime"

	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

// SetupCommand tries set ergo as the proxy on networking settings.
// For now, this feature is only supported for:
//
//   - OSX
//   - Linux-gnome
//   - Windows
//
// Usage:
//
// `ergo setup osx`
//
type SetupCommand struct {
	System string
	Remove bool
}

var usage = `
List of supported systems:

-linux-gnome
-osx
-windows

For more support please open an issue on: https://github.com/cristianoliveira/ergo
`

// Execute apply the SetupCommand
func (c SetupCommand) Execute(config *proxy.Config) (string, error) {
	fmt.Println("Current detected system: " + runtime.GOOS)

	proxyURL := config.GetProxyPacURL()

	configurator := setup.GetConfigurator(c.System)

	if configurator == nil {
		return "", fmt.Errorf(usage)
	}

	var err error
	if c.Remove {
		err = configurator.SetDown()
	} else {
		err = configurator.SetUp(proxyURL)
	}

	if err != nil {
		return "", fmt.Errorf("Setup failed cause %s", err)
	}

	return "Setup executed", nil
}
