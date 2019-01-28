package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Commit struct {
	BaseCommitHash string
	Commits        []string
	Diff           string
}

type GhostBranchSpec interface {
	CreateBranch() (GhostBranch, error)
}

type LocalBaseBranchSpec struct {
	SrcDir            string
	DstDir            string
	Prefix            string
	RemoteBaseRefspec string
	LocalBaseRefspec  string
}

type LocalModBranchSpec struct {
	SrcDir           string
	DstDir           string
	Prefix           string
	LocalBaseRefspec string
}

func (bs LocalBaseBranchSpec) CreateBranch() (GhostBranch, error) {
	err := git.ValidateRefspec(bs.SrcDir, bs.RemoteBaseRefspec)
	if err != nil {
		return nil, err
	}
	remoteBaseCommit, err := git.ResolveRefspec(bs.SrcDir, bs.RemoteBaseRefspec)
	if err != nil {
		return nil, err
	}
	err = git.ValidateRefspec(bs.SrcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(bs.SrcDir, bs.LocalBaseRefspec)
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
	err = git.CreateDiffBundleFile(bs.SrcDir, tmpFile.Name(), remoteBaseCommit, localBaseCommit)
	if err != nil {
		return nil, err
	}
	err = os.Rename(tmpFile.Name(), filepath.Join(bs.DstDir, branch.FileName()))
	if err != nil {
		return nil, err
	}
	err = git.CreateOrphanBranch(bs.DstDir, branch.BranchName())
	if err != nil {
		return nil, err
	}
	err = git.CommitFile(bs.DstDir, branch.FileName(), fmt.Sprintf("%s..%s", remoteBaseCommit, localBaseCommit))
	if err != nil {
		return nil, err
	}
	return &branch, err
}

func (bs LocalModBranchSpec) CreateBranch() (GhostBranch, error) {
	err := git.ValidateRefspec(bs.SrcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}
	localBaseCommit, err := git.ResolveRefspec(bs.SrcDir, bs.LocalBaseRefspec)
	if err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	err = git.CreateDiffPatchFile(bs.SrcDir, tmpFile.Name(), localBaseCommit)
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
	err = os.Rename(tmpFile.Name(), filepath.Join(bs.DstDir, branch.FileName()))
	if err != nil {
		return nil, err
	}
	err = git.CreateOrphanBranch(bs.DstDir, branch.BranchName())
	if err != nil {
		return nil, err
	}
	err = git.CommitFile(bs.DstDir, branch.FileName(), fmt.Sprintf("%s/%s", localBaseCommit, hash))
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

type GhostBranch interface {
	BranchName() string
	FileName() string
}

type LocalBaseBranch struct {
	Prefix           string
	RemoteBaseCommit string
	LocalBaseCommit  string
}

type LocalModBranch struct {
	Prefix          string
	LocalBaseCommit string
	LocalModHash    string
}

func (b LocalBaseBranch) BranchName() string {
	return fmt.Sprintf("%s/%s-%s", b.Prefix, b.RemoteBaseCommit, b.LocalBaseCommit)
}

func (b LocalBaseBranch) FileName() string {
	return "commits.patch"
}

func (b LocalModBranch) BranchName() string {
	return fmt.Sprintf("%s/%s/%s", b.Prefix, b.LocalBaseCommit, b.LocalModHash)
}

func (b LocalModBranch) FileName() string {
	return "local-mod.patch"
}

type PushOptions struct {
	SrcDir      string
	DstDir      string
	GhostPrefix string
	GhostRepo   string
	RemoteBase  string
	LocalBase   string
}

type PushResult struct {
	LocalBaseBranch *LocalBaseBranch
	LocalModBranch  *LocalModBranch
}
