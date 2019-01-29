package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"path"

	log "github.com/Sirupsen/logrus"
)

type PullOptions struct {
	WorkingEnvSpec
	GhostSpec
	ForceApply bool
}

func Pull(options PullOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	// pull command assumed pwd is the src directory to apply ghost commits.
	srcDir := options.SrcDir
	workingEnv, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	ghostDir := workingEnv.GhostDir
	defer workingEnv.clean()

	localBaseBranch, localModBranch, err := options.GhostSpec.validateAndCreateGhostBranches(options.WorkingEnvSpec)
	if err != nil {
		return err
	}

	applyGhostBranchToSrc := func(ghost GhostBranch) error {
		log.WithFields(util.MergeFields(util.ToFields(ghost), log.Fields{"ghostDir": ghostDir, "srcDir": srcDir})).Info("applying ghost branch")

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
			message := "HEAD is not equal to remote-base"
			if options.ForceApply {
				log.WithFields(
					util.MergeFields(util.ToFields(*localBaseBranch),
						log.Fields{"HEAD": srcHead, "srcDir": srcDir}),
				).Warnf("%s. Applying local base branch will be failed.", message)
			} else {
				return fmt.Errorf("abort because %s (HEAD=%s, remote-base=%s)", message, srcHead, localBaseBranch.RemoteBaseCommit)
			}
		}
		err = applyGhostBranchToSrc(localBaseBranch)
		if err != nil {
			return err
		}
	} else {
		log.WithFields(log.Fields{
			"RemoteBaseCommit": localModBranch.LocalBaseCommit,
			"LocalBaseCommit":  localModBranch.LocalBaseCommit,
		}).Info("skipping to apply local base branch because local-base equals to remote-base")

		srcHead, err := git.ResolveRefspec(srcDir, "HEAD")
		if err != nil {
			return err
		}

		if srcHead != localModBranch.LocalBaseCommit {
			message := "HEAD is not equal to local-base"
			if options.ForceApply {
				log.WithFields(util.MergeFields(
					util.ToFields(*localModBranch),
					log.Fields{"HEAD": srcHead, "srcDir": srcDir}),
				).Warnf("%s. Applying local mod branch will be failed.", message)
			} else {
				return fmt.Errorf("abort because %s (HEAD=%s, local-base=%s)", message, srcHead, localModBranch.LocalBaseCommit)
			}
		}
	}
	return applyGhostBranchToSrc(localModBranch)
}
