package proxy

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWhenHasErgoFile(t *testing.T) {
	config := NewConfig()
	config.ConfigFile = "../.ergo"
	err := config.LoadServices()
	if err != nil {
		t.Fatal("could not load required configuration file for tests")
	}

	t.Run("Wildcard feature", func(t *testing.T) {
		cases := []struct {
			host     string
			expected bool
		}{
			{"foobar.wildcard.dev", true},
			{"other.wildcard.dev", true},
		}

		for _, c := range cases {
			t.Run("It matches the service host with wildcard "+c.host, func(t *testing.T) {
				service, err := config.GetService(c.host)
				if err != nil {
					t.Errorf("Expected no error. Got: %s", err)
				}

				if service.Empty() {
					t.Errorf("Expected service to not be nil")
				}
			})
		}
	})

	t.Run("It loads the services redirections", func(t *testing.T) {
		expected := 11
		service := len(config.Services)

		if expected != service {
			t.Errorf("Expected to get %d services, but got %d\r\n%q",
				expected, service, config.Services)
		}
	})

	t.Run("It match the service host mysite", func(t *testing.T) {
		service, err := config.GetService("mysite.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It match the service host mylocalsite", func(t *testing.T) {
		service, err := config.GetService("mylocalsite.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It match the service host bla.dev", func(t *testing.T) {
		service, err := config.GetService("bla.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It match the service host foo.dev", func(t *testing.T) {
		service, err := config.GetService("foo.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It match the service host withspaces.dev", func(t *testing.T) {
		service, err := config.GetService("withspaces.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It does not match the service host", func(t *testing.T) {
		service, err := config.GetService("undefined.dev")
		if service != nil {
			t.Errorf("Expected service to be nil got: %#v", service)
		}

		if err == nil {
			t.Errorf("Expected error. Got none")
		}
	})

	t.Run("It does match other protocols than http", func(t *testing.T) {
		service, err := config.GetService("redislocal.dev")
		if err != nil {
			t.Errorf("Expected no error. Got: %s", err)
		}

		if service.Empty() {
			t.Errorf("Expected service to not be nil")
		}
	})

	t.Run("It match subdomains", func(tt *testing.T) {
		tt.Run("for one.domain.dev", func(tt *testing.T) {
			service, err := config.GetService("one.domain.dev")
			if err != nil {
				t.Errorf("Expected no error. Got: %s", err)
			}
			if service.Empty() {
				tt.Errorf("Expected service to not be nil")
			}
		})

		tt.Run("for two.domain.dev", func(tt *testing.T) {
			service, err := config.GetService("two.domain.dev")
			if err != nil {
				t.Errorf("Expected no error. Got: %s", err)
			}

			if service.Empty() {
				tt.Errorf("Expected service to not be nil")
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

		service := UnsafeNewService("testservice", "http://localhost:8080")

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

		service := UnsafeNewService("servicetoberemoved", "http://localhost:8083")

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

		expected := []byte(`http://localhost:80 mysite
			http://localhost:8282 mylocalsite
			http://localhost:8082  two.domain
			foo http://localhost:3000
			bla http://localhost:5000
			withspaces       http://localhost:8080
			one.domain       http://localhost:8081
			two.domain       http://localhost:8082
			redislocal       redis://localhost:6543
			*.wildcard http://localhost:4000
			withextraspace http://localhost:2222  
		`)

		if !bytes.Equal(expected, newFileContent) {
			tt.Errorf("Expected only to remove servicetoberemoved. Got %s\n", newFileContent)
		}
	})

	t.Run("It fails to remove service with invalid path", func(tt *testing.T) {
		service := UnsafeNewService("testservice", "http://localhost:8080")

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

		service := UnsafeNewService("service-without-url", "")
		AddService(config.ConfigFile, service)
		err = config.LoadServices()
		if err == nil {
			tt.Error("Expected LoadServices to fail")
		}
	})
}

func TestChangingByEnvironmentVariable(t *testing.T) {
	cases := []struct {
		title       string
		varName     string
		value       string
		expectation func(config *Config) bool
	}{
		{
			title:   "Variable env " + PortEnv + " changes port config",
			varName: PortEnv,
			value:   "3000",
			expectation: func(config *Config) bool {
				return config.Port == "3000"
			},
		},
		{
			title:   "Variable env " + DomainEnv + " changes port config",
			varName: DomainEnv,
			value:   ".new",
			expectation: func(config *Config) bool {
				return config.Domain == ".new"
			},
		},
		{
			title:   "Variable env " + VerboseEnv + " changes port config",
			varName: VerboseEnv,
			value:   "1",
			expectation: func(config *Config) bool {
				return config.Verbose == true
			},
		},
		{
			title:   "Variable env " + ConfigFileEnv + " changes port config",
			varName: ConfigFileEnv,
			value:   "/tmp/dev.txt",
			expectation: func(config *Config) bool {
				return config.ConfigFile == "/tmp/dev.txt"
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(tt *testing.T) {
			os.Setenv(c.varName, c.value)

			config := NewConfig()

			if !c.expectation(config) {
				tt.Fatal("Expected config to be changed", config)
			}

			os.Setenv(c.varName, "")
		})
	}
}

func TestOverrideBy(t *testing.T) {
	cases := []struct {
		title       string
		newConfig   *Config
		expectation func(config *Config) bool
	}{
		{
			title: "It changes only the default value domain",
			newConfig: &Config{
				Domain: ".foo",
			},
			expectation: func(config *Config) bool {
				return config.Domain == ".foo" ||
					config.Port != PortDefault ||
					config.ConfigFile != ConfigFilePathDefault ||
					config.Verbose != false
			},
		},
		{
			title: "It changes only the default value port",
			newConfig: &Config{
				Port: "2111",
			},
			expectation: func(config *Config) bool {
				return config.Port == "2111" ||
					config.Domain != DomainDefault ||
					config.ConfigFile != ConfigFilePathDefault ||
					config.Verbose != false
			},
		},
		{
			title: "It changes only the default value config file",
			newConfig: &Config{
				ConfigFile: "/tmp/ergo",
			},
			expectation: func(config *Config) bool {
				return config.ConfigFile == "/tmp/ergo" ||
					config.Domain != DomainDefault ||
					config.Port != PortDefault ||
					config.Verbose != false
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(tt *testing.T) {
			config := NewConfig()
			config.OverrideBy(c.newConfig)
			if !c.expectation(config) {
				tt.Fatal("Expected config to be changed", config)
			}
		})
	}
}
