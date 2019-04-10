// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"fmt"
	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
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

	lines := strings.Split(string(output), "\n")
	branchNames := make([]string, 0, len(lines))
	for _, line := range lines {
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
