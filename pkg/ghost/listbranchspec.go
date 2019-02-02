package ghost

import (
	"git-ghost/pkg/ghost/git"

	log "github.com/Sirupsen/logrus"
)

type ListCommitsBranchSpec struct {
	HashFrom string
	HashTo   string
}

type ListDiffBranchSpec struct {
	HashFrom string
	HashTo   string
}

func (options ListCommitsBranchSpec) Resolve(srcDir string) *ListCommitsBranchSpec {
	var newOptions ListCommitsBranchSpec
	newOptions.HashFrom = resolveRefspecOrIgnore(srcDir, options.HashFrom)
	newOptions.HashTo = resolveRefspecOrIgnore(srcDir, options.HashTo)
	return &newOptions
}

func (options *ListCommitsBranchSpec) GetBranches(repo, prefix string) (LocalBaseBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, prefix, options.HashFrom, options.HashTo)
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

func (options ListDiffBranchSpec) Resolve(srcDir string) *ListDiffBranchSpec {
	var newOptions ListDiffBranchSpec
	newOptions.HashFrom = resolveRefspecOrIgnore(srcDir, options.HashFrom)
	newOptions.HashTo = resolveRefspecOrIgnore(srcDir, options.HashTo)
	return &newOptions
}

func (options *ListDiffBranchSpec) GetBranches(repo, prefix string) (LocalModBranches, error) {
	branchNames, err := git.ListGhostBranchNames(repo, prefix, options.HashFrom, options.HashTo)
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

func resolveRefspecOrIgnore(dir, refspec string) string {
	commit, err := git.ResolveRefspec(dir, refspec)
	if err != nil {
		log.WithFields(log.Fields{
			"refspec": refspec,
		}).Warning("Failed to resolve refspec")
		commit = refspec
	}
	return commit
}
