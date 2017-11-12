package commands

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestRunCommand(t *testing.T) {
	command := RunCommand{}

	t.Run("when domain config has wrong domain format", func(tt *testing.T) {
		config := buildConfig([]proxy.Service{
			{Name: "test.dev", URL: "localhost:9999"},
		})

		config.Domain = "foobar"

		out, err := command.Execute(&config)
		if err == nil {
			tt.Errorf("Expected error got none. Output: %s", out)
		}

		if !strings.Contains(err.Error(), "Domain has a wrong format") {
			tt.Errorf("Received error is different than expected. Error: %s", err)
		}
	})
}
