package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"os/exec"
	"strings"
)

// ListRemoteBranchNames returns remote branch names
func ListRemoteBranchNames(repo string, branchnames []string) ([]string, errors.GitGhostError) {
	if len(branchnames) == 0 {
		return []string{}, nil
	}

	branchNamesToSearch := []string{}
	for _, b := range branchnames {
		prefixed := b
		if !strings.HasPrefix(b, "refs/heads/") {
			prefixed = fmt.Sprintf("%s%s", "refs/heads/", b)
		}
		branchNamesToSearch = append(branchNamesToSearch, prefixed)
	}
	opts := append([]string{"ls-remote", "-q", "--heads", "--refs", repo}, branchNamesToSearch...)
	output, err := util.JustOutputCmd(exec.Command("git", opts...))
	if err != nil {
		return []string{}, errors.WithStack(err)
	}

	var branchNames []string
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) != 2 {
			return []string{}, errors.Errorf("Got unexpected line: %s", line)
		}
		// Assume it starts from "refs/heads/"
		name := tokens[1][11:]
		branchNames = append(branchNames, name)
	}
	return branchNames, nil
}
