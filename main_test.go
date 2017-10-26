package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testExec = os.Args[0]

func TestMain(t *testing.T) {
	t.Run("it shows usage", func(tt *testing.T) {
		args := []string{"ergo", "-h"}
		os.Args = args

		main()
		// Output: USAGE
	})
}

// TestMissingCommand will cause the main function to return a non-zero exit code
// so we need to do some wrapping to make the test not fail.
// For more info check ou to uset: https://talks.golang.org/2014/testing.slide#23
func TestMissingCommand(t *testing.T) {
	t.Run("missing a command so it shows usage", func(tt *testing.T) {
		if os.Getenv("TEST_MISSING_COMMAND") == "1" {
			args := []string{"ergo"}
			os.Args = args

			main()
			// Output: USAGE and exit with an error code
		}

		cmd := exec.Command(testExec, "-test.run=TestMissingCommand")
		cmd.Env = append(os.Environ(), "TEST_MISSING_COMMAND=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
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

func TestRemoveCommand(t *testing.T) {
	t.Run("it returns no command", func(tt *testing.T) {
		args := []string{"ergo", "remove"}
		os.Args = args

		result := command()
		if result == nil {
			t.Errorf("Expected result not to be nil")
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		outC := make(chan string)
		var buf bytes.Buffer

		go func() {
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		result()

		w.Close()

		os.Stdout = old

		out := <-outC

		if !strings.Contains(out, "Usage: ergo remove <name|url>") {
			t.Fatalf("Expected Remove to show usage. Got %s", out)
		}
	})
}
