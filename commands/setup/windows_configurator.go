package setup

// WindowsConfigurator implements Configurator for windows
type WindowsConfigurator struct{}

// SetUp is responsible for setting up the ergo as proxy
func (c *WindowsConfigurator) SetUp(proxyURL string) error {
	err := runner.Run(
		`reg add HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings /v AutoConfigURL /t REG_SZ /d ` + proxyURL + ` /f`)

	return err
}

// SetDown is responsible for remove the ergo as proxy
func (c *WindowsConfigurator) SetDown() error {
	err := runner.Run(
		`reg delete HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings /v AutoConfigURL /f`)

	return err
}
