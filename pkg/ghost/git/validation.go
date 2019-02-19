package git

import (
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"os/exec"
)

// ValidateGit check the environment has 'git' command or not.
func ValidateGit() errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "version"),
	)
}

// ValidateComittish check comittish is valid on dir
func ValidateComittish(dir, comittish string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "cat-file", "-e", comittish),
	)
}
