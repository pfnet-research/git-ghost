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

	"github.com/pfnet-research/git-ghost/pkg/ghost"
	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	"github.com/spf13/cobra"
)

type deleteFlags struct {
	hashFrom string
	hashTo   string
	all      bool
	dryrun   bool
}

func NewDeleteCommand() *cobra.Command {
	var (
		deleteFlags deleteFlags
	)

	var command = &cobra.Command{
		Use:   "delete",
		Short: "delete ghost branches of diffs.",
		Long:  "delete ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runDeleteDiffCommand(&deleteFlags),
	}
	command.AddCommand(&cobra.Command{
		Use:   "commits",
		Short: "delete ghost branches of commits.",
		Long:  "delete ghost branches of commits.",
		Args:  cobra.NoArgs,
		Run:   runDeleteCommitsCommand(&deleteFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "diff",
		Short: "delete ghost branches of diffs.",
		Long:  "delete ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runDeleteDiffCommand(&deleteFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "delete ghost branches of all types.",
		Long:  "delete ghost branches of all types.",
		Args:  cobra.NoArgs,
		Run:   runDeleteAllCommand(&deleteFlags),
	})
	command.PersistentFlags().StringVar(&deleteFlags.hashFrom, "from", "", "commit or diff hash to which ghost branches are deleted.")
	command.PersistentFlags().StringVar(&deleteFlags.hashTo, "to", "", "commit or diff hash from which ghost branches are deleted.")
	command.PersistentFlags().BoolVar(&deleteFlags.all, "all", false, "flag to ensure multiple ghost branches.")
	command.PersistentFlags().BoolVar(&deleteFlags.dryrun, "dry-run", false, "If true, only print the branch names that would be deleted, without deleting them.")
	return command
}

func runDeleteCommitsCommand(flags *deleteFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString())
	}
}

func runDeleteDiffCommand(flags *deleteFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString())
	}
}

func runDeleteAllCommand(flags *deleteFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Print(res.PrettyString())
	}
}

func (flags deleteFlags) validate() errors.GitGhostError {
	if (flags.hashFrom == "" || flags.hashTo == "") && !flags.all && !flags.dryrun {
		return errors.Errorf("all must be set if multiple ghosts branches are deleted")
	}
	return nil
}
