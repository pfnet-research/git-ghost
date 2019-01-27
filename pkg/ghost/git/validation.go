package git

import (
	"os/exec"
)

func ValidateGit() error {
	gitCmd := exec.Command("git", "version")
	err := gitCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ValidateCommitish(commitish string) error {
	gitCmd := exec.Command("git", "cat-file", "-e", commitish)
	err := gitCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
