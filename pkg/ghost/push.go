package ghost

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
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
func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")

	var result PushResult
	if options.CommitsBranchSpec != nil {
		branch, err := pushGhostBranch(options.CommitsBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		commitsBranch, _ := branch.(*types.CommitsBranch)
		result.CommitsBranch = commitsBranch
	}

	if options.DiffBranchSpec != nil {
		branch, err := pushGhostBranch(options.DiffBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		diffBranch, _ := branch.(*types.DiffBranch)
		result.DiffBranch = diffBranch
	}

	return &result, nil
}

func pushGhostBranch(branchSpec types.GhostBranchSpec, workingEnvSpec types.WorkingEnvSpec) (types.GhostBranch, error) {
	workingEnv, err := workingEnvSpec.Initialize()
	if err != nil {
		return nil, err
	}
	defer util.LogDeferredError(workingEnv.Clean)
	dstDir := workingEnv.GhostDir
	branch, err := branchSpec.CreateBranch(*workingEnv)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, nil
	}
	existence, err := git.ValidateRemoteBranchExistence(
		workingEnv.GhostRepo,
		branch.BranchName(),
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return branch, nil
}
