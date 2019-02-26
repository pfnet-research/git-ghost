package errors_test

import (
	"fmt"
	"git-ghost/pkg/util/errors"
	"testing"

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
