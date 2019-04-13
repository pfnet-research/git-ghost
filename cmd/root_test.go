package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	assert.NotNil(t, cmd)
}
