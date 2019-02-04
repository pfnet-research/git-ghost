package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type ListOptions struct {
	WorkingEnvSpec
	GhostPrefix string
	BaseCommit  string
}

type ListResult struct {
	LocalBaseBranches LocalBaseBranches
	LocalModBranches  LocalModBranches
}

func List(options ListOptions) (*ListResult, error) {
	log.WithFields(util.ToFields(options)).Debug("list command with")

	baseCommit, err := git.ResolveComittish(options.SrcDir, options.BaseCommit)
	if err != nil {
		return nil, err
	}
	branchNames, err := git.ListGhostBranchNames(options.GhostRepo, options.GhostPrefix, baseCommit, "")
	if err != nil {
		return nil, err
	}

	var res ListResult
	for _, name := range branchNames {
		branch := CreateGhostBranchByName(name)
		if br, ok := branch.(*LocalBaseBranch); ok {
			res.LocalBaseBranches = append(res.LocalBaseBranches, *br)
			continue
		}
		if br, ok := branch.(*LocalModBranch); ok {
			res.LocalModBranches = append(res.LocalModBranches, *br)
			continue
		}
		log.Warning(fmt.Sprintf("unknown branch: %s", name))
	}
	res.LocalBaseBranches.Sort()
	res.LocalModBranches.Sort()

	return &res, nil
}

func (res *ListResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("Local Base Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalBaseBranches {
		buffer.WriteString(fmt.Sprintf("%s\n", branch.BranchName()))
	}
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("\n")
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("Local Mod Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalModBranches {
		buffer.WriteString(fmt.Sprintf("%s\n", branch.BranchName()))
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("\n")
	}
	return buffer.String()
}
