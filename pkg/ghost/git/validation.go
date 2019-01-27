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

func ValidateRefspec(dir, refspec string) error {
	cmd := exec.Command("git", "-C", dir, "cat-file", "-e", refspec)
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
