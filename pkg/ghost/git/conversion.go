package git

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func ResolveRefspec(dir, refspec string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-list", "-1", refspec)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	commit, err := cmd.Output()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return "", errors.New(s)
		}
		return "", err
	}
	return strings.TrimRight(string(commit), "\r\n"), nil
}
