package cmd

import (
	"git-ghost/pkg/ghost"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewShowCommand())
}

type showFlags struct {
	localBase string
}

func NewShowCommand() *cobra.Command {
	var (
		showFlags showFlags
	)
	var command = &cobra.Command{
		Use:   "show [hash]",
		Short: "show ghost commits on remote repository.",
		Long:  "show ghost commits on remote repository.",
		Run: func(cmd *cobra.Command, args []string) {
			hashArg := args[0]
			opts := ghost.ShowOptions{
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
			}

			if showFlags.localBase == "" {
				opts.GhostSpec.LocalBase = globalOpts.baseCommit
			} else {
				opts.GhostSpec.LocalBase = showFlags.localBase
			}

			if err := ghost.Show(opts); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&showFlags.localBase, "local-base", "", "git refspec used to create a local modification patch from (default \"value of --base-commit\")")
	return command
}
