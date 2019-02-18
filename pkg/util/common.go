package util

import (
	log "github.com/Sirupsen/logrus"
)

// LogDeferredError calls a given function and log an error according to the result
func LogDeferredError(f func() error) {
	err := f()
	if err != nil {
		log.Debugf("Error during defered call: %s", err)
	}
}
