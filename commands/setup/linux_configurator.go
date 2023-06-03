package setup

// LinuxConfigurator implements Configurator for windows
type LinuxConfigurator struct{}

// SetUp is responsible for setting up the ergo as proxy
func (c *LinuxConfigurator) SetUp(proxyURL string) error {
	_, err := RunnerDefault.Run("gsettings", "set", "org.gnome.system.proxy", "mode", "'auto'")

	if err != nil {
		return err
	}

	_, err = RunnerDefault.Run("gsettings", "set", "org.gnome.system.proxy", "autoconfig-url", `'`+proxyURL+`'`)
	return err
}

// SetDown is responsible for remove the ergo as proxy
func (c *LinuxConfigurator) SetDown() error {
	_, err := RunnerDefault.Run("gsettings", "set", "org.gnome.system.proxy", "mode", "'none'")

	if err != nil {
		return err
	}

	_, err = RunnerDefault.Run("gsettings", "set", "org.gnome.system.proxy", "autoconfig-url", "''")
	return err
}
