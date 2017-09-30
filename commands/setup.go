package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func Setup(system string, remove bool, config *proxy.Config) {
	fmt.Println("Current detected system: " + runtime.GOOS)
	proxyURL := "http://127.0.0.1:" + config.Port + "/proxy.pac"
	script := ""
	cmd := exec.Command("/bin/sh")

	switch system {
	case "linux-gnome":
		if remove {
			script = `
				gsettings set org.gnome.system.proxy mode 'none'
				gsettings set org.gnome.system.proxy autoconfig-url ''`

		} else {
			script = `
				gsettings set org.gnome.system.proxy mode 'auto'
				gsettings set org.gnome.system.proxy autoconfig-url '` + proxyURL + `'`

			fmt.Println(`To configure the proxy on your terminal execute:
			export http_proxy=http://127.0.0.1:` + config.Port)
		}

		cmd.Stdin = strings.NewReader(script)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

	case "osx":
		if remove {
			script = `sudo networksetup -setautoproxyurl "Wi-Fi" ""`
		} else {
			script = `sudo networksetup -setautoproxyurl "Wi-Fi" "` + proxyURL + `"`

			fmt.Println(`To configure the proxy on your terminal execute:
			export http_proxy=http://127.0.0.1:` + config.Port)
		}

		cmd.Stdin = strings.NewReader(script)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println(`
List of supported system

-linux-gnome
-osx

For more support please open an issue on: https://github.com/cristianoliveira/ergo
		`)
	}

	fmt.Println("Ergo proxy setup executed.")
}
