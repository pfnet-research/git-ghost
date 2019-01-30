package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"regexp"
)

type GhostSpec struct {
	GhostPrefix  string
	RemoteBase   string
	LocalBase    string
	LocalModHash string
}

func (gs GhostSpec) validateAndCreateGhostBranches(weSpec WorkingEnvSpec) (*LocalBaseBranch, *LocalModBranch, error) {
	var err error

	// resolve HEAD If necessary
	remoteBaseResolved := gs.RemoteBase
	localBaseResolved := gs.LocalBase
	sha1Regex := regexp.MustCompile(`\b[0-9a-f]{5,40}\b`)

	if !sha1Regex.MatchString(gs.RemoteBase) {
		remoteBaseResolved, err = git.ResolveRefspec(weSpec.SrcDir, gs.RemoteBase)
		if err != nil {
			return nil, nil, err
		}
	}
	if !sha1Regex.MatchString(gs.RemoteBase) {
		localBaseResolved, err = git.ResolveRefspec(weSpec.SrcDir, gs.LocalBase)
		if err != nil {
			return nil, nil, err
		}
	}

	// ghost branch validations and create ghost branches
	var localBaseBranch *LocalBaseBranch
	if remoteBaseResolved != localBaseResolved {
		// TODO warning when srcDir is on remoteBaseResolved.
		localBaseBranch = &LocalBaseBranch{
			Prefix:           gs.GhostPrefix,
			RemoteBaseCommit: remoteBaseResolved,
			LocalBaseCommit:  localBaseResolved,
		}

		existence, err := git.ValidateRemoteBranchExistence(weSpec.GhostRepo, localBaseBranch.BranchName())
		if err != nil {
			return nil, nil, err
		}
		if !existence {
			return nil, nil, fmt.Errorf("can't resolve local base branch on %s: %+v", weSpec.GhostRepo, localBaseBranch)
		}
	}

	localModBranch := &LocalModBranch{
		Prefix:          gs.GhostPrefix,
		LocalBaseCommit: localBaseResolved,
		LocalModHash:    gs.LocalModHash,
	}
	existence, err := git.ValidateRemoteBranchExistence(weSpec.GhostRepo, localModBranch.BranchName())
	if err != nil {
		return nil, nil, err
	}
	if !existence {
		return nil, nil, fmt.Errorf("can't resolve local mod branch on %s: %+v", weSpec.GhostRepo, localModBranch)
	}

	return localBaseBranch, localModBranch, nil
}
