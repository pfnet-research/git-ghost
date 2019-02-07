package ghost

import (
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

// PullOptions represents arg for Pull func
type PullOptions struct {
	types.WorkingEnvSpec
	*types.CommitsBranchSpec
	*types.PullableDiffBranchSpec
}

func pullAndApply(spec types.PullableGhostBranchSpec, we types.WorkingEnv) error {
	pulledBranch, err := spec.PullBranch(we)
	if err != nil {
		return err
	}
	return pulledBranch.Apply(we)
}

// Pull pulls ghost branches and apply to workind directory
func Pull(options PullOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.Initialize()
	if err != nil {
		return err
	}
	defer we.Clean()

	if options.CommitsBranchSpec != nil {
		err := pullAndApply(*options.CommitsBranchSpec, *we)
		if err != nil {
			return err
		}
	}

	if options.PullableDiffBranchSpec != nil {
		return pullAndApply(*options.PullableDiffBranchSpec, *we)
	}

	log.WithFields(util.ToFields(options)).Warn("pull command has nothing to do with")
	return nil
}
