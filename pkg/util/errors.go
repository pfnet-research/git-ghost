package util

import (
	"github.com/pkg/errors"
)

type GitGhostError interface {
	Error() string
	ShortError() string
	StackTrace() errors.StackTrace
}

type gitghosterr struct {
	message string
	wrappedError
}

// stackTracer is interface for error types that have a stack trace
type wrappedError interface {
	Error() string
	Cause() error
	StackTrace() errors.StackTrace
}

func (e gitghosterr) Error() string {
	return e.wrappedError.Error()
}

func (e gitghosterr) ShortError() string {
	return e.message
}

func (e gitghosterr) StackTrace() errors.StackTrace {
	return e.wrappedError.StackTrace()
}

func WrapError(err error, message string) GitGhostError {
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, message)
	return gitghosterr{message, err.(wrappedError)}
}
