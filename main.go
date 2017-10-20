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
  ergo setup [options] [linux-gnome|osx|windows] [-remove]
  ergo add [options] <service-name> <host:port>

Options:
  -h      Shows this message.
  -v      Shows ergo's version.
  -config     Set the config file to the proxy.

run:
  -p          Set ports to proxy.
  -V          Set verbosity on output.

setup:
  -remove     Set remove proxy configurations.
`

func command() func() {
	config := proxy.NewConfig()
	command := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	configFile := command.String("config", "./.ergo", "Set the services file")
	command.Parse(os.Args[2:])

	services, err := proxy.LoadServices(*configFile)
	if err != nil {
		log.Printf("Could not load services: %v\n", err)
	}
	config.Services = services
	config.ConfigFile = *configFile

	switch os.Args[1] {
	case "list":
		return func() {
			commands.List(config)
		}

	case "list-names":
		return func() {
			commands.ListNames(config)
		}

	case "setup":
		if len(os.Args) <= 2 {
			return nil
		}

		system := command.Args()[0]
		setupRemove := command.Bool("remove", false, "Set remove proxy configurations.")
		command.Parse(command.Args()[1:])

		return func() {
			commands.Setup(system, *setupRemove, config)
		}

	case "url":
		if len(os.Args) != 3 {
			return nil
		}

		name := os.Args[2]
		return func() {
			commands.URL(name, config)
		}

	case "run":
		command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")
		command.BoolVar(&config.Verbose, "V", false, "Set verbosity on proxy output")

		command.Parse(os.Args[2:])

		return func() {
			commands.Run(config)
		}
	case "add":
		if len(os.Args) <= 3 {
			return nil
		}

		name := os.Args[2]
		url := os.Args[3]
		service := proxy.NewService(name, url)

		return func() {
			commands.AddService(config, service, *configFile)
		}
	}

	return nil
}

func main() {
	help := flag.Bool("h", false, "Shows ergo's help.")
	version := flag.Bool("v", false, "Shows ergo's version.")

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	cmd := command()
	showUsage := *help || len(os.Args) == 1 || cmd == nil

	if showUsage {
		fmt.Println(USAGE)
	} else {
		cmd()
	}
}
