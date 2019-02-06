package cmd

import (
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/types"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewShowCommand())
}

func NewShowCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "show [from-hash(default=HEAD)] [diff-hash]",
		Short: "show commits(hash1...hash2), diff(hash...current state) in ghost repo",
		Long:  "show commits or diff or all from ghost repo.  If you didn't specify any subcommand, this commands works as an alias for 'show diff' command.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowDiffCommand,
	}
	command.AddCommand(&cobra.Command{
		Use:   "diff [diff-from-hash(default=HEAD)] [diff-hash]",
		Short: "show diff in ghost repo ",
		Long:  "show diff from [diff-from-hash] to [diff-hash] in ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowDiffCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   "commits [from-hash(default=HEAD)] [to-hash]",
		Short: "show commits in ghost repo",
		Long:  "show commits from [from-hash] to [to-hash] in ghost repo",
		Args:  cobra.RangeArgs(1, 2),
		Run:   runShowCommitsCommand,
	})
	command.AddCommand(&cobra.Command{
		Use:   "all [from-hash(default=HEAD)] [to-hash] [diff-hash]",
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

func (arg showCommitsArg) validate() error {
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
		log.Error(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		LocalBaseBranchSpec: &types.LocalBaseBranchSpec{
			Prefix:              globalOpts.ghostPrefix,
			RemoteBaseCommitish: arg.commitsFrom,
			LocalBaseCommitish:  arg.commitsTo,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

type showDiffArg struct {
	diffFrom string
	diffHash string
}

func newShowDiffArg(args []string) showDiffArg {
	arg := showDiffArg{
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

func (arg showDiffArg) validate() error {
	if err := nonEmpty("diff-from-hash", arg.diffFrom); err != nil {
		return err
	}
	if err := nonEmpty("diff-hash", arg.diffHash); err != nil {
		return err
	}
	return nil
}

func runShowDiffCommand(cmd *cobra.Command, args []string) {
	arg := newShowDiffArg(args)
	if err := arg.validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		PullableLocalModBranchSpec: &types.PullableLocalModBranchSpec{
			LocalModBranchSpec: types.LocalModBranchSpec{
				Prefix:             globalOpts.ghostPrefix,
				LocalBaseCommitish: arg.diffFrom,
			},
			LocalModHash: arg.diffHash,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func runShowAllCommand(cmd *cobra.Command, args []string) {
	var showCommitsArg showCommitsArg
	var showDiffArg showDiffArg

	switch len(args) {
	case 3:
		showCommitsArg = newShowCommitsArg(args[0:2])
		showDiffArg = newShowDiffArg(args[1:])
	case 2:
		showCommitsArg = newShowCommitsArg(args[0:1])
		showDiffArg = newShowDiffArg(args)
	default:
		log.Error(cmd.Args(cmd, args))
		os.Exit(1)
	}

	if err := showCommitsArg.validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := showDiffArg.validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	options := ghost.ShowOptions{
		WorkingEnvSpec: globalOpts.WorkingEnvSpec(),
		LocalBaseBranchSpec: &types.LocalBaseBranchSpec{
			Prefix:              globalOpts.ghostPrefix,
			RemoteBaseCommitish: showCommitsArg.commitsFrom,
			LocalBaseCommitish:  showCommitsArg.commitsTo,
		},
		PullableLocalModBranchSpec: &types.PullableLocalModBranchSpec{
			LocalModBranchSpec: types.LocalModBranchSpec{
				Prefix:             globalOpts.ghostPrefix,
				LocalBaseCommitish: showDiffArg.diffFrom,
			},
			LocalModHash: showDiffArg.diffHash,
		},
		Writer: os.Stdout,
	}

	err := ghost.Show(options)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
