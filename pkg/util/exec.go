package util

import (
	"bytes"
	"errors"
	"os/exec"
)

func JustOutputCmd(cmd *exec.Cmd) ([]byte, error) {
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	bytes, err := cmd.Output()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return []byte{}, errors.New(s)
		}
		return []byte{}, err
	}
	return bytes, err
}

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
