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

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pfnet-research/git-ghost/pkg/ghost"
	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	"github.com/spf13/cobra"
)

var outputTypes = []string{"only-from", "only-to"}
var regexpOutputPattern = regexp.MustCompile("^(|" + strings.Join(outputTypes, "|") + ")$")

type listFlags struct {
	hashFrom  string
	hashTo    string
	noHeaders bool
	output    string
}

func NewListCommand() *cobra.Command {
	var (
		listFlags listFlags
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "list ghost branches of diffs.",
		Long:  "list ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runListDiffCommand(&listFlags),
	}
	command.AddCommand(&cobra.Command{
		Use:   "commits",
		Short: "list ghost branches of commits.",
		Long:  "list ghost branches of commits.",
		Args:  cobra.NoArgs,
		Run:   runListCommitsCommand(&listFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "diff",
		Short: "list ghost branches of diffs.",
		Long:  "list ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runListDiffCommand(&listFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "list ghost branches of all types.",
		Long:  "list ghost branches of all types.",
		Args:  cobra.NoArgs,
		Run:   runListAllCommand(&listFlags),
	})
	command.PersistentFlags().StringVar(&listFlags.hashFrom, "from", "", "commit or diff hash to which ghost branches are listed.")
	command.PersistentFlags().StringVar(&listFlags.hashTo, "to", "", "commit or diff hash from which ghost branches are listed.")
	command.PersistentFlags().BoolVar(&listFlags.noHeaders, "no-headers", false, "When using the default, only-from or only-to output format, don't print headers (default print headers).")
	command.PersistentFlags().StringVarP(&listFlags.output, "output", "o", "", "Output format. One of: only-from|only-to")
	return command
}

func runListCommitsCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &types.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func runListDiffCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListDiffBranchSpec: &types.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func runListAllCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &types.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
			ListDiffBranchSpec: &types.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func (flags listFlags) validate() errors.GitGhostError {
	if !regexpOutputPattern.MatchString(flags.output) {
		return errors.Errorf("output must be one of %v", outputTypes)
	}
	return nil
}
