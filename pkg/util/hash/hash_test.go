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
