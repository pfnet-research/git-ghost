package util

import (
	log "github.com/Sirupsen/logrus"
)

// LogError calls a given function and log errors according to results
func LogError(f func() error) {
	err := f()
	if err != nil {
		log.Errorf("Error during defered call: %s", err)
	}
}
