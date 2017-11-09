package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "-h"}
		os.Args = args

		main()
		// Output: USAGE
	})
}

func TestListCommand(t *testing.T) {
	t.Run("it is list command", func(tt *testing.T) {
		args := []string{"ergo", "list"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result to not be nil")
		}

		result()
	})
}

func TestListNamesCommand(t *testing.T) {
	t.Run("it is list-names command", func(tt *testing.T) {
		args := []string{"ergo", "list-names"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result to not be nil")
		}

		result()
	})
}

func TestSetupCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "setup"}

		result := command(args)
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is setup command", func(tt *testing.T) {
		args := []string{"ergo", "setup", "osx"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}
	})
}

func TestUrlCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "url"}

		result := command(args)
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "url", "foo"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}

		result()
	})
}

func TestRunCommand(t *testing.T) {
	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "run"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}
	})
}

// TestInvalidRunCommandOptions will cause the main function to return a non-zero exit code
// so we need to do some wrapping to make the test not fail.
// For more info check ou to uset: https://talks.golang.org/2014/testing.slide#23
func TestInvalidRunCommandOptions(t *testing.T) {
	t.Run("invalid run command options so it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "run", "-domain=foo"}

		if command(args) == nil {
			t.Errorf("Expected result to not be nil")
		}
	})
}

func TestAddCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "add"}

		result := command(args)
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "add", "foo", "bar"}

		result := command(args)
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}

		result()
	})
}
