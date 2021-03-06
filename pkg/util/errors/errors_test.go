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

package errors_test

import (
	"fmt"
	"testing"

	"github.com/pfnet-research/git-ghost/pkg/util/errors"

	"github.com/stretchr/testify/assert"
)

func TestLogErrorWithStack(t *testing.T) {
	errors.LogErrorWithStack(errors.New("foo"))
}

func TestNew(t *testing.T) {
	err := errors.New("foo")
	assert.Equal(t, "foo", err.Error())
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("%s", "foo")
	assert.Equal(t, "foo", err.Error())
}

func TestWithStack(t *testing.T) {
	// for normal error
	err := errors.WithStack(fmt.Errorf("foo"))
	assert.Equal(t, "foo", err.Error())
	// for GitGhostError
	err = errors.WithStack(errors.New("foo"))
	assert.Equal(t, "foo", err.Error())
}
