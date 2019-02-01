package ghost

import (
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type PullOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	LocalModBranchSpec
	LocalModHash string
	ForceApply   bool
}

type PullCommitsOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	ForceApply bool
}
type PullDiffOptions struct {
	WorkingEnvSpec
	LocalModBranchSpec
	LocalModHash string
	ForceApply   bool
}

func PullCommits(options PullCommitsOptions, workingEnv *WorkingEnv) error {
	log.WithFields(util.ToFields(options)).Debug("pull commits command with")

	we, initialized, err := initializeWorkingEnvIfRequired(options.WorkingEnvSpec, workingEnv)
	if err != nil {
		return err
	}
	if initialized {
		defer we.clean()
	}

	resolved, err := options.LocalBaseBranchSpec.resolve(*we)
	if err != nil {
		return err
	}
	branch := LocalBaseBranch{
		Prefix:           resolved.Prefix,
		RemoteBaseCommit: resolved.RemoteBaseCommitish,
		LocalBaseCommit:  resolved.LocalBaseCommitish,
	}

	return branch.Apply(*we, options.ForceApply)
}

func PullDiff(options PullDiffOptions, workingEnv *WorkingEnv) error {
	log.WithFields(util.ToFields(options)).Debug("pull diff command with")
	we, initialized, err := initializeWorkingEnvIfRequired(options.WorkingEnvSpec, workingEnv)
	if err != nil {
		return err
	}
	if initialized {
		defer we.clean()
	}

	resolved, err := options.LocalModBranchSpec.resolve(*we)
	if err != nil {
		return err
	}

	branch := LocalModBranch{
		Prefix:          resolved.Prefix,
		LocalBaseCommit: resolved.LocalBaseCommitish,
		LocalModHash:    options.LocalModHash,
	}

	return branch.Apply(*we, options.ForceApply)
}

func Pull(options PullOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()

	err = PullCommits(PullCommitsOptions{
		WorkingEnvSpec:      options.WorkingEnvSpec,
		LocalBaseBranchSpec: options.LocalBaseBranchSpec,
		ForceApply:          options.ForceApply,
	}, we)
	if err != nil {
		return err
	}

	err = PullDiff(PullDiffOptions{
		WorkingEnvSpec:     options.WorkingEnvSpec,
		LocalModBranchSpec: options.LocalModBranchSpec,
		LocalModHash:       options.LocalModHash,
		ForceApply:         options.ForceApply,
	}, we)
	if err != nil {
		return err
	}

	return nil
}

func initializeWorkingEnvIfRequired(spec WorkingEnvSpec, we *WorkingEnv) (*WorkingEnv, bool, error) {
	if we == nil {
		we, err := spec.initialize()
		if err != nil {
			return nil, false, err
		}
		return we, true, nil
	}
	return we, false, nil
}
