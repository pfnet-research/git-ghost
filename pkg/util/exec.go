package util

import (
	"bytes"
	"git-ghost/pkg/util/errors"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func JustOutputCmd(cmd *exec.Cmd) ([]byte, errors.GitGhostError) {
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
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func JustStreamOutputCmd(cmd *exec.Cmd, writer io.Writer) errors.GitGhostError {
	wd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"pwd":     wd,
		"command": strings.Join(cmd.Args, " "),
	}).Debug("exec")
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	cmd.Stdout = writer
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return errors.WithStack(err)
	}
	return nil
}

func JustRunCmd(cmd *exec.Cmd) errors.GitGhostError {
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
		return errors.WithStack(err)
	}
	return nil
}
