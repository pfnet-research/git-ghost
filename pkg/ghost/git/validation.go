package git

import (
	"bytes"
	"errors"
	"os/exec"
)

func ValidateGit() error {
	cmd := exec.Command("git", "version")
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}

func ValidateCommitish(commitish string) error {
	cmd := exec.Command("git", "cat-file", "-e", commitish)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}
