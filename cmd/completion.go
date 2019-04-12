// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	# TODO: Support second and third argument completion
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
		git-ghost_list_diff | git-ghost_list_commits | git-ghost_list_all | \
		git-ghost_delete_diff | git-ghost_delete_commits | git-ghost_delete_all )
			# TODO: Support --from and --to completion
			return
			;;
		*)
			;;
	esac
}
	`
)

func NewCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
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
			rootCmd.BashCompletionFunction = bashCompletionFunc
			availableCompletions := map[string]func(io.Writer) error{
				"bash": rootCmd.GenBashCompletion,
				"zsh":  rootCmd.GenZshCompletion,
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
}
