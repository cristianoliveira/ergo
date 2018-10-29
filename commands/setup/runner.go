package setup

import (
	"os/exec"
)

// Runner is a process runner
type Runner interface {
	Run(string, string) error
}

// DefaultRunner implements the default runner for ergo
type DefaultRunner struct{}

// Run a given command script
func (r DefaultRunner) Run(command, args string) error {
	return exec.Command(command, args).Run()
}

// RunnerDefault run given commands
var RunnerDefault Runner = &DefaultRunner{}
