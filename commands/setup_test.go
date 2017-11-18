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
		{Name: "test.dev", URL: "localhost:999"},
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
}

func (r *TestRunner) Run(command string) error {
	if r.ExpectToInclude == "" {
		fmt.Println("No expectation")
		return nil
	}

	if !strings.Contains(command, r.ExpectToInclude) {
		r.Test.Fatalf(
			"Expected command to include '%s' but it is not present.\n Command: %s",
			r.ExpectToInclude,
			command,
		)
	}

	return nil
}

func TestSetupLinuxGnome(t *testing.T) {
	config := buildConfig([]proxy.Service{
		{Name: "test.dev", URL: "localhost:999"},
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title: "expect to run with sh",
				CommandExpectToInclude: "/bin/sh -c",
			},
			{
				Title: "expect to set networking mode auto",
				CommandExpectToInclude: "mode 'auto'",
			},
			{
				Title: "expect to set networking url",
				CommandExpectToInclude: `autoconfig-url '` + config.GetProxyPacURL() + `'`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

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
				Title: "expect to run with sh",
				CommandExpectToInclude: "/bin/sh -c",
			},
			{
				Title: "expect to set networking mode none",
				CommandExpectToInclude: "mode 'none'",
			},
			{
				Title: "expect to set networking no url",
				CommandExpectToInclude: `autoconfig-url ''`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

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
		{Name: "test.dev", URL: "localhost:999"},
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title: "expect to run with sh",
				CommandExpectToInclude: "/bin/sh -c",
			},
			{
				Title: "expect to set networking proxy pac url",
				CommandExpectToInclude: `-setautoproxyurl "Wi-Fi" "` + config.GetProxyPacURL() + `"`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

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
				Title: "expect to run with sh",
				CommandExpectToInclude: "/bin/sh -c",
			},
			{
				Title: "expect to set networking wi-fi to none",
				CommandExpectToInclude: `-setautoproxyurl "Wi-Fi" ""`,
			},
		}

		for _, c := range cases {
			t.Run(c.Title, func(tt *testing.T) {
				setup.RunnerDefault = &TestRunner{
					Test:            tt,
					ExpectToInclude: c.CommandExpectToInclude,
				}

				command := SetupCommand{System: "osx", Remove: true}
				_, err := command.Execute(config)
				if err != nil {
					t.Fatalf(err.Error())
				}
			})
		}
	})
}

func TestSetupWindows(t *testing.T) {
	config := buildConfig([]proxy.Service{
		{Name: "test.dev", URL: "localhost:999"},
	})

	t.Run("when setting up", func(t *testing.T) {
		var cases = []struct {
			Title                  string
			CommandExpectToInclude string
		}{
			{
				Title: "expect to add a new register",
				CommandExpectToInclude: "reg add",
			},
			{
				Title: "expect to set networking proxy pac url",
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
				Title: "expect to delete the register",
				CommandExpectToInclude: "reg delete",
			},
			{
				Title: "expect to set networking wi-fi to none",
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
