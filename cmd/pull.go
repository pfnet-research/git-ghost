package cmd

import (
	"errors"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/git"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewPullCommand())
}

type pullFlags struct {
	localBase  string
	baseCommit string
	force      bool
}

func NewPullCommand() *cobra.Command {
	var (
		flags pullFlags
	)

	var command = &cobra.Command{
		Use:   "pull [hash]",
		Short: "pull a ghost commit from remote repository and apply to your working git repository.",
		Long:  "pull a ghost commit from remote repository and apply to your working git repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := flags.Validate()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			hashArg := args[0]
			opts := ghost.PullOptions{
				WorkingEnvSpec: ghost.WorkingEnvSpec{
					SrcDir:          globalOpts.srcDir,
					GhostWorkingDir: globalOpts.ghostWorkDir,
					GhostRepo:       globalOpts.ghostRepo,
				},
				GhostSpec: ghost.GhostSpec{
					GhostPrefix:  globalOpts.ghostPrefix,
					RemoteBase:   flags.baseCommit,
					LocalModHash: hashArg,
				},
				ForceApply: flags.force,
			}

			if flags.localBase == "" {
				opts.GhostSpec.LocalBase = flags.baseCommit
			} else {
				opts.GhostSpec.LocalBase = flags.localBase
			}

			if err := ghost.Pull(opts); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&flags.baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	command.PersistentFlags().StringVar(&flags.localBase, "local-base", "", "git refspec used to create a local modification patch from (default \"value of --base-commit\")")
	command.PersistentFlags().BoolVar(&flags.force, "force", false, "try applying patch even when your working repository checked out different base commit")
	return command
}

func (flags pullFlags) Validate() error {
	if flags.baseCommit != "" {
		err := git.ValidateRefspec(".", flags.baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	return nil
}
