package util

import (
	"bytes"
	"errors"
	"os/exec"
)

func JustRunCmd(cmd *exec.Cmd) error {
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
