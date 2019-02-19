package errors

import (
	"fmt"

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
	var fields log.Fields
	if log.GetLevel() == log.TraceLevel {
		fields = log.Fields{"stacktrace": fmt.Sprintf("%+v", err)}
	}
	log.WithFields(fields).Error(err)
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
