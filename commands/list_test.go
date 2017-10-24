package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestList(t *testing.T) {

	config := buildConfig([]proxy.Service{
		proxy.Service{Name: "test.dev", URL: "localhost:9999"},
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

	List(&config)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "Ergo Proxy current list:") {
		t.Fatalf("Expected List to return something containing\"Ergo Proxy current list:\". Got %s.", out)
	}

	if !strings.Contains(out, "- http://test.dev -> localhost:9999") {
		t.Fatalf("Expected List to return something containing\"- http://test.dev -> localhost:9999\". Got %s.", out)
	}
}
