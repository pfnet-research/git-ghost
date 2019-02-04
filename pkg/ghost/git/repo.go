package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"os/exec"
)

var (
	ORIGIN string = "origin"
)

func InitializeGitDir(dir, repo, branch string) error {
	args := []string{"clone", "-q", "-o", ORIGIN}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repo, dir)
	cmd := exec.Command("git", args...)
	return util.JustRunCmd(cmd)
}

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

func DeleteRemoteBranches(dir string, branchNames ...string) error {
	args := []string{"-C", dir, "push", "origin"}
	for _, name := range branchNames {
		args = append(args, fmt.Sprintf(":%s", name))
	}
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

func Push(dir string, comittishes ...string) error {
	args := []string{"-C", dir, "push", "origin"}
	args = append(args, comittishes...)
	return util.JustRunCmd(
		exec.Command("git", args...),
	)
}

func Pull(dir, comittish string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "pull", "origin", comittish),
	)
}

func CreateOrphanBranch(dir, branch string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "checkout", "--orphan", branch),
	)
}

func ResetHardToBranch(dir, branch string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "reset", "--hard", branch),
	)
}
