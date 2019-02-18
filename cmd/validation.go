package cmd

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util/errors"
)

func nonEmpty(name, value string) errors.GitGhostError {
	if value == "" {
		return errors.Errorf("%s must not be empty", name)
	}
	return nil
}

func isValidComittish(name, comittish string) errors.GitGhostError {
	err := git.ValidateComittish(globalOpts.srcDir, comittish)
	if err != nil {
		return errors.Errorf("%s is not a valid object", name)
	}
	return nil
}
