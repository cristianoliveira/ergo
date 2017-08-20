package proxy

import (
	"testing"
)

func TestWhenHasErgoFile(t *testing.T) {
	config := LoadConfig("../.ergo")

	t.Run("It loads the services redirections", func(t *testing.T) {
		expected := 2
		result := len(config.Services)

		if expected != result {
			t.Errorf("Expected %s got %s", expected, result)
		}
	})

	t.Run("It match the service host bla.dev", func(t *testing.T) {
		url := "bla.dev"
		result := config.GetService(url)

		if result == nil {
			t.Errorf("Expected result to not be nil", result)
		}
	})

	t.Run("It match the service host foo.dev", func(t *testing.T) {
		url := "foo.dev"
		result := config.GetService(url)

		if result == nil {
			t.Errorf("Expected result to not be nil", result)
		}
	})

	t.Run("It does not match the service host", func(t *testing.T) {
		url := "undefined.dev"
		result := config.GetService(url)

		if result != nil {
			t.Errorf("Expected result to be nil got: ", result)
		}
	})
}
