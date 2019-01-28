package git

import (
	"bytes"
	"errors"
	"os/exec"
)

func ValidateRemoteBranchExistence(repo, branch string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", repo, branch)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	output, err := cmd.Output()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return false, errors.New(s)
		}
		return false, err
	}
	return string(output) != "", nil
}
