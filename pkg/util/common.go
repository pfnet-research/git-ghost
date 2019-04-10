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

package util

import (
	"git-ghost/pkg/util/errors"

	log "github.com/Sirupsen/logrus"
)

const (
	CommitStartFromInit = "_"
)

// LogDeferredError calls a given function and log an error according to the result
func LogDeferredError(f func() error) {
	err := f()
	if err != nil {
		log.Debugf("Error during defered call: %s", err)
	}
}

// LogDeferredGitGhostError calls a given function and log an GitGhostError according to the result
func LogDeferredGitGhostError(f func() errors.GitGhostError) {
	err := f()
	if err != nil {
		log.Errorf("Error during defered call: %s", err)
	}
}
