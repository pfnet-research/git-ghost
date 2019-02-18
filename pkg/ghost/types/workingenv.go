package types

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util/errors"
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

func (weSpec WorkingEnvSpec) Initialize() (*WorkingEnv, errors.GitGhostError) {
	ghostDir, err := ioutil.TempDir(weSpec.GhostWorkingDir, "git-ghost-")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ggerr := git.InitializeGitDir(ghostDir, weSpec.GhostRepo, "")
	if ggerr != nil {
		return nil, ggerr
	}

	log.WithFields(log.Fields{
		"dir": ghostDir,
	}).Debug("ghost repo was cloned")

	return &WorkingEnv{
		WorkingEnvSpec: weSpec,
		GhostDir:       ghostDir,
	}, nil
}

func (weSpec WorkingEnv) Clean() errors.GitGhostError {
	return errors.WithStack(os.RemoveAll(weSpec.GhostDir))
}
