package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/git"
	"os"

	"github.com/spf13/cobra"
)

type pushFlags struct {
	localBase string
}

func init() {
	RootCmd.AddCommand(NewPushCommand())
}

func NewPushCommand() *cobra.Command {
	var (
		pushOpts pushFlags
	)
	command := &cobra.Command{
		Use:   "push",
		Short: "generate and push a ghost commit to remote repository",
		Long:  "generate and push a ghost commit to remote repository",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			hash, err := ghost.Push(ghost.PushOptions{
				SrcDir:      globalOpts.srcDir,
				GhostPrefix: globalOpts.ghostPrefix,
				GhostRepo:   globalOpts.ghostRepo,
				RemoteBase:  globalOpts.baseCommit,
				LocalBase:   pushOpts.localBase,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(hash)
		},
	}
	command.PersistentFlags().StringVar(&pushOpts.localBase, "local-base", "HEAD", "git refspec used to create a local modification patch from")
	return command
}

func (flags pushFlags) Validate() error {
	if flags.localBase == "" {
		return errors.New("local-base must be specified")
	}
	err := git.ValidateRefspec(".", flags.localBase)
	if err != nil {
		return errors.New("local-base is not a valid object")
	}
	return nil
}
