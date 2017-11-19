package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/commands"
	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

func TestHelperFlags(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "-h"}
		os.Args = args

		main()
		// Output: USAGE
	})

	t.Run("it shows version", func(tt *testing.T) {
		args := []string{"ergo", "-v"}
		os.Args = args

		main()
		// Output: USAGE
	})
}

func TestShowUsage(t *testing.T) {
	t.Run("missing a command so it shows usage", func(tt *testing.T) {
		args := []string{"ergo"}
		os.Args = args

		result, config := prepareSubCommand(args)
		if result != nil || config != nil {
			t.Errorf("Expected result to not be nil got %s and %v", result, config)
		}
	})

	t.Run("when wrong command so it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "foobar"}
		os.Args = args

		result, config := prepareSubCommand(args)
		if result != nil || config != nil {
			t.Errorf("Expected result to not be nil got %s and %v", result, config)
		}
	})
}

func TestCommandsWithDefaultConfigs(t *testing.T) {
	cases := []struct {
		title       string
		args        []string
		commandName string
		output      string
		err         error
	}{
		{
			title:       "it list names",
			args:        []string{"ergo", "list-names"},
			commandName: "ListNameCommand",
			output:      "Ergo Proxy current list",
		},
		{
			title:       "it list services",
			args:        []string{"ergo", "list"},
			commandName: "ListCommand",
			output:      "Ergo Proxy current list",
		},
		{
			title:       "it is url command",
			args:        []string{"ergo", "url", "foo"},
			commandName: "URLCommand",
			output:      "http://foo.dev",
		},
		{
			title:       "when add command include service",
			args:        []string{"ergo", "add", "bla", "http://localhost:3030"},
			commandName: "AddServiceCommand",
			output:      "Service added successfully",
		},
		{
			title:       "when add command try insert existent service",
			args:        []string{"ergo", "add", "foo", "http://localhost:3030"},
			commandName: "AddServiceCommand",
			output:      "",
			err:         fmt.Errorf("Service already present"),
		},
		{
			title:       "when remove command remove a service",
			args:        []string{"ergo", "remove", "foo"},
			commandName: "RemoveServiceCommand",
			output:      "Service Removed",
		},
		{
			title:       "when remove command remove a service",
			args:        []string{"ergo", "remove", "bla"},
			commandName: "RemoveServiceCommand",
			output:      "",
			err:         fmt.Errorf("Service bla not found"),
		},
	}

	services := map[string]proxy.Service{
		"foo": {Name: "foo", URL: "http://localhost:3000"},
		"bar": {Name: "bar", URL: "http://localhost:5000"},
	}

	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		t.Fatalf("Error creating tempfile %s", err.Error())
	}

	for _, c := range cases {
		t.Run(c.title, func(tt *testing.T) {
			command, cfg := prepareSubCommand(c.args)
			cfg.ConfigFile = tmpfile.Name()
			cfg.Services = services

			if command == nil {
				tt.Fatal("Expected command to not be nil")
			}

			commandType := reflect.TypeOf(command).Name()
			if commandType != c.commandName {
				tt.Fatalf("Expected result to be equal %s got %s", c.commandName, commandType)
			}

			if cfg == nil {
				tt.Fatal("Expected config to not be nil")
			}

			output, err := command.Execute(cfg)
			if c.err == nil && err != nil {
				tt.Errorf("Not expecting an error got %s", err)
			} else if c.err != nil && err.Error() != c.err.Error() {
				tt.Errorf("Expecting an error %s got %s", c.err, err)
			}

			if !strings.Contains(output, c.output) {
				tt.Errorf("Expected result to contain '%s' got '%s'", c.output, output)
			}
		})
	}
}

type TestRunner struct{}

func (r *TestRunner) Run(command string) error {
	return nil
}

func TestSetupCommand(t *testing.T) {
	setup.RunnerDefault = &TestRunner{}

	t.Run("when no system given", func(tt *testing.T) {
		args := []string{"ergo", "setup"}

		command, config := prepareSubCommand(args)
		if command != nil || config != nil {
			tt.Fatal("Expected result to not be nil")
		}
	})

	cases := []struct {
		title      string
		args       []string
		output     string
		expectFail bool
	}{
		{
			title:  "when system given is osx",
			args:   []string{"ergo", "setup", "osx"},
			output: "Setup executed",
		},
		{
			title:  "when system given is windows",
			args:   []string{"ergo", "setup", "windows"},
			output: "Setup executed",
		},
		{
			title:  "when system given is linux-gnome",
			args:   []string{"ergo", "setup", "linux-gnome"},
			output: "Setup executed",
		},
		{
			title:      "when system given is unkown",
			args:       []string{"ergo", "setup", "unkown"},
			output:     "",
			expectFail: true,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(tt *testing.T) {
			command, config := prepareSubCommand(c.args)
			if command == nil || config == nil {
				t.Fatal("Expected result to not be nil")
			}

			setup := command.(commands.SetupCommand)

			output, err := setup.Execute(config)
			if c.expectFail && err == nil {
				t.Errorf("Expected error got none")
			}
			if !c.expectFail && err != nil {
				t.Errorf("Expected no error got %s", err)
			}

			if !strings.Contains(output, c.output) {
				t.Errorf("Expected output '%s' got '%s'", c.output, output)
			}

			if setup.Remove {
				t.Error("Expected to not have remove setup", setup)
			}
		})
	}

	t.Run("when removing setup", func(tt *testing.T) {
		args := []string{"ergo", "setup", "osx", "-remove"}

		command, config := prepareSubCommand(args)
		if command == nil || config == nil {
			tt.Fatal("Expected result to not be nil")
		}

		setup := command.(commands.SetupCommand)

		if !setup.Remove {
			t.Error("Expected to have remove setup", setup)
		}
		if setup.System != "osx" {
			t.Error("Expected to setup osx", setup)
		}
	})
}

func TestCommandFlags(t *testing.T) {
	commands := []string{"list", "list-names", "run"}

	for _, cmd := range commands {
		cases := []struct {
			title  string
			args   []string
			config *proxy.Config
		}{
			{
				title: "when -domain flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "-domain", ".foo"},
				config: &proxy.Config{
					Domain: ".foo",
				},
			},
			{
				title: "when -config flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "-config", ".file"},
				config: &proxy.Config{
					ConfigFile: ".file",
				},
			},
			{
				title: "when -p flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "-p", "2002"},
				config: &proxy.Config{
					Port: "2002",
				},
			},
			{
				title: "when -V flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "-V"},
				config: &proxy.Config{
					Verbose: true,
				},
			},
		}

		for _, c := range cases {
			t.Run(c.title, func(tt *testing.T) {
				_, config := prepareSubCommand(c.args)
				if config == nil {
					tt.Fatal("Expected config to not be nil")
				}

				if c.config.Domain != "" && config.Domain != c.config.Domain {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Domain, config.Domain)
				}

				if c.config.ConfigFile != "" && config.ConfigFile != c.config.ConfigFile {
					tt.Errorf("Expected config to have domain %s got %s", c.config.ConfigFile, config.ConfigFile)
				}

				if c.config.Port != "" && config.Port != c.config.Port {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Port, config.Port)
				}

				if !c.config.Verbose && config.Verbose != c.config.Verbose {
					tt.Errorf("Expected config to have domain %v got %v", c.config.Verbose, config.Verbose)
				}
			})
		}
	}
}

func TestSubCommandFlags(t *testing.T) {
	t.Run("setup", func(tt *testing.T) {
		t.Run("when removing setup", func(tt *testing.T) {
			args := []string{"ergo", "setup", "osx", "-remove"}

			command, config := prepareSubCommand(args)
			if command == nil || config == nil {
				tt.Fatal("Expected result to not be nil")
			}

			setup := command.(commands.SetupCommand)

			if !setup.Remove {
				t.Error("Expected to have remove setup", setup)
			}
			if setup.System != "osx" {
				t.Error("Expected to setup osx", setup)
			}
		})
	})

	t.Run("url", func(tt *testing.T) {
		cmd := "url"

		cases := []struct {
			title  string
			args   []string
			config *proxy.Config
		}{
			{
				title: "when -domain flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-domain", ".foo"},
				config: &proxy.Config{
					Domain: ".foo",
				},
			},
			{
				title: "when -config flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-config", ".file"},
				config: &proxy.Config{
					ConfigFile: ".file",
				},
			},
			{
				title: "when -p flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-p", "2002"},
				config: &proxy.Config{
					Port: "2002",
				},
			},
			{
				title: "when -V flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-V"},
				config: &proxy.Config{
					Verbose: true,
				},
			},
		}

		for _, c := range cases {
			t.Run(c.title, func(tt *testing.T) {
				log.Println("tests", c.args)
				_, config := prepareSubCommand(c.args)
				if config == nil {
					tt.Fatal("Expected config to not be nil", c.args)
				}

				if c.config.Domain != "" && config.Domain != c.config.Domain {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Domain, config.Domain)
				}

				if c.config.ConfigFile != "" && config.ConfigFile != c.config.ConfigFile {
					tt.Errorf("Expected config to have domain %s got %s", c.config.ConfigFile, config.ConfigFile)
				}

				if c.config.Port != "" && config.Port != c.config.Port {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Port, config.Port)
				}

				if !c.config.Verbose && config.Verbose != c.config.Verbose {
					tt.Errorf("Expected config to have domain %v got %v", c.config.Verbose, config.Verbose)
				}
			})
		}
	})

	t.Run("add", func(tt *testing.T) {
		cmd := "add"

		cases := []struct {
			title  string
			args   []string
			config *proxy.Config
		}{
			{
				title: "when -domain flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "bla", "-domain", ".foo"},
				config: &proxy.Config{
					Domain: ".foo",
				},
			},
			{
				title: "when -config flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "bla", "-config", ".file"},
				config: &proxy.Config{
					ConfigFile: ".file",
				},
			},
			{
				title: "when -p flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "bla", "-p", "2002"},
				config: &proxy.Config{
					Port: "2002",
				},
			},
			{
				title: "when -V flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "bla", "-V"},
				config: &proxy.Config{
					Verbose: true,
				},
			},
		}

		for _, c := range cases {
			t.Run(c.title, func(tt *testing.T) {
				log.Println("tests", c.args)
				_, config := prepareSubCommand(c.args)
				if config == nil {
					tt.Fatal("Expected config to not be nil", c.args)
				}

				if c.config.Domain != "" && config.Domain != c.config.Domain {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Domain, config.Domain)
				}

				if c.config.ConfigFile != "" && config.ConfigFile != c.config.ConfigFile {
					tt.Errorf("Expected config to have domain %s got %s", c.config.ConfigFile, config.ConfigFile)
				}

				if c.config.Port != "" && config.Port != c.config.Port {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Port, config.Port)
				}

				if !c.config.Verbose && config.Verbose != c.config.Verbose {
					tt.Errorf("Expected config to have domain %v got %v", c.config.Verbose, config.Verbose)
				}
			})
		}
	})

	t.Run("remove", func(tt *testing.T) {
		cmd := "remove"

		cases := []struct {
			title  string
			args   []string
			config *proxy.Config
		}{
			{
				title: "when -domain flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-domain", ".foo"},
				config: &proxy.Config{
					Domain: ".foo",
				},
			},
			{
				title: "when -config flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-config", ".file"},
				config: &proxy.Config{
					ConfigFile: ".file",
				},
			},
			{
				title: "when -p flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-p", "2002"},
				config: &proxy.Config{
					Port: "2002",
				},
			},
			{
				title: "when -V flag is passed for " + cmd,
				args:  []string{"ergo", cmd, "foo", "-V"},
				config: &proxy.Config{
					Verbose: true,
				},
			},
		}

		for _, c := range cases {
			t.Run(c.title, func(tt *testing.T) {
				log.Println("tests", c.args)
				_, config := prepareSubCommand(c.args)
				if config == nil {
					tt.Fatal("Expected config to not be nil", c.args)
				}

				if c.config.Domain != "" && config.Domain != c.config.Domain {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Domain, config.Domain)
				}

				if c.config.ConfigFile != "" && config.ConfigFile != c.config.ConfigFile {
					tt.Errorf("Expected config to have domain %s got %s", c.config.ConfigFile, config.ConfigFile)
				}

				if c.config.Port != "" && config.Port != c.config.Port {
					tt.Errorf("Expected config to have domain %s got %s", c.config.Port, config.Port)
				}

				if !c.config.Verbose && config.Verbose != c.config.Verbose {
					tt.Errorf("Expected config to have domain %v got %v", c.config.Verbose, config.Verbose)
				}
			})
		}
	})
}
