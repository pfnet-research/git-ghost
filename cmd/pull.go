package cmd

import (
	"git-ghost/pkg/ghost"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewPullCommand())
}

type pullFlags struct {
	forceApply bool
}

func NewPullCommand() *cobra.Command {
	var (
		flags pullFlags
	)
	command := &cobra.Command{
		Use:   "pull [from-hash(default=HEAD)] [diff-hash]",
		Short: "pull commits(hash1...hash2), diff(hash...current state) from ghost repo and apply them to working dir",
		Long:  "pull commits or diff or all from ghost repo and apply them to working dir.  If you didn't specify any subcommand, this commands works as an alias for 'pull diff' command.",
		Args:  cobra.MinimumNArgs(1),
		Run:   runPullDiffCommand(flags),
	}
	command.PersistentFlags().BoolVarP(&flags.forceApply, "force", "f", true, "force apply pulled ghost branches to working dir")

	command.AddCommand(&cobra.Command{
		Use:   "diff [diff-from-hash(default=HEAD)] [diff-hash]",
		Short: "pull diff from ghost repo and apply it to working dir",
		Long:  "pull diff from [diff-from-hash] to [diff-hash] from your ghost repo and apply it to working dir",
		Args:  cobra.MinimumNArgs(1),
		Run:   runPullDiffCommand(flags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "commits [from-hash(default=HEAD)] [to-hash]",
		Short: "pull commits from ghost repo and apply it to working dir",
		Long:  "pull commits from [from-hash] to [to-hash] from your ghost repo and apply it to working dir",
		Args:  cobra.MinimumNArgs(1),
		Run:   runPullCommitsCommand(flags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "all [from-hash(default=HEAD)] [to-hash] [diff-hash]",
		Short: "pull both commits and diff from ghost repo and apply them to working dir sequentially",
		Long:  "pull commits([from-hash]...[to-hash]) and diff([to-hash]...[diff-hash]) and apply them to working dir sequentially",
		Args:  cobra.MinimumNArgs(2),
		Run:   runPullAllCommand(flags),
	})
	return command
}

type pullCommitsArg struct {
	commitsFrom string
	commitsTo   string
}

func newPullCommitsArg(args []string) pullCommitsArg {
	arg := pullCommitsArg{
		commitsFrom: "HEAD",
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

func (arg pullCommitsArg) validate() error {
	if err := nonEmpty("commit-from", arg.commitsFrom); err != nil {
		return err
	}
	if err := nonEmpty("commit-to", arg.commitsTo); err != nil {
		return err
	}
	return nil
}

func runPullCommitsCommand(flags pullFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		arg := newPullCommitsArg(args)
		if err := arg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		options := ghost.PullCommitsOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			LocalBaseBranchSpec: ghost.LocalBaseBranchSpec{
				Prefix:              globalOpts.ghostPrefix,
				RemoteBaseCommitish: arg.commitsFrom,
				LocalBaseCommitish:  arg.commitsTo,
			},
			ForceApply: flags.forceApply,
		}

		err := ghost.PullCommits(options, nil)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

type pullDiffArg struct {
	diffFrom string
	diffHash string
}

func newPullDiffArg(args []string) pullDiffArg {
	arg := pullDiffArg{
		diffFrom: "HEAD",
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

func (arg pullDiffArg) validate() error {
	if err := nonEmpty("diff-from-hash", arg.diffFrom); err != nil {
		return err
	}
	if err := nonEmpty("diff-hash", arg.diffHash); err != nil {
		return err
	}
	return nil
}

func runPullDiffCommand(flags pullFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		arg := newPullDiffArg(args)
		if err := arg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		options := ghost.PullDiffOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			PullableLocalModBranchSpec: ghost.PullableLocalModBranchSpec{
				LocalModBranchSpec: ghost.LocalModBranchSpec{
					Prefix:             globalOpts.ghostPrefix,
					LocalBaseCommitish: arg.diffFrom,
				},
				LocalModHash: arg.diffHash,
			},
			ForceApply: flags.forceApply,
		}

		err := ghost.PullDiff(options, nil)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func runPullAllCommand(flags pullFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var pullCommitsArg pullCommitsArg
		var pullDiffArg pullDiffArg

		if len(args) >= 3 {
			pullCommitsArg = newPullCommitsArg(args[0:2])
			pullDiffArg = newPullDiffArg(args[1:])
		} else { // len(args) == 2
			pullCommitsArg = newPullCommitsArg(args[0:1])
			pullDiffArg = newPullDiffArg(args)
		}

		if err := pullCommitsArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		if err := pullDiffArg.validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		options := ghost.PullOptions{
			WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
			LocalBaseBranchSpec: ghost.LocalBaseBranchSpec{
				Prefix:              globalOpts.ghostPrefix,
				RemoteBaseCommitish: pullCommitsArg.commitsFrom,
				LocalBaseCommitish:  pullCommitsArg.commitsTo,
			},
			PullableLocalModBranchSpec: ghost.PullableLocalModBranchSpec{
				LocalModBranchSpec: ghost.LocalModBranchSpec{
					Prefix:             globalOpts.ghostPrefix,
					LocalBaseCommitish: pullDiffArg.diffFrom,
				},
				LocalModHash: pullDiffArg.diffHash,
			},
			ForceApply: flags.forceApply,
		}

		err := ghost.Pull(options)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}
