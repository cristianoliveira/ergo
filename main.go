package main

import (
	"flag"
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
	"net/http"
	"os"
)

const VERSION = "0.0.4"

const USAGE = `
Ergo proxy.
The local proxy agent for multiple services development.

Usage:
  ergo [options]
  ergo run [options]

Options:
  -h      Shows this message.
  -v      Shows ergs's version.
`

func main() {
	run := flag.Bool("run", true, "Starts the proxy service")
	help := flag.Bool("-h", false, "Shows ergs's help.")
	version := flag.Bool("-v", false, "Shows ergs's version.")

	if *help {
		fmt.Println(USAGE)
		os.Exit(0)
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *run {
		config := proxy.LoadConfig("./.ergo")

		http.HandleFunc("/proxy.pac", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./resources/proxy.pac")
		})

		proxy := proxy.NewErgoProxy(config)
		http.Handle("/", proxy)

		fmt.Println("Ergo Proxy listening on port: 2000")
		http.ListenAndServe(":2000", nil)
		return
	}
}
