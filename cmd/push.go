package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/types"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pushFlags struct {
	includedFilepaths []string
}

func init() {
	RootCmd.AddCommand(NewPushCommand())
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

func runPushCommitsCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushArg := newPushCommitsArg(args)
		if err := pushArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			CommitsBranchSpec: &types.CommitsBranchSpec{
				Prefix:        globalOpts.ghostPrefix,
				CommitishFrom: pushArg.commitsFrom,
				CommitishTo:   pushArg.commitsTo,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if result.CommitsBranch != nil {
			fmt.Printf(
				"%s %s\n",
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

func (arg pushDiffArg) validate() error {
	if err := nonEmpty("diff-from", arg.diffFrom); err != nil {
		return err
	}
	if err := isValidComittish("diff-from", arg.diffFrom); err != nil {
		return err
	}
	return nil
}

func runPushDiffCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushArg := newPushDiffArg(args)
		if err := pushArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			DiffBranchSpec: &types.DiffBranchSpec{
				Prefix:            globalOpts.ghostPrefix,
				ComittishFrom:     pushArg.diffFrom,
				IncludedFilepaths: flags.includedFilepaths,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if result.DiffBranch != nil {
			fmt.Printf(result.DiffBranch.DiffHash)
		}
	}
}

func runPushAllCommand(flags *pushFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pushCommitsArg := newPushCommitsArg(args[0:1])
		if err := pushCommitsArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		pushDiffArg := newPushDiffArg(args[1:])
		if err := pushDiffArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		options := ghost.PushOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			CommitsBranchSpec: &types.CommitsBranchSpec{
				Prefix:        globalOpts.ghostPrefix,
				CommitishFrom: pushCommitsArg.commitsFrom,
				CommitishTo:   pushCommitsArg.commitsTo,
			},
			DiffBranchSpec: &types.DiffBranchSpec{
				Prefix:            globalOpts.ghostPrefix,
				ComittishFrom:     pushDiffArg.diffFrom,
				IncludedFilepaths: flags.includedFilepaths,
			},
		}

		result, err := ghost.Push(options)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if result.CommitsBranch != nil {
			fmt.Printf(
				"%s %s\n",
				result.CommitsBranch.CommitHashFrom,
				result.CommitsBranch.CommitHashTo,
			)
		}

		if result.DiffBranch != nil {
			fmt.Printf(result.DiffBranch.DiffHash)
		}
	}
}
