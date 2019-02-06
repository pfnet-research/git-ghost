package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

// ListOptions represents arg for List func
type ListOptions struct {
	WorkingEnvSpec
	*ListCommitsBranchSpec
	*ListDiffBranchSpec
}

// ListResult contains results of List func

type ListResult struct {
	LocalBaseBranches LocalBaseBranches
	LocalModBranches  LocalModBranches
}

// List returns ghost branches list per ghost branch type
func List(options ListOptions) (*ListResult, error) {
	log.WithFields(util.ToFields(options)).Debug("list command with")

	res := ListResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		res.LocalBaseBranches = branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		res.LocalModBranches = branches
	}

	res.LocalBaseBranches.Sort()
	res.LocalModBranches.Sort()

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *ListResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("Local Base Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalBaseBranches {
		buffer.WriteString(fmt.Sprintf("%s => %s\n", branch.RemoteBaseCommit, branch.LocalBaseCommit))
	}
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("\n")
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("Local Mod Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalModBranches {
		buffer.WriteString(fmt.Sprintf("%s -> %s\n", branch.LocalBaseCommit, branch.LocalModHash))
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("\n")
	}
	return buffer.String()
}
