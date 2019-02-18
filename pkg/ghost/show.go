package ghost

import (
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"
	"io"

	log "github.com/Sirupsen/logrus"
)

// ShowOptions represents arg for Pull func
type ShowOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.PullableDiffBranchSpec
	// if you want to consume and transform the output of `ghost.Show()`,
	// Please use `io.Pipe()` as below,
	// ```
	// r, w := io.Pipe()
	// go func() { ghost.Show(ShowOptions{ Writer: w }); w.Close()}
	// ````
	// Then, you can read the output from `r` and transform them as you like.
	Writer io.Writer
}

func pullAndshow(branchSpec types.PullableGhostBranchSpec, we types.WorkingEnv, writer io.Writer) error {
	branch, err := branchSpec.PullBranch(we)
	if err != nil {
		return err
	}
	if branch != nil {
		return branch.Show(we, writer)
	}
	return nil
}

// Show writes ghost branches contents to option.Writer
func Show(options ShowOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	if options.CommitsBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer util.LogDeferredError(we.Clean)
		err = pullAndshow(options.CommitsBranchSpec, *we, options.Writer)
		if err != nil {
			return err
		}
	}

	if options.PullableDiffBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer util.LogDeferredError(we.Clean)
		return pullAndshow(options.PullableDiffBranchSpec, *we, options.Writer)
	}

	log.WithFields(util.ToFields(options)).Warn("show command has nothing to do with")
	return nil
}
