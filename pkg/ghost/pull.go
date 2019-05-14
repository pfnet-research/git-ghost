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

	log "github.com/sirupsen/logrus"
)

// PullOptions represents arg for Pull func
type PullOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.PullableDiffBranchSpec
}

func pullAndApply(spec types.PullableGhostBranchSpec, we types.WorkingEnv) errors.GitGhostError {
	pulledBranch, err := spec.PullBranch(we)
	if err != nil {
		return errors.WithStack(err)
	}
	return pulledBranch.Apply(we)
}

// Pull pulls ghost branches and apply to workind directory
func Pull(options PullOptions) errors.GitGhostError {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.LogDeferredGitGhostError(we.Clean)

	if options.CommitsBranchSpec != nil {
		err := pullAndApply(*options.CommitsBranchSpec, *we)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if options.PullableDiffBranchSpec != nil {
		err := pullAndApply(*options.PullableDiffBranchSpec, *we)
		return errors.WithStack(err)
	}

	log.WithFields(util.ToFields(options)).Warn("pull command has nothing to do with")
	return nil
}
