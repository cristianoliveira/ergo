package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/commands/setup"
	"github.com/cristianoliveira/ergo/proxy"
)

func initialize() (proxy.Config, error) {
	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		return proxy.Config{}, fmt.Errorf("Error creating tempfile: %s", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		return proxy.Config{}, fmt.Errorf("Error writing to temporary file: %s", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		return proxy.Config{}, fmt.Errorf("Error closing temp file: %s", err.Error())
	}

	if err != nil {
		return proxy.Config{}, fmt.Errorf("No error expected while initializing Config file. Got %s", err.Error())
	}
	config := proxy.Config{}
	config.ConfigFile = tmpfile.Name()
	config.Services, err = proxy.LoadServices(config.ConfigFile)

	if err != nil {
		return proxy.Config{}, fmt.Errorf("No error expected while loading services from config file. Got %s", err.Error())
	}

	return config, nil
}

func TestSetup(t *testing.T) {

	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	service := proxy.Service{}
	service.Name = config.Services[0].Name

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	Setup("inexistent-os", false, &config)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "List of supported system") {
		t.Fatalf("Expected Setup to tell us about the supported systems if we ask it to run an unsupported system. Got %s.", out)
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
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

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

				Setup("linux-gnome", false, &config)
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

				Setup("linux-gnome", true, &config)
			})
		}
	})
}

func TestSetupOSX(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

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

				Setup("osx", false, &config)
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

				Setup("osx", true, &config)
			})
		}
	})
}

func TestSetupWindows(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

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

				Setup("windows", false, &config)
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

				Setup("windows", true, &config)
			})
		}
	})
}
