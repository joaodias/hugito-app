package mocks

import (
	"errors"
)

type BuildEngine struct {
	IsError         bool
	BuildSiteCalled bool
}

func (be *BuildEngine) BuildSite(source string) error {
	be.BuildSiteCalled = true
	if be.IsError {
		return errors.New("Some error")
	}
	return nil
}
