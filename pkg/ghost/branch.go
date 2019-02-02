package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"sort"

	log "github.com/Sirupsen/logrus"
)

type GhostBranch interface {
	BranchName() string
	FileName() string
	Show(we WorkingEnv, writer io.Writer) error
	Apply(we WorkingEnv, forceApply bool) error
}

// interface assetions
var _ GhostBranch = LocalBaseBranch{}
var _ GhostBranch = LocalModBranch{}

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

func show(ghost GhostBranch, we WorkingEnv, writer io.Writer) error {
	cmd := exec.Command("git", "-C", we.GhostDir, "--no-pager", "cat-file", "-p", fmt.Sprintf("HEAD:%s", ghost.FileName()))
	cmd.Stdout = writer
	return util.JustRunCmd(cmd)
}

func apply(ghost GhostBranch, we WorkingEnv, expectedSrcHead string, forceApply bool) error {
	log.WithFields(util.MergeFields(
		util.ToFields(ghost),
		log.Fields{
			"ghostDir":        we.GhostDir,
			"srcDir":          we.SrcDir,
			"expectedSrcHead": expectedSrcHead,
			"forceApply":      forceApply,
		},
	)).Info("applying ghost branch")

	srcHead, err := git.ResolveComittish(we.SrcDir, "HEAD")
	if err != nil {
		return err
	}

	if srcHead != expectedSrcHead {
		message := "HEAD is not equal to expected"
		if forceApply {
			log.WithFields(util.MergeFields(
				util.ToFields(ghost),
				log.Fields{
					"actualSrcHead":   srcHead,
					"expectedSrcHead": expectedSrcHead,
					"srcDir":          we.SrcDir,
				},
			),
			).Warnf("%s. Applying ghost branch will be failed.", message)
		} else {
			return fmt.Errorf("abort because %s (actual=%s, expected=%s)", message, srcHead, expectedSrcHead)
		}
	}

	// TODO make this instance methods.
	switch ghost.(type) {
	case LocalBaseBranch:
		err = git.ApplyDiffBundleFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))
	case LocalModBranch:
		err = git.ApplyDiffPatchFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))
	default:
		return fmt.Errorf("not supported on type = %+v", reflect.TypeOf(ghost))
	}

	if err != nil {
		return err
	}
	return nil
}

func (bs LocalBaseBranch) Show(we WorkingEnv, writer io.Writer) error {
	return show(bs, we, writer)
}

func (bs LocalBaseBranch) Apply(we WorkingEnv, forceApply bool) error {
	err := apply(bs, we, bs.RemoteBaseCommit, forceApply)
	if err != nil {
		return err
	}
	return nil
}

func (bs LocalModBranch) Show(we WorkingEnv, writer io.Writer) error {
	return show(bs, we, writer)
}

func (bs LocalModBranch) Apply(we WorkingEnv, forceApply bool) error {
	err := apply(bs, we, bs.LocalBaseCommit, forceApply)
	if err != nil {
		return err
	}
	return nil
}
