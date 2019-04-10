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

package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"os/exec"
	"strings"
)

var (
	ORIGIN string = "origin"
)

// InitializeGitDir clone repo to dir.
// if you set empty branchname, it will checkout default branch of repo.
func InitializeGitDir(dir, repo, branch string) errors.GitGhostError {
	args := []string{"clone", "-q", "-o", ORIGIN}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repo, dir)
	cmd := exec.Command("git", args...)
	return util.JustRunCmd(cmd)
}

// CopyUserConfig copies user config from source directory to destination directory.
func CopyUserConfig(srcDir, dstDir string) errors.GitGhostError {
	// Get user config from src
	userNameBytes, err := util.JustOutputCmd(exec.Command("git", "-C", srcDir, "config", "user.name"))
	if err != nil {
		return errors.WithStack(err)
	}
	userName := strings.TrimSuffix(string(userNameBytes), "\n")
	userEmailBytes, err := util.JustOutputCmd(exec.Command("git", "-C", srcDir, "config", "user.email"))
	if err != nil {
		return errors.WithStack(err)
	}
	userEmail := strings.TrimSuffix(string(userEmailBytes), "\n")
	// Copy the user config to dst
	err = util.JustRunCmd(exec.Command("git", "-C", dstDir, "config", "user.name", fmt.Sprintf("\"%s\"", userName)))
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(util.JustRunCmd(exec.Command("git", "-C", dstDir, "config", "user.email", fmt.Sprintf("\"%s\"", userEmail))))
}

// CommitAndPush commits and push to its origin
func CommitAndPush(dir, filename, message, committish string) errors.GitGhostError {
	err := CommitFile(dir, filename, message)
	if err != nil {
		return errors.WithStack(err)
	}
	err = Push(dir, committish)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// CommitFile commits a file
func CommitFile(dir, filename, message string) errors.GitGhostError {
	err := util.JustRunCmd(
		exec.Command("git", "-C", dir, "add", filename),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "commit", "-q", filename, "-m", message),
	)
}

// DeleteRemoteBranches delete branches from its origin
func DeleteRemoteBranches(dir string, branchNames ...string) errors.GitGhostError {
	args := []string{"-C", dir, "push", "origin"}
	for _, name := range branchNames {
		args = append(args, fmt.Sprintf(":%s", name))
	}
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

// Push pushes current HEAD to its origin
func Push(dir string, committishes ...string) errors.GitGhostError {
	args := []string{"-C", dir, "push", "origin"}
	args = append(args, committishes...)
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

// Pull pulls committish from its origin
func Pull(dir, committish string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "pull", "origin", committish),
	)
}

// CreateOrphanBranch creates an orphan branch on dir
func CreateOrphanBranch(dir, branch string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "checkout", "--orphan", branch),
	)
}

// ResetHardToBranch reset dir to branch with --hard option
func ResetHardToBranch(dir, branch string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "reset", "--hard", branch),
	)
}

// Init initializes a git repo
func Init(dir string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "init"),
	)
}
