//go:build windows

package setup

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

//ShowInternetOptions will display the Internet Options
//control panel. If the user clicks OK here, then the modifications
//will be visible without a restart
//this is just a fallback for the InetRefresh function
func ShowInternetOptions() {
	cmd := exec.Command("control", "inetcpl.cpl,Connections,4")
	fmt.Println("Starting the Internet Options control panel")
	err := cmd.Run()
	if err != nil {
		log.Printf("Internet Options control panel could not be started.\r\n%s\r\n", err.Error())
	}
}

//InetRefresh will inform windows that a proxy change was performed.
//if windows is not informed about the change, then the modification will only be noticed
//after a restart
func InetRefresh() {
	wininet := syscall.MustLoadDLL("wininet.dll")
	inetsetoption := wininet.MustFindProc("InternetSetOptionW")
	//Option 39 is INTERNET_OPTION_SETTINGS_CHANGED
	ret, _, callErr := inetsetoption.Call(0, 39, 0, 0)

	if ret == 0 {
		log.Printf(`
			Could not auto refresh the proxy settings.
			Step 1
			The received error is: %s
			Please press OK on the next dialog.\r\n`,
			callErr.Error())

		ShowInternetOptions()
	}
	//Option 37 is INTERNET_OPTION_REFRESH
	ret, _, callErr = inetsetoption.Call(0, 37, 0, 0)
	if ret == 0 {
		log.Printf(`
			Could not auto refresh the proxy settings.
			Step 2
			The received error is: %s
			Please press OK on the next dialog.\r\n`,
			callErr.Error())
		ShowInternetOptions()
	}

	return
}
