package commands

import (
	"fmt"
	"log"

	"github.com/cristianoliveira/ergo/proxy"
)

// Run command starts the ergo proxy server.
//
// Usage:
// `ergo run`
func Run(config *proxy.Config) {

	fmt.Println("Ergo Proxy listening on port: " + config.Port)
	fmt.Println("Ergo Proxy listening for domain: " + config.Domain)
	err := proxy.ServeProxy(config)
	if err != nil {
		log.Fatal(err)
	}
}
