package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	bashCompletionFunc = `
# Need _git_ghost function to autocomplete on 'git ghost' instead of 'git-ghost'
_git_ghost ()
{
	__start_git-ghost
}
__git-ghost_get_hash() {
	local ghost_out
	if ghost_out=$(git-ghost list -o only-from --no-headers --from "$1*" | uniq 2>/dev/null); then
	    __git-ghost_debug "${FUNCNAME[0]}: ${ghost_out} -- $cur"
	    COMPREPLY+=( $( compgen -W "${ghost_out[*]}" -- "$cur" ) )
	fi
}
__git-ghost_custom_func() {
	case ${last_command} in
		git-ghost_push_diff | git-ghost_push_commits | git-ghost_push_all | \
		git-ghost_pull_diff | git-ghost_pull_commits | git-ghost_pull_all | \
		git-ghost_show_diff | git-ghost_show_commits | git-ghost_show_all )
			__git-ghost_get_hash
			return
			;;
		*)
			;;
	esac
}
	`
)

func init() {
	RootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion SHELL",
	Short: "output shell completion code for the specified shell (bash or zsh)",
	Long: `Write bash or zsh shell completion code to standard output.

	For bash, ensure you have bash completions installed and enabled.
	To access completions in your current shell, run
	$ source <(git-ghost completion bash)
	Alternatively, write it to a file and source in .bash_profile

	For zsh, output to a file in a directory referenced by the $fpath shell
	variable.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		RootCmd.BashCompletionFunction = bashCompletionFunc
		availableCompletions := map[string]func(io.Writer) error{
			"bash": RootCmd.GenBashCompletion,
			"zsh":  RootCmd.GenZshCompletion,
		}
		completion, ok := availableCompletions[shell]
		if !ok {
			fmt.Printf("Invalid shell '%s'. The supported shells are bash and zsh.\n", shell)
			os.Exit(1)
		}
		if err := completion(os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}
