package cmd

import (
	"errors"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/ghost/types"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type globalFlags struct {
	srcDir       string
	ghostWorkDir string
	ghostPrefix  string
	ghostRepo    string
	verbose      bool
}

func (gf globalFlags) WorkingEnvSpec() types.WorkingEnvSpec {
	return types.WorkingEnvSpec{
		SrcDir:          gf.srcDir,
		GhostWorkingDir: gf.ghostWorkDir,
		GhostRepo:       gf.ghostRepo,
	}
}

var (
	Version  string
	Revision string
)

var RootCmd = &cobra.Command{
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
		err = globalOpts.Validate()
		if err != nil {
			return err
		}
		if globalOpts.verbose {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
}

var globalOpts globalFlags

func init() {
	cobra.OnInitialize()
	currentDir := os.Getenv("PWD")
	RootCmd.PersistentFlags().StringVar(&globalOpts.srcDir, "src-dir", currentDir, "source directory which you create ghost from")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostWorkDir, "ghost-working-dir", os.TempDir(), "local root directory for git-ghost interacting with ghost repository")
	ghostPrefixEnv := os.Getenv("GHOST_PREFIX")
	if ghostPrefixEnv == "" {
		ghostPrefixEnv = "ghost"
	}
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostPrefix, "ghost-prefix", ghostPrefixEnv, "prefix of ghost branch name")
	ghostRepoEnv := os.Getenv("GHOST_REPO")
	RootCmd.PersistentFlags().StringVar(&globalOpts.ghostRepo, "ghost-repo", ghostRepoEnv, "git remote url for ghosts repository")
	RootCmd.PersistentFlags().BoolVarP(&globalOpts.verbose, "verbose", "v", false, "verbose mode")
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

func (flags *globalFlags) Validate() error {
	if flags.srcDir == "" {
		return errors.New("src-dir must be specified")
	}
	_, err := os.Stat(flags.ghostWorkDir)
	if err != nil {
		return fmt.Errorf("ghost-working-dir is not found (value: %v)", flags.ghostWorkDir)
	}
	if flags.ghostPrefix == "" {
		return errors.New("ghost-prefix must be specified")
	}
	if flags.ghostRepo == "" {
		return errors.New("ghost-repo must be specified")
	}
	return nil
}
