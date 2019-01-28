package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"io/ioutil"
	"os"
	"path"
)

func Pull(options PullOptions) error {
	fmt.Fprintf(os.Stderr, "pull command with %+v\n", options)

	// pull command assumed pwd is the src directory to apply ghost commits.
	srcDir := options.SrcDir
	ghostDir, err := ioutil.TempDir(options.GhostWorkingDir, "git-ghost-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(ghostDir)
	err = git.InitializeGitDir(ghostDir, options.GhostRepo, "")
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "tmpdir: %s\n", ghostDir)

	localBaseBranch, localModBranch, err := validateAndCreateGhostBranches(options)
	if err != nil {
		return err
	}

	applyGhostBranchToSrc := func(ghost GhostBranch) error {
		fmt.Fprintf(os.Stderr, "applying ghost branch %+v on %s to %s\n", ghost, ghostDir, srcDir)

		err = git.ResetHardToBranch(ghostDir, git.ORIGIN+"/"+ghost.BranchName())
		if err != nil {
			return err
		}

		// TODO make this instance methods.
		switch t := ghost.(type) {
		case *LocalBaseBranch:
			err = git.ApplyDiffBundleFile(srcDir, path.Join(ghostDir, ghost.FileName()))
		case *LocalModBranch:
			err = git.ApplyDiffPatchFile(srcDir, path.Join(ghostDir, ghost.FileName()))
		default:
			return fmt.Errorf("not supported on type = %+v", t)
		}

		if err != nil {
			return err
		}
		return nil
	}

	if localBaseBranch != nil {
		srcHead, err := git.ResolveRefspec(srcDir, "HEAD")
		if err != nil {
			return err
		}
		if srcHead != localBaseBranch.RemoteBaseCommit {
			fmt.Fprintf(os.Stderr, "WARNING: HEAD is not equal to remote-base.  applying local base branch will be failed. HEAD=%s, remote-base=%s", srcHead, localBaseBranch.RemoteBaseCommit)
		}
		err = applyGhostBranchToSrc(localBaseBranch)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "Skipping to apply local base branch because --local-base equals to --remote-base (%s)", localModBranch.LocalBaseCommit)
		srcHead, err := git.ResolveRefspec(srcDir, "HEAD")
		if err != nil {
			return err
		}
		if srcHead != localModBranch.LocalBaseCommit {
			fmt.Fprintf(os.Stderr, "WARNING: HEAD is not equal to local-base.  applying local mod branch will be failed. HEAD=%s, local-base=%s", srcHead, localModBranch.LocalBaseCommit)
		}
	}
	return applyGhostBranchToSrc(localModBranch)
}

func validateAndCreateGhostBranches(options PullOptions) (*LocalBaseBranch, *LocalModBranch, error) {
	var err error

	// resolve HEAD If necessary
	remoteBaseResolved := options.RemoteBase
	localBaseResolved := options.LocalBase
	if options.RemoteBase == "HEAD" {
		remoteBaseResolved, err = git.ResolveRefspec(options.SrcDir, options.RemoteBase)
		if err != nil {
			return nil, nil, err
		}
	}
	if options.LocalBase == "HEAD" {
		localBaseResolved, err = git.ResolveRefspec(options.SrcDir, options.LocalBase)
		if err != nil {
			return nil, nil, err
		}
	}

	// ghost branch validations and create ghost branches
	var localBaseBranch *LocalBaseBranch
	if remoteBaseResolved != localBaseResolved {
		// TODO warning when srcDir is on remoteBaseResolved.
		localBaseBranch = &LocalBaseBranch{
			Prefix:           options.GhostPrefix,
			RemoteBaseCommit: remoteBaseResolved,
			LocalBaseCommit:  localBaseResolved,
		}

		existence, err := git.ValidateRemoteBranchExistence(options.GhostRepo, localBaseBranch.BranchName())
		if err != nil {
			return nil, nil, err
		}
		if !existence {
			return nil, nil, fmt.Errorf("can't resolve local base branch on %s: %+v", options.GhostRepo, localBaseBranch)
		}
	}

	localModBranch := &LocalModBranch{
		Prefix:          options.GhostPrefix,
		LocalBaseCommit: localBaseResolved,
		LocalModHash:    options.Hash,
	}
	existence, err := git.ValidateRemoteBranchExistence(options.GhostRepo, localModBranch.BranchName())
	if err != nil {
		return nil, nil, err
	}
	if !existence {
		return nil, nil, fmt.Errorf("can't resolve local mod branch on %s: %+v", options.GhostRepo, localModBranch)
	}

	return localBaseBranch, localModBranch, nil
}
