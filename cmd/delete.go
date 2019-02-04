package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewDeleteCommand())
}

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
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func runDeleteDiffCommand(flags *deleteFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func runDeleteAllCommand(flags *deleteFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		opts := ghost.DeleteOptions{
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
			Dryrun: flags.dryrun,
		}

		res, err := ghost.Delete(opts)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString())
	}
}

func (flags deleteFlags) validate() error {
	if (flags.hashFrom == "" || flags.hashTo == "") && !flags.all && !flags.dryrun {
		return fmt.Errorf("all must be set if multiple ghosts branches are deleted")
	}
	return nil
}
