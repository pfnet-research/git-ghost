package util

import (
	"git-ghost/pkg/util/errors"
	"os/exec"
	"strings"
)

func GenerateFileContentHash(filepath string) (string, errors.GitGhostError) {
	// TODO: Use appropriate hash algorithm
	cmd := exec.Command("sha1sum", "-b", filepath)
	output, err := cmd.Output()
	if err != nil {
		return "", errors.WithStack(err)
	}
	hash := strings.Split(string(output), " ")[0]
	return hash, nil
}
