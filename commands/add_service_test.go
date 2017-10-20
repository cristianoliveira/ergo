package commands

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestAddServiceAllreadyThere(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("Error writing to temporary file: %s", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temp file: %s", err.Error())
	}

	if err != nil {
		t.Fatalf("No error expected while initializing config file. Got %s.", err.Error())
	}
	config := proxy.Config{}
	config.ConfigFile = tmpfile.Name()
	config.Services, err = proxy.LoadServices(config.ConfigFile)

	if err != nil {
		t.Fatalf("No error expected while loading services from config file. Got %s.", err.Error())
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

	AddService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service already present") {
		t.Fatalf("Expected AddService to refuse to add an existing service. Got %s.", out)
	}
}

func TestAddServiceAddOK(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("Error writing to temporary file: %s", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temp file: %s", err.Error())
	}

	if err != nil {
		t.Fatalf("No error expected while initializing config file. Got %s.", err.Error())
	}
	config := proxy.Config{}
	config.ConfigFile = tmpfile.Name()
	config.Services, err = proxy.LoadServices(config.ConfigFile)

	if err != nil {
		t.Fatalf("No error expected while loading services from config file. Got %s.", err.Error())
	}

	service := proxy.Service{}
	service.Name = "newservice"
	service.URL = "http://localhost:3333"

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	AddService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service added successfully") {
		t.Fatalf("Expected AddService add a service. Got %s.", out)
	}
}

func TestAddServiceFileNOK(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "testaddservice")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write([]byte("test.dev localhost:9999")); err != nil {
		t.Fatalf("Error writing to temporary file: %s", err.Error())
	}

	if err = tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temp file: %s", err.Error())
	}

	if err != nil {
		t.Fatalf("No error expected while initializing config file. Got %s.", err.Error())
	}
	config := proxy.Config{}
	config.ConfigFile = tmpfile.Name()
	config.Services, err = proxy.LoadServices(config.ConfigFile)

	if err != nil {
		t.Fatalf("No error expected while loading services from config file. Got %s.", err.Error())
	}

	config.ConfigFile = "anyfilethatdoesnotexist.here"

	service := proxy.Service{}
	service.Name = "newservice"
	service.URL = "http://localhost:3333"

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	AddService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Error while adding new service") {
		t.Fatalf("Expected AddService to refuse to add an service in an unexisting file. Got %s.", out)
	}
}
