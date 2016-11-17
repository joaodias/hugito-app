package hugo

import (
	"os/exec"
)

const (
	hugoCommand = "hugo"
)

// Runner is something that runs. In our case the command executor that runs.
type Runner interface {
	Run() error
}

// Executor is something that executes. In this case an executor execute a
// command with given arguments.
type Executor interface {
	Execute(string, []string) Runner
}

// CommandExecutor is the command executor in the command line.
type CommandExecutor struct{}

// BuildSite builds an hugo project using the command line hugo tool. In the
// future a library use should be studied. Maybe this is better to be done when
// the HUGO project gets a more usable API.
func BuildSite(command string, source string, commandExecutor Executor) error {
	arguments := []string{"--source", source}
	runner := commandExecutor.Execute(command, arguments)
	err := runner.Run()
	if err != nil {
		return err
	}
	return nil
}

// Execute executes a command line command. It does not run it. It returns a
// runner that can run the command.
func (c *CommandExecutor) Execute(command string, arguments []string) Runner {
	return exec.Command(command, arguments...)
}
