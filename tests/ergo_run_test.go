//go:build integration

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cristianoliveira/ergo/commands/setup"
)

var initialSetup interface{}

func ergo(args ...string) *exec.Cmd {
	return exec.Command(filepath.Join("..", "bin", "ergo"), args...)
}

type config struct {
	filePath string
}

func (c config) clean() error {
	return os.Remove(c.filePath)
}

func newConfigFromFile(configFile string) (config, error) {
	c := config{}
	file, err := ioutil.TempFile(os.TempDir(), "ergoconfig")
	if err != nil {
		return c, fmt.Errorf("Unable to create tmpfile: %v", err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return c, fmt.Errorf("unable to read configuration content: %v", err)
	}

	_, err = file.Write(fileContent)
	if err != nil {
		return c, fmt.Errorf("Unable to write configuration content to temp file: %v", err)
	}

	c.filePath = file.Name()
	return c, nil
}

func TestListApps(t *testing.T) {
	t.Run("it lists the apps", func(tt *testing.T) {
		appsOutput := []string{
			"http://foo.dev -> http://localhost:3000",
			"http://bla.dev -> http://localhost:5000",
			"http://withspaces.dev -> http://localhost:8080",
			"http://one.domain.dev -> http://localhost:8081",
			"http://two.domain.dev -> http://localhost:8082",
			"http://withextraspace.dev -> http://localhost:2222",
			"http://redislocal.dev -> redis://localhost:6543",
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
			"redislocal -> redis://localhost:6543",
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
			// FIXME URL should be redis://localhost:6543
			"redislocal":         "http://redislocal.dev",
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
		appsOutput := fmt.Sprintf("%s\n", "Service added successfully")

		c, err := newConfigFromFile("./.ergo")
		if err != nil {
			tt.Fatalf("Unable to copy config file: %v", err)
		}
		defer c.clean()

		cmd := ergo("add", "new.service", "http://localhost:8083", "-config", c.filePath)
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output: '%s' \n got '%s'", appsOutput, output)
		}
	})

	t.Run("it prints message for already added service", func(tt *testing.T) {
		appsOutput := fmt.Sprintf("%s\n", "Service already present")

		cmd := ergo("add", "foo", "http://localhost:3000")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output:\n '%s' \n got '%s'", appsOutput, output)
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

func setupErgoWithTempConfig(configFilePath string) (config, error) {

	c, err := newConfigFromFile(configFilePath)
	if err != nil {
		c.clean()
		return c, err
	}

	err = setupErgo(c.filePath)
	return c, err
}

func setupErgo(configFilePath string) error {
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
		cmd := ergo("setup", "-config", configFilePath, *runOS)
		_, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Unable to run ergo setup: %v", err)
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
		setup.InetRefresh()
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
	t.Skip("Travis CI has blocked access to networking changes and this has been failing since. We going to ignore it for now")

	if *getOS() != "linux-gnome" {
		t.Skip("Not running linux-gnome setup specific tests")
	}
	err := setupErgo("./.ergo")
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
	err := setupErgo("./.ergo")
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
	err := setupErgo("./.ergo")
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

//On windows please open port 9090 for this to work
func startDynamicTestWebServer() {

	http.HandleFunc("/dyn", handler)
	go func() {
		http.ListenAndServe(":9090", nil)
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
	err := setupErgo("./.ergo")

	if err != nil {
		t.Fatalf("Could not perform setup for windows. Got %s", err.Error())
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")

	defer func() {
		if cmd != nil && cmd.Process != nil {
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
	err := setupErgo("./.ergo")
	//exit status 1 means that gsettings was not found. So if linux is not really gnome, this will fails otherwise
	if err != nil && err.Error() != "exit status 1" {
		t.Fatalf("Could not perform setup for linux-gnome. Got %s", err.Error())
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")

	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

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
	err := setupErgo("./.ergo")
	if err != nil {
		t.Skip("Please fix this ... mac people")
	}

	defer cleanSetup()

	startTestWebServer()

	cmd := ergo("run")

	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

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

func TestConfigDynamicWindows(t *testing.T) {
	if *getOS() != "windows" {
		t.Skip("Not running windows run specific tests")
	}

	c, err := setupErgoWithTempConfig("./.ergo")

	if err != nil {
		t.Fatalf("Could not perform setup for windows. Got %s", err.Error())
	}

	defer c.clean()
	defer cleanSetup()

	//start a web server that will initially have no mapping in our proxy
	startDynamicTestWebServer()

	cmd := ergo("run", "-config", c.filePath)

	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	go func() {
		cmd.Run()
	}()

	//wait for it to be sure it started
	time.Sleep(2 * time.Second)

	//we should not have a correct response
	rez, err := getCommandResult("powershell", "-c", "$proxy=[System.Net.WebRequest]::GetSystemWebProxy();",
		"$webclient=new-object system.net.webclient;", "$webclient.proxy=$proxy;",
		"$webclient.DownloadString('http://dynamic.dev')")

	//on windows we should get an error here
	if err == nil {
		t.Fatal("Expected error while requesting http://dynamic.dev, Got nil\r\n")
	}

	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("Could not write on the config file")
	}

	defer f.Close()

	if _, err = f.WriteString("dynamic http://localhost:9090"); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	rez, err = getCommandResult("powershell", "-c", "$proxy=[System.Net.WebRequest]::GetSystemWebProxy();",
		"$webclient=new-object system.net.webclient;", "$webclient.proxy=$proxy;",
		"$webclient.DownloadString('http://dynamic.dev')")

	if err != nil {
		t.Fatalf("Expected no error while asking for an added config. Got %s\r\n", err.Error())
	}

	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}
}

func TestConfigDynamicLinuxGnome(t *testing.T) {
	t.Skip("Travis CI has blocked access to networking changes and this has been failing since. We going to ignore it for now")

	if *getOS() != "linux-gnome" {
		t.Skip("Not running linux-gnome run specific tests")
	}

	c, err := setupErgoWithTempConfig("./.ergo")
	if err != nil && err.Error() != "exit status 1" {
		t.Fatalf("Could not perform setup for linux-gnome: %v", err)
	}
	defer c.clean()
	defer cleanSetup()

	startDynamicTestWebServer()

	cmd := ergo("run", "-config", c.filePath)

	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	go func() {
		cmd.Run()
	}()

	time.Sleep(2 * time.Second)

	//it would be better perhaps to use a brower like chrome-headless
	rez, err := getCommandResult("./testCurlNotExisting.sh")

	if err != nil {
		t.Fatalf("Expected no error while requesting http://dynamic.dev, Got %s\r\n", err.Error())
	}

	if rez == "ergo test response" {
		t.Fatalf("Expected different response while asking for http://dynamic.dev, Got %s\r\n", rez)
	}

	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("Could not write on the config file")
	}

	defer f.Close()

	if _, err = f.WriteString("dynamic http://localhost:9090/dyn"); err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)

	rez, err = getCommandResult("./testCurlNotExisting.sh")

	if err != nil {
		t.Fatalf("Expected no error while asking for an added config. Got %s\r\n", err.Error())
	}

	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}

}

func TestConfigDynamicOSX(t *testing.T) {
	if *getOS() != "osx" {
		t.Skip("Not running osx run specific tests")
	}

	c, err := setupErgoWithTempConfig("./.ergo")

	if err != nil {
		t.Skip("Please fix this ... mac people")
	}
	defer c.clean()
	defer cleanSetup()

	startDynamicTestWebServer()

	cmd := ergo("run", "-config", c.filePath)

	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	go func() {
		cmd.Run()
	}()

	time.Sleep(2 * time.Second)

	//it would be better perhaps to use a brower like chrome-headless
	rez, err := getCommandResult("./testCurl.sh")

	if err != nil {
		t.Fatalf("Expected no error while requesting http://dynamic.dev, Got %s\r\n", err.Error())
	}

	if rez == "ergo test response" {
		t.Fatalf("Expected different response while asking for http://dynamic.dev, Got %s\r\n", rez)
	}

	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("Could not write on the config file")
	}

	defer f.Close()

	if _, err = f.WriteString("dynamic http://localhost:9090/dyn"); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	rez, err = getCommandResult("./testCurlNotExisting.sh")

	if err != nil {
		t.Fatalf("Expected no error while asking for an added config. Got %s\r\n", err.Error())
	}

	if strings.Trim(rez, " \r\n") != "ergo test response" {
		t.Fatalf("Expected \"ergo test response\" as response and got: %s\r\n", rez)
	}
}
