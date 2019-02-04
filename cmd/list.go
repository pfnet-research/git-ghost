package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewListCommand())
}

type listFlags struct {
	hashFrom string
	hashTo   string
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
	return command
}

func runListCommitsCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: ghost.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &ghost.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func runListDiffCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: ghost.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListDiffBranchSpec: &ghost.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func runListAllCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: ghost.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &ghost.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
			ListDiffBranchSpec: &ghost.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func (flags listFlags) validate() error {
	return nil
}
