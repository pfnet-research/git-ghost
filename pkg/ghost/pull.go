package ghost

import (
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

// PullOptions represents arg for Pull func
type PullOptions struct {
	WorkingEnvSpec
	*LocalBaseBranchSpec
	*PullableLocalModBranchSpec
}

func pullAndApply(spec PullableGhostBranchSpec, we WorkingEnv) error {
	pulledBranch, err := spec.PullBranch(we)
	if err != nil {
		return err
	}
	return pulledBranch.Apply(we)
}

// Pull pulls ghost branches and apply to workind directory
func Pull(options PullOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()

	if options.LocalBaseBranchSpec != nil {
		err := pullAndApply(*options.LocalBaseBranchSpec, *we)
		if err != nil {
			return err
		}
	}

	if options.PullableLocalModBranchSpec != nil {
		return pullAndApply(*options.PullableLocalModBranchSpec, *we)
	}

	log.WithFields(util.ToFields(options)).Warn("pull command has nothing to do with")
	return nil
}
