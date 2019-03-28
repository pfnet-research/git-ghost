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
