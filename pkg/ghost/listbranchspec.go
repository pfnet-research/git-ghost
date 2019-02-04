package ghost

import (
	"git-ghost/pkg/ghost/git"
)

type ListCommitsBranchSpec struct {
	HashFrom string
	HashTo   string
}

type ListDiffBranchSpec struct {
	HashFrom string
	HashTo   string
}

func (ls ListCommitsBranchSpec) Resolve(srcDir string) *ListCommitsBranchSpec {
	var newOptions ListCommitsBranchSpec
	if ls.HashFrom != "" {
		newOptions.HashFrom = resolveComittishOr(srcDir, ls.HashFrom)
	}
	if ls.HashTo != "" {
		newOptions.HashTo = resolveComittishOr(srcDir, ls.HashTo)
	}
	return &newOptions
}

func (ls *ListCommitsBranchSpec) GetBranches(repo, prefix string) (LocalBaseBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, err
	}
	var branches LocalBaseBranches
	for _, name := range branchNames {
		branch := CreateGhostBranchByName(name)
		if br, ok := branch.(*LocalBaseBranch); ok {
			branches = append(branches, *br)
		}
	}
	return branches, nil
}

func (ls ListDiffBranchSpec) Resolve(srcDir string) *ListDiffBranchSpec {
	var newOptions ListDiffBranchSpec
	if ls.HashFrom != "" {
		newOptions.HashFrom = resolveComittishOr(srcDir, ls.HashFrom)
	}
	return &newOptions
}

func (ls *ListDiffBranchSpec) GetBranches(repo, prefix string) (LocalModBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, prefix, ls.HashFrom, ls.HashTo)
	if err != nil {
		return nil, err
	}
	var branches LocalModBranches
	for _, name := range branchNames {
		branch := CreateGhostBranchByName(name)
		if br, ok := branch.(*LocalModBranch); ok {
			branches = append(branches, *br)
		}
	}
	return branches, nil
}
