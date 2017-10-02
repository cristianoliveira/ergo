// +build integration

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
	"path/filepath"
)

func ergo(args ...string) *exec.Cmd {
	return exec.Command(filepath.Join("..","bin","ergo"), args...)
}

func TestListApps(t *testing.T) {	

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
}

func TestListAppNames(t *testing.T){

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
}
func TestShowUrlForName(t *testing.T){ 

	t.Run("it shows the url for a given name", func(tt *testing.T) {
		appsOutput := map[string]string{
			"foo":                "http://foo.dev",
			"bla":                "http://bla.dev",
			"withspaces":         "http://withspaces.dev",
			"one.domain":         "http://one.domain.dev",
			"two.domain":         "http://two.domain.dev",
			"redis://redislocal": "http://redis://redislocal.dev",
		}

		for name, url := range appsOutput {
			cmd := ergo("url", name)
			bs, err := cmd.Output()
			if err != nil {
				tt.Fatal(err)
			}

			output := string(bs)			
			if strings.Trim(output," \r\n")!=url {
				tt.Errorf("Expected output:\n [%s] \n got [%s]", url, strings.Trim(output," \r\n"))
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
