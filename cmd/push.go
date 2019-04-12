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

type pushFlags struct {
	includedFilepaths []string
	followSymlinks    bool
}

func NewPushCommand() *cobra.Command {
	var (
		flags pushFlags
	)
	command := &cobra.Command{
		Use:   "push [from-hash(default=HEAD)]",
		Short: "push commits(hash1...hash2), diff(hash...current state) to your ghost repo",
		Long:  "push commits or diff or all to your ghost repo.  If you didn't specify any subcommand, this commands works as an alias for 'push diff' command.",
		Args:  cobra.RangeArgs(0, 1),
		Run:   runPushDiffCommand(&flags),
	}
	command.AddCommand(&cobra.Command{
		Use:   "commits [from-hash] [to-hash(default=HEAD)]",
		Short: "push commits between two commits to your ghost repo",
		Long:  "push all the commits between [from-hash]...[to-hash] to your ghost repo.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runPushCommitsCommand(&flags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "diff [from-hash(default=HEAD)]",
		Short: "push diff from a commit to current state of your working dir to your ghost repo",
		Long:  "push diff from [from-hash] to current state of your working dir to your ghost repo.  please be noted that this pushes only diff, which means that it doesn't save any commits information.",
		Args:  cobra.RangeArgs(0, 1),
		Run:   runPushDiffCommand(&flags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "all [commits-from-hash] [diff-from-hash(default=HEAD)]",
		Short: "push both commits and diff to your ghost repo",
		Long:  "push both commits([commits-from-hash]...[diff-from-hash]) and diff([diff-from-hash]...current state) to your ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runPushAllCommand(&flags),
	})

	command.PersistentFlags().StringSliceVarP(&flags.includedFilepaths, "include", "I", []string{}, "include a non-indexed file, this flag can be repeated to specify multiple files.")
	command.PersistentFlags().BoolVar(&flags.followSymlinks, "follow-symlinks", false, "follow symlinks inside the repository.")

	return command
}

type pushCommitsArg struct {
	commitsFrom string
	commitsTo   string
}

func newPushCommitsArg(args []string) pushCommitsArg {
	pushCommitsArg := pushCommitsArg{
		commitsFrom: "",
		commitsTo:   "HEAD",
	}
	if len(args) >= 1 {
		pushCommitsArg.commitsFrom = args[0]
	}
	if len(args) >= 2 {
		pushCommitsArg.commitsTo = args[1]
	}
	return pushCommitsArg
}

func (arg pushCommitsArg) validate() errors.GitGhostError {
	if err := nonEmpty("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := nonEmpty("commit-to", arg.commitsTo); err != nil {
		return err
	}
	if err := isValidCommittish("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := isValidCommittish("commit-to", arg.commitsTo); err != nil {
		return err
	}
	return nil
}

func runPushCommitsCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushArg := newPushCommitsArg(args)
		if err := pushArg.validate(); err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			CommitsBranchSpec: &types.CommitsBranchSpec{
				Prefix:         globalOpts.ghostPrefix,
				CommittishFrom: pushArg.commitsFrom,
				CommittishTo:   pushArg.commitsTo,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}

		if result.CommitsBranch != nil {
			fmt.Printf(
				"%s %s",
				result.CommitsBranch.CommitHashFrom,
				result.CommitsBranch.CommitHashTo,
			)
		}
	}
}

type pushDiffArg struct {
	diffFrom string
}

func newPushDiffArg(args []string) pushDiffArg {
	pushDiffArg := pushDiffArg{
		diffFrom: "HEAD",
	}
	if len(args) >= 1 {
		pushDiffArg.diffFrom = args[0]
	}
	return pushDiffArg
}

func (arg pushDiffArg) validate() errors.GitGhostError {
	if err := nonEmpty("diff-from", arg.diffFrom); err != nil {
		return err
	}
	if err := isValidCommittish("diff-from", arg.diffFrom); err != nil {
		return err
	}
	return nil
}

func runPushDiffCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushArg := newPushDiffArg(args)
		if err := pushArg.validate(); err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			DiffBranchSpec: &types.DiffBranchSpec{
				Prefix:            globalOpts.ghostPrefix,
				CommittishFrom:    pushArg.diffFrom,
				IncludedFilepaths: flags.includedFilepaths,
				FollowSymlinks:    flags.followSymlinks,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}

		if result.DiffBranch != nil {
			fmt.Printf(
				"%s %s",
				result.DiffBranch.CommitHashFrom,
				result.DiffBranch.DiffHash,
			)
			fmt.Print("\n")
		}
	}
}

func runPushAllCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushCommitsArg := newPushCommitsArg(args[0:1])
		if err := pushCommitsArg.validate(); err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}

		pushDiffArg := newPushDiffArg(args[1:])
		if err := pushDiffArg.validate(); err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}

		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			CommitsBranchSpec: &types.CommitsBranchSpec{
				Prefix:         globalOpts.ghostPrefix,
				CommittishFrom: pushCommitsArg.commitsFrom,
				CommittishTo:   pushCommitsArg.commitsTo,
			},
			DiffBranchSpec: &types.DiffBranchSpec{
				Prefix:            globalOpts.ghostPrefix,
				CommittishFrom:    pushDiffArg.diffFrom,
				IncludedFilepaths: flags.includedFilepaths,
				FollowSymlinks:    flags.followSymlinks,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}

		if result.CommitsBranch != nil {
			fmt.Printf(
				"%s %s",
				result.CommitsBranch.CommitHashFrom,
				result.CommitsBranch.CommitHashTo,
			)
			if result.DiffBranch != nil {
				fmt.Print("\n")
			}
		}

		if result.DiffBranch != nil {
			fmt.Printf(
				"%s %s",
				result.DiffBranch.CommitHashFrom,
				result.DiffBranch.DiffHash,
			)
			fmt.Print("\n")
		}
	}
}
