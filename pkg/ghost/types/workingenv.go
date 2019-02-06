package types

import (
	"git-ghost/pkg/ghost/git"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// WorkingEnvSpec abstract an environment git-ghost works with
type WorkingEnvSpec struct {
	// SrcDir is local git directory
	SrcDir string
	// GhostWorkingDir is a root directory which git-ghost creates temporary directories
	GhostWorkingDir string
	// GhostRepo is a repository url git-ghost works with
	GhostRepo string
}

// WorkingEnv is initialized environment containing temporary local ghost repository
type WorkingEnv struct {
	WorkingEnvSpec
	GhostDir string
}

func (weSpec WorkingEnvSpec) Initialize() (*WorkingEnv, error) {
	ghostDir, err := ioutil.TempDir(weSpec.GhostWorkingDir, "git-ghost-")
	if err != nil {
		return nil, err
	}
	err = git.InitializeGitDir(ghostDir, weSpec.GhostRepo, "")
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"dir": ghostDir,
	}).Debug("ghost repo was cloned")

	return &WorkingEnv{
		WorkingEnvSpec: weSpec,
		GhostDir:       ghostDir,
	}, nil
}

func (weSpec WorkingEnv) Clean() {
	os.RemoveAll(weSpec.GhostDir)
}
