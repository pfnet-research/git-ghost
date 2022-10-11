// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type WorkDir struct {
	Dir string
	Env map[string]string
}

type CommandError struct {
	InternalError error
	Stdout        string
	Stderr        string
}

func (ce *CommandError) Error() string {
	return fmt.Sprintf("%s\n\nstdout:\n%s\n\nstderr:\n%s", ce.InternalError, ce.Stdout, ce.Stderr)
}

func CloneWorkDir(base *WorkDir) (*WorkDir, error) {
	wd, err := CreateWorkDir()
	if err != nil {
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "clone", base.Dir, wd.Dir)
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.email", "you@example.com")
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.name", "\"Your Name\"")
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	return wd, nil
}

func CreateGitWorkDir() (*WorkDir, error) {
	wd, err := CreateWorkDir()
	if err != nil {
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "init")
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.email", "you@example.com")
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	_, _, err = wd.RunCommmand("git", "config", "user.name", "\"Your Name\"")
	if err != nil {
		_ = wd.Remove()
		return nil, err
	}
	return wd, nil
}

func CreateWorkDir() (*WorkDir, error) {
	dir, err := os.MkdirTemp("", "git-ghost-e2e-test-")
	if err != nil {
		return nil, err
	}
	return &WorkDir{Dir: dir}, nil
}

func (wd *WorkDir) Remove() error {
	return os.RemoveAll(wd.Dir)
}

func (wd *WorkDir) RunGitGhostCommmand(args ...string) (string, string, error) {
	newArgs := []string{"ghost"}
	debug := os.Getenv("DEBUG")
	if debug != "" {
		newArgs = append(newArgs, "-vvv")
	}
	newArgs = append(newArgs, args...)
	return wd.RunCommmand("git", newArgs...)
}

func (wd *WorkDir) RunCommmand(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.Dir = wd.Dir
	env := make([]string, 0, len(os.Environ())+len(wd.Env)+1)
	env = append(env, os.Environ()...)
	for key, val := range wd.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	cmd.Env = env
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		err = &CommandError{
			InternalError: err,
			Stdout:        stdout.String(),
			Stderr:        stderr.String(),
		}
	}
	return stdout.String(), stderr.String(), err
}
