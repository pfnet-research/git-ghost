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
var _ GhostBranchSpec = CommitsBranchSpec{}
var _ GhostBranchSpec = DiffBranchSpec{}
var _ PullableGhostBranchSpec = CommitsBranchSpec{}
var _ PullableGhostBranchSpec = PullableDiffBranchSpec{}

// CommitsBranchSpec is a spec for creating local base branch
type CommitsBranchSpec struct {
	Prefix        string
	CommitishFrom string
	CommitishTo   string
}

// DiffBranchSpec is a spec for creating local mod branch
type DiffBranchSpec struct {
	Prefix        string
	ComittishFrom string
}

// PullableDiffBranchSpec is a spec for pulling local base branch
type PullableDiffBranchSpec struct {
	DiffBranchSpec
	DiffHash string
}

// Resolve resolves comittish in DiffBranchSpec as full commit hash values
func (bs CommitsBranchSpec) Resolve(srcDir string) (*CommitsBranchSpec, error) {
	err := git.ValidateComittish(srcDir, bs.CommitishFrom)
	if err != nil {
		return nil, err
	}
	commitHashFrom := resolveComittishOr(srcDir, bs.CommitishFrom)
	err = git.ValidateComittish(srcDir, bs.CommitishTo)
	if err != nil {
		return nil, err
	}
	commitHashTo := resolveComittishOr(srcDir, bs.CommitishTo)
	branch := &CommitsBranchSpec{
		Prefix:        bs.Prefix,
		CommitishFrom: commitHashFrom,
		CommitishTo:   commitHashTo,
	}
	return branch, nil
}

// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
func (bs CommitsBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, error) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}

	branch := &CommitsBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: resolved.CommitishFrom,
		CommitHashTo:   resolved.CommitishTo,
	}
	if branch.CommitHashFrom == branch.CommitHashTo {
		log.WithFields(log.Fields{
			"from": branch.CommitHashFrom,
			"to":   branch.CommitHashTo,
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
func (bs CommitsBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}

	commitHashFrom := resolved.CommitishFrom
	commitHashTo := resolved.CommitishTo
	if commitHashFrom == commitHashTo {
		return nil, nil
	}

	branch := CommitsBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: commitHashFrom,
		CommitHashTo:   commitHashTo,
	}
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-base")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	defer util.LogError(func() error { return os.Remove(tmpFile.Name()) })
	err = git.CreateDiffBundleFile(srcDir, tmpFile.Name(), commitHashFrom, commitHashTo)
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

// Resolve resolves comittish in DiffBranchSpec as full commit hash values
func (bs DiffBranchSpec) Resolve(srcDir string) (*DiffBranchSpec, error) {
	err := git.ValidateComittish(srcDir, bs.ComittishFrom)
	if err != nil {
		return nil, err
	}
	commitHashFrom := resolveComittishOr(srcDir, bs.ComittishFrom)
	return &DiffBranchSpec{
		Prefix:        bs.Prefix,
		ComittishFrom: commitHashFrom,
	}, nil
}

// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
func (bs DiffBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, error) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}
	commitHashFrom := resolved.ComittishFrom
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	defer util.LogError(func() error { return os.Remove(tmpFile.Name()) })
	err = git.CreateDiffPatchFile(srcDir, tmpFile.Name(), commitHashFrom)
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
	branch := DiffBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: commitHashFrom,
		DiffHash:       hash,
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
func (bs PullableDiffBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, error) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}
	branch := &DiffBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: resolved.ComittishFrom,
		DiffHash:       bs.DiffHash,
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
