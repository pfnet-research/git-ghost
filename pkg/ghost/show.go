package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type ShowOptions struct {
	WorkingEnvSpec
	GhostSpec
	Writer io.Writer
}

func Show(options ShowOptions) error {
	log.WithFields(util.ToFields(options)).Debug("show command with")

	localBaseBranch, localModBranch, err := options.GhostSpec.validateAndCreateGhostBranches(options.WorkingEnvSpec)
	if err != nil {
		return err
	}

	checkoutGhostBranch := func(gb GhostBranch) (*WorkingEnv, error) {
		workingEnv, err := options.WorkingEnvSpec.initialize()
		if err != nil {
			return nil, err
		}

		err = git.ResetHardToBranch(workingEnv.GhostDir, git.ORIGIN+"/"+gb.BranchName())
		if err != nil {
			return nil, err
		}
		return workingEnv, nil
	}
	genGitCatFileCommand := func(we *WorkingEnv, gb GhostBranch) []string {
		return []string{"git", "-C", we.GhostDir, "cat-file", "-p", fmt.Sprintf("HEAD:%s", gb.FileName())}
	}

	commandStr := []string{}
	if localBaseBranch != nil {
		we, err := checkoutGhostBranch(localBaseBranch)
		if err != nil {
			return err
		}
		defer we.clean()
		commandStr = append(commandStr, genGitCatFileCommand(we, localBaseBranch)...)
		commandStr = append(commandStr, "&&")
	}
	we, err := checkoutGhostBranch(localModBranch)
	if err != nil {
		return err
	}
	defer we.clean()
	commandStr = append(commandStr, genGitCatFileCommand(we, localModBranch)...)
	cmd := exec.Command("/bin/sh", "-c", strings.Join(commandStr, " "))
	cmd.Stdout = options.Writer
	return util.JustRunCmd(cmd)
}
