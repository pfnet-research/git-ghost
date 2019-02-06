package types

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// GhostBranchSpec is an interface
//
// GhostBranchSpec is a specification for creating ghost branch
type GhostBranchSpec interface {
	// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
	CreateBranch(we WorkingEnv) (GhostBranch, error)
}

// PullableGhostBranchSpec is an interface
//
// PullableGhostBranchSpec is a specification for pulling ghost branch from ghost repo
type PullableGhostBranchSpec interface {
	// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
	PullBranch(we WorkingEnv) (GhostBranch, error)
}

// ensuring interfaces
var _ GhostBranchSpec = LocalBaseBranchSpec{}
var _ GhostBranchSpec = LocalModBranchSpec{}
var _ PullableGhostBranchSpec = LocalBaseBranchSpec{}
var _ PullableGhostBranchSpec = PullableLocalModBranchSpec{}

// LocalBaseBranchSpec is a spec for creating local base branch
type LocalBaseBranchSpec struct {
	Prefix              string
	RemoteBaseCommitish string
	LocalBaseCommitish  string
}

// LocalModBranchSpec is a spec for creating local mod branch
type LocalModBranchSpec struct {
	Prefix             string
	LocalBaseCommitish string
}

// PullableLocalModBranchSpec is a spec for pulling local base branch
type PullableLocalModBranchSpec struct {
	LocalModBranchSpec
	LocalModHash string
}

// Resolve resolves comittish in LocalModBranchSpec as full commit hash values
func (bs LocalBaseBranchSpec) Resolve(srcDir string) (*LocalBaseBranchSpec, error) {
	err := git.ValidateComittish(srcDir, bs.RemoteBaseCommitish)
	if err != nil {
		return nil, err
	}
	remoteBaseCommit := resolveComittishOr(srcDir, bs.RemoteBaseCommitish)
	err = git.ValidateComittish(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		return nil, err
	}
	localBaseCommit := resolveComittishOr(srcDir, bs.LocalBaseCommitish)
	branch := &LocalBaseBranchSpec{
		Prefix:              bs.Prefix,
		RemoteBaseCommitish: remoteBaseCommit,
		LocalBaseCommitish:  localBaseCommit,
	}
	return branch, nil
}

// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
func (bs LocalBaseBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, error) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}

	branch := &LocalBaseBranch{
		Prefix:           resolved.Prefix,
		RemoteBaseCommit: resolved.RemoteBaseCommitish,
		LocalBaseCommit:  resolved.LocalBaseCommitish,
	}
	if branch.RemoteBaseCommit == branch.LocalBaseCommit {
		log.WithFields(log.Fields{
			"from": branch.RemoteBaseCommit,
			"to":   branch.LocalBaseCommit,
		}).Warn("skipping pull and apply ghost commits branch because from-hash and to-hash is the same.")
		return nil, nil
	}

	err = pull(branch, we)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
func (bs LocalBaseBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.Resolve(we.SrcDir)
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

// Resolve resolves comittish in LocalModBranchSpec as full commit hash values
func (bs LocalModBranchSpec) Resolve(srcDir string) (*LocalModBranchSpec, error) {
	err := git.ValidateComittish(srcDir, bs.LocalBaseCommitish)
	if err != nil {
		return nil, err
	}
	localBaseCommit := resolveComittishOr(srcDir, bs.LocalBaseCommitish)
	return &LocalModBranchSpec{
		Prefix:             bs.Prefix,
		LocalBaseCommitish: localBaseCommit,
	}, nil
}

// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
func (bs LocalModBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}
	localBaseCommit := resolved.LocalBaseCommitish
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

// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
func (bs PullableLocalModBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, error) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}
	branch := &LocalModBranch{
		Prefix:          resolved.Prefix,
		LocalBaseCommit: resolved.LocalBaseCommitish,
		LocalModHash:    bs.LocalModHash,
	}
	err = pull(branch, we)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func pull(ghost GhostBranch, we WorkingEnv) error {
	return git.ResetHardToBranch(we.GhostDir, git.ORIGIN+"/"+ghost.BranchName())
}

func resolveComittishOr(srcDir string, commitishToResolve string) string {
	resolved, err := git.ResolveComittish(srcDir, commitishToResolve)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": srcDir,
			"specified":  commitishToResolve,
		}).Warn("can't resolve commit-ish value on local git repository.  specified commit-ish value will be used.")
		return commitishToResolve
	}
	return resolved
}
