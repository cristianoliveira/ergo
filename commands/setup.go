package commands

import (
	"fmt"
	"runtime"

	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

// Setup command tries set ergo as the proxy on networking settings.
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
func Setup(system string, remove bool, config *proxy.Config) {

	fmt.Println("Current detected system: " + runtime.GOOS)

	proxyURL := config.GetProxyPacURL()

	c := setup.GetConfigurator(system)

	if c == nil {
		fmt.Println(`
List of supported systems:

-linux-gnome
-osx
-windows

For more support please open an issue on: https://github.com/cristianoliveira/ergo
`)

		return
	}

	var err error
	if remove {
		err = c.SetDown()
	} else {
		err = c.SetUp(proxyURL)
	}

	if err != nil {
		fmt.Printf(`Setup failed cause %s`, err)
	}
}
