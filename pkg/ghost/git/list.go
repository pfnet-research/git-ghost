package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"os/exec"
	"strings"
)

func ListGhostBranchNames(repo, prefix, fromRefspec, toRefspec string) ([]string, error) {
	fromPattern := "*"
	toPattern := "*"
	if fromRefspec != "" {
		fromPattern = fromRefspec
	}
	if toRefspec != "" {
		toPattern = toRefspec
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
