package ghost

import (
	"fmt"
)

type Commit struct {
	BaseCommitHash string
	Commits        []string
	Diff           string
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
	SrcDir          string
	GhostWorkingDir string
	GhostPrefix     string
	GhostRepo       string
	RemoteBase      string
	LocalBase       string
}

type PushResult struct {
	LocalBaseBranch *LocalBaseBranch
	LocalModBranch  *LocalModBranch
}

type PullOptions struct {
	SrcDir          string
	GhostWorkingDir string
	GhostPrefix     string
	GhostRepo       string
	RemoteBase      string
	LocalBase       string
	Hash            string
}
