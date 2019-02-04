package ghost

import (
	"git-ghost/pkg/ghost/git"
)

type ListCommitsBranchSpec struct {
	Prefix   string
	HashFrom string
	HashTo   string
}

type ListDiffBranchSpec struct {
	Prefix   string
	HashFrom string
	HashTo   string
}

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

func (ls *ListCommitsBranchSpec) GetBranches(repo string) (LocalBaseBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
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

func (ls *ListDiffBranchSpec) Resolve(srcDir string) *ListDiffBranchSpec {
	newSpec := *ls
	if ls.HashFrom != "" {
		newSpec.HashFrom = resolveComittishOr(srcDir, ls.HashFrom)
	}
	return &newSpec
}

func (ls *ListDiffBranchSpec) GetBranches(repo string) (LocalModBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, ls.Prefix, ls.HashFrom, ls.HashTo)
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
