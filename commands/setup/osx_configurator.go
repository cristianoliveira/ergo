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

func checkSupportedVersion() error {
	cmd := `sw_vers -productVersion`
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return err
	}

	var majorVersionNumber int
	outputString := string(output)
	if outputString == "" {
		return fmt.Errorf("checking the current osx version failed")
	}

	if strings.Contains(outputString, ".") {
		majorVersion := strings.Split(string(output), ".")[0]
		majorVersionNumber, err = strconv.Atoi(majorVersion)
	} else {
		majorVersionNumber, err = strconv.Atoi(outputString)
	}

	if err != nil {
		return err
	}

	if majorVersionNumber >= SUPPORTED_OSX_VERSION {
		fmt.Println("The ergo setup is not supported for the current osx version.")
		fmt.Println("Supported versions Catalina or below.")
		fmt.Println("Please, consider setting up ergo as proxy manually.")
		return fmt.Errorf("unsupported osx version")
	}

	return nil
}

// SetUp is responsible for setting up the ergo as proxy
func (c *OSXConfigurator) SetUp(proxyURL string) error {
	err := checkSupportedVersion()
	if err != nil {
		return err
	}

	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" "`+proxyURL+`"`)
}

// SetDown is responsible for remove the ergo as proxy
func (c *OSXConfigurator) SetDown() error {
	err := checkSupportedVersion()
	if err != nil {
		return err
	}

	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" ""`)
}
