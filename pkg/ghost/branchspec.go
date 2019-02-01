package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type GhostBranchSpec interface {
	CreateBranch(we WorkingEnv) (GhostBranch, error)
}

type LocalBaseBranchSpec struct {
	Prefix              string
	RemoteBaseCommitish string
	LocalBaseCommitish  string
}

type LocalModBranchSpec struct {
	Prefix             string
	LocalBaseCommitish string
}

func (bs LocalBaseBranchSpec) resolve(we WorkingEnv) (*LocalBaseBranchSpec, error) {
	srcDir := we.SrcDir
	err := git.ValidateRefspec(srcDir, bs.RemoteBaseCommitish)
	if err != nil {
		return nil, err
	}

	remoteBaseCommit, err := git.ResolveRefspec(srcDir, bs.RemoteBaseCommitish)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": srcDir,
			"specified":  bs.RemoteBaseCommitish,
		}).Warn("can't resolve commit-ish value on local git repository.  specified commit-ish value will be used.")
		remoteBaseCommit = bs.RemoteBaseCommitish
	}

	err = git.ValidateRefspec(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": srcDir,
			"specified":  bs.LocalBaseCommitish,
		}).Warn("can't resolve commit-ish value on local git repository.  specified commit-ish value will be used.")
		localBaseCommit = bs.LocalBaseCommitish
	}

	return &LocalBaseBranchSpec{
		Prefix:              bs.Prefix,
		RemoteBaseCommitish: remoteBaseCommit,
		LocalBaseCommitish:  localBaseCommit,
	}, nil
}

func (bs LocalBaseBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.resolve(we)
	if err != nil {
		return nil, err
	}

	remoteBaseCommit := resolved.RemoteBaseCommitish
	localBaseCommit := resolved.LocalBaseCommitish
	if localBaseCommit == remoteBaseCommit {
		return nil, nil
	}

	branch := LocalBaseBranch{
		Prefix:           resolved.Prefix,
		LocalBaseCommit:  localBaseCommit,
		RemoteBaseCommit: remoteBaseCommit,
	}
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-base")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	err = git.CreateDiffBundleFile(srcDir, tmpFile.Name(), remoteBaseCommit, localBaseCommit)
	if err != nil {
		return nil, err
	}
	err = os.Rename(tmpFile.Name(), filepath.Join(dstDir, branch.FileName()))
	if err != nil {
		return nil, err
	}

	err = git.CreateOrphanBranch(dstDir, branch.BranchName())
	if err != nil {
		return nil, err
	}
	err = git.CommitFile(dstDir, branch.FileName(), fmt.Sprintf("Create ghost commit"))
	if err != nil {
		return nil, err
	}

	return &branch, nil
}

func (bs LocalModBranchSpec) resolve(we WorkingEnv) (*LocalModBranchSpec, error) {
	srcDir := we.SrcDir
	err := git.ValidateRefspec(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": srcDir,
			"specified":  bs.LocalBaseCommitish,
		}).Warn("can't resolve commit-ish value on local git repository.  specified commit-ish value will be used.")
		localBaseCommit = bs.LocalBaseCommitish
	}
	return &LocalModBranchSpec{
		Prefix:             bs.Prefix,
		LocalBaseCommitish: localBaseCommit,
	}, nil
}

func (bs LocalModBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.resolve(we)
	localBaseCommit := resolved.LocalBaseCommitish
	if err != nil {
		return nil, err
	}
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	err = git.CreateDiffPatchFile(srcDir, tmpFile.Name(), localBaseCommit)
	if err != nil {
		return nil, err
	}
	size, err := util.FileSize(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	if size == 0 {
		return nil, nil
	}

	hash, err := util.GenerateFileContentHash(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	branch := LocalModBranch{
		Prefix:          resolved.Prefix,
		LocalBaseCommit: localBaseCommit,
		LocalModHash:    hash,
	}
	err = os.Rename(tmpFile.Name(), filepath.Join(dstDir, branch.FileName()))
	if err != nil {
		return nil, err
	}

	err = git.CreateOrphanBranch(dstDir, branch.BranchName())
	if err != nil {
		return nil, err
	}
	err = git.CommitFile(dstDir, branch.FileName(), fmt.Sprintf("Create ghost commit"))
	if err != nil {
		return nil, err
	}

	return &branch, nil
}
