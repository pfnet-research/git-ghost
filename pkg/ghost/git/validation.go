package git

import (
	"git-ghost/pkg/util"
	"os/exec"
)

// ValidateGit check the environment has 'git' command or not.
func ValidateGit() error {
	return util.JustRunCmd(
		exec.Command("git", "version"),
	)
}

// ValidateComittish check comittish is valid on dir
func ValidateComittish(dir, comittish string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "cat-file", "-e", comittish),
	)
}
