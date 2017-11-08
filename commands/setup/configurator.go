package setup

// Configurator is a interface of setup ergo configuration
type Configurator interface {
	SetUp(proxyURL string) error
	SetDown() error
}
