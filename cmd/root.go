package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version  string
	Revision string
)

var RootCmd = &cobra.Command{
	Use:   "git-ghost",
	Short: "git-ghost",
}

// Global Flags
var ghostRemote string
var baseCommit string

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringVar(&ghostRemote, "ghost-remote", "", "git refspec for ghost commits repository")
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
