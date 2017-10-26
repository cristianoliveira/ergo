package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestRemoveOK(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:999"},
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

	RemoveService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service Removed") {
		t.Fatalf("Expected RemoveService to remove an existing service. Got %s", out)
	}
}

func TestRemoveFailIfNotExists(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:999"},
	})

	service := proxy.Service{Name: "doesntexist.dev"}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	RemoveService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service doesntexist.dev not found") {
		t.Fatalf("Expected RemoveService to fail removing non-existing service. Got %s", out)
	}
}

func TestRemoveOKUrl(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:999"},
	})

	service := proxy.Service{URL: "localhost:999"}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	RemoveService(&config, service, config.ConfigFile)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Service Removed") {
		t.Fatalf("Expected RemoveService to remove based on url. Got %s", out)
	}
}
