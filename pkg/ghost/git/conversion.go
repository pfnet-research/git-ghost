package git

import (
	"git-ghost/pkg/util"
	"os/exec"
	"strings"
)

func ResolveRefspec(dir, refspec string) (string, error) {
	commit, err := util.JustOutputCmd(
		exec.Command("git", "-C", dir, "rev-list", "-1", refspec),
	)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(commit), "\r\n"), nil
}
