package git

import (
	"git-ghost/pkg/util"
	"os/exec"
)

func ValidateGit() error {
	return util.JustRunCmd(
		exec.Command("git", "version"),
	)
}

func ValidateRefspec(dir, refspec string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "cat-file", "-e", refspec),
	)
}
