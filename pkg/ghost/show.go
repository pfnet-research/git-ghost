package ghost

import (
	"git-ghost/pkg/util"
	"io"

	log "github.com/Sirupsen/logrus"
)

type ShowOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	LocalModBranchSpec
	LocalModHash string
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
	LocalModBranchSpec
	LocalModHash string
	Writer       io.Writer
}

func prepareShowCommits(options ShowCommitsOptions, we WorkingEnv) (*LocalBaseBranch, error) {
	resolved, err := options.LocalBaseBranchSpec.resolve(we)
	if err != nil {
		return nil, err
	}
	branch := LocalBaseBranch{
		Prefix:           resolved.Prefix,
		RemoteBaseCommit: resolved.RemoteBaseCommitish,
		LocalBaseCommit:  resolved.LocalBaseCommitish,
	}
	if err := branch.Prepare(we); err != nil {
		return nil, err
	}
	return &branch, nil
}

func ShowCommits(options ShowCommitsOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull commits command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()
	branch, err := prepareShowCommits(options, *we)
	if err != nil {
		return err
	}
	return branch.Show(*we, options.Writer)
}

func prepareShowDiff(options ShowDiffOptions, we WorkingEnv) (*LocalModBranch, error) {
	resolved, err := options.LocalModBranchSpec.resolve(we)
	if err != nil {
		return nil, err
	}
	branch := LocalModBranch{
		Prefix:          resolved.Prefix,
		LocalBaseCommit: resolved.LocalBaseCommitish,
		LocalModHash:    options.LocalModHash,
	}
	err = branch.Prepare(we)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func ShowDiff(options ShowDiffOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull diff command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()
	branch, err := prepareShowDiff(options, *we)
	if err != nil {
		return err
	}
	return branch.Show(*we, options.Writer)
}

func ShowAll(options ShowOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	showCommitsOpt := ShowCommitsOptions{
		WorkingEnvSpec:      options.WorkingEnvSpec,
		LocalBaseBranchSpec: options.LocalBaseBranchSpec,
		Writer:              options.Writer,
	}
	showDiffOpt := ShowDiffOptions{
		WorkingEnvSpec:     options.WorkingEnvSpec,
		LocalModBranchSpec: options.LocalModBranchSpec,
		LocalModHash:       options.LocalModHash,
		Writer:             options.Writer,
	}

	weForCommits, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer weForCommits.clean()
	localBaseBranch, err := prepareShowCommits(showCommitsOpt, *weForCommits)
	if err != nil {
		return err
	}

	weForDiff, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer weForDiff.clean()
	localModBranch, err := prepareShowDiff(showDiffOpt, *weForDiff)
	if err != nil {
		return err
	}

	err = localBaseBranch.Show(*weForCommits, options.Writer)
	if err != nil {
		return err
	}
	err = localModBranch.Show(*weForDiff, options.Writer)
	if err != nil {
		return err
	}

	return nil
}
