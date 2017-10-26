package proxy

import (
	"io/ioutil"
	"testing"
)

func TestWhenHasErgoFile(t *testing.T) {
	config := NewConfig()
	services, err := LoadServices("../.ergo")
	if err != nil {
		t.Fatal("could not load requied configuration file for tests")
	}

	config.Services = services

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

		if result == nil {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It match the service host foo.dev", func(t *testing.T) {
		result := config.GetService("foo.dev")

		if result == nil {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It match the service host withspaces.dev", func(t *testing.T) {
		result := config.GetService("withspaces.dev")

		if result == nil {
			t.Errorf("Expected result to not be nil")
		}
	})

	t.Run("It does not match the service host", func(t *testing.T) {
		result := config.GetService("undefined.dev")

		if result != nil {
			t.Errorf("Expected result to be nil got: %#v", result)
		}
	})

	t.Run("It does match the other protocols than http", func(t *testing.T) {
		if result := config.GetService("redis://redislocal.dev"); result == nil {
			t.Errorf("Expected  result to not be nil")
		}
	})

	t.Run("It match subdomains", func(tt *testing.T) {
		tt.Run("for one.domain.dev", func(tt *testing.T) {
			if result := config.GetService("one.domain.dev"); result == nil {
				tt.Errorf("Expected  result to not be nil")
			}
		})

		tt.Run("for two.domain.dev", func(tt *testing.T) {
			if result := config.GetService("two.domain.dev"); result == nil {
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

		service := NewService("testservice", "http://localhost:8080")

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
	})

	t.Run("It returns an error if there's an invalid declaration", func(tt *testing.T) {
		fileContent, err := ioutil.ReadFile("../.ergo")

		if err != nil {
			tt.Skip("Could not load initial .ergo file")
		}

		defer ioutil.WriteFile("../.ergo", fileContent, 0755)

		service := NewService("service-without-url", "")
		AddService("../.ergo", service)
		_, err = LoadServices("../.ergo")
		if err == nil {
			tt.Error("Expected LoadServices to fail")
		}
	})
}
