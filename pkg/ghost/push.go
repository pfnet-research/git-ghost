package ghost

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

type PushOptions struct {
	SrcDir          string
	GhostWorkingDir string
	GhostPrefix     string
	GhostRepo       string
	RemoteBase      string
	LocalBase       string
}

type PushResult struct {
	LocalBaseBranch *LocalBaseBranch
	LocalModBranch  *LocalModBranch
}

func Push(options PushOptions) (*PushResult, error) {
	log.WithFields(util.ToFields(options)).Debug("push command with")
	branchSpecs := []GhostBranchSpec{
		LocalBaseBranchSpec{
			Repo:              options.GhostRepo,
			Prefix:            options.GhostPrefix,
			SrcDir:            options.SrcDir,
			RemoteBaseRefspec: options.RemoteBase,
			LocalBaseRefspec:  options.LocalBase,
		},
		LocalModBranchSpec{
			Repo:             options.GhostRepo,
			Prefix:           options.GhostPrefix,
			SrcDir:           options.SrcDir,
			LocalBaseRefspec: options.LocalBase,
		},
	}

	branches := []GhostBranch{}
	for _, branchSpec := range branchSpecs {
		dstDir, err := ioutil.TempDir(options.GhostWorkingDir, "git-ghost-")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(dstDir)
		branch, err := branchSpec.CreateBranch(dstDir)
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
