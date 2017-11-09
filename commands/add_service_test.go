package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestAddServiceAllreadyThere(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:9999"},
	})

	service := proxy.Service{Name: "test.dev"}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	AddService(&config, service)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service already present") {
		t.Fatalf("Expected AddService to refuse to add an existing service. Got %s.", out)
	}
}

func TestAddServiceAddOK(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:9999"},
	})

	service := proxy.Service{
		Name: "newtest.dev",
		URL:  "http://localhost:3333",
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	AddService(&config, service)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service added successfully") {
		t.Fatalf("Expected AddService add a service. Got %s.", out)
	}
}

func TestAddServiceAddFileNotFound(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:9999"},
	})

	service := proxy.Service{
		Name: "newtest.dev",
		URL:  "http://localhost:3333",
	}

	newConfig := proxy.Config{
		Services:   config.Services,
		ConfigFile: "anyfilethatdoesnotexist.here",
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	AddService(&newConfig, service)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Error while adding new service") {
		t.Fatalf("Expected AddService to refuse to add an service in an unexisting file. Got %s.", out)
	}
}
