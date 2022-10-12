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

package hash

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/pfnet-research/git-ghost/pkg/util/errors"
)

func GenerateFileContentHash(filepath string) (string, errors.GitGhostError) {
	// ref: https://pkg.go.dev/crypto/sha1#example-New-File
	f, err := os.Open(filepath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.WithStack(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
