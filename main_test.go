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

func TestMain(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo"}
		os.Args = args

		main()
		// Output: USAGE
	})

	t.Run("it shows help", func(tt *testing.T) {
		args := []string{"ergo", "-h"}
		os.Args = args

		main()
		// Output: USAGE
	})

	t.Run("it shows version", func(tt *testing.T) {
		args := []string{"ergo", "-v"}
		os.Args = args

		main()
	})

	t.Run("it list version", func(tt *testing.T) {
		args := []string{"ergo", "list"}
		os.Args = args

		main()
	})
}

func TestShowUsage(t *testing.T) {
	t.Run("when missing a command", func(tt *testing.T) {
		args := []string{"ergo"}
		os.Args = args

		result, config := prepareSubCommand(args)
		if result != nil || config != nil {
			t.Errorf("Expected result to not be nil got %s and %v", result, config)
		}
	})

	t.Run("when wrong unknown command", func(tt *testing.T) {
		args := []string{"ergo", "foobar"}
		os.Args = args

		result, config := prepareSubCommand(args)
		if result != nil || config != nil {
			t.Errorf("Expected result to not be nil got %s and %v", result, config)
		}
	})

	t.Run("for add command", func(tt *testing.T) {
		t.Run("when missing arguments", func(tt *testing.T) {
			args := []string{"ergo", "add"}
			os.Args = args

			result, config := prepareSubCommand(args)
			if result != nil || config != nil {
				t.Errorf("Expected result to not be nil got %s and %v", result, config)
			}
		})

		t.Run("when missing url", func(tt *testing.T) {
			args := []string{"ergo", "add", "foo"}
			os.Args = args

			result, config := prepareSubCommand(args)
			if result != nil || config != nil {
				t.Errorf("Expected result to not be nil got %s and %v", result, config)
			}
		})
	})
}

type TestRunner struct{}

func (r *TestRunner) Run(command, args string) error {
	return nil
}

func TestListCommand(t *testing.T) {
	services := map[string]proxy.Service{
		"foo": {Name: "foo", URL: "http://localhost:3000"},
		"bar": {Name: "bar", URL: "http://localhost:5000"},
	}

	t.Run("list services", func(tt *testing.T) {
		args := []string{"ergo", "list"}
		expectedOutput := "Ergo Proxy current list"

		command, config := prepareSubCommand(args)
		config.Services = services

		if command == nil {
			tt.Fatal("Expected command to not be nil")
		}

		output, err := command.Execute(config)
		if err != nil {
			tt.Errorf("Not expecting an error got %s", err)
		}

		if !strings.Contains(output, expectedOutput) {
			tt.Errorf("Expected result to contain '%s' got '%s'", expectedOutput, output)
		}
	})
}

func TestListNamesCommand(t *testing.T) {
	services := map[string]proxy.Service{
		"foo": {Name: "foo", URL: "http://localhost:3000"},
		"bar": {Name: "bar", URL: "http://localhost:5000"},
	}

	t.Run("list names services", func(tt *testing.T) {
		args := []string{"ergo", "list-names"}
		expectedOutput := "Ergo Proxy current list"

		command, config := prepareSubCommand(args)
		config.Services = services

		if command == nil {
			tt.Fatal("Expected command to not be nil")
		}

		output, err := command.Execute(config)
		if err != nil {
			tt.Errorf("Not expecting an error got %s", err)
		}

		if !strings.Contains(output, expectedOutput) {
			tt.Errorf("Expected result to contain '%s' got '%s'", expectedOutput, output)
		}
	})
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

func TestUrlCommand(t *testing.T) {
	commandName := "URLCommand"

	services := map[string]proxy.Service{
		"foo": {Name: "foo", URL: "http://localhost:3000"},
		"bar": {Name: "bar", URL: "http://localhost:5000"},
	}

	t.Run("when there is the service", func(tt *testing.T) {
		expectedOutput := "http://foo.dev"
		args := []string{"ergo", "url", "foo"}

		command, argconfig := prepareSubCommand(args)

		config := proxy.NewConfig()
		config.OverrideBy(argconfig)
		config.Services = services

		if command == nil || config == nil {
			tt.Fatalf("Expected command to be nil got %v", command)
		}

		commandType := reflect.TypeOf(command).Name()
		if commandType != commandName {
			tt.Fatalf("Expected result to be equal %s got %s", commandName, commandType)
		}

		output, err := command.Execute(config)
		if err != nil {
			tt.Errorf("Not expecting an error got %s", err)
		}

		if !strings.Contains(output, expectedOutput) {
			tt.Errorf("Expected result to contain '%s' got '%s'", expectedOutput, output)
		}
	})

	t.Run("when there is no service", func(tt *testing.T) {
		args := []string{"ergo", "url", "foobla"}

		command, argconfig := prepareSubCommand(args)

		config := proxy.NewConfig()
		config.OverrideBy(argconfig)
		config.Services = services

		if command == nil || config == nil {
			tt.Fatalf("Expected command to be nil got %v", command)
		}

		commandType := reflect.TypeOf(command).Name()
		if commandType != commandName {
			tt.Fatalf("Expected result to be equal %s got %s", commandName, commandType)
		}

		_, err := command.Execute(config)
		if err == nil {
			tt.Errorf("Not expecting an error got %s", err)
		}
	})
}

func TestRunCommand(t *testing.T) {
	commandName := "RunCommand"

	t.Run("runs", func(tt *testing.T) {
		args := []string{"ergo", "run"}

		command, config := prepareSubCommand(args)

		if command == nil || config == nil {
			tt.Fatalf("Expected command to be nil got %v", command)
		}

		commandType := reflect.TypeOf(command).Name()
		if commandType != commandName {
			tt.Fatalf("Expected result to be equal %s got %s", commandName, commandType)
		}
	})
}

func TestAddCommand(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		services := map[string]proxy.Service{
			"foo": {Name: "foo", URL: "http://localhost:3000"},
			"bar": {Name: "bar", URL: "http://localhost:5000"},
		}

		tmpfile, err := ioutil.TempFile("", "testaddservice")
		if err != nil {
			t.Fatalf("Error creating tempfile %s", err.Error())
		}

		cases := []struct {
			title       string
			args        []string
			commandName string
			output      string
			err         error
		}{
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
	})

	t.Run("fail", func(tt *testing.T) {
		t.Run("when missing arguments", func(tt *testing.T) {
			args := []string{"ergo", "add"}

			command, config := prepareSubCommand(args)
			if command != nil && config != nil {
				tt.Fatalf("Expected command to be nil got %s", command)
			}
		})

		t.Run("when missing url", func(tt *testing.T) {
			args := []string{"ergo", "add", "fo"}

			command, config := prepareSubCommand(args)
			if command != nil && config != nil {
				tt.Fatalf("Expected command to be nil got %s", command)
			}
		})
	})
}

func TestRemoveCommand(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
		services := map[string]proxy.Service{
			"foo": {Name: "foo", URL: "http://localhost:3000"},
			"bar": {Name: "bar", URL: "http://localhost:5000"},
		}

		tmpfile, err := ioutil.TempFile("", "testaddservice")
		if err != nil {
			t.Fatalf("Error creating tempfile %s", err.Error())
		}

		cases := []struct {
			title       string
			args        []string
			commandName string
			output      string
			err         error
		}{
			{
				title:       "when remove by url",
				args:        []string{"ergo", "remove", "http://localhost:5000"},
				commandName: "RemoveServiceCommand",
				output:      "Service Removed",
			},
			{
				title:       "when remove by name",
				args:        []string{"ergo", "remove", "foo"},
				commandName: "RemoveServiceCommand",
				output:      "Service Removed",
			},
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
				if fmt.Sprintf("%s", c.err) != fmt.Sprintf("%s", err) {
					tt.Errorf("Expecting an error %s got %s", c.err, err)
				}

				if !strings.Contains(output, c.output) {
					tt.Errorf("Expected result to contain '%s' got '%s'", c.output, output)
				}
			})
		}
	})

	t.Run("fail", func(tt *testing.T) {
		t.Run("when missing arguments", func(tt *testing.T) {
			args := []string{"ergo", "remove"}

			command, config := prepareSubCommand(args)
			if command != nil && config != nil {
				tt.Fatalf("Expected command to be nil got %s", command)
			}
		})
	})
}

func TestCommandLineFlags(t *testing.T) {
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
