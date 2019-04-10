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
	"os/exec"

	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
)

// ValidateRemoteBranchExistence checks repo has branch or not.
func ValidateRemoteBranchExistence(repo, branch string) (bool, errors.GitGhostError) {
	output, err := util.JustOutputCmd(
		exec.Command("git", "ls-remote", "--heads", repo, branch),
	)
	if err != nil {
		return false, err
	}
	return string(output) != "", nil
}
