package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull [hash]",
	Short: "pull a ghost commit from remote repository and apply to your working git repository.",
	Long:  "pull a ghost commit from remote repository and apply to your working git repository.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull command")
	},
}
