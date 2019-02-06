package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"os/exec"
	"strings"
)

// ListGhostBranchNames returns all ghost branchnames for fromComittish..toComittish.
// You can use wildcard in fromComittish and toComittish
func ListGhostBranchNames(repo, prefix, fromComittish, toComittish string) ([]string, error) {
	fromPattern := "*"
	toPattern := "*"
	if fromComittish != "" {
		fromPattern = fromComittish
	}
	if toComittish != "" {
		toPattern = toComittish
	}
	output, err := util.JustOutputCmd(exec.Command("git",
		"ls-remote", "-q", "--heads", "--refs", repo,
		fmt.Sprintf("refs/heads/%s/%s-%s", prefix, fromPattern, toPattern),
		fmt.Sprintf("refs/heads/%s/%s/%s", prefix, fromPattern, toPattern),
	))
	if err != nil {
		return []string{}, err
	}
	var branchNames []string
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) != 2 {
			return []string{}, fmt.Errorf("Got unexpected line: %s", line)
		}
		// Assume it starts from "refs/heads/"
		name := tokens[1][11:]
		branchNames = append(branchNames, name)
	}
	return branchNames, nil
}
