package infrastructure

import (
	"os/exec"
)

const (
	hugoCommand = "hugo"
)

// Hugo holds the Hugo engine functionality
type Hugo struct{}

// BuildSite builds an hugo project using the command line hugo tool. In the
// future a library use should be studied. Maybe this is better to be done when
// the HUGO project gets a more usable API.
func (hugo *Hugo) BuildSite(source string) error {
	arguments := []string{"--source", source}
	err := exec.Command(hugoCommand, arguments...).Run()
	if err != nil {
		return err
	}
	return nil
}
