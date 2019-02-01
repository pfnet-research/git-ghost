package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
)

func nonEmpty(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty", name)
	}
	return nil
}

func isValidComittish(name, comittish string) error {
	err := git.ValidateRefspec(globalOpts.srcDir, comittish)
	if err != nil {
		return fmt.Errorf("%s is not a valid object", name)
	}
	return nil
}
