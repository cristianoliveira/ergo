package commands

import (
	"bytes"
	"errors"
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
	ExpectedCommand string
}

func (t *TestRunner) Run(command string) error {
	if t.ExpectedCommand == "" {
		return nil
	}

	if command != t.ExpectedCommand {
		return errors.New("Expected command" + t.ExpectedCommand + "received: " + command)
	}

	return nil
}

func TestSetupLinuxGnome(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{
		ExpectedCommand: "Foo",
	}

	Setup("linux-gnome", false, &config)

}

func TestSetupLinuxGnomeRemove(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{
		ExpectedCommand: "Foo",
	}

	Setup("linux-gnome", true, &config)
}

func TestSetupOsx(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{}

	Setup("osx", false, &config)
}

func TestSetupOsxRemove(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{}

	Setup("osx", true, &config)
}

func TestSetupWindows(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{}

	Setup("windows", false, &config)
}

func TestSetupRemoveWindows(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	setup.RunnerDefault = &TestRunner{}

	Setup("windows", true, &config)
}
