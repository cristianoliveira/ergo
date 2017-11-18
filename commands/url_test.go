package commands

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestURLCommand(t *testing.T) {
	config := buildConfig([]proxy.Service{
		{Name: "test.dev", URL: "localhost:9999"},
	})

	t.Run("when found the service", func(tt *testing.T) {
		command := URLCommand{FilterName: "test.dev"}

		out, err := command.Execute(config)

		if err != nil {
			t.Fatalf("Expected no error. Got: %s", err)
		}

		if !strings.Contains(out, "http://test.dev") {
			t.Fatalf("Expected URL to return something containing\"http://test.dev\". Got %s.", out)
		}
	})

	t.Run("when doesnt found the service", func(tt *testing.T) {
		command := URLCommand{FilterName: "undefined"}

		_, err := command.Execute(config)

		if err == nil {
			t.Fatalf("Expected error. Got: %s", err)
		}
	})
}
