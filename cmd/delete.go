package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete ghost commits from remote repository.",
	Long:  "delete ghost commits from remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete command")
	},
}
