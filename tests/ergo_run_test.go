package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func ergo(args ...string) *exec.Cmd {
	return exec.Command("../bin/ergo", args...)
}

func integrationDisabled() bool {
	value, found := os.LookupEnv("INTEGRATION_TEST")
	if found == false || value == "false" {
		return true
	}

	return false
}
func TestIntegration(t *testing.T) {
	// TODO: create tests for windows
	if runtime.GOOS == "windows" {
		t.Skipf("skipping test on %q", runtime.GOOS)
	}

	if integrationDisabled() {
		t.Skip("skipping integration tests - run with INTEGRATION_TEST=true to include")
	}

	t.Run("it lists the apps", func(tt *testing.T) {
		appsOutput := []string{
			"http://foo.dev -> http://localhost:3000",
			"http://bla.dev -> http://localhost:5000",
			"http://withspaces.dev -> http://localhost:8080",
			"http://one.domain.dev -> http://localhost:8081",
			"http://two.domain.dev -> http://localhost:8082",
			"http://redis://redislocal.dev -> redis://localhost:6543",
		}

		cmd := ergo("list")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		for _, app := range appsOutput {
			if !strings.Contains(output, app) {
				tt.Errorf("Expected output:\n %s \n got %s", output, app)
			}
		}
	})

	t.Run("it lists the app names", func(tt *testing.T) {
		appsOutput := []string{
			"foo -> http://localhost:3000",
			"bla -> http://localhost:5000",
			"withspaces -> http://localhost:8080",
			"one.domain -> http://localhost:8081",
			"two.domain -> http://localhost:8082",
			"redis://redislocal -> redis://localhost:6543",
		}

		cmd := ergo("list-names", "foo")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		for _, app := range appsOutput {
			if !strings.Contains(output, app) {
				tt.Errorf("Expected output:\n %s \n got %s", output, app)
			}
		}
	})

	t.Run("it shows the url for a given name", func(tt *testing.T) {
		appsOutput := map[string]string{
			"foo":                "http://localhost:3000",
			"bla":                "http://localhost:5000",
			"withspaces":         "http://localhost:8080",
			"one.domain":         "http://localhost:8081",
			"two.domain":         "http://localhost:8082",
			"redis://redislocal": "redis://localhost:6543",
		}

		for name, url := range appsOutput {
			cmd := ergo("list-names", "foo")
			bs, err := cmd.Output()
			if err != nil {
				tt.Fatal(err)
			}

			output := string(bs)
			if !strings.Contains(output, url) {
				tt.Errorf("Expected output:\n %s \n got %s", output, name)
			}
		}
	})

	t.Run("it adds new service if not present", func(tt *testing.T) {
		appsOutput := fmt.Sprintf("%s\n", "Service added successfully!")

		cmd := ergo("add", "new.service", "http://localhost:8083")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output:\n %s \n got %s", appsOutput, output)
		}
	})

	t.Run("it prints message for already added service", func(tt *testing.T) {
		appsOutput := fmt.Sprintf("%s\n", "Service already present!")

		cmd := ergo("add", "foo", "http://localhost:3000")
		bs, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		output := string(bs)

		if strings.Compare(output, appsOutput) != 0 {
			tt.Errorf("Expected output:\n %s \n got %s", appsOutput, output)
		}
	})
	// TODO: Add tests for server
	//
	// t.Run("it runs binding the sites", func(tt *testing.T) {
	// 	cmd := ergo("run", "-p", "25000")
	// 	defer cmd.Process.Kill()
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		tt.Fatal(err)
	// 	}
	// })
}
