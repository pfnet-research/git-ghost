package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show ghost commits on remote repository.",
	Long:  "show ghost commits on remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull command")
	},
}
