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

package cmd

import (
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util/errors"
)

func nonEmpty(name, value string) errors.GitGhostError {
	if value == "" {
		return errors.Errorf("%s must not be empty", name)
	}
	return nil
}

func isValidComittish(name, comittish string) errors.GitGhostError {
	err := git.ValidateComittish(globalOpts.srcDir, comittish)
	if err != nil {
		return errors.Errorf("%s is not a valid object", name)
	}
	return nil
}
