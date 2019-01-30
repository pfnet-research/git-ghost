package ghost

import (
	"fmt"
	"regexp"
	"sort"
)

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

type LocalBaseBranches []LocalBaseBranch
type LocalModBranches []LocalModBranch

var localBaseBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)-([a-f0-9]+)$`)
var localModBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)/([a-f0-9]+)$`)

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

func CreateGhostBranchByName(branchName string) GhostBranch {
	m := localBaseBranchNamePattern.FindStringSubmatch(branchName)
	if len(m) > 0 {
		return &LocalBaseBranch{
			Prefix:           m[1],
			RemoteBaseCommit: m[2],
			LocalBaseCommit:  m[3],
		}
	}
	m = localModBranchNamePattern.FindStringSubmatch(branchName)
	if len(m) > 0 {
		return &LocalModBranch{
			Prefix:          m[1],
			LocalBaseCommit: m[2],
			LocalModHash:    m[3],
		}
	}
	return nil
}

func (branches LocalBaseBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}

func (branches LocalModBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}
