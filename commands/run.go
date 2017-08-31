package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
	"log"
)

func Run(config *proxy.Config) {
	err := proxy.ServeProxy(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ergo Proxy listening on port: " + config.Port)
}
