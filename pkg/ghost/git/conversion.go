package git

import (
	"git-ghost/pkg/util"
	"os/exec"
	"strings"
)

// ResolveComittish resolves comittish as full commit hash on dir
func ResolveComittish(dir, comittish string) (string, error) {
	commit, err := util.JustOutputCmd(
		exec.Command("git", "-C", dir, "rev-list", "-1", comittish),
	)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(commit), "\r\n"), nil
}
