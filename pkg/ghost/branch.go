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

// GhostBranch is an interface representing a ghost branch.
//
// It is created from GhostBranchSpec/PullableGhostBranchSpec
type GhostBranch interface {
	// BranchName returns its full branch name on git repository
	BranchName() string
	// FileName returns a file name contained in the GhostBranch
	FileName() string
	// Show writes contents of this ghost branch on passed working env to writer
	Show(we WorkingEnv, writer io.Writer) error
	// Apply applies contents(diff or patch) of this ghost branch on passed working env
	Apply(we WorkingEnv) error
}

// interface assetions
var _ GhostBranch = LocalBaseBranch{}
var _ GhostBranch = LocalModBranch{}

// LocalBaseBranch represents a local base branch
//
// This contains patches for RemoteBaseCommit..LocalBaseCommit
type LocalBaseBranch struct {
	Prefix           string
	RemoteBaseCommit string
	LocalBaseCommit  string
}

// LocalModBranch represents a local mod branch
//
// This contains diff
// - whose content hash value is LocalModHash
// - which is generated on LocalBaseCommit
type LocalModBranch struct {
	// Prefix is a prefix of branch name
	Prefix string
	// LocalBaseCommit is full commit hash to which this local mod branch's diff contains
	LocalBaseCommit string
	// LocalModHash is a hash value of its diff
	LocalModHash string
}

// LocalBaseBranches is an alias for []LocalBaseBranch
type LocalBaseBranches []LocalBaseBranch

// LocalModBranches is an alias for []LocalModBranch
type LocalModBranches []LocalModBranch

var localBaseBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)-([a-f0-9]+)$`)
var localModBranchNamePattern = regexp.MustCompile(`^([a-z0-9]+)/([a-f0-9]+)/([a-f0-9]+)$`)

// BranchName returns its full branch name on git repository
func (b LocalBaseBranch) BranchName() string {
	return fmt.Sprintf("%s/%s-%s", b.Prefix, b.RemoteBaseCommit, b.LocalBaseCommit)
}

// FileName returns a file name containing this GhostBranch
func (b LocalBaseBranch) FileName() string {
	return "commits.patch"
}

// BranchName returns its full branch name on git repository
func (b LocalModBranch) BranchName() string {
	return fmt.Sprintf("%s/%s/%s", b.Prefix, b.LocalBaseCommit, b.LocalModHash)
}

// FileName returns a file name containing this GhostBranch
func (b LocalModBranch) FileName() string {
	return "local-mod.patch"
}

// CreateGhostBranchByName instantiates GhostBranch object from branchname
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

// Sort sorts passed branches in lexicographic order of BranchName()
func (branches LocalBaseBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}

// AsGhostBranches just lifts item type to GhostBranch
func (branches LocalBaseBranches) AsGhostBranches() []GhostBranch {
	ghostBranches := make([]GhostBranch, len(branches))
	for i, branch := range branches {
		ghostBranches[i] = branch
	}
	return ghostBranches
}

// Sort sorts passed branches in lexicographic order of BranchName()
func (branches LocalModBranches) Sort() {
	sortFunc := func(i, j int) bool {
		return branches[i].BranchName() < branches[j].BranchName()
	}
	sort.Slice(branches, sortFunc)
}

// AsGhostBranches just lifts item type to GhostBranch
func (branches LocalModBranches) AsGhostBranches() []GhostBranch {
	ghostBranches := make([]GhostBranch, len(branches))
	for i, branch := range branches {
		ghostBranches[i] = branch
	}
	return ghostBranches
}

func show(ghost GhostBranch, we WorkingEnv, writer io.Writer) error {
	cmd := exec.Command("git", "-C", we.GhostDir, "--no-pager", "cat-file", "-p", fmt.Sprintf("HEAD:%s", ghost.FileName()))
	cmd.Stdout = writer
	return util.JustRunCmd(cmd)
}

func apply(ghost GhostBranch, we WorkingEnv, expectedSrcHead string) error {
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
	case LocalBaseBranch:
		return git.ApplyDiffBundleFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))
	case LocalModBranch:
		return git.ApplyDiffPatchFile(we.SrcDir, path.Join(we.GhostDir, ghost.FileName()))

	default:
		return fmt.Errorf("not supported on type = %+v", reflect.TypeOf(ghost))
	}
}

// Show writes contents of this ghost branch on passed working env to writer
func (bs LocalBaseBranch) Show(we WorkingEnv, writer io.Writer) error {
	return show(bs, we, writer)
}

// Apply applies contents(diff or patch) of this ghost branch on passed working env
func (bs LocalBaseBranch) Apply(we WorkingEnv) error {
	err := apply(bs, we, bs.RemoteBaseCommit)
	if err != nil {
		return err
	}
	return nil
}

// Show writes contents of this ghost branch on passed working env to writer
func (bs LocalModBranch) Show(we WorkingEnv, writer io.Writer) error {
	return show(bs, we, writer)
}

// Apply applies contents(diff or patch) of this ghost branch on passed working env
func (bs LocalModBranch) Apply(we WorkingEnv) error {
	err := apply(bs, we, bs.LocalBaseCommit)
	if err != nil {
		return err
	}
	return nil
}
