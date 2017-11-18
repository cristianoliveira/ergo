package commands

import (
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestAddServiceAllreadyThere(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test", URL: "localhost:9999"},
	})

	service := proxy.Service{Name: "test"}

	command := AddServiceCommand{Service: service}
	result, err := command.Execute(config)
	if err == nil {
		t.Fatalf("Expected to receive error. Result: %s", result)
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

	command := AddServiceCommand{Service: service}
	result, err := command.Execute(config)
	if err != nil {
		t.Fatalf("Expected to not receive error. Got: %s", err)
	}

	if result != "Service added successfully" {
		t.Fatalf("Expected AddServiceCommand to add service. Got %s.", result)
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

	config.ConfigFile = "anyfilethatdoesnotexist.here"

	command := AddServiceCommand{Service: service}
	result, err := command.Execute(config)
	if err == nil {
		t.Fatalf("Expected to not receive error. Got: %s", result)
	}
}
