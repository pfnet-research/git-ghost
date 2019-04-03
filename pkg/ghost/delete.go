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
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"

	log "github.com/Sirupsen/logrus"
)

// DeleteOptions represents arg for Delete func
type DeleteOptions struct {
	types.WorkingEnvSpec
	*types.ListCommitsBranchSpec
	*types.ListDiffBranchSpec
	Dryrun bool
}

// DeleteResult contains deleted ghost branches in Delete func
type DeleteResult struct {
	*types.CommitsBranches
	*types.DiffBranches
}

// Delete deletes ghost branches from ghost repo and returns deleted branches
func Delete(options DeleteOptions) (*DeleteResult, errors.GitGhostError) {
	log.WithFields(util.ToFields(options)).Debug("delete command with")

	res := DeleteResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res.CommitsBranches = &branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res.DiffBranches = &branches
	}

	workingEnv, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer util.LogDeferredGitGhostError(workingEnv.Clean)

	deleteBranches := func(branches []types.GhostBranch, dryrun bool) errors.GitGhostError {
		if len(branches) == 0 {
			return nil
		}
		var branchNames []string
		for _, branch := range branches {
			branchNames = append(branchNames, branch.BranchName())
		}
		log.WithFields(log.Fields{
			"branches": branchNames,
		}).Info("Delete branch")
		if dryrun {
			return nil
		}
		err := git.DeleteRemoteBranches(workingEnv.GhostDir, branchNames...)
		return errors.WithStack(err)
	}

	if res.CommitsBranches != nil {
		err := deleteBranches(res.CommitsBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if res.DiffBranches != nil {
		err := deleteBranches(res.DiffBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *DeleteResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if res.CommitsBranches != nil {
		buffer.WriteString("Deleted Local Base Branches:\n")
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("%-40s %-40s\n", "Remote Base", "Local Base"))
		branches := *res.CommitsBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s %s\n", branch.CommitHashFrom, branch.CommitHashTo))
		}
		buffer.WriteString("\n")
	}
	if res.DiffBranches != nil {
		buffer.WriteString("Deleted Local Mod Branches:\n")
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("%-40s %-40s\n", "Local Base", "Local Mod"))
		branches := *res.DiffBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s %s\n", branch.CommitHashFrom, branch.DiffHash))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
