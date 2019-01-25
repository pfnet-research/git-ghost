package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(gcCmd)
}

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "gc ghost commits from remote repository.",
	Long:  "gc ghost commits from remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gc command")
	},
}
