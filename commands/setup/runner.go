package setup

import (
	"fmt"
	"os/exec"
)

// Runner is a process runner
type Runner interface {
	Run(string, ...string) ([]byte, error)
}

// DefaultRunner implements the default runner for ergo
type DefaultRunner struct{}

// Run a given command script
func (r DefaultRunner) Run(command string, args ...string) ([]byte, error) {
	out, err := exec.Command(command, args...).Output()
	if len(out) != 0 {
		fmt.Printf("Command: %s\n", out)
	}
	return out, err
}

// RunnerDefault run given commands
var RunnerDefault Runner = &DefaultRunner{}
