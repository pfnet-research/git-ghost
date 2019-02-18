package errors

import (
	"github.com/pkg/errors"

	log "github.com/Sirupsen/logrus"
)

type GitGhostError interface {
	StackTrace() errors.StackTrace
	Error() string
	Cause() error
}

// LogErrorWithStack emits a log message with errors.GitGhostError level and stack trace with debug level
func LogErrorWithStack(err GitGhostError) {
	log.Error(err)
	log.Tracef("%+v", err)
}

func Errorf(s string, args ...interface{}) GitGhostError {
	return errors.Errorf(s, args...).(GitGhostError)
}

func New(s string) GitGhostError {
	return errors.New(s).(GitGhostError)
}

func WithStack(err error) GitGhostError {
	if err == nil {
		return nil
	}
	if gitghosterr, ok := err.(GitGhostError); ok {
		return gitghosterr
	}
	return errors.WithStack(err).(GitGhostError)
}
