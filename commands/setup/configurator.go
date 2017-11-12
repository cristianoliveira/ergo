package setup

// Configurator is a interface of setup ergo configuration
type Configurator interface {
	SetUp(proxyURL string) error
	SetDown() error
}

// GetConfigurator gets the right configurator strategy for a given system
func GetConfigurator(system string) Configurator {
	switch system {
	case "windows":
		return &WindowsConfigurator{}
	case "osx":
		return &OSXConfigurator{}
	case "linux-gnome":
		return &LinuxConfigurator{}
	default:
		return nil
	}
}
