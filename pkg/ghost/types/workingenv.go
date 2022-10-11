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
	"os"

	"github.com/pfnet-research/git-ghost/pkg/ghost/git"
	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	log "github.com/sirupsen/logrus"
)

const (
	defaultGhostUserName  = "Git Ghost"
	defaultGhostUserEmail = "git-ghost@example.com"
)

// WorkingEnvSpec abstract an environment git-ghost works with
type WorkingEnvSpec struct {
	// SrcDir is local git directory
	SrcDir string
	// GhostWorkingDir is a root directory which git-ghost creates temporary directories
	GhostWorkingDir string
	// GhostRepo is a repository url git-ghost works with
	GhostRepo string
	// GhostUserName is a user name which is used in ghost working directories.
	GhostUserName string
	// GhostUserEmail is a user email which is used in ghost working directories.
	GhostUserEmail string
}

// WorkingEnv is initialized environment containing temporary local ghost repository
type WorkingEnv struct {
	WorkingEnvSpec
	GhostDir string
}

func (weSpec WorkingEnvSpec) Initialize() (*WorkingEnv, errors.GitGhostError) {
	ghostDir, err := os.MkdirTemp(weSpec.GhostWorkingDir, "git-ghost-")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ggerr := git.InitializeGitDir(ghostDir, weSpec.GhostRepo, "")
	if ggerr != nil {
		return nil, ggerr
	}
	ghostUserName := defaultGhostUserName
	if weSpec.GhostUserName != "" {
		ghostUserName = weSpec.GhostUserName
	}
	ghostUserEmail := defaultGhostUserEmail
	if weSpec.GhostUserEmail != "" {
		ghostUserEmail = weSpec.GhostUserEmail
	}
	ggerr = git.SetUserConfig(ghostDir, ghostUserName, ghostUserEmail)
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
