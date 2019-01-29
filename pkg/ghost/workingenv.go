package ghost

import (
	"git-ghost/pkg/ghost/git"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

type WorkingEnvSpec struct {
	SrcDir          string
	GhostWorkingDir string
	GhostRepo       string
}

type WorkingEnv struct {
	WorkingEnvSpec
	GhostDir string
}

func (weSpec WorkingEnvSpec) initialize() (*WorkingEnv, error) {
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
	}).Info("ghost repo was cloned")

	return &WorkingEnv{
		WorkingEnvSpec: weSpec,
		GhostDir:       ghostDir,
	}, nil
}

func (weSpec WorkingEnv) clean() {
	os.RemoveAll(weSpec.GhostDir)
}
