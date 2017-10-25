package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
  -domain     Set a custom domain for services.


run:
  -p          Set ports to proxy.
  -V          Set verbosity on output.

setup:
  -remove     Set remove proxy configurations.
`

func command() func() {
	// Fail fast if we didn't receive a command argument
	if len(os.Args) == 1 {
		return nil
	}

	command := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	configFile := command.String("config", "./.ergo.toml", "Set the config file")
	config := proxy.NewConfig(*configFile)
	domain := command.String("domain", ".dev", "Set a custom domain for services")
	command.Parse(os.Args[2:])

	services := config.Services
	config.Services = services
	config.ConfigFile = *configFile
	config.Domain = *domain

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
		if !strings.HasPrefix(config.Domain, ".") {
			return nil
		}

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

	if *help {
		fmt.Println(USAGE)
		return
	}

	cmd := command()

	if cmd == nil {
		fmt.Println(USAGE)
		os.Exit(1)
	} else {
		cmd()
	}
}
