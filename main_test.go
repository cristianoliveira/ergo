package main

import (
	"testing"
)

func TestShowingUsage(t *testing.T) {
	t.Run("it shows usage when missing command", func(tt *testing.T) {
		args := []string{"ergo"}

		cmd := command(args)
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})

	t.Run("it shows usage when pass h flag", func(tt *testing.T) {
		args := []string{"ergo", "-h"}

		cmd := command(args)
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})

	t.Run("it shows usage when unknown command", func(tt *testing.T) {
		args := []string{"ergo", "foobar"}

		cmd := command(args)
		if cmd != nil {
			t.Errorf("Expected cmd to be nil")
		}
		// Output: USAGE
	})
}

func TestShowingVersion(t *testing.T) {
	t.Run("it shows usage when missing command", func(tt *testing.T) {
		args := []string{"ergo", "-v"}

		cmd := command(args)
		if cmd == nil {
			t.Errorf("Expected cmd to not be nil")
		}

		cmd()
		// Output: USAGE
	})
}

// TestMissingCommand will cause the main function to return a non-zero exit code
// so we need to do some wrapping to make the test not fail.
// For more info check ou to uset: https://talks.golang.org/2014/testing.slide#23
func TestMissingCommand(t *testing.T) {
	t.Run("missing a command so it shows usage", func(tt *testing.T) {
		args := []string{"ergo"}

		cmd := command(args)

		if cmd != nil {
			t.Errorf("Expected cmd to not be nil")
		}

		// Output: USAGE and exit with an error code
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
