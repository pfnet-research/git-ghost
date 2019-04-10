// Copyright 2019 Preferred Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
)

// GhostBranchSpec is an interface
//
// GhostBranchSpec is a specification for creating ghost branch
type GhostBranchSpec interface {
	// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
	CreateBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError)
}

// PullableGhostBranchSpec is an interface
//
// PullableGhostBranchSpec is a specification for pulling ghost branch from ghost repo
type PullableGhostBranchSpec interface {
	// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
	PullBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError)
}

// ensuring interfaces
var _ GhostBranchSpec = CommitsBranchSpec{}
var _ GhostBranchSpec = DiffBranchSpec{}
var _ PullableGhostBranchSpec = CommitsBranchSpec{}
var _ PullableGhostBranchSpec = PullableDiffBranchSpec{}

// Constants
const maxSymlinkDepth = 3

// CommitsBranchSpec is a spec for creating local base branch
type CommitsBranchSpec struct {
	Prefix         string
	CommittishFrom string
	CommittishTo   string
}

// DiffBranchSpec is a spec for creating local mod branch
type DiffBranchSpec struct {
	Prefix            string
	CommittishFrom    string
	IncludedFilepaths []string
	FollowSymlinks    bool
}

// PullableDiffBranchSpec is a spec for pulling local base branch
type PullableDiffBranchSpec struct {
	Prefix         string
	CommittishFrom string
	DiffHash       string
}

// Resolve resolves committish in DiffBranchSpec as full commit hash values.
// The special character "_" is allowed to indicate full commmits.
func (bs CommitsBranchSpec) Resolve(srcDir string) (*CommitsBranchSpec, errors.GitGhostError) {
	commitHashFrom := bs.CommittishFrom
	if bs.CommittishFrom != util.CommitStartFromInit {
		// CommittishFrom must be a valid existing committish
		err := git.ValidateCommittish(srcDir, bs.CommittishFrom)
		if err != nil {
			return nil, err
		}
		commitHashFrom = resolveCommittishOr(srcDir, bs.CommittishFrom)
		err = git.ValidateCommittish(srcDir, bs.CommittishTo)
		if err != nil {
			return nil, err
		}
	}
	commitHashTo := resolveCommittishOr(srcDir, bs.CommittishTo)
	branch := &CommitsBranchSpec{
		Prefix:         bs.Prefix,
		CommittishFrom: commitHashFrom,
		CommittishTo:   commitHashTo,
	}
	return branch, nil
}

// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
func (bs CommitsBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}

	branch := &CommitsBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: resolved.CommittishFrom,
		CommitHashTo:   resolved.CommittishTo,
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
func (bs CommitsBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, ggerr := bs.Resolve(we.SrcDir)
	if ggerr != nil {
		return nil, ggerr
	}

	commitHashFrom := resolved.CommittishFrom
	commitHashTo := resolved.CommittishTo
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
		return nil, errors.WithStack(err)
	}
	util.LogDeferredError(tmpFile.Close)
	defer util.LogDeferredError(func() error { return os.Remove(tmpFile.Name()) })
	if commitHashFrom == util.CommitStartFromInit {
		ggerr = git.CreateFullBundleFile(srcDir, tmpFile.Name(), commitHashTo)
		if ggerr != nil {
			return nil, ggerr
		}
	} else {
		ggerr = git.CreateDiffBundleFile(srcDir, tmpFile.Name(), commitHashFrom, commitHashTo)
		if ggerr != nil {
			return nil, ggerr
		}
	}
	err = os.Rename(tmpFile.Name(), filepath.Join(dstDir, branch.FileName()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ggerr = git.CreateOrphanBranch(dstDir, branch.BranchName())
	if ggerr != nil {
		return nil, ggerr
	}
	ggerr = git.CommitFile(dstDir, branch.FileName(), fmt.Sprintf("Create ghost commit"))
	if ggerr != nil {
		return nil, ggerr
	}

	return &branch, nil
}

// Resolve resolves committish in DiffBranchSpec as full commit hash values
func (bs DiffBranchSpec) Resolve(srcDir string) (*DiffBranchSpec, errors.GitGhostError) {
	// CommittishFrom must be a valid existing committish
	err := git.ValidateCommittish(srcDir, bs.CommittishFrom)
	if err != nil {
		return nil, err
	}
	commitHashFrom := resolveCommittishOr(srcDir, bs.CommittishFrom)

	var errs error
	includedFilepaths := make([]string, 0, len(bs.IncludedFilepaths))
	for _, p := range bs.IncludedFilepaths {
		resolved, err := resolveFilepath(srcDir, p)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		includedFilepaths = append(includedFilepaths, resolved)

		if bs.FollowSymlinks {
			islink, err := util.IsSymlink(p)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			if islink {
				err := util.WalkSymlink(srcDir, p, func(paths []string, pp string) errors.GitGhostError {
					if len(paths) > maxSymlinkDepth {
						return errors.Errorf("symlink is too deep (< %d): %s", maxSymlinkDepth, strings.Join(paths, " -> "))
					}
					if filepath.IsAbs(pp) {
						return errors.Errorf("symlink to absolute path is not supported: %s -> %s", strings.Join(paths, " -> "), pp)
					}
					resolved, err := resolveFilepath(srcDir, pp)
					if err != nil {
						return errors.WithStack(err)
					}
					includedFilepaths = append(includedFilepaths, resolved)
					return nil
				})
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
			}
		}
	}
	if errs != nil {
		return nil, errors.WithStack(errs)
	}
	if len(includedFilepaths) > 0 {
		includedFilepaths = util.UniqueStringSlice(includedFilepaths)
	}

	return &DiffBranchSpec{
		Prefix:            bs.Prefix,
		CommittishFrom:    commitHashFrom,
		IncludedFilepaths: includedFilepaths,
	}, nil
}

// CreateBranch create a ghost branch on WorkingEnv and returns a GhostBranch object
func (bs DiffBranchSpec) CreateBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError) {
	dstDir := we.GhostDir
	srcDir := we.SrcDir
	resolved, ggerr := bs.Resolve(we.SrcDir)
	if ggerr != nil {
		return nil, ggerr
	}
	commitHashFrom := resolved.CommittishFrom
	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	util.LogDeferredError(tmpFile.Close)
	defer util.LogDeferredError(func() error { return os.Remove(tmpFile.Name()) })
	err = git.CreateDiffPatchFile(srcDir, tmpFile.Name(), commitHashFrom)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(bs.IncludedFilepaths) > 0 {
		err = git.AppendNonIndexedDiffFiles(srcDir, tmpFile.Name(), resolved.IncludedFilepaths)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	size, err := util.FileSize(tmpFile.Name())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if size == 0 {
		return nil, nil
	}

	hash, err := util.GenerateFileContentHash(tmpFile.Name())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	branch := DiffBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: commitHashFrom,
		DiffHash:       hash,
	}
	err = os.Rename(tmpFile.Name(), filepath.Join(dstDir, branch.FileName()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = git.CreateOrphanBranch(dstDir, branch.BranchName())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = git.CommitFile(dstDir, branch.FileName(), fmt.Sprintf("Create ghost commit"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &branch, nil
}

// Resolve resolves committish in PullableDiffBranchSpec as full commit hash values
func (bs PullableDiffBranchSpec) Resolve(srcDir string) (*PullableDiffBranchSpec, errors.GitGhostError) {
	err := git.ValidateCommittish(srcDir, bs.CommittishFrom)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	commitHashFrom := resolveCommittishOr(srcDir, bs.CommittishFrom)

	return &PullableDiffBranchSpec{
		Prefix:         bs.Prefix,
		CommittishFrom: commitHashFrom,
		DiffHash:       bs.DiffHash,
	}, nil
}

// PullBranch pulls a ghost branch on from ghost repo in WorkingEnv and returns a GhostBranch object
func (bs PullableDiffBranchSpec) PullBranch(we WorkingEnv) (GhostBranch, errors.GitGhostError) {
	resolved, err := bs.Resolve(we.SrcDir)
	if err != nil {
		return nil, err
	}
	branch := &DiffBranch{
		Prefix:         resolved.Prefix,
		CommitHashFrom: resolved.CommittishFrom,
		DiffHash:       bs.DiffHash,
	}
	err = pull(branch, we)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func pull(ghost GhostBranch, we WorkingEnv) errors.GitGhostError {
	return git.ResetHardToBranch(we.GhostDir, git.ORIGIN+"/"+ghost.BranchName())
}

func resolveCommittishOr(srcDir string, committishToResolve string) string {
	resolved, err := git.ResolveCommittish(srcDir, committishToResolve)
	if err != nil {
		log.WithFields(log.Fields{
			"repository": srcDir,
			"specified":  committishToResolve,
		}).Warn("can't resolve commit-ish value on local git repository.  specified commit-ish value will be used.")
		return committishToResolve
	}
	return resolved
}

func resolveFilepath(dir, p string) (string, errors.GitGhostError) {
	absp := p
	if !filepath.IsAbs(p) {
		absp = filepath.Clean(filepath.Join(dir, p))
	}
	relp, err := filepath.Rel(dir, absp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	log.WithFields(log.Fields{
		"dir":  dir,
		"path": p,
		"absp": absp,
		"relp": relp,
	}).Debugf("resolved path")
	if strings.HasPrefix(relp, "../") {
		return "", errors.Errorf("%s is not located in the source directory", p)
	}
	isdir, err := util.IsDir(relp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if isdir {
		return "", errors.Errorf("directory diff is not supported: %s", p)
	}
	return relp, nil
}
