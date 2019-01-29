package ghost

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type PushOptions struct {
	WorkingEnvSpec
	GhostPrefix string
	RemoteBase  string
	LocalBase   string
}

type PushResult struct {
	LocalBaseBranch *LocalBaseBranch
	LocalModBranch  *LocalModBranch
}

func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")
	branchSpecs := []GhostBranchSpec{
		LocalBaseBranchSpec{
			Prefix:            options.GhostPrefix,
			RemoteBaseRefspec: options.RemoteBase,
			LocalBaseRefspec:  options.LocalBase,
		},
		LocalModBranchSpec{
			Prefix:           options.GhostPrefix,
			LocalBaseRefspec: options.LocalBase,
		},
	}

	branches := []GhostBranch{}
	for _, branchSpec := range branchSpecs {
		workingEnv, err := options.WorkingEnvSpec.initialize()
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
			continue
		}
		branches = append(branches, branch)
		existence, err := git.ValidateRemoteBranchExistence(options.GhostRepo, branch.BranchName())
		if err != nil {
			return nil, err
		}
		if existence {
			log.WithFields(log.Fields{
				"branch":    branch.BranchName(),
				"ghostRepo": options.GhostRepo,
			}).Info("skipped pushing existing branch")
			continue
		}

		log.WithFields(log.Fields{
			"branch":    branch.BranchName(),
			"ghostRepo": options.GhostRepo,
		}).Info("pushing branch")
		err = git.Push(dstDir, branch.BranchName())
		if err != nil {
			return nil, err
		}
	}

	resp := PushResult{}
	for _, branch := range branches {
		if br, ok := branch.(*LocalBaseBranch); ok {
			resp.LocalBaseBranch = br
		}
		if br, ok := branch.(*LocalModBranch); ok {
			resp.LocalModBranch = br
		}
	}

	return &resp, nil
}
