package setup

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// OSXConfigurator implements Configurator for windows
type OSXConfigurator struct{}

const SUPPORTED_OSX_VERSION = 10 // up to Catalina

func checkSupportedVersion() (bool, error) {
	cmd := `sw_vers -productVersion`
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return false, err
	}

	var majorVersionNumber int
	outputString := string(output)
	if strings.Contains(outputString, ".") {
		majorVersion := strings.Split(string(output), ".")[0]
		majorVersionNumber, err = strconv.Atoi(majorVersion)
	} else {
		majorVersionNumber, err = strconv.Atoi(outputString)
	}

	if err != nil {
		return false, err
	}

	if majorVersionNumber >= SUPPORTED_OSX_VERSION {
		fmt.Println("The ergo setup is not supported for the current osx version.")
		return false, nil
	}

	return true, nil
}

// SetUp is responsible for setting up the ergo as proxy
func (c *OSXConfigurator) SetUp(proxyURL string) error {
	isSupported, err := checkSupportedVersion()
	if err != nil {
		fmt.Println("Error while checking the current osx version: ", err)
		return err
	}

	if !isSupported {
		fmt.Println("The ergo setup is not supported for the current osx version.")
		return nil
	}

	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" "`+proxyURL+`"`)
}

// SetDown is responsible for remove the ergo as proxy
func (c *OSXConfigurator) SetDown() error {
	isSupported, err := checkSupportedVersion()
	if err != nil {
		fmt.Println("Error while checking the current osx version: ", err)
		return err
	}

	if !isSupported {
		fmt.Println("The ergo setup is not supported for the current osx version.")
		return nil
	}

	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" ""`)
}
