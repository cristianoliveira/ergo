package commands

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestList(t *testing.T) {

	config := buildConfig([]proxy.Service{
		{Name: "test.dev", URL: "localhost:9999"},
	})

	out, _ := ListCommand{}.Execute(&config)

	if !strings.Contains(out, "Ergo Proxy current list:") {
		t.Fatalf("Expected List to return something containing\"Ergo Proxy current list:\". Got %s.", out)
	}

	if !strings.Contains(out, "- http://test.dev -> localhost:9999") {
		t.Fatalf("Expected List to return something containing\"- http://test.dev -> localhost:9999\". Got %s.", out)
	}
}
