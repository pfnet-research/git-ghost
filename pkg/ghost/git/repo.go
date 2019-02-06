package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"os/exec"
)

var (
	ORIGIN string = "origin"
)

// InitializeGitDir clone repo to dir.
// if you set empty branchname, it will checkout default branch of repo.
func InitializeGitDir(dir, repo, branch string) error {
	args := []string{"clone", "-q", "-o", ORIGIN}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repo, dir)
	cmd := exec.Command("git", args...)
	return util.JustRunCmd(cmd)
}

// CommitAndPush commits and push to its origin
func CommitAndPush(dir, filename, message, comittish string) error {
	err := CommitFile(dir, filename, message)
	if err != nil {
		return err
	}
	err = Push(dir, comittish)
	if err != nil {
		return err
	}
	return nil
}

// CommitFile commits a file
func CommitFile(dir, filename, message string) error {
	err := util.JustRunCmd(
		exec.Command("git", "-C", dir, "add", filename),
	)
	if err != nil {
		return err
	}
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "commit", "-q", filename, "-m", message),
	)
}

// DeleteRemoteBranches delete branches from its origin
func DeleteRemoteBranches(dir string, branchNames ...string) error {
	args := []string{"-C", dir, "push", "origin"}
	for _, name := range branchNames {
		args = append(args, fmt.Sprintf(":%s", name))
	}
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

// Push pushes current HEAD to its origin
func Push(dir string, comittishes ...string) error {
	args := []string{"-C", dir, "push", "origin"}
	args = append(args, comittishes...)
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

// Pull pulls comittish from its origin
func Pull(dir, comittish string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "pull", "origin", comittish),
	)
}

// CreateOrphanBranch creates an orphan branch on dir
func CreateOrphanBranch(dir, branch string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "checkout", "--orphan", branch),
	)
}

// ResetHardToBranch reset dir to branch with --hard option
func ResetHardToBranch(dir, branch string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "reset", "--hard", branch),
	)
}
