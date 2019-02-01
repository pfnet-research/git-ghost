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

type pushFlags struct {
	baseCommit string
	localBase  string
}

func init() {
	RootCmd.AddCommand(NewPushCommand())
}

func NewPushCommand() *cobra.Command {
	var (
		flags pushFlags
	)
	command := &cobra.Command{
		Use:   "push",
		Short: "generate and push a ghost commit to remote repository",
		Long:  "generate and push a ghost commit to remote repository",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := flags.Validate()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			resp, err := ghost.Push(ghost.PushOptions{
				WorkingEnvSpec: ghost.WorkingEnvSpec{
					SrcDir:          globalOpts.srcDir,
					GhostWorkingDir: globalOpts.ghostWorkDir,
					GhostRepo:       globalOpts.ghostRepo,
				},
				GhostPrefix: globalOpts.ghostPrefix,
				RemoteBase:  flags.baseCommit,
				LocalBase:   flags.localBase,
			})
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if resp.LocalModBranch != nil {
				fmt.Println(resp.LocalModBranch.LocalModHash)
			}
		},
	}
	command.PersistentFlags().StringVar(&flags.baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	command.PersistentFlags().StringVar(&flags.localBase, "local-base", "HEAD", "git refspec used to create a local modification patch from")
	return command
}

func (flags pushFlags) Validate() error {
	if flags.baseCommit != "" {
		err := git.ValidateRefspec(".", flags.baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	if flags.localBase == "" {
		return errors.New("local-base must be specified")
	}
	err := git.ValidateRefspec(".", flags.localBase)
	if err != nil {
		return errors.New("local-base is not a valid object")
	}
	return nil
}
