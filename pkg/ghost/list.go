// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ghost

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	log "github.com/sirupsen/logrus"
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
func List(options ListOptions) (*ListResult, errors.GitGhostError) {
	log.WithFields(util.ToFields(options)).Debug("list command with")

	res := ListResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res.CommitsBranches = &branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, errors.WithStack(err)
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
