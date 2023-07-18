package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

func TestSetup(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.UnsafeNewService("test.dev", "localhost:9999"),
	})

	command := SetupCommand{System: "inexistent-os", Remove: false}
	_, err := command.Execute(config)

	if !strings.Contains(err.Error(), "List of supported system") {
		t.Fatalf("Expected Setup to tell us about the supported systems if we ask"+
			" it to run an unsupported system. Got %s.", err.Error())
	}
}

type TestRunner struct {
	Test            *testing.T
	ExpectToInclude string
	History         string
}

func (r *TestRunner) Run(command string, args ...string) ([]byte, error) {
	if r.ExpectToInclude == "" {
		fmt.Println("No expectation")
		return []byte{}, nil
	}

	r.History = r.History + " > " + command + " " + strings.Join(args, " ")
	if !strings.Contains(r.History, r.ExpectToInclude) {
		r.Test.Fatalf(
			"Expected command to include '%s' but it is not present.\n Command: %s",
			r.ExpectToInclude,
			r.History,
		)
	}

	return []byte{}, nil
}

type TestRunnerWithOutput struct {
	Test         *testing.T
	History      string
	MockedOutput map[string][]byte
}

func (r *TestRunnerWithOutput) Mock(command string, output []byte) {
	// parse command to a unique key
	key := strings.Join(strings.Split(command, " "), "_")
	r.MockedOutput[key] = output
}

func (r *TestRunnerWithOutput) Run(command string, args ...string) ([]byte, error) {
	commandWithArgs := command + " " + strings.Join(args, " ")
	key := strings.Join(strings.Split(commandWithArgs, " "), "_")
	mockedOutput, ok := r.MockedOutput[key]
	if !ok {
		r.Test.Fatalf("No more expectation for command %s", commandWithArgs)
		return []byte{}, nil
	}

	return mockedOutput, nil
}

func TestSetupLinuxGnome(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.UnsafeNewService("test.dev", "localhost:9999"),
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to set networking mode auto",
				CommandExpectToInclude: "mode 'auto'",
			},
			{
				Title:                  "expect to set networking url",
				CommandExpectToInclude: `autoconfig-url '` + config.GetProxyPacURL() + `'`,
			},
		}

		testRunner := &TestRunner{}
		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				testRunner.Test = tt
				testRunner.ExpectToInclude = c.CommandExpectToInclude

				setup.RunnerDefault = testRunner
				command := SetupCommand{System: "linux-gnome", Remove: false}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})

	t.Run("when setting down", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to set networking mode none",
				CommandExpectToInclude: "mode 'none'",
			},
			{
				Title:                  "expect to set networking no url",
				CommandExpectToInclude: `autoconfig-url ''`,
			},
		}

		testRunner := &TestRunner{}
		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				testRunner.Test = tt
				testRunner.ExpectToInclude = c.CommandExpectToInclude

				setup.RunnerDefault = testRunner
				command := SetupCommand{System: "linux-gnome", Remove: true}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})
}

func TestSetupOSX(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.UnsafeNewService("test.dev", "localhost:999"),
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to set networking proxy pac url",
				CommandExpectToInclude: `networksetup -setautoproxyurl "Wi-Fi" "` + config.GetProxyPacURL() + `"`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				mockedRunner := &TestRunnerWithOutput{
					Test:         tt,
					MockedOutput: map[string][]byte{},
				}

				mockedRunner.Mock("/bin/sh -c sw_vers -productVersion", []byte("10.11.6"))
				mockedRunner.Mock(c.CommandExpectToInclude, []byte{})

				setup.RunnerDefault = mockedRunner
				command := SetupCommand{System: "osx", Remove: false}
				_, err := command.Execute(config)

				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})

	t.Run("when setting down", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to set networking wi-fi to none",
				CommandExpectToInclude: `networksetup -setautoproxyurl "Wi-Fi" ""`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				mockedRunner := &TestRunnerWithOutput{
					Test:         tt,
					MockedOutput: map[string][]byte{},
				}

				mockedRunner.Mock("/bin/sh -c sw_vers -productVersion", []byte("10.11.6"))
				mockedRunner.Mock(c.CommandExpectToInclude, []byte{})
				setup.RunnerDefault = mockedRunner

				command := SetupCommand{System: "osx", Remove: true}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})

	t.Run("it does not work after Catalina version * > 10", func(tt *testing.T) {
		mockedRunner := &TestRunnerWithOutput{
			Test:         tt,
			MockedOutput: map[string][]byte{},
		}

		mockedRunner.Mock("/bin/sh -c sw_vers -productVersion", []byte("11.11.6"))
		setup.RunnerDefault = mockedRunner

		command := SetupCommand{System: "osx", Remove: true}
		_, err := command.Execute(config)
		// check if the error message contains Setup failed
		if !strings.Contains(err.Error(), "Setup failed cause unsupported osx version") {
			t.Fatalf(err.Error())
		}
	})
}

func TestSetupWindows(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.UnsafeNewService("test.dev", "localhost:999"),
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to add a new register",
				CommandExpectToInclude: "add",
			},
			{
				Title:                  "expect to set networking proxy pac url",
				CommandExpectToInclude: `AutoConfigURL /t REG_SZ /d ` + config.GetProxyPacURL(),
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

				command := SetupCommand{System: "windows", Remove: false}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})

	t.Run("when setting down", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title:                  "expect to delete the register",
				CommandExpectToInclude: "delete",
			},
			{
				Title:                  "expect to set networking wi-fi to none",
				CommandExpectToInclude: "AutoConfigURL /f",
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

				command := SetupCommand{System: "windows", Remove: true}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})
}
