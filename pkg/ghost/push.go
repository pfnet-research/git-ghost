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
	"github.com/pfnet-research/git-ghost/pkg/ghost/git"
	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	log "github.com/sirupsen/logrus"
)

// PushOptions represents arg for Push func
type PushOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.DiffBranchSpec
}

// PushResult contains resultant ghost branches of Push func
type PushResult struct {
	*types.CommitsBranch
	*types.DiffBranch
}

// Push pushes create ghost branches and push them to remote ghost repository
func Push(options PushOptions) (*PushResult, errors.GitGhostError) {
	log.WithFields(util.ToFields(options)).Debug("push command with")

	var result PushResult
	if options.CommitsBranchSpec != nil {
		branch, err := pushGhostBranch(options.CommitsBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		commitsBranch, _ := branch.(*types.CommitsBranch)
		result.CommitsBranch = commitsBranch
	}

	if options.DiffBranchSpec != nil {
		branch, err := pushGhostBranch(options.DiffBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		diffBranch, _ := branch.(*types.DiffBranch)
		result.DiffBranch = diffBranch
	}

	return &result, nil
}

func pushGhostBranch(branchSpec types.GhostBranchSpec, workingEnvSpec types.WorkingEnvSpec) (types.GhostBranch, errors.GitGhostError) {
	workingEnv, err := workingEnvSpec.Initialize()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer util.LogDeferredGitGhostError(workingEnv.Clean)
	dstDir := workingEnv.GhostDir
	branch, err := branchSpec.CreateBranch(*workingEnv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if branch == nil {
		return nil, nil
	}
	existence, err := git.ValidateRemoteBranchExistence(
		workingEnv.GhostRepo,
		branch.BranchName(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if existence {
		log.WithFields(log.Fields{
			"branch":    branch.BranchName(),
			"ghostRepo": workingEnv.GhostRepo,
		}).Info("skipped pushing existing branch")
		return branch, nil
	}

	log.WithFields(log.Fields{
		"branch":    branch.BranchName(),
		"ghostRepo": workingEnv.GhostRepo,
	}).Info("pushing branch")
	err = git.Push(dstDir, branch.BranchName())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return branch, nil
}
