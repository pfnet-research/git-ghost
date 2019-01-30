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
}

func NewListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "list ghost commits on remote repository.",
		Long:  "list ghost commits on remote repository.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			opts := ghost.ListOptions{
				WorkingEnvSpec: ghost.WorkingEnvSpec{
					SrcDir:          globalOpts.srcDir,
					GhostWorkingDir: globalOpts.ghostWorkDir,
					GhostRepo:       globalOpts.ghostRepo,
				},
				GhostPrefix: globalOpts.ghostPrefix,
				BaseCommit:  globalOpts.baseCommit,
			}

			res, err := ghost.List(opts)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Printf(res.PrettyString())
		},
	}
	return command
}
