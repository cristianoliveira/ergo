package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cristianoliveira/ergo/commands"
	"github.com/cristianoliveira/ergo/proxy"
)

//VERSION of ergo.
var VERSION string

//USAGE details the usage for ergo proxy.
const USAGE = `
Ergo proxy.
The local proxy agent for multiple services development.

Usage:
  ergo run [options]
  ergo list
  ergo list-names
  ergo url <name>
  ergo setup [linux-gnome|osx|windows]
  ergo add <service-name> <host:port>

Options:
  -h      Shows this message.
  -v      Shows ergo's version.

run:
  -p          Set ports to proxy.
  -V          Set verbosity on output.
  -config     Set the config file to the proxy.

setup:
  -remove     Set remove proxy configurations.
`

func main() {
	help := flag.Bool("h", false, "Shows ergo's help.")
	version := flag.Bool("v", false, "Shows ergo's version.")

	flag.Parse()

	if *help || len(os.Args) == 1 {
		fmt.Println(USAGE)
		os.Exit(0)
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	config := proxy.NewConfig()
	command := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	configFile := command.String("config", "./.ergo", "Set the services file")
	command.Parse(os.Args[2:])

	services, err := proxy.LoadServices(*configFile)
	if err != nil {
		log.Printf("Could not load services: %v\n", err)
	}
	config.Services = services

	switch os.Args[1] {
	case "list":
		commands.List(config)

	case "list-names":
		commands.ListNames(config)

	case "setup":
		if len(os.Args) <= 2 {
			fmt.Println(USAGE)
			os.Exit(0)
		}

		system := os.Args[2]
		setupRemove := command.Bool("remove", false, "Set remove proxy configurations.")
		command.Parse(os.Args[3:])

		commands.Setup(system, *setupRemove, config)

	case "url":
		if len(os.Args) != 3 {
			fmt.Println(USAGE)
			os.Exit(0)
		}

		name := os.Args[2]
		commands.URL(name, config)

	case "run":
		command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")
		command.BoolVar(&config.Verbose, "V", false, "Set verbosity on proxy output")

		command.Parse(os.Args[2:])
		commands.Run(config)
	case "add":
		if len(os.Args) <= 3 {
			fmt.Println(USAGE)
			os.Exit(0)
		}

		name := os.Args[2]
		url := os.Args[3]
		service := proxy.NewService(name, url)

		commands.AddService(config, service, *configFile)

	default:
		fmt.Println(USAGE)
	}
}
