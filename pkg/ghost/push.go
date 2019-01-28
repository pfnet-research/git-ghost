package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"io/ioutil"
	"os"
)

func Push(options PushOptions) (*PushResult, error) {
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
		dstDir, err := ioutil.TempDir(options.DstDir, "git-ghost-")
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
			fmt.Fprintf(os.Stderr, "Skipped pushing existing branch '%s' in %s\n", branch.BranchName(), options.GhostRepo)
			continue
		}
		fmt.Fprintf(os.Stderr, "Pushing branch %s to %s\n", branch.BranchName(), options.GhostRepo)
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
