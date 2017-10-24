package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/ergo/proxy"
)

func TestURL(t *testing.T) {

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

	URL("test.dev", &config)

	w.Close()

	os.Stdout = old

	out := <-outC

	if !strings.Contains(out, "http://test.dev") {
		t.Fatalf("Expected URL to return something containing\"http://test.dev\". Got %s.", out)
	}
}
