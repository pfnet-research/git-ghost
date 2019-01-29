package git

import (
	"git-ghost/pkg/util"
	"os/exec"
)

func ValidateRemoteBranchExistence(repo, branch string) (bool, error) {
	output, err := util.JustOutputCmd(
		exec.Command("git", "ls-remote", "--heads", repo, branch),
	)
	if err != nil {
		return false, err
	}
	return string(output) != "", nil
}
