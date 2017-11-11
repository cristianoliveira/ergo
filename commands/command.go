package commands

import "github.com/cristianoliveira/ergo/proxy"

// Command is the interface for a command in ergo
// returns a string with the result of the execution or a error
type Command interface {
	Execute(config *proxy.Config) (string, error)
}
