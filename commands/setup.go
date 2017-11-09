package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

var ergoCmd = exec.Command

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
	proxyURL := "http://127.0.0.1:" + config.Port + "/proxy.pac"

	c := setup.GetConfigurator(system)

	if c != nil {

		if remove {
			c.SetDown()
		} else {
			c.SetUp(proxyURL)
		}

		fmt.Println("Ergo proxy setup executed.")

	} else {
		fmt.Println(`
List of supported systems:

-linux-gnome
-osx
-windows

For more support please open an issue on: https://github.com/cristianoliveira/ergo
		`)
	}
}
