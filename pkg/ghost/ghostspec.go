package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"

	log "github.com/Sirupsen/logrus"
)

type GhostSpec struct {
	GhostPrefix  string
	RemoteBase   string
	LocalBase    string
	LocalModHash string
}

func (gs GhostSpec) validateAndCreateGhostBranches(weSpec WorkingEnvSpec) (*LocalBaseBranch, *LocalModBranch, error) {
	var err error
	remoteBaseResolved := gs.RemoteBase
	localBaseResolved := gs.LocalBase

	// try to resolve specified commit-ish values
	// if failed, we will use specified values.
	triedToBeResolved, err := git.ResolveRefspec(weSpec.SrcDir, gs.RemoteBase)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": weSpec.SrcDir,
			"specified":  gs.RemoteBase,
		}).Warn("can't resolve --base-commit on local git repository.  specified commit-ish value will be used.")
	} else {
		remoteBaseResolved = triedToBeResolved
	}

	triedToBeResolved, err = git.ResolveRefspec(weSpec.SrcDir, gs.LocalBase)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": weSpec.SrcDir,
			"specified":  gs.LocalBase,
		}).Warn("can't resolve --local-base on local git repository.  specified commit-ish value will be used.")
	} else {
		localBaseResolved = triedToBeResolved
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
