package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

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

func fakeCommandLinuxGnome(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_LINUX_GNOME=1"}
	return cmd
}

func fakeCommandLinuxGnomeRemove(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_LINUX_GNOME_REMOVE=1"}
	return cmd
}

func fakeCommandOsx(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_OSX=1"}
	return cmd
}

func fakeCommandOsxRemove(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_OSX_REMOVE=1"}
	return cmd
}

func fakeCommandWindows(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_WINDOWS=1"}
	return cmd
}

func fakeCommandWindowsRemove(executable string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestMeWantHelp", "--", executable}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"ME_WANT_HELP=1", "TEST_FOR_WINDOWS_REMOVE=1"}
	return cmd
}

func TestMeWantHelp(t *testing.T) {
	log.Println("ME_WANT_HELP", os.Getenv("ME_WANT_HELP"))
	if os.Getenv("ME_WANT_HELP") != "1" {
		return
	}

	fmt.Fprintf(os.Stderr, "%s,%v", "deci", os.Args)

	if os.Getenv("TEST_FOR_LINUX_GNOME") == "1" {
		//we should also check for arguments
		if os.Args[3] != "/bin/bash" {
			t.Fatalf("expected \"/bin/bash\". got %s", os.Args[1])
		}
	} else if os.Getenv("TEST_FOR_LINUX_GNOME_REMOVE") == "1" {
		//we should also check for arguments
		if os.Args[3] != "/bin/bash" {
			t.Fatalf("expected \"/bin/bash\". got %s", os.Args[1])
		}
	} else if os.Getenv("TEST_FOR_OSX") == "1" {
		//we should also check for arguments
		if os.Args[3] != "/bin/bash" {
			t.Fatalf("expected \"/bin/bash\". got %s", os.Args[1])
		}
	} else if os.Getenv("TEST_FOR_OSX_REMOVE") == "1" {
		//we should also check for arguments
		if os.Args[3] != "/bin/bash" {
			t.Fatalf("expected \"/bin/bash\". got %s", os.Args[1])
		}
	} else if os.Getenv("TEST_FOR_WINDOWS") == "1" {
		rez := strings.Join(os.Args[3:], " ")
		if rez != "reg add HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings /v AutoConfigURL /t REG_SZ /d http://127.0.0.1:/proxy.pac /f" {
			t.Fatalf("expected \"reg\". got %s", os.Args[1])
		}
	} else if os.Getenv("TEST_FOR_WINDOWS_REMOVE") == "1" {
		rez := strings.Join(os.Args[3:], " ")
		if rez != "reg delete HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings /v AutoConfigURL /f" {
			t.Fatalf("expected \"reg\". got %s", os.Args[1])
		}
	}

	os.Exit(0)
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

func TestSetupLinuxGnome(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandLinuxGnome
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("linux-gnome", false, &config)

}

func TestSetupLinuxGnomeRemove(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandLinuxGnomeRemove
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("linux-gnome", true, &config)

}

func TestSetupOsx(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandOsx
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("osx", false, &config)
}

func TestSetupOsxRemove(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandOsxRemove
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("osx", true, &config)
}

func TestSetupWindows(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandWindows
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("windows", false, &config)
}

func TestSetupRemoveWindows(t *testing.T) {
	config, err := initialize()

	if err != nil {
		t.Fatalf(err.Error())
	}

	ergoCmd = fakeCommandWindowsRemove
	defer func() {
		ergoCmd = exec.Command
	}()

	Setup("windows", true, &config)
}
