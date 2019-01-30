package util

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func JustOutputCmd(cmd *exec.Cmd) ([]byte, error) {
	wd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"pwd":     wd,
		"command": strings.Join(cmd.Args, " "),
	}).Debug("exec")
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
	wd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"pwd":     wd,
		"command": strings.Join(cmd.Args, " "),
	}).Debug("exec")
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
