package mocks

import (
	"errors"
	"github.com/joaodias/hugito-app/hugo"
)

type CommandRunner struct {
	IsError bool
}

type CommandExecutor struct {
	IsError bool
}

func (c *CommandExecutor) Execute(cmd string, args []string) hugo.Runner {
	if c.IsError {
		return &CommandRunner{
			IsError: true,
		}
	}
	return &CommandRunner{
		IsError: false,
	}
}

func (c *CommandRunner) Run() error {
	if c.IsError {
		return errors.New("Error building HUGO site.")
	}
	return nil
}
