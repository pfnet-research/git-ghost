package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// ListOptions represents arg for List func
type ListOptions struct {
	types.WorkingEnvSpec
	*types.ListCommitsBranchSpec
	*types.ListDiffBranchSpec
}

// ListResult contains results of List func

type ListResult struct {
	*types.CommitsBranches
	*types.DiffBranches
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

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *ListResult) PrettyString(headers bool, output string) string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if res.CommitsBranches != nil {
		branches := *res.CommitsBranches
		branches.Sort()
		if headers {
			buffer.WriteString("Local Base Branches:\n")
			buffer.WriteString("\n")
			columns := []string{}
			switch output {
			case "only-from":
				columns = append(columns, fmt.Sprintf("%-40s", "Remote Base"))
			case "only-to":
				columns = append(columns, fmt.Sprintf("%-40s", "Local Base"))
			default:
				columns = append(columns, fmt.Sprintf("%-40s", "Remote Base"))
				columns = append(columns, fmt.Sprintf("%-40s", "Local Base"))
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		for _, branch := range branches {
			columns := []string{}
			switch output {
			case "only-from":
				columns = append(columns, branch.CommitHashFrom)
			case "only-to":
				columns = append(columns, branch.CommitHashTo)
			default:
				columns = append(columns, branch.CommitHashFrom)
				columns = append(columns, branch.CommitHashTo)
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		if headers {
			buffer.WriteString("\n")
		}
	}
	if res.DiffBranches != nil {
		branches := *res.DiffBranches
		branches.Sort()
		if headers {
			buffer.WriteString("Local Mod Branches:\n")
			buffer.WriteString("\n")
			columns := []string{}
			switch output {
			case "only-from":
				columns = append(columns, fmt.Sprintf("%-40s", "Local Base"))
			case "only-to":
				columns = append(columns, fmt.Sprintf("%-40s", "Local Mod"))
			default:
				columns = append(columns, fmt.Sprintf("%-40s", "Local Base"))
				columns = append(columns, fmt.Sprintf("%-40s", "Local Mod"))
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		for _, branch := range branches {
			columns := []string{}
			switch output {
			case "only-from":
				columns = append(columns, branch.CommitHashFrom)
			case "only-to":
				columns = append(columns, branch.DiffHash)
			default:
				columns = append(columns, branch.CommitHashFrom)
				columns = append(columns, branch.DiffHash)
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		if headers {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
