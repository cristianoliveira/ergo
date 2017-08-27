package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
)

func Run(config *proxy.Config) {
	fmt.Println("Ergo Proxy listening on port: " + config.Port)
	proxy.ServeProxy(config)
}
