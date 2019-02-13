package util_test

import (
	"errors"
	"fmt"
	"testing"

	"git-ghost/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	err := errors.New("foo")

	// wrap once
	gitghosterr := util.WrapError(err, "bar")
	assert.Equal(t, "bar", gitghosterr.ShortError())
	assert.Equal(t, "bar: foo", gitghosterr.Error())
	assert.Equal(t, "errors.go", fmt.Sprintf("%s", gitghosterr.StackTrace()[0]))

	// wrap twice
	gitghosterr = util.WrapError(gitghosterr, "wow")
	assert.Equal(t, "wow", gitghosterr.ShortError())
	assert.Equal(t, "wow: bar: foo", gitghosterr.Error())
	assert.Equal(t, "errors.go", fmt.Sprintf("%s", gitghosterr.StackTrace()[0]))
}
