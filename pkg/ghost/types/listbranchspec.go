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

package types

import (
	"fmt"

	"github.com/pfnet-research/git-ghost/pkg/ghost/git"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
)

// ListCommitsBranchSpec is spec for list commits branch
type ListCommitsBranchSpec struct {
	// Prefix is a prefix of branch name
	Prefix string
	// HashFrom is committish value to list HashFrom..HashTo
	HashFrom string
	// HashTo is a committish value to list HashFrom..HashTo
	HashTo string
}

// ListCommitsBranchSpec is spec for list diff branch
type ListDiffBranchSpec struct {
	Prefix string
	// HashFrom is committish value to list HashFrom..HashTo
	HashFrom string
	// HashTo is a committish value to list HashFrom..HashTo
	HashTo string
}

// Resolve resolves committish values in ListCommitsBranchSpec as full commit hash
func (ls *ListCommitsBranchSpec) Resolve(srcDir string) *ListCommitsBranchSpec {
	newSpec := *ls
	if ls.HashFrom != "" {
		newSpec.HashFrom = resolveCommittishOr(srcDir, ls.HashFrom)
	}
	if ls.HashTo != "" {
		newSpec.HashTo = resolveCommittishOr(srcDir, ls.HashTo)
	}
	return &newSpec
}

// GetBranches returns CommitsBranches from spec
func (ls *ListCommitsBranchSpec) GetBranches(repo string) (CommitsBranches, errors.GitGhostError) {
	branchNames, err := listGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var branches CommitsBranches
	for _, name := range branchNames {
		branch := CreateGhostBranchByName(name)
		if br, ok := branch.(*CommitsBranch); ok {
			branches = append(branches, *br)
		}
	}
	return branches, nil
}

// Resolve resolves committish values in ListDiffBranchSpec as full commit hash
func (ls *ListDiffBranchSpec) Resolve(srcDir string) *ListDiffBranchSpec {
	newSpec := *ls
	if ls.HashFrom != "" {
		newSpec.HashFrom = resolveCommittishOr(srcDir, ls.HashFrom)
	}
	return &newSpec
}

// GetBranches returns DiffBranches from spec
func (ls *ListDiffBranchSpec) GetBranches(repo string) (DiffBranches, errors.GitGhostError) {
	branchNames, err := listGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var branches DiffBranches
	for _, name := range branchNames {
		branch := CreateGhostBranchByName(name)
		if br, ok := branch.(*DiffBranch); ok {
			branches = append(branches, *br)
		}
	}
	return branches, nil
}

func listGhostBranchNames(repo, prefix, fromCommittish, toCommittish string) ([]string, error) {
	fromPattern := "*"
	toPattern := "*"
	if fromCommittish != "" {
		fromPattern = fromCommittish
	}
	if toCommittish != "" {
		toPattern = toCommittish
	}

	branchNames, err := git.ListRemoteBranchNames(repo, []string{
		fmt.Sprintf("%s/%s-%s", prefix, fromPattern, toPattern),
		fmt.Sprintf("%s/%s/%s", prefix, fromPattern, toPattern),
	})
	if err != nil {
		return []string{}, errors.WithStack(err)
	}

	return branchNames, nil
}
