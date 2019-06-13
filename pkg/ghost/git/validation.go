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

// ValidateGit check the environment has 'git' command or not.
func ValidateGit() errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "version"),
	)
}

// ValidateCommittish check committish is valid on dir
func ValidateCommittish(dir, committish string) errors.GitGhostError {
	output, err := util.JustOutputCmd(
		exec.Command("git", "-C", dir, "cat-file", "-e", committish),
	)
	if err != nil && util.GetExitCode(err.Cause()) == 1 && len(output) == 0 {
		// exit 1 is for unexisting committish.
		return errors.Errorf("%s does not exist", committish)
	}
	return err
}
