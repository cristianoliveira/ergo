package setup

// Configurator is a interface of setup ergo configuration
type Configurator interface {
	SetUp(proxyURL string) error
	SetDown() error
}

func GetCofigurator(system string) Configurator {
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
