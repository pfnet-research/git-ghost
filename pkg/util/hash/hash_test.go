package hash_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pfnet-research/git-ghost/pkg/util/hash"
	"github.com/stretchr/testify/assert"
)

func CalculateHashWithCommand(filepath string) (string, error) {
	cmd := exec.Command("sha1sum", "-b", filepath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	hash := strings.Split(string(output), " ")[0]
	return hash, nil
}

func TestHashCompatibility(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "tempfile-test-")
	if err != nil {
		t.Fatal(err)
	}
	oldHash, err := CalculateHashWithCommand(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	newHash, err := hash.GenerateFileContentHash(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, oldHash, newHash)
}
