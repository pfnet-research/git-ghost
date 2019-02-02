package ghost

import (
	"git-ghost/pkg/util"
	"io"

	log "github.com/Sirupsen/logrus"
)

type ShowOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	PullableLocalModBranchSpec
	// if you want to consume and transform the output of `ghost.Show()`,
	// Please use `io.Pipe()` as below,
	// ```
	// r, w := io.Pipe()
	// go func() { ghost.Show(ShowOptions{ Writer: w }); w.Close()}
	// ````
	// Then, you can read the output from `r` and transform them as you like.
	Writer io.Writer
}

type ShowCommitsOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	Writer io.Writer
}
type ShowDiffOptions struct {
	WorkingEnvSpec
	PullableLocalModBranchSpec
	Writer io.Writer
}

func ShowCommits(options ShowCommitsOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull commits command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()

	branch, err := options.LocalBaseBranchSpec.PullBranch(*we)
	if err != nil {
		return err
	}
	if branch != nil {
		return branch.Show(*we, options.Writer)
	}
	return nil
}

func ShowDiff(options ShowDiffOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull diff command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()

	branch, err := options.PullableLocalModBranchSpec.PullBranch(*we)
	if err != nil {
		return err
	}
	return branch.Show(*we, options.Writer)
}

func ShowAll(options ShowOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	weForCommits, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer weForCommits.clean()
	localBaseBranch, err := options.LocalBaseBranchSpec.PullBranch(*weForCommits)
	if err != nil {
		return err
	}

	weForDiff, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer weForDiff.clean()
	localModBranch, err := options.PullableLocalModBranchSpec.PullBranch(*weForDiff)
	if err != nil {
		return err
	}

	if localBaseBranch != nil {
		err = localBaseBranch.Show(*weForCommits, options.Writer)
		if err != nil {
			return err
		}
	}
	err = localModBranch.Show(*weForDiff, options.Writer)
	if err != nil {
		return err
	}

	return nil
}
