package main

import (
	"os"
	"testing"
)

func TestShowingUsage(t *testing.T) {
	t.Run("it shows usage when missing command", func(tt *testing.T) {
		os.Args = []string{"ergo"}

		cmd := command()
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})

	t.Run("it shows usage when pass h flag", func(tt *testing.T) {
		os.Args = []string{"ergo", "-h"}

		cmd := command()
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})

	t.Run("it shows usage when unknown command", func(tt *testing.T) {
		os.Args = []string{"ergo", "foobar"}

		cmd := command()
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})
}

func TestShowingVersion(t *testing.T) {
	t.Run("it shows usage when missing command", func(tt *testing.T) {
		os.Args = []string{"ergo", "-v"}

		cmd := command()
		if cmd == nil {
			t.Errorf("Expected cmd to not be nil")
		}

		cmd()
		// Output: USAGE
	})
}

func TestListCommand(t *testing.T) {
	t.Run("it is list command", func(tt *testing.T) {
		args := []string{"ergo", "list"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result to not be nil")
		}

		result()
	})
}

func TestListNamesCommand(t *testing.T) {
	t.Run("it is list-names command", func(tt *testing.T) {
		args := []string{"ergo", "list-names"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result to not be nil")
		}

		result()
	})
}

func TestSetupCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "setup"}
		os.Args = args

		result := command()
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is setup command", func(tt *testing.T) {
		args := []string{"ergo", "setup", "osx"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}
	})
}

func TestUrlCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "url"}
		os.Args = args

		result := command()
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "url", "foo"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}

		result()
	})
}

func TestRunCommand(t *testing.T) {
	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "run"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}
	})
}

func TestAddCommand(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "add"}
		os.Args = args

		result := command()
		if result != nil {
			t.Errorf("Expected result to be nil")
		}
	})

	t.Run("it is url command", func(tt *testing.T) {
		args := []string{"ergo", "add", "foo", "bar"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}

		result()
	})
}
