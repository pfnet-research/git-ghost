package ghost

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type PushOptions struct {
	PushCommitsOptions
	PushDiffOptions
}

type PushResult struct {
	*PushCommitsResult
	*PushDiffResult
}

type PushCommitsOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
}
type PushCommitsResult struct {
	LocalBaseBranch *LocalBaseBranch
}
type PushDiffOptions struct {
	WorkingEnvSpec
	LocalModBranchSpec
}

type PushDiffResult struct {
	LocalModBranch *LocalModBranch
}

func PushDiff(options PushDiffOptions) (*PushDiffResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push diff command with")
	branch, err := pushGhostBranch(options.LocalModBranchSpec, options.WorkingEnvSpec)
	if err != nil {
		return nil, err
	}
	localModBranch, _ := branch.(*LocalModBranch)
	return &PushDiffResult{localModBranch}, nil
}

func PushCommits(options PushCommitsOptions) (*PushCommitsResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push commits command with")
	branch, err := pushGhostBranch(options.LocalBaseBranchSpec, options.WorkingEnvSpec)
	if err != nil {
		return nil, err
	}
	localBaseBranch, _ := branch.(*LocalBaseBranch)
	return &PushCommitsResult{localBaseBranch}, nil
}

func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")

	var result *PushResult
	pushCommitsResult, err := PushCommits(options.PushCommitsOptions)
	if err != nil {
		return result, err
	}
	result = &PushResult{
		PushCommitsResult: pushCommitsResult,
	}

	pushDiffResult, err := PushDiff(options.PushDiffOptions)
	if err != nil {
		return result, err
	}
	result.PushDiffResult = pushDiffResult

	return result, nil
}

func pushGhostBranch(branchSpec GhostBranchSpec, workingEnvSpec WorkingEnvSpec) (GhostBranch, error) {
	workingEnv, err := workingEnvSpec.initialize()
	if err != nil {
		return nil, err
	}
	defer workingEnv.clean()
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
		return nil, nil
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
