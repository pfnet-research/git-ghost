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
	localBase string
	force     bool
}

func NewPullCommand() *cobra.Command {
	var (
		pullFlags pullFlags
	)

	var command = &cobra.Command{
		Use:   "pull [hash]",
		Short: "pull a ghost commit from remote repository and apply to your working git repository.",
		Long:  "pull a ghost commit from remote repository and apply to your working git repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hashArg := args[0]
			opts := ghost.PullOptions{
				WorkingEnvSpec: ghost.WorkingEnvSpec{
					SrcDir:          globalOpts.srcDir,
					GhostWorkingDir: globalOpts.ghostWorkDir,
					GhostRepo:       globalOpts.ghostRepo,
				},
				GhostSpec: ghost.GhostSpec{
					GhostPrefix:  globalOpts.ghostPrefix,
					RemoteBase:   globalOpts.baseCommit,
					LocalModHash: hashArg,
				},
				ForceApply: pullFlags.force,
			}

			if pullFlags.localBase == "" {
				opts.GhostSpec.LocalBase = globalOpts.baseCommit
			} else {
				opts.GhostSpec.LocalBase = pullFlags.localBase
			}

			if err := ghost.Pull(opts); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&pullFlags.localBase, "local-base", "", "git refspec used to create a local modification patch from (default \"value of --base-commit\")")
	command.PersistentFlags().BoolVar(&pullFlags.force, "force", false, "try applying patch even when your working repository checked out different base commit")
	return command
}
