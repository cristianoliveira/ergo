package commands

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cristianoliveira/ergo/proxy"
)

// Setup command tries set ergo as the proxy on networking settings.
// For now, this feature is only supported for:
//   - OSX
//   - Gnome (tested on Linux and FreeBSD)
//   - Windows
//
// Usage:
// `ergo setup osx`
func Setup(system string, remove bool, config *proxy.Config) {

	fmt.Println("Current detected system: " + runtime.GOOS)
	proxyURL := "http://127.0.0.1:" + config.Port + "/proxy.pac"
	script := ""
	cmd := exec.Command("/bin/sh")

	switch system {
	case "gnome":
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

	case "windows":
		if remove {
			cmd = exec.Command("reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "AutoConfigURL", "/f")
		} else {
			cmd = exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "AutoConfigURL", "/t", "REG_SZ", "/d", proxyURL, "/f")
		}
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		//this is because windows needs to be told about the change
		InetRefresh()
	default:
		fmt.Println(`
List of supported system

-gnome (tested on linux and freebsd)
-osx
-windows

For more support please open an issue on: https://github.com/cristianoliveira/ergo
		`)
	}

	fmt.Println("Ergo proxy setup executed.")
}
