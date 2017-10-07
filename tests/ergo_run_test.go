// +build integration

package main

import (
	"fmt"
	"github.com/cristianoliveira/ergo/commands"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var initialSetup interface{}

func ergo(args ...string) *exec.Cmd {
	return exec.Command(filepath.Join("..", "bin", "ergo"), args...)
}

func TestListApps(t *testing.T) {

	t.Run("it lists the apps", func(tt *testing.T) {
		appsOutput := []string{
			"http://foo.dev -> http://localhost:3000",
			"http://bla.dev -> http://localhost:5000",
			"http://withspaces.dev -> http://localhost:8080",
			"http://one.domain.dev -> http://localhost:8081",
			"http://two.domain.dev -> http://localhost:8082",
			"http://redis://redislocal.dev -> redis://localhost:6543",
		}

		cmd := ergo("list")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		for _, app := range appsOutput {
			if !strings.Contains(output, app) {
				tt.Errorf("Expected output:\n %s \n got %s", output, app)
			}
		}
	})
}

func TestListAppNames(t *testing.T) {

	t.Run("it lists the app names", func(tt *testing.T) {
		appsOutput := []string{
			"foo -> http://localhost:3000",
			"bla -> http://localhost:5000",
			"withspaces -> http://localhost:8080",
			"one.domain -> http://localhost:8081",
			"two.domain -> http://localhost:8082",
			"redis://redislocal -> redis://localhost:6543",
		}

		cmd := ergo("list-names", "foo")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		for _, app := range appsOutput {
			if !strings.Contains(output, app) {
				tt.Errorf("Expected output:\n %s \n got %s", output, app)
			}
		}
	})
}
func TestShowUrlForName(t *testing.T) {

	t.Run("it shows the url for a given name", func(tt *testing.T) {
		appsOutput := map[string]string{
			"foo":                "http://foo.dev",
			"bla":                "http://bla.dev",
			"withspaces":         "http://withspaces.dev",
			"one.domain":         "http://one.domain.dev",
			"two.domain":         "http://two.domain.dev",
			"redis://redislocal": "redis://redislocal.dev",
		}

		for name, url := range appsOutput {
			cmd := ergo("url", name)
			bs, err := cmd.Output()
			if err != nil {
				tt.Fatal(err)
			}

			output := string(bs)
			if strings.Trim(output, " \r\n") != url {
				tt.Errorf("Expected output:\n [%s] \n got [%s]", url, strings.Trim(output, " \r\n"))
			}
		}
	})

}

func TestAddService(t *testing.T) {
	t.Run("it adds new service if not present", func(tt *testing.T) {
		appsOutput := fmt.Sprintf("%s\n", "Service added successfully!")

		fileContent, err := ioutil.ReadFile("./.ergo")

		if err == nil {
			//we clean after the test. Otherwise the next test will fail
			defer ioutil.WriteFile("./.ergo", fileContent, 0755)
		}

		cmd := ergo("add", "new.service", "http://localhost:8083")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output:\n %s \n got %s", appsOutput, output)
		}
	})

	t.Run("it prints message for already added service", func(tt *testing.T) {
		appsOutput := fmt.Sprintf("%s\n", "Service already present!")

		cmd := ergo("add", "foo", "http://localhost:3000")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output:\n %s \n got %s", appsOutput, output)
		}
	})
}

// functions needed for setup and run testing
func getOS() *string {
	var runOS string
	//we setup the right runtime
	if runtime.GOOS == "windows" {
		runOS = "windows"
	} else if runtime.GOOS == "linux" {
		//here we only have the tests for gnome
		runOS = "linux-gnome"
	} else if runtime.GOOS == "darwin" {
		runOS = "osx"
	}

	return &runOS
}

func setupErgo() error {
	var err error
	runOS := getOS()

	//we need to store the initial values to make sure they are restored
	//linux-gnome
	if *runOS == "linux-gnome" {
		var mode string
		var autoconfigURL string
		mode, err = getLinuxGnomeProxyMode()
		if err != nil {
			return err
		}
		autoconfigURL, err = getLinuxGnomeProxyAutoConfig()
		if err != nil {
			return err
		}
		initialSetup = struct {
			mode          string
			autoconfigURL string
		}{mode: mode, autoconfigURL: autoconfigURL}
	}
	//osx
	if *runOS == "osx" {
		var autoconfigURL string
		autoconfigURL, err = getDarwinProxyAutoURL()
		if err != nil {
			return err
		}
		initialSetup = struct{ autoconfigURL string }{autoconfigURL: autoconfigURL}
	}
	//windows
	if *runOS == "windows" {
		var autoconfigURL string
		autoconfigURL, err = getWindowsProxyAutoURL()
		if err != nil {
			return err
		}
		initialSetup = struct{ autoconfigURL string }{autoconfigURL: autoconfigURL}
	}

	if runOS != nil {
		log.Println("starting setup " + *runOS)
		cmd := ergo("setup", *runOS)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

func removeSetupErgo() *exec.Cmd {
	runOS := getOS()
	if runOS != nil {
		return ergo("setup", *runOS, "-remove")

	}
	return nil
}

func cleanSetup() error {
	runOS := getOS()

	//if initialSetup is not set, then the initial setup failed and there is no point
	//trying to clean after it, anyway we can't since we do not have the initial info
	if initialSetup == nil {
		return nil
	}

	//we need to store the initial values to make sure they are restored
	//linux-gnome
	if *runOS == "linux-gnome" {
		s := initialSetup.(struct {
			mode          string
			autoconfigURL string
		})
		return clearLinuxSetup(s.mode, s.autoconfigURL)
	}
	//osx
	if *runOS == "osx" {
		s := initialSetup.(struct{ autoconfigURL string })
		return clearDarwinSetup(s.autoconfigURL)
	}
	//windows
	if *runOS == "windows" {
		s := initialSetup.(struct{ autoconfigURL string })
		err := clearWindowsSetup(s.autoconfigURL)
		if err != nil {
			return err
		}
		commands.InetRefresh()
	}
	return nil
}

func getCommandResult(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	rez, err := cmd.Output()

	return string(rez), err
}

func getLinuxGnomeProxyMode() (string, error) {
	return getCommandResult("gsettings", "get", "org.gnome.system.proxy", "mode")
}

func getLinuxGnomeProxyAutoConfig() (string, error) {
	return getCommandResult("gsettings", "get", "org.gnome.system.proxy", "autoconfig-url")
}

func clearLinuxSetup(mode string, autoconfigURL string) error {
	_, err := getCommandResult("gsettings", "set", "org.gnome.system.proxy", "mode", "'"+mode+"'")
	if err != nil {
		if err.Error() != "exit status 1" {
			return err
		} else {
			return nil
		}
	}
	_, err = getCommandResult("gsettings", "set", "org.gnome.system.proxy", "autoconfig-url", "'"+autoconfigURL+"'")

	return err
}

func getDarwinProxyAutoURL() (string, error) {
	return getCommandResult("sudo", "networksetup", "-getautoproxyurl", "\"Wi-Fi\"")
}

func clearDarwinSetup(autoconfigURL string) error {
	_, err := getCommandResult("sudo", "networksetup", "-setautoproxyurl", "\"Wi-Fi\"", "'"+autoconfigURL+"'")

	return err
}

func getWindowsProxyAutoURL() (string, error) {
	rez, err := getCommandResult("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "AutoconfigURL")

	//it is nromal to get an exit status 1 if we do not yet have the key in the registry
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}
	return rez, nil
}

func clearWindowsSetup(autoconfigURL string) error {
	var err error

	if autoconfigURL != "" {
		_, err = getCommandResult("reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "AutoConfigURL", "/f")
	} else {
		_, err = getCommandResult("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "AutoConfigURL", "/t", "REG_SZ", "/d", autoconfigURL, "/f")
	}

	return err
}
func TestSetupLinuxGnome(t *testing.T) {

	if *getOS() != "linux-gnome" {
		t.Skip("Not running linux-gnome setup specific tests")
	}
	err := setupErgo()
	if err != nil {
		if err.Error() != "exit status 1" {
			t.Fatalf("Could not perform setup for linux-gnome. Got %s", err.Error())
		}
		t.Skipf("Skipping test as gsettings does not exist on this system.")
	}

	defer cleanSetup()

	ac, err := getLinuxGnomeProxyAutoConfig()
	if err != nil {
		t.Fatalf("No error expected while getting AutoconfigURL from gsettings. Got: %s, %s\r\n", err.Error(), ac)
	}

	if !strings.Contains(ac, "http://127.0.0.1:2000/proxy.pac") {
		t.Fatalf("Expected to find \"http://127.0.0.1:2000/proxy.pac\" as part of the AutoconfigURL. Got \"%s\"\r\n", ac)
	}

	mode, err := getLinuxGnomeProxyMode()
	if err != nil {
		t.Fatalf("No error expected while getting \"mode\" from gsettings. Got: %s, %s\r\n", err.Error(), mode)
	}

	if strings.Contains(mode, "none") {
		t.Fatalf("Expected to find a mode different then \"none\". Got \"%s\"\r\n", mode)
	}

}

func TestSetupOSX(t *testing.T) {
	if *getOS() != "osx" {
		t.Skip("Not running osx setup specific tests")
	}
	err := setupErgo()
	if err != nil {
		t.Skipf("Please fix this for osx")
		//t.Fatalf("Could not perform setup for osx. Got %s", err.Error())
	}

	defer cleanSetup()

	ac, err := getDarwinProxyAutoURL()
	if err != nil {
		t.Fatalf("No error expected while getting AutoconfigURL from osx. Got: %s, %s\r\n", err.Error(), ac)
	}

	if !strings.Contains(ac, "http://127.0.0.1:2000/proxy.pac") {
		t.Fatalf("Expected to find \"http://127.0.0.1:2000/proxy.pac\" as part of the AutoconfigURL. Got \"%s\"\r\n", ac)
	}

}

func TestSetupWindows(t *testing.T) {
	if *getOS() != "windows" {
		t.Skip("Not running windows setup specific tests")
	}
	err := setupErgo()
	if err != nil {
		t.Fatalf("Could not perform setup for windows. Got %s", err.Error())
	}

	defer cleanSetup()

	rez, err := getWindowsProxyAutoURL()
	if err != nil {
		t.Fatalf("No error expected while getting AutoconfigURL from registry. Got: %s, %s\r\n", err.Error(), rez)
	}

	if !strings.Contains(rez, "http://127.0.0.1:2000/proxy.pac") {
		t.Fatalf("Expected to find \"http://127.0.0.1:2000/proxy.pac\" as part of the AutoconfigURL. Got \"%s\"\r\n", rez)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ergo test response")
}

//please change the ports if you change .ergo file
//the server will be closed when the test finishes
//not programatically
//Perhaps a programatic approach should be used later
func startTestWebServer() {
	http.HandleFunc("/", handler)
	go func() {
		http.ListenAndServe(":5000", nil)
	}()
}

//Please be aware that for windows you should have ports
//2000 : ergo proxy
// and
//5000 : test web server
//opened in
//the firewall for this test to work
func TestRunWindows(t *testing.T) {
	if *getOS() != "windows" {
		t.Skip("Not running windows run specific tests")
	}
	err := setupErgo()

	if err != nil {
		t.Fatalf("Could not perform setup for windows. Got %s", err.Error())
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")
	defer func() {
		if cmd == nil || cmd.ProcessState == nil {
			return
		}

		if !cmd.ProcessState.Exited() {
			cmd.Process.Kill()
		}
	}()

	go func() {
		cmd.Run()
	}()

	time.Sleep(2 * time.Second)

	//it would be much easier to use Invoke-WebRequest, but using this makes it compatible with
	//powershell 2.0 (windows 7)
	rez, err := getCommandResult("powershell", "-c", "$proxy=[System.Net.WebRequest]::GetSystemWebProxy();",
		"$webclient=new-object system.net.webclient;", "$webclient.proxy=$proxy;",
		"$webclient.DownloadString('http://bla.dev')")

	if err != nil {
		t.Fatalf("Expected no error while requesting http://bla.dev, Got %s\r\n", err.Error())
	}
	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}
}

func TestRunLinuxGnome(t *testing.T) {
	if *getOS() != "linux-gnome" {
		t.Skip("Not running linux-gnome run specific tests")
	}
	err := setupErgo()
	//exit status 1 means that gsettings was not found. So if linux is not really gnome, this will fails otherwise
	if err != nil && err.Error() != "exit status 1" {
		t.Fatalf("Could not perform setup for linux-gnome. Got %s", err.Error())
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")
	if cmd != nil {
		defer func() {
			if cmd != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				cmd.Process.Kill()
			}
		}()
	}
	go func() {
		cmd.Run()
	}()

	time.Sleep(2 * time.Second)

	//it would be better perhaps to use a brower like chrome-headless
	rez, err := getCommandResult("./testCurl.sh")

	if err != nil {
		t.Fatalf("Expected no error while requesting http://bla.dev, Got %s\r\n", err.Error())
	}
	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}
}

func TestRunOSX(t *testing.T) {
	if *getOS() != "osx" {
		t.Skip("Not running osx run specific tests")
	}
	err := setupErgo()
	if err != nil {
		t.Skip("Please fix this ... mac people")
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")
	if cmd != nil {
		defer func() {
			if cmd != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				cmd.Process.Kill()
			}
		}()
	}
	go func() {
		cmd.Run()
	}()

	time.Sleep(2 * time.Second)

	//it would be better perhaps to use a brower like chrome-headless
	rez, err := getCommandResult("./testCurl.sh")

	if err != nil {
		t.Fatalf("Expected no error while requesting http://bla.dev, Got %s\r\n", err.Error())
	}
	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}
}
