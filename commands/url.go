package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
)

func Url(name string, config *proxy.Config) {
	for _, s := range config.Services {
		if name == s.Name {
			localUrl := `http://` + s.Name + config.Domain
			fmt.Println(localUrl)
		}
	}
}
