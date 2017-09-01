package proxy

import (
	"testing"
)

func TestWhenHasErgoFile(t *testing.T) {
	config := NewConfig()
	config.Services = LoadConfig("../.ergo")

	t.Run("It loads the services redirections", func(t *testing.T) {
		expected := 6
		result := len(config.Services)

		if expected != result {
			t.Errorf("Expected %s got %s", expected, result)
		}
	})

	t.Run("It match the service host bla.dev", func(t *testing.T) {
		result := config.GetService("bla.dev")

		if result == nil {
			t.Errorf("Expected result to not be nil", result)
		}
	})

	t.Run("It match the service host foo.dev", func(t *testing.T) {
		result := config.GetService("foo.dev")

		if result == nil {
			t.Errorf("Expected result to not be nil", result)
		}
	})

	t.Run("It match the service host withspaces.dev", func(t *testing.T) {
		result := config.GetService("withspaces.dev")

		if result == nil {
			t.Errorf("Expected result to not be nil", result)
		}
	})

	t.Run("It does not match the service host", func(t *testing.T) {
		result := config.GetService("undefined.dev")

		if result != nil {
			t.Errorf("Expected result to be nil got: ", result)
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
}
