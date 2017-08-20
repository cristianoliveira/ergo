package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
)

func Run(config *proxy.Config) {
	fmt.Println("Ergo Proxy listening on port: 2000")
	proxy.ServeProxy(config)
}
