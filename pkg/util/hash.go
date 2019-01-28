package util

import (
	"os/exec"
	"strings"
)

func GenerateFileContentHash(filepath string) (string, error) {
	// TODO: Use appropriate hash algorithm
	cmd := exec.Command("sha1sum", "-b", filepath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	hash := strings.Split(string(output), " ")[0]
	return hash, nil
}
