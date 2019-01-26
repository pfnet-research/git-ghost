package git

import (
	"os/exec"
)

// For test mocks
var (
	ValidateGit       = validateGit
	ValidateCommitish = validateCommitish
)

func validateGit() error {
	gitCmd := exec.Command("git", "version")
	err := gitCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func validateCommitish(commitish string) error {
	gitCmd := exec.Command("git", "cat-file", "-e", commitish)
	err := gitCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
