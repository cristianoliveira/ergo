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

// TestInvalidRunCommandOptions will cause the main function to return a non-zero exit code
// so we need to do some wrapping to make the test not fail.
// For more info check ou to uset: https://talks.golang.org/2014/testing.slide#23
func TestInvalidRunCommandOptions(t *testing.T) {
	t.Run("invalid run command options so it shows usage", func(tt *testing.T) {
		if os.Getenv("TEST_INVALID_RUN_OPTION") == "1" {
			args := []string{"ergo", "run", "-domain=foo"}
			os.Args = args

			main()
			// Output: USAGE and exit with an error code
		}

		cmd := exec.Command(testExec, "-test.run=TestInvalidRunCommandOptions")
		cmd.Env = append(os.Environ(), "TEST_INVALID_RUN_OPTION=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
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
