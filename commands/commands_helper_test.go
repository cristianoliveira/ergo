package commands

import (
	"io/ioutil"

	"github.com/cristianoliveira/ergo/proxy"
)

func buildConfig(services []proxy.Service) proxy.Config {
	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		panic("Error creating tempfile" + err.Error())
	}

	return proxy.Config{
		ConfigFile: tmpfile.Name(),
		Services:   services,
	}
}
