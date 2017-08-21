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

Options:
  -h      Shows this message.
  -v      Shows ergs's version.
`

func main() {
	var command string = "run"

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	help := flag.Bool("h", false, "Shows ergs's help.")
	version := flag.Bool("v", false, "Shows ergs's version.")

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
		os.Exit(0)

	case "list-names":
		commands.ListNames(config)
		os.Exit(0)

	case "run":
		commands.Run(config)

	default:
		fmt.Println(USAGE)
		os.Exit(0)
	}
}
