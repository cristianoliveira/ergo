package setup

// OSXConfigurator implements Configurator for windows
type OSXConfigurator struct{}

// SetUp is responsible for setting up the ergo as proxy
func (c *OSXConfigurator) SetUp(proxyURL string) error {
	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" "`+proxyURL+`"`)
}

// SetDown is responsible for remove the ergo as proxy
func (c *OSXConfigurator) SetDown() error {
	return RunnerDefault.Run(`networksetup`, `-setautoproxyurl "Wi-Fi" ""`)
}
