package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

type ShowOptions struct {
	WorkingEnvSpec
	GhostSpec
	// if you want to consume and transform the output of `ghost.Show()`,
	// Please use `io.Pipe()` as below,
	// ```
	// r, w := io.Pipe()
	// go func() { ghost.Show(ShowOptions{ Writer: w }); w.Close()}
	// ````
	// Then, you can read the output from `r` and transform them as you like.
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
	execShow := func(we *WorkingEnv, gb GhostBranch) error {
		cmd := exec.Command("git", "-C", we.GhostDir, "--no-pager", "cat-file", "-p", fmt.Sprintf("HEAD:%s", gb.FileName()))
		cmd.Stdout = options.Writer
		return util.JustRunCmd(cmd)
	}

	if localBaseBranch != nil {
		we, err := checkoutGhostBranch(localBaseBranch)
		if err != nil {
			return err
		}
		defer we.clean()
		err = execShow(we, localBaseBranch)
		if err != nil {
			return err
		}
	}
	we, err := checkoutGhostBranch(localModBranch)
	if err != nil {
		return err
	}
	defer we.clean()
	return execShow(we, localModBranch)
}
