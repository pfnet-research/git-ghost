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

package ghost

import (
	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
	"io"

	log "github.com/Sirupsen/logrus"
)

// ShowOptions represents arg for Pull func
type ShowOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.PullableDiffBranchSpec
	// if you want to consume and transform the output of `ghost.Show()`,
	// Please use `io.Pipe()` as below,
	// ```
	// r, w := io.Pipe()
	// go func() { ghost.Show(ShowOptions{ Writer: w }); w.Close()}
	// ````
	// Then, you can read the output from `r` and transform them as you like.
	Writer io.Writer
}

func pullAndshow(branchSpec types.PullableGhostBranchSpec, we types.WorkingEnv, writer io.Writer) errors.GitGhostError {
	branch, err := branchSpec.PullBranch(we)
	if err != nil {
		return err
	}
	if branch != nil {
		err := branch.Show(we, writer)
		return err
	}
	return nil
}

// Show writes ghost branches contents to option.Writer
func Show(options ShowOptions) errors.GitGhostError {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	if options.CommitsBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer util.LogDeferredGitGhostError(we.Clean)
		err = pullAndshow(options.CommitsBranchSpec, *we, options.Writer)
		if err != nil {
			return err
		}
	}

	if options.PullableDiffBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer util.LogDeferredGitGhostError(we.Clean)
		return pullAndshow(options.PullableDiffBranchSpec, *we, options.Writer)
	}

	log.WithFields(util.ToFields(options)).Warn("show command has nothing to do with")
	return nil
}
