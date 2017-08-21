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
`

func main() {
	var command string = "run"

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

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

	config := proxy.LoadConfig("./.ergo")

	switch command {
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
