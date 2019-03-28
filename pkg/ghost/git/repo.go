package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"os/exec"
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

// CommitAndPush commits and push to its origin
func CommitAndPush(dir, filename, message, comittish string) errors.GitGhostError {
	err := CommitFile(dir, filename, message)
	if err != nil {
		return errors.WithStack(err)
	}
	err = Push(dir, comittish)
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
func Push(dir string, comittishes ...string) errors.GitGhostError {
	args := []string{"-C", dir, "push", "origin"}
	args = append(args, comittishes...)
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

// Pull pulls comittish from its origin
func Pull(dir, comittish string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "pull", "origin", comittish),
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
