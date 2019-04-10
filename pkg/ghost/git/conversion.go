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
	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
	"os/exec"
	"strings"
)

// ResolveCommittish resolves committish as full commit hash on dir
func ResolveCommittish(dir, committish string) (string, errors.GitGhostError) {
	commit, err := util.JustOutputCmd(
		exec.Command("git", "-C", dir, "rev-list", "-1", committish),
	)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(commit), "\r\n"), nil
}
