package commands

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
	"log"
)

// Run command starts the ergo proxy server.
//
// Usage:
// `ergo run`
func Run(config *proxy.Config) {

	fmt.Println("Ergo Proxy listening on port: " + config.Port)
	err := proxy.ServeProxy(config)
	if err != nil {
		log.Fatal(err)
	}
}
