package types

import (
	"git-ghost/pkg/ghost/git"
)

// ListCommitsBranchSpec is spec for list commits branch
type ListCommitsBranchSpec struct {
	// Prefix is a prefix of branch name
	Prefix string
	// HashFrom is comittish value to list HashFrom..HashTo
	HashFrom string
	// HashTo is a comittish value to list HashFrom..HashTo
	HashTo string
}

// ListCommitsBranchSpec is spec for list diff branch
type ListDiffBranchSpec struct {
	Prefix string
	// HashFrom is comittish value to list HashFrom..HashTo
	HashFrom string
	// HashTo is a comittish value to list HashFrom..HashTo
	HashTo string
}

// Resolve resolves commitish values in ListCommitsBranchSpec as full commit hash
func (ls *ListCommitsBranchSpec) Resolve(srcDir string) *ListCommitsBranchSpec {
	newSpec := *ls
	if ls.HashFrom != "" {
		newSpec.HashFrom = resolveComittishOr(srcDir, ls.HashFrom)
	}
	if ls.HashTo != "" {
		newSpec.HashTo = resolveComittishOr(srcDir, ls.HashTo)
	}
	return &newSpec
}

// GetBranches returns CommitsBranches from spec
func (ls *ListCommitsBranchSpec) GetBranches(repo string) (CommitsBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, err
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

// Resolve resolves commitish values in ListDiffBranchSpec as full commit hash
func (ls *ListDiffBranchSpec) Resolve(srcDir string) *ListDiffBranchSpec {
	newSpec := *ls
	if ls.HashFrom != "" {
		newSpec.HashFrom = resolveComittishOr(srcDir, ls.HashFrom)
	}
	return &newSpec
}

// GetBranches returns DiffBranches from spec
func (ls *ListDiffBranchSpec) GetBranches(repo string) (DiffBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, err
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
