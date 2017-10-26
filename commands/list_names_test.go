package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestListNames(t *testing.T) {

	config := buildConfig([]proxy.Service{
		{Name: "test.dev", URL: "localhost:9999"},
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	ListNames(&config)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Ergo Proxy current list:") {
		t.Fatalf("Expected ListNames to return something containing\"Ergo Proxy current list:\". Got %s.", out)
	}

	if !strings.Contains(out, "- test.dev -> localhost:9999") {
		t.Fatalf("Expected ListNames to return something containing\"- test.dev -> localhost:9999\". Got %s.", out)
	}
}
