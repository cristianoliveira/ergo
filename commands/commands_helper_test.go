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

	config := proxy.NewConfig()
	config.ConfigFile = tmpfile.Name()

	for _, s := range services {
		config.AddService(s)
	}

	return config
}
