package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"os"
)

func Push(options PushOptions) (*PushResult, error) {
	repo := options.GhostRepo
	dstDir, err := git.CreateTempGitDir(options.DstDir, options.GhostRepo, "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dstDir)

	branchSpecs := []GhostBranchSpec{
		LocalBaseBranchSpec{
			Prefix:            options.GhostPrefix,
			SrcDir:            options.SrcDir,
			DstDir:            dstDir,
			RemoteBaseRefspec: options.RemoteBase,
			LocalBaseRefspec:  options.LocalBase,
		},
		LocalModBranchSpec{
			Prefix:           options.GhostPrefix,
			SrcDir:           options.SrcDir,
			DstDir:           dstDir,
			LocalBaseRefspec: options.LocalBase,
		},
	}
	branches := []GhostBranch{}
	for _, branchSpec := range branchSpecs {
		branch, err := branchSpec.CreateBranch()
		if err != nil {
			return nil, err
		}
		existence, err := git.ValidateRemoteBranchExistence(repo, branch.BranchName())
		if err != nil {
			return nil, err
		}
		if existence {
			fmt.Printf("Skipped pushing existing branch '%s' in %s\n", branch.BranchName(), repo)
		} else {
			fmt.Printf("Created branch '%s' in %s\n", branch.BranchName(), repo)
		}
		branches = append(branches, branch)
	}

	branchNames := []string{}
	for _, branch := range branches {
		branchNames = append(branchNames, branch.BranchName())
	}
	if len(branchNames) > 0 {
		fmt.Printf("Pushing branches %v\n", branchNames)
		err := git.Push(dstDir, branchNames...)
		if err != nil {
			return nil, err
		}
	}

	resp := PushResult{}
	for _, branch := range branches {
		if br, ok := branch.(LocalBaseBranch); ok {
			resp.LocalBaseBranch = &br
		}
		if br, ok := branch.(LocalModBranch); ok {
			resp.LocalModBranch = &br
		}
	}

	return &resp, nil
}
