package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GhostBranchSpec interface {
	CreateBranch(we WorkingEnv) (GhostBranch, error)
}

type LocalBaseBranchSpec struct {
	Prefix            string
	RemoteBaseRefspec string
	LocalBaseRefspec  string
}

type LocalModBranchSpec struct {
	Prefix           string
	LocalBaseRefspec string
}

func (bs LocalBaseBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	err := git.InitializeGitDir(dstDir, we.GhostRepo, "")

	if err != nil {
		return nil, err
	}
	err = git.ValidateRefspec(srcDir, bs.RemoteBaseRefspec)
	if err != nil {
		return nil, err
	}
	remoteBaseCommit, err := git.ResolveRefspec(srcDir, bs.RemoteBaseRefspec)
	if err != nil {
		return nil, err
	}
	err = git.ValidateRefspec(srcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(srcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}

	if localBaseCommit == remoteBaseCommit {
		return nil, nil
	}

	branch := LocalBaseBranch{
		Prefix:           bs.Prefix,
		LocalBaseCommit:  localBaseCommit,
		RemoteBaseCommit: remoteBaseCommit,
	}
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
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

func (bs LocalModBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	err := git.InitializeGitDir(dstDir, we.GhostRepo, "")
	if err != nil {
		return nil, err
	}
	err = git.ValidateRefspec(srcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(srcDir, bs.LocalBaseRefspec)
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
		Prefix:          bs.Prefix,
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
