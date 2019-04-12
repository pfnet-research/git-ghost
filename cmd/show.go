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

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   fmt.Sprintf("show [from-hash(default=%s)] [diff-hash]", ghostDefaultPushFromHash),
		Short: "show commits(hash1...hash2), diff(hash...current state) in ghost repo",
		Long:  "show commits or diff or all from ghost repo.  If you didn't specify any subcommand, this commands works as an alias for 'show diff' command.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowDiffCommand,
	}
	command.AddCommand(&cobra.Command{
		Use:   fmt.Sprintf("diff [diff-from-hash(default=%s)] [diff-hash]", ghostDefaultPushFromHash),
		Short: "show diff in ghost repo ",
		Long:  "show diff from [diff-from-hash] to [diff-hash] in ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowDiffCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   fmt.Sprintf("commits [from-hash(default=%s)] [to-hash]", ghostDefaultPushFromHash),
		Short: "show commits in ghost repo",
		Long:  "show commits from [from-hash] to [to-hash] in ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowCommitsCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   fmt.Sprintf("all [from-hash(default=%s)] [to-hash] [diff-hash]", ghostDefaultPushFromHash),
		Short: "show both commits and diff in ghost repo",
		Long:  "show commits([from-hash]...[to-hash]) and diff([to-hash]...[diff-hash]) in ghost repo",
		Args:  cobra.RangeArgs(2, 3),
		Run:   runShowAllCommand,
	})
	return command
}

type showCommitsArg struct {
	commitsFrom string
	commitsTo   string
}

func newShowCommitsArg(args []string) showCommitsArg {
	arg := showCommitsArg{
		commitsFrom: ghostDefaultPushFromHash,
		commitsTo:   "",
	}

	if len(args) >= 2 {
		arg.commitsFrom = args[0]
		arg.commitsTo = args[1]
		return arg
	}

	if len(args) >= 1 {
		arg.commitsTo = args[0]
		return arg
	}

	return arg
}

func (arg showCommitsArg) validate() errors.GitGhostError {
	if err := nonEmpty("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := nonEmpty("commit-to", arg.commitsTo); err != nil {
		return err
	}
	return nil
}

func runShowCommitsCommand(cmd *cobra.Command, args []string) {
	arg := newShowCommitsArg(args)
	if err := arg.validate(); err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		CommitsBranchSpec: &types.CommitsBranchSpec{
			Prefix:         globalOpts.ghostPrefix,
			CommittishFrom: arg.commitsFrom,
			CommittishTo:   arg.commitsTo,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}
}

type showDiffArg struct {
	diffFrom string
	diffHash string
}

func newShowDiffArg(defaultFromHash string, args []string) showDiffArg {
	arg := showDiffArg{
		diffFrom: defaultFromHash,
		diffHash: "",
	}

	if len(args) >= 2 {
		arg.diffFrom = args[0]
		arg.diffHash = args[1]
		return arg
	}

	if len(args) >= 1 {
		arg.diffHash = args[0]
		return arg
	}

	return arg
}

func (arg showDiffArg) validate() errors.GitGhostError {
	if err := nonEmpty("diff-from-hash", arg.diffFrom); err != nil {
		return err
	}
	if err := nonEmpty("diff-hash", arg.diffHash); err != nil {
		return err
	}
	return nil
}

func runShowDiffCommand(cmd *cobra.Command, args []string) {
	arg := newShowDiffArg(ghostDefaultPushFromHash, args)
	if err := arg.validate(); err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		PullableDiffBranchSpec: &types.PullableDiffBranchSpec{
			Prefix:         globalOpts.ghostPrefix,
			CommittishFrom: arg.diffFrom,
			DiffHash:       arg.diffHash,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}
}

func runShowAllCommand(cmd *cobra.Command, args []string) {
	var showCommitsArg showCommitsArg
	var showDiffArg showDiffArg

	switch len(args) {
	case 3:
		showCommitsArg = newShowCommitsArg(args[0:2])
		showDiffArg = newShowDiffArg("HEAD", args[1:])
	case 2:
		showCommitsArg = newShowCommitsArg(args[0:1])
		showDiffArg = newShowDiffArg("HEAD", args)
	default:
		log.Error(cmd.Args(cmd, args))
		os.Exit(1)
	}

	if err := showCommitsArg.validate(); err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}
	if err := showDiffArg.validate(); err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		CommitsBranchSpec: &types.CommitsBranchSpec{
			Prefix:         globalOpts.ghostPrefix,
			CommittishFrom: showCommitsArg.commitsFrom,
			CommittishTo:   showCommitsArg.commitsTo,
		},
		PullableDiffBranchSpec: &types.PullableDiffBranchSpec{
			Prefix:         globalOpts.ghostPrefix,
			CommittishFrom: showDiffArg.diffFrom,
			DiffHash:       showDiffArg.diffHash,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		errors.LogErrorWithStack(err)
		os.Exit(1)
	}
}
