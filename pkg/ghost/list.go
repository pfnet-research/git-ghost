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
	*types.LocalBaseBranches
	*types.LocalModBranches
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

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *ListResult) PrettyString(headers bool, output string) string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if res.LocalBaseBranches != nil {
		branches := *res.LocalBaseBranches
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
				columns = append(columns, branch.RemoteBaseCommit)
			case "only-to":
				columns = append(columns, branch.LocalBaseCommit)
			default:
				columns = append(columns, branch.RemoteBaseCommit)
				columns = append(columns, branch.LocalBaseCommit)
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		if headers {
			buffer.WriteString("\n")
		}
	}
	if res.LocalModBranches != nil {
		branches := *res.LocalModBranches
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
				columns = append(columns, branch.LocalBaseCommit)
			case "only-to":
				columns = append(columns, branch.LocalModHash)
			default:
				columns = append(columns, branch.LocalBaseCommit)
				columns = append(columns, branch.LocalModHash)
			}
			buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(columns, " ")))
		}
		if headers {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
