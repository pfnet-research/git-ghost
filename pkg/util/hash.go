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

package util

import (
	"crypto/sha1"
	"io"
	"strings"

	"github.com/pfnet-research/git-ghost/pkg/util/errors"
)

func GenerateFileContentHash(filepath string) (string, errors.GitGhostError) {
	input := strings.NewReader(filepath)

	hash := sha1.New()
	if _, err := io.Copy(hash, input); err != nil {
		return "", errors.WithStack(err)
	}
	sum := hash.Sum(nil)
	return string(sum), nil
}
