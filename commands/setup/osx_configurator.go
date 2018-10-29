package setup

// OSXConfigurator implements Configurator for windows
type OSXConfigurator struct{}

// SetUp is responsible for setting up the ergo as proxy
func (c *OSXConfigurator) SetUp(proxyURL string) error {
	script := `sudo networksetup -setautoproxyurl "Wi-Fi" "` + proxyURL + `"`

	return RunnerDefault.Run(`/bin/sh`, ` -c '`+script+`'`)
}

// SetDown is responsible for remove the ergo as proxy
func (c *OSXConfigurator) SetDown() error {
	script := `sudo networksetup -setautoproxyurl "Wi-Fi" ""`

	return RunnerDefault.Run(`/bin/sh`, ` -c '`+script+`'`)
}
