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
	return errors.WithStack(fmt.Errorf(s, args...)).(GitGhostError)
}

func New(s string) GitGhostError {
	return errors.WithStack(fmt.Errorf(s)).(GitGhostError)
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
