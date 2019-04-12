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
	"os"

	"github.com/pfnet-research/git-ghost/pkg/ghost/git"
	"github.com/pfnet-research/git-ghost/pkg/ghost/types"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	srcDir       string
	ghostWorkDir string
	ghostPrefix  string
	ghostRepo    string
	verbose      int
}

func (gf globalFlags) WorkingEnvSpec() types.WorkingEnvSpec {
	workingEnvSpec := types.WorkingEnvSpec{
		SrcDir:          gf.srcDir,
		GhostWorkingDir: gf.ghostWorkDir,
		GhostRepo:       gf.ghostRepo,
	}
	userName, userEmail, err := git.GetUserConfig(globalOpts.srcDir)
	if err == nil {
		workingEnvSpec.GhostUserName = userName
		workingEnvSpec.GhostUserEmail = userEmail
	} else {
		log.Debug("failed to get user name and email of the source directory")
	}
	return workingEnvSpec
}

var (
	Version  string
	Revision string
)

var globalOpts globalFlags

func NewRootCmd() *cobra.Command {
	cobra.OnInitialize()
	rootCmd := &cobra.Command{
		Use:           "git-ghost",
		Short:         "git-ghost",
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}
			err := validateEnvironment()
			if err != nil {
				return err
			}
			err = globalOpts.SetDefaults()
			if err != nil {
				return err
			}
			err = globalOpts.Validate()
			if err != nil {
				return err
			}
			switch globalOpts.verbose {
			case 0:
				log.SetLevel(log.ErrorLevel)
			case 1:
				log.SetLevel(log.InfoLevel)
			case 2:
				log.SetLevel(log.DebugLevel)
			case 3:
				log.SetLevel(log.TraceLevel)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&globalOpts.srcDir, "src-dir", "", "source directory which you create ghost from (default to PWD env)")
	rootCmd.PersistentFlags().StringVar(&globalOpts.ghostWorkDir, "ghost-working-dir", "", "local root directory for git-ghost interacting with ghost repository (default to a temporary directory)")
	rootCmd.PersistentFlags().StringVar(&globalOpts.ghostPrefix, "ghost-prefix", "", "prefix of ghost branch name (default to GIT_GHOST_PREFIX env, or ghost)")
	rootCmd.PersistentFlags().StringVar(&globalOpts.ghostRepo, "ghost-repo", "", "git remote url for ghosts repository (default to GIT_GHOST_REPO env)")
	rootCmd.PersistentFlags().CountVarP(&globalOpts.verbose, "verbose", "v", "verbose mode")
	rootCmd.AddCommand(NewPushCommand())
	rootCmd.AddCommand(NewPullCommand())
	rootCmd.AddCommand(NewShowCommand())
	rootCmd.AddCommand(NewListCommand())
	rootCmd.AddCommand(NewDeleteCommand())
	rootCmd.AddCommand(NewGCCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewCompletionCommand(rootCmd))
	return rootCmd
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of git-ghost",
		Long:  `Print the version number of git-ghost`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("git-ghost %s (revision: %s)", Version, Revision)
		},
	}
}

func validateEnvironment() errors.GitGhostError {
	err := git.ValidateGit()
	if err != nil {
		return errors.New("git is required")
	}
	return nil
}

func (flags *globalFlags) SetDefaults() errors.GitGhostError {
	if globalOpts.srcDir == "" {
		globalOpts.srcDir = os.Getenv("PWD")
	}
	if globalOpts.ghostWorkDir == "" {
		globalOpts.ghostWorkDir = os.TempDir()
	}
	if globalOpts.ghostPrefix == "" {
		ghostPrefixEnv := os.Getenv("GIT_GHOST_PREFIX")
		if ghostPrefixEnv == "" {
			ghostPrefixEnv = "ghost"
		}
		globalOpts.ghostPrefix = ghostPrefixEnv
	}
	if globalOpts.ghostRepo == "" {
		globalOpts.ghostRepo = os.Getenv("GIT_GHOST_REPO")
	}
	return nil
}

func (flags *globalFlags) Validate() errors.GitGhostError {
	if flags.srcDir == "" {
		return errors.New("src-dir must be specified")
	}
	_, err := os.Stat(flags.ghostWorkDir)
	if err != nil {
		return errors.Errorf("ghost-working-dir is not found (value: %v)", flags.ghostWorkDir)
	}
	if flags.ghostPrefix == "" {
		return errors.New("ghost-prefix must be specified")
	}
	if flags.ghostRepo == "" {
		return errors.New("ghost-repo must be specified")
	}
	return nil
}
