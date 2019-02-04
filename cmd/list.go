package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/git"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewListCommand())
}

type listFlags struct {
	baseCommit string
}

func NewListCommand() *cobra.Command {
	var (
		flags listFlags
	)
	var command = &cobra.Command{
		Use:   "list",
		Short: "list ghost commits on remote repository.",
		Long:  "list ghost commits on remote repository.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := flags.Validate()
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
				GhostPrefix: globalOpts.ghostPrefix,
				BaseCommit:  flags.baseCommit,
			}

			res, err := ghost.List(opts)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Printf(res.PrettyString())
		},
	}
	command.PersistentFlags().StringVar(&flags.baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	return command
}

func (flags listFlags) Validate() error {
	if flags.baseCommit != "" {
		err := git.ValidateComittish(".", flags.baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	return nil
}
