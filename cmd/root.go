package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version  string
	Revision string
)

var RootCmd = &cobra.Command{
	Use:   "git-ghost",
	Short: "git-ghost",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "version" {
			return nil
		}
		err := validateEnvironment()
		if err != nil {
			return err
		}
		err = validateGlobalFlags()
		if err != nil {
			return err
		}
		return nil
	},
}

// Global Flags
var ghostRepo string
var baseCommit string

func init() {
	cobra.OnInitialize()
	ghostRepoEnv := os.Getenv("GHOST_REPO")
	RootCmd.PersistentFlags().StringVar(&ghostRepo, "ghost-repo", ghostRepoEnv, "git refspec for ghost commits repository")
	RootCmd.PersistentFlags().StringVar(&baseCommit, "base-commit", "HEAD", "base commit hash for generating ghost commit.")
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of git-ghost",
	Long:  `Print the version number of git-ghost`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-ghost %s (revision: %s)", Version, Revision)
	},
}

func validateEnvironment() error {
	err := git.ValidateGit()
	if err != nil {
		return errors.New("git is required")
	}
	return nil
}

func validateGlobalFlags() error {
	if ghostRepo == "" {
		return errors.New("ghost-repo must be specified")
	}
	if baseCommit != "" {
		err := git.ValidateCommitish(baseCommit)
		if err != nil {
			return errors.New("base-commit is not a valid object")
		}
	}
	return nil
}
