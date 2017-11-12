package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cristianoliveira/ergo/commands"
	"github.com/cristianoliveira/ergo/proxy"
)

// VERSION of ergo
// When ergo is build without a proper tag/release it is named as `unofficial version`.
// For instance, installing through `go get github.com/cristianoliveira/ergo` or `go build`.
var VERSION = "unofficial version"

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
  ergo remove [options] <service-name|host:port>

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

	config := proxy.NewConfig()
	command := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	command.StringVar(&config.ConfigFile, "config", "./.ergo", "Set the services file")
	command.StringVar(&config.Domain, "domain", ".dev", "Set a custom domain for services")
	command.Parse(os.Args[2:])

	err := config.LoadServices()
	if err != nil {
		log.Fatalf("Could not load services: %v\n", err)
		return nil
	}

	switch os.Args[1] {
	case "list":
		return execute(commands.ListCommand{}, config)

	case "list-names":
		return execute(commands.ListNameCommand{}, config)

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

		return execute(commands.URLCommand{FilterName: name}, config)

	case "run":
		command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")
		command.BoolVar(&config.Verbose, "V", false, "Set verbosity on proxy output")

		command.Parse(os.Args[2:])

		return execute(commands.RunCommand{}, config)
	case "add":
		if len(os.Args) <= 3 {
			return nil
		}

		name := os.Args[2]
		url := os.Args[3]
		service := proxy.NewService(name, url)

		command = flag.NewFlagSet(os.Args[1], flag.ExitOnError)
		command.StringVar(&config.ConfigFile, "config", "./.ergo", "Set the services file")
		command.Parse(os.Args[4:])

		err := config.LoadServices()
		if err != nil {
			log.Fatalf("Could not load services: %v\n", err)
		}

		return execute(commands.AddServiceCommand{Service: service}, config)

	case "remove":
		if len(os.Args) <= 2 {
			return nil
		}

		nameOrURL := os.Args[2]

		service := proxy.NewService(nameOrURL, nameOrURL)

		return execute(commands.RemoveServiceCommand{Service: service}, config)
	}

	return nil
}

func execute(command commands.Command, config *proxy.Config) func() {
	result, err := command.Execute(config)
	return func() {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
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
