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
	"io"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"sort"

	log "github.com/Sirupsen/logrus"
)

// GhostBranch is an interface representing a ghost branch.
//
// It is created from GhostBranchSpec/PullableGhostBranchSpec
type GhostBranch interface {
	// BranchName returns its full branch name on git repository
	BranchName() string
	// FileName returns a file name contained in the GhostBranch
	FileName() string
	// Show writes contents of this ghost branch on passed working env to writer
	Show(we WorkingEnv, writer io.Writer) errors.GitGhostError
	// Apply applies contents(diff or patch) of this ghost branch on passed working env
	Apply(we WorkingEnv) errors.GitGhostError
}

// interface assetions
var _ GhostBranch = CommitsBranch{}
var _ GhostBranch = DiffBranch{}

// CommitsBranch represents a local base branch
//
// This contains patches for CommitHashFrom..CommitHashTo
type CommitsBranch struct {
	Prefix         string
	CommitHashFrom string
	CommitHashTo   string
}

// DiffBranch represents a local mod branch
//
// This contains diff
// - whose content hash value is DiffHash
// - which is generated on CommitHashFrom
type DiffBranch struct {
	// Prefix is a prefix of branch name
	Prefix string
	// CommitHashFrom is full commit hash to which this local mod branch's diff contains
	CommitHashFrom string
	// DiffHash is a hash value of its diff
	DiffHash string
}

// CommitsBranches is an alias for []CommitsBranch
type CommitsBranches []CommitsBranch

// DiffBranches is an alias for []DiffBranch
type DiffBranches []DiffBranch

var commitsBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)-([a-f0-9]+)$`)
var diffBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)/([a-f0-9]+)$`)

// BranchName returns its full branch name on git repository
func (b CommitsBranch) BranchName() string {
	return fmt.Sprintf("%s/%s-%s", b.Prefix, b.CommitHashFrom, b.CommitHashTo)
}

// FileName returns a file name containing this GhostBranch
func (b CommitsBranch) FileName() string {
	return "commits.patch"
}

// BranchName returns its full branch name on git repository
func (b DiffBranch) BranchName() string {
	return fmt.Sprintf("%s/%s/%s", b.Prefix, b.CommitHashFrom, b.DiffHash)
}

// FileName returns a file name containing this GhostBranch
func (b DiffBranch) FileName() string {
	return "local-mod.patch"
}

// CreateGhostBranchByName instantiates GhostBranch object from branchname
func CreateGhostBranchByName(branchName string) GhostBranch {
	m := commitsBranchNamePattern.FindStringSubmatch(branchName)
	if len(m) > 0 {
		return &CommitsBranch{
			Prefix:         m[1],
			CommitHashFrom: m[2],
			CommitHashTo:   m[3],
		}
	}
	m = diffBranchNamePattern.FindStringSubmatch(branchName)
	if len(m) > 0 {
		return &DiffBranch{
			Prefix:         m[1],
			CommitHashFrom: m[2],
			DiffHash:       m[3],
		}
	}
	return nil
}

// Sort sorts passed branches in lexicographic order of BranchName()
func (branches CommitsBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}

// AsGhostBranches just lifts item type to GhostBranch
func (branches CommitsBranches) AsGhostBranches() []GhostBranch {
	ghostBranches := make([]GhostBranch, len(branches))
	for i, branch := range branches {
		ghostBranches[i] = branch
	}
	return ghostBranches
}

// Sort sorts passed branches in lexicographic order of BranchName()
func (branches DiffBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}

// AsGhostBranches just lifts item type to GhostBranch
func (branches DiffBranches) AsGhostBranches() []GhostBranch {
	ghostBranches := make([]GhostBranch, len(branches))
	for i, branch := range branches {
		ghostBranches[i] = branch
	}
	return ghostBranches
}

func show(ghost GhostBranch, we WorkingEnv, writer io.Writer) errors.GitGhostError {
	cmd := exec.Command("git", "-C", we.GhostDir, "--no-pager", "cat-file", "-p", fmt.Sprintf("HEAD:%s", ghost.FileName()))
	cmd.Stdout = writer
	return util.JustRunCmd(cmd)
}

func apply(ghost GhostBranch, we WorkingEnv, expectedSrcHead string) errors.GitGhostError {
	log.WithFields(util.MergeFields(
		util.ToFields(ghost),
		log.Fields{
			"ghostDir":        we.GhostDir,
			"srcDir":          we.SrcDir,
			"expectedSrcHead": expectedSrcHead,
		},
	)).Info("applying ghost branch")

	srcHead, err := git.ResolveComittish(we.SrcDir, "HEAD")
	if err != nil {
		return err
	}

	if srcHead != expectedSrcHead {
		message := "HEAD is not equal to expected"
		log.WithFields(util.MergeFields(
			util.ToFields(ghost),
			log.Fields{
				"actualSrcHead":   srcHead,
				"expectedSrcHead": expectedSrcHead,
				"srcDir":          we.SrcDir,
			},
		),
		).Warnf("%s. Applying ghost branch might be failed.", message)
	}

	// TODO make this instance methods.
	switch ghost.(type) {
	case CommitsBranch:
		return git.ApplyDiffBundleFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))
	case DiffBranch:
		return git.ApplyDiffPatchFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))

	default:
		return errors.Errorf("not supported on type = %+v", reflect.TypeOf(ghost))
	}
}

// Show writes contents of this ghost branch on passed working env to writer
func (bs CommitsBranch) Show(we WorkingEnv, writer io.Writer) errors.GitGhostError {
	return show(bs, we, writer)
}

// Apply applies contents(diff or patch) of this ghost branch on passed working env
func (bs CommitsBranch) Apply(we WorkingEnv) errors.GitGhostError {
	err := apply(bs, we, bs.CommitHashFrom)
	if err != nil {
		return err
	}
	return nil
}

// Show writes contents of this ghost branch on passed working env to writer
func (bs DiffBranch) Show(we WorkingEnv, writer io.Writer) errors.GitGhostError {
	return show(bs, we, writer)
}

// Apply applies contents(diff or patch) of this ghost branch on passed working env
func (bs DiffBranch) Apply(we WorkingEnv) errors.GitGhostError {
	err := apply(bs, we, bs.CommitHashFrom)
	if err != nil {
		return err
	}
	return nil
}
