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
	RootCmd.AddCommand(NewShowCommand())
}

type showFlags struct {
	baseCommit string
	localBase  string
}

func NewShowCommand() *cobra.Command {
	var (
		flags showFlags
	)
	var command = &cobra.Command{
		Use:   "show [hash]",
		Short: "show ghost commits on remote repository.",
		Long:  "show ghost commits on remote repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := flags.Validate()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			hashArg := args[0]
			opts := ghost.ShowOptions{
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
				Writer: os.Stdout,
			}

			if flags.localBase == "" {
				opts.GhostSpec.LocalBase = flags.baseCommit
			} else {
				opts.GhostSpec.LocalBase = flags.localBase
			}

			if err := ghost.Show(opts); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&flags.baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	command.PersistentFlags().StringVar(&flags.localBase, "local-base", "", "git refspec used to create a local modification patch from (default \"value of --base-commit\")")
	return command
}

func (flags showFlags) Validate() error {
	if flags.baseCommit != "" {
		err := git.ValidateRefspec(".", flags.baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	return nil
}
