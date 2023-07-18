package commands

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestRemove(t *testing.T) {
	config := buildConfig([]proxy.Service{
		proxy.UnsafeNewService("test.dev", "localhost:999"),
		proxy.UnsafeNewService("test2.dev", "localhost:9292"),
	})

	t.Run("when remove service", func(tt *testing.T) {
		service := proxy.UnsafeNewService("test.dev", "")

		command := RemoveServiceCommand{Service: service}
		out, err := command.Execute(config)
		if err != nil {
			t.Fatalf("Expected no error got: %s", err)
		}

		if !strings.Contains(out, "Service Removed") {
			t.Fatalf("Expected RemoveService to remove an existing service. Got %s", out)
		}
	})

	t.Run("when service not found", func(tt *testing.T) {
		service := proxy.UnsafeNewService("doesntexist.dev", "")

		command := RemoveServiceCommand{Service: service}
		_, err := command.Execute(config)
		if err == nil {
			t.Fatalf("Expected error got: %s", err)
		}
	})

	t.Run("when config file not found", func(tt *testing.T) {
		service := proxy.UnsafeNewService("test.dev", "")
		config.ConfigFile = "undefined"

		command := RemoveServiceCommand{Service: service}
		_, err := command.Execute(config)
		if err == nil {
			t.Fatalf("Expected error got: %s", err)
		}
	})
}
