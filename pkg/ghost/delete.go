package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

// DeleteOptions represents arg for Delete func
type DeleteOptions struct {
	types.WorkingEnvSpec
	*types.ListCommitsBranchSpec
	*types.ListDiffBranchSpec
	Dryrun bool
}

// DeleteResult contains deleted ghost branches in Delete func
type DeleteResult struct {
	*types.LocalBaseBranches
	*types.LocalModBranches
}

// Delete deletes ghost branches from ghost repo and returns deleted branches
func Delete(options DeleteOptions) (*DeleteResult, error) {
	log.WithFields(util.ToFields(options)).Debug("delete command with")

	res := DeleteResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		res.LocalBaseBranches = &branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		res.LocalModBranches = &branches
	}

	workingEnv, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return nil, err
	}
	defer workingEnv.Clean()

	deleteBranches := func(branches []types.GhostBranch, dryrun bool) error {
		var branchNames []string
		for _, branch := range branches {
			branchNames = append(branchNames, branch.BranchName())
		}
		log.WithFields(log.Fields{
			"branches": branchNames,
		}).Info("Delete branch")
		if dryrun {
			return nil
		}
		return git.DeleteRemoteBranches(workingEnv.GhostDir, branchNames...)
	}

	if res.LocalBaseBranches != nil {
		err := deleteBranches(res.LocalBaseBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, err
		}
	}

	if res.LocalModBranches != nil {
		err := deleteBranches(res.LocalModBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *DeleteResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if res.LocalBaseBranches != nil {
		buffer.WriteString("Deleted Local Base Branches:\n")
		buffer.WriteString("\n")
		branches := *res.LocalBaseBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s => %s\n", branch.RemoteBaseCommit, branch.LocalBaseCommit))
		}
		buffer.WriteString("\n")
	}
	if res.LocalModBranches != nil {
		buffer.WriteString("Deleted Local Mod Branches:\n")
		buffer.WriteString("\n")
		branches := *res.LocalModBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s -> %s\n", branch.LocalBaseCommit, branch.LocalModHash))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
