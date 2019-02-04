package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewPushCommand())
}

func NewPushCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "push [from-hash(default=HEAD)]",
		Short: "push commits(hash1...hash2), diff(hash...current state) to your ghost repo",
		Long:  "push commits or diff or all to your ghost repo.  If you didn't specify any subcommand, this commands works as an alias for 'push diff' command.",
		Args:  cobra.RangeArgs(0, 1),
		Run:   runPushDiffCommand,
	}
	command.AddCommand(&cobra.Command{
		Use:   "commits [from-hash] [to-hash(default=HEAD)]",
		Short: "push commits between two commits to your ghost repo",
		Long:  "push all the commits between [from-hash]...[to-hash] to your ghost repo.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runPushCommitsCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   "diff [from-hash(default=HEAD)]",
		Short: "push diff from a commit to current state of your working dir to your ghost repo",
		Long:  "push diff from [from-hash] to current state of your working dir to your ghost repo.  please be noted that this pushes only diff, which means that it doesn't save any commits information.",
		Args:  cobra.RangeArgs(0, 1),
		Run:   runPushDiffCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   "all [commits-from-hash] [diff-from-hash(default=HEAD)]",
		Short: "push both commits and diff to your ghost repo",
		Long:  "push both commits([commits-from-hash]...[diff-from-hash]) and diff([diff-from-hash]...current state) to your ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			runPushCommitsCommand(cmd, args)
			runPushDiffCommand(cmd, args[1:])
		},
	})
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

func (arg pushCommitsArg) validate() error {
	if err := nonEmpty("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := nonEmpty("commit-to", arg.commitsTo); err != nil {
		return err
	}
	if err := isValidComittish("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := isValidComittish("commit-to", arg.commitsTo); err != nil {
		return err
	}
	return nil
}

func runPushCommitsCommand(cmd *cobra.Command, args []string) {
	pushArg := newPushCommitsArg(args)
	if err := pushArg.validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	options := ghost.PushCommitsOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		LocalBaseBranchSpec: ghost.LocalBaseBranchSpec{
			Prefix:              globalOpts.ghostPrefix,
			RemoteBaseCommitish: pushArg.commitsFrom,
			LocalBaseCommitish:  pushArg.commitsTo,
		},
	}

	result, err := ghost.PushCommits(options)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if result.LocalBaseBranch != nil {
		fmt.Printf(
			"%s %s\n",
			result.LocalBaseBranch.RemoteBaseCommit,
			result.LocalBaseBranch.LocalBaseCommit,
		)
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
func (arg pushDiffArg) validate() error {
	if err := nonEmpty("diff-from", arg.diffFrom); err != nil {
		return err
	}
	if err := isValidComittish("diff-from", arg.diffFrom); err != nil {
		return err
	}
	return nil
}

func runPushDiffCommand(cmd *cobra.Command, args []string) {
	pushArg := newPushDiffArg(args)
	if err := pushArg.validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	options := ghost.PushDiffOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		LocalModBranchSpec: ghost.LocalModBranchSpec{
			Prefix:             globalOpts.ghostPrefix,
			LocalBaseCommitish: pushArg.diffFrom,
		},
	}

	result, err := ghost.PushDiff(options)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if result.LocalModBranch != nil {
		fmt.Printf(result.LocalModBranch.LocalModHash)
	}
}
