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
	*types.CommitsBranches
	*types.DiffBranches
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
		res.CommitsBranches = &branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		res.DiffBranches = &branches
	}

	workingEnv, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return nil, err
	}
	defer util.LogError(workingEnv.Clean)

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

	if res.CommitsBranches != nil {
		err := deleteBranches(res.CommitsBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, err
		}
	}

	if res.DiffBranches != nil {
		err := deleteBranches(res.DiffBranches.AsGhostBranches(), options.Dryrun)
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
	if res.CommitsBranches != nil {
		buffer.WriteString("Deleted Local Base Branches:\n")
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("%-40s %-40s\n", "Remote Base", "Local Base"))
		branches := *res.CommitsBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s %s\n", branch.CommitHashFrom, branch.CommitHashTo))
		}
		buffer.WriteString("\n")
	}
	if res.DiffBranches != nil {
		buffer.WriteString("Deleted Local Mod Branches:\n")
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("%-40s %-40s\n", "Local Base", "Local Mod"))
		branches := *res.DiffBranches
		branches.Sort()
		for _, branch := range branches {
			buffer.WriteString(fmt.Sprintf("%s %s\n", branch.CommitHashFrom, branch.DiffHash))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
