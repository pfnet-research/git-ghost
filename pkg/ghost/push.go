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
	*types.LocalBaseBranchSpec
	*types.LocalModBranchSpec
}

// PushResult contains resultant ghost branches of Push func
type PushResult struct {
	*types.LocalBaseBranch
	*types.LocalModBranch
}

// Push pushes create ghost branches and push them to remote ghost repository
func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")

	var result PushResult
	if options.LocalBaseBranchSpec != nil {
		branch, err := pushGhostBranch(options.LocalBaseBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		localBaseBranch, _ := branch.(*types.LocalBaseBranch)
		result.LocalBaseBranch = localBaseBranch
	}

	if options.LocalModBranchSpec != nil {
		branch, err := pushGhostBranch(options.LocalModBranchSpec, options.WorkingEnvSpec)
		if err != nil {
			return nil, err
		}
		localModBranch, _ := branch.(*types.LocalModBranch)
		result.LocalModBranch = localModBranch
	}

	return &result, nil
}

func pushGhostBranch(branchSpec types.GhostBranchSpec, workingEnvSpec types.WorkingEnvSpec) (types.GhostBranch, error) {
	workingEnv, err := workingEnvSpec.Initialize()
	if err != nil {
		return nil, err
	}
	defer workingEnv.Clean()
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
