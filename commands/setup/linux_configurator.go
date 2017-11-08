package setup

// LinuxConfigurator implements Configurator for windows
type LinuxConfigurator struct{}

// SetUp is responsible for setting up the ergo as proxy
func (c *LinuxConfigurator) SetUp(proxyURL string) error {
	script := `gsettings set org.gnome.system.proxy mode 'auto' &&
	gsettings set org.gnome.system.proxy autoconfig-url '` + proxyURL + `'`

	return runner.Run(
		`/bin/sh -c '` + script + `'`)

}

// SetDown is responsible for remove the ergo as proxy
func (c *LinuxConfigurator) SetDown() error {
	script := `gsettings set org.gnome.system.proxy mode 'none' &&
	gsettings set org.gnome.system.proxy autoconfig-url ''`

	return runner.Run(
		`/bin/sh -c '` + script + `'`)
}
