package cmd

import (
	"fmt"
	"os"

	"git-ghost/pkg/ghost"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewPullCommand())
}

type pullFlags struct {
	localBase string
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
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostPrefix:     globalOpts.ghostPrefix,
				GhostRepo:       globalOpts.ghostRepo,
				RemoteBase:      globalOpts.baseCommit,
				Hash:            hashArg,
			}

			if pullFlags.localBase == "" {
				opts.LocalBase = globalOpts.baseCommit
			} else {
				opts.LocalBase = pullFlags.localBase
			}

			err := ghost.Pull(opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&pullFlags.localBase, "local-base", "", "git refspec used to create a local modification patch from (default \"value of --base-commit\")")
	return command
}
