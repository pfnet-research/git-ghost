package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list ghost commits on remote repository.",
	Long:  "list ghost commits on remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list command")
	},
}
