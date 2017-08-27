package main

import (
	"flag"
	"fmt"
	"github.com/cristianoliveira/ergo/commands"
	"github.com/cristianoliveira/ergo/proxy"
	"os"
)

var VERSION string

const USAGE = `
Ergo proxy.
The local proxy agent for multiple services development.

Usage:
  ergo [options]
  ergo run [options]
  ergo list
  ergo list-names
  ergo url <name>

Options:
  -h      Shows this message.
  -v      Shows ergo's version.

  Run:
	-p      Set ports to proxy
`

func main() {
	help := flag.Bool("h", false, "Shows ergo's help.")
	version := flag.Bool("v", false, "Shows ergo's version.")
	flag.Parse()

	if *help {
		fmt.Println(USAGE)
		os.Exit(0)
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	config := proxy.NewConfig()

	command := flag.NewFlagSet(os.Args[1], flag.ExitOnError)

	command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")
	command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")

	config.Services = proxy.LoadConfig("./.ergo")
	switch os.Args[1] {
	case "list":
		commands.List(config)

	case "list-names":
		commands.ListNames(config)

	case "url":
		if len(os.Args) != 3 {
			fmt.Println(USAGE)
			os.Exit(0)
		}

		name := os.Args[2]
		commands.Url(name, config)

	case "run":
		commands.Run(config)

	default:
		fmt.Println(USAGE)
	}
}
