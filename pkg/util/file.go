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
	"github.com/pfnet-research/git-ghost/pkg/util/errors"
	"os"
	"path/filepath"
)

// FileSize returns file size of a given file
func FileSize(filepath string) (int64, errors.GitGhostError) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return fi.Size(), nil
}

// WalkSymlink reads a symlink and call a given callback until the resolved path is not a symlink.WalkSymlink
func WalkSymlink(dir, path string, cb func([]string, string) errors.GitGhostError) errors.GitGhostError {
	abspath := path
	if !filepath.IsAbs(path) {
		abspath = filepath.Clean(filepath.Join(dir, path))
	}
	islink, err := IsSymlink(abspath)
	if err != nil {
		return err
	}
	if !islink {
		return errors.Errorf("%s is not a symlink", abspath)
	}

	resolved := abspath
	paths := []string{path}
	for {
		abspath = resolved
		if !filepath.IsAbs(resolved) {
			abspath = filepath.Clean(filepath.Join(dir, resolved))
		}
		islink, ggerr := IsSymlink(abspath)
		if ggerr != nil {
			return ggerr
		}
		if !islink {
			break
		}
		path, err := os.Readlink(abspath)
		if err != nil {
			return errors.WithStack(err)
		}
		ggerr = cb(paths, path)
		if ggerr != nil {
			return errors.WithStack(ggerr)
		}
		resolved = path
		paths = append(paths, path)
	}
	return nil
}

// IsSymlink returns whether a given file is a directory
func IsDir(path string) (bool, errors.GitGhostError) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return stat.IsDir(), nil
}

// IsSymlink returns whether a given file is a symlink
func IsSymlink(path string) (bool, errors.GitGhostError) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}
