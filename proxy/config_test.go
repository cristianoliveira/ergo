package proxy

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestWhenHasErgoFile(t *testing.T) {
	config := NewConfig()
	config.ConfigFile = "../.ergo"
	err := config.LoadServices()
	if err != nil {
		t.Fatal("could not load required configuration file for tests")
	}

	t.Run("It loads the services redirections", func(t *testing.T) {
		expected := 6
		result := len(config.Services)

		if expected != result {
			t.Errorf("Expected to get %d services, but got %d\r\n%q",
				expected, result, config.Services)
		}
	})

	t.Run("It match the service host bla.dev", func(t *testing.T) {
		result := config.GetService("bla.dev")

		if result.Empty() {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It match the service host foo.dev", func(t *testing.T) {
		result := config.GetService("foo.dev")

		if result.Empty() {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It match the service host withspaces.dev", func(t *testing.T) {
		result := config.GetService("withspaces.dev")

		if result.Empty() {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It does not match the service host", func(t *testing.T) {
		result := config.GetService("undefined.dev")

		if !result.Empty() {
			t.Errorf("Expected result to be nil got: %#v", result)
		}
	})

	t.Run("It does match the other protocols than http", func(t *testing.T) {
		if result := config.GetService("redis://redislocal.dev"); result.Empty() {
			t.Errorf("Expected  result to not be nil")
		}
	})

	t.Run("It match subdomains", func(tt *testing.T) {
		tt.Run("for one.domain.dev", func(tt *testing.T) {
			if result := config.GetService("one.domain.dev"); result.Empty() {
				tt.Errorf("Expected  result to not be nil")
			}
		})

		tt.Run("for two.domain.dev", func(tt *testing.T) {
			if result := config.GetService("two.domain.dev"); result.Empty() {
				tt.Errorf("Expected  result to not be nil")
			}
		})
	})

	t.Run("It adds new service", func(tt *testing.T) {
		fileContent, err := ioutil.ReadFile("../.ergo")

		if err != nil {
			tt.Skipf("Could not load initial .ergo file")
		}

		//we clean after the test. Otherwise the next test will fail
		defer ioutil.WriteFile("../.ergo", fileContent, 0755)

		service := NewService("testservice", "http://localhost:8080")

		if err := AddService("../.ergo", service); err != nil {
			tt.Errorf("Expected service to be added")
		}

	})

	t.Run("It removes a service", func(tt *testing.T) {
		fileContent, err := ioutil.ReadFile("../ergo")

		if err != nil {
			tt.Skipf("Could not load initial .ergo file")
		}

		//we clean after the test. Otherwise the next test will fail
		defer ioutil.WriteFile("../.ergo", fileContent, 0755)

		service := NewService("servicetoberemoved", "http://localhost:8083")

		if err := RemoveService("../.ergo", service); err != nil {
			tt.Errorf("Expected no error while removing service. Got %v\n", err)
		}

		newFileContent, err := ioutil.ReadFile("../ergo")
		if err != nil {
			tt.Skip("Could not load initial .ergo file")
		}

		if !(len(fileContent) > len(newFileContent)) {
			tt.Errorf("Expected service to be removed")
		}

		expected := []byte(`foo http://localhost:3000
		bla http://localhost:5000
		withspaces       http://localhost:8080
		one.domain       http://localhost:8081
		two.domain       http://localhost:8082
		redis://redislocal       redis://localhost:6543
		`)

		if !bytes.Equal(expected, newFileContent) {
			tt.Errorf("Expected only to remove servicetoberemoved. Got %s\n", newFileContent)
		}
	})

	t.Run("It fails to remove service with invalid path", func(tt *testing.T) {

		service := NewService("testservice", "http://localhost:8080")

		if err := RemoveService("foobarinvalid", service); err == nil {
			tt.Errorf("Expected failure to read invalid path.\n")
		}
	})

	t.Run("It returns an error if there's an invalid declaration", func(tt *testing.T) {
		fileContent, err := ioutil.ReadFile("../.ergo")

		if err != nil {
			tt.Skip("Could not load initial .ergo file")
		}

		defer ioutil.WriteFile("../.ergo", fileContent, 0755)

		service := NewService("service-without-url", "")
		AddService(config.ConfigFile, service)
		err = config.LoadServices()
		if err == nil {
			tt.Error("Expected LoadServices to fail")
		}
	})
}
