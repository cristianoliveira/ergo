package setup

import (
	"os/exec"
	"strings"
)

// Runner is a process runner
type Runner interface {
	Run(string) error
}

// DefaultRunner implements the default runner for ergo
type DefaultRunner struct{}

// Run a given command script
func (r DefaultRunner) Run(command string) error {
	args := strings.Split(command, " ")
	return exec.Command(args[0], args[1:]...).Run()
}

// RunnerDefault run given commands
var RunnerDefault Runner = &DefaultRunner{}
