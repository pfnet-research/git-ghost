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

package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"git-ghost/test/util"

	"github.com/stretchr/testify/assert"
)

var (
	ghostDir *util.WorkDir
)

func setup() error {
	dir, err := util.CreateGitWorkDir()
	if err != nil {
		return err
	}
	ghostDir = dir
	return nil
}

func teardown() error {
	if ghostDir != nil {
		defer ghostDir.Remove()
	}
	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		teardown()
		panic(err)
	}

	result := m.Run()

	err = teardown()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	os.Exit(result)
}

func TestTypeDefault(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	// Make one modification
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	stdout, _, err := srcDir.RunCommmand("git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	baseCommit := strings.TrimRight(stdout, "\n")

	stdout, _, err = srcDir.RunGitGhostCommmand("push")
	if err != nil {
		t.Fatal(err)
	}
	hashes := strings.Split(strings.TrimRight(stdout, "\n"), " ")
	assert.Equal(t, 2, len(hashes))
	diffBaseCommit := hashes[0]
	diffHash := hashes[1]
	assert.NotEqual(t, "", diffBaseCommit)
	assert.NotEqual(t, "", diffHash)

	stdout, _, err = srcDir.RunGitGhostCommmand("show", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, "-b\n+c\n")

	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "b\n", stdout)
	_, _, err = dstDir.RunGitGhostCommmand("pull", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "c\n", stdout)

	stdout, _, err = dstDir.RunGitGhostCommmand("list")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("delete", "--all")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("list")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))
}

func TestTypeCommits(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	stdout, _, err := srcDir.RunGitGhostCommmand("push", "commits", "HEAD~1")
	if err != nil {
		t.Fatal(err)
	}
	hashes := strings.Split(stdout, " ")
	assert.Equal(t, 2, len(hashes))
	baseCommit := hashes[0]
	targetCommit := hashes[1]
	assert.NotEqual(t, "", baseCommit)
	assert.NotEqual(t, "", targetCommit)

	stdout, _, err = srcDir.RunGitGhostCommmand("show", "commits", baseCommit, targetCommit)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, "-a\n+b\n")

	_, _, err = dstDir.RunCommmand("git", "checkout", baseCommit)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "a\n", stdout)
	_, _, err = dstDir.RunGitGhostCommmand("pull", "commits", baseCommit, targetCommit)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "b\n", stdout)

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "commits")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))

	stdout, _, err = dstDir.RunGitGhostCommmand("delete", "commits", "--all")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "commits")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
}

func TestTypeDiff(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	// Make one modification
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	stdout, _, err := srcDir.RunCommmand("git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	baseCommit := strings.TrimRight(stdout, "\n")

	stdout, _, err = srcDir.RunGitGhostCommmand("push", "diff")
	if err != nil {
		t.Fatal(err)
	}
	hashes := strings.Split(strings.TrimRight(stdout, "\n"), " ")
	assert.Equal(t, 2, len(hashes))
	diffBaseCommit := hashes[0]
	diffHash := hashes[1]
	assert.NotEqual(t, "", diffBaseCommit)
	assert.NotEqual(t, "", diffHash)

	stdout, _, err = srcDir.RunGitGhostCommmand("show", "diff", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, "-b\n+c\n")

	_, _, err = dstDir.RunGitGhostCommmand("pull", "diff", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "c\n", stdout)

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "diff")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("delete", "diff", "--all")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "diff")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))
}

func TestTypeAll(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	// Make one modification
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	stdout, _, err := srcDir.RunGitGhostCommmand("push", "all", "HEAD~1")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(stdout, "\n")
	hashes := strings.Split(lines[0], " ")
	assert.Equal(t, 2, len(hashes))
	baseCommit := hashes[0]
	targetCommit := hashes[1]
	assert.NotEqual(t, "", baseCommit)
	assert.NotEqual(t, "", targetCommit)

	hashes = strings.Split(lines[1], " ")
	assert.Equal(t, 2, len(hashes))
	diffBaseCommit := hashes[0]
	diffHash := hashes[1]
	assert.NotEqual(t, "", diffBaseCommit)
	assert.NotEqual(t, "", diffHash)

	stdout, _, err = srcDir.RunGitGhostCommmand("show", "all", baseCommit, targetCommit, diffHash)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, "-a\n+b\n")
	assert.Contains(t, stdout, "-b\n+c\n")

	_, _, err = dstDir.RunCommmand("git", "checkout", baseCommit)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "a\n", stdout)
	_, _, err = dstDir.RunGitGhostCommmand("pull", "all", baseCommit, targetCommit, diffHash)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "c\n", stdout)

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "all")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("delete", "all", "-v", "--all")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
	assert.Contains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))

	stdout, _, err = dstDir.RunGitGhostCommmand("list", "all")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
	assert.NotContains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))
}

func TestIncludeFile(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	// Make one modification
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Make a file
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo 'this is an included file' > included_file")
	if err != nil {
		t.Fatal(err)
	}

	stdout, _, err := srcDir.RunGitGhostCommmand("push", "--include", "included_file")
	if err != nil {
		t.Fatal(err)
	}
	hashes := strings.Split(strings.TrimRight(stdout, "\n"), " ")
	assert.Equal(t, 2, len(hashes))
	diffBaseCommit := hashes[0]
	diffHash := hashes[1]
	assert.NotEqual(t, "", diffBaseCommit)
	assert.NotEqual(t, "", diffHash)

	_, _, err = dstDir.RunGitGhostCommmand("pull", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "included_file")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "this is an included file\n", stdout)
}

func TestIncludeLinkFile(t *testing.T) {
	srcDir, dstDir, err := setupBasicEnv(ghostDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Remove()
	defer dstDir.Remove()

	// Make one modification
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Make a file
	_, _, err = srcDir.RunCommmand("bash", "-c", "echo 'this is an included file' > included_file")
	if err != nil {
		t.Fatal(err)
	}

	// Make a symlink to the file above
	_, _, err = srcDir.RunCommmand("bash", "-c", "ln -ns included_file included_link")
	if err != nil {
		t.Fatal(err)
	}

	stdout, _, err := srcDir.RunGitGhostCommmand("push", "--include", "included_link", "--follow-symlinks")
	if err != nil {
		t.Fatal(err)
	}
	hashes := strings.Split(strings.TrimRight(stdout, "\n"), " ")
	assert.Equal(t, 2, len(hashes))
	diffBaseCommit := hashes[0]
	diffHash := hashes[1]
	assert.NotEqual(t, "", diffBaseCommit)
	assert.NotEqual(t, "", diffHash)

	_, _, err = dstDir.RunGitGhostCommmand("pull", diffHash)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _, err = dstDir.RunCommmand("cat", "included_file")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "this is an included file\n", stdout)
	stdout, _, err = dstDir.RunCommmand("cat", "included_link")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "this is an included file\n", stdout)
}

func setupBasicEnv(workDir *util.WorkDir) (*util.WorkDir, *util.WorkDir, error) {
	srcDir, err := util.CreateGitWorkDir()
	if err != nil {
		return nil, nil, err
	}

	err = setupBasicGitRepo(srcDir)
	if err != nil {
		srcDir.Remove()
		return nil, nil, err
	}
	srcDir.Env = map[string]string{
		"GIT_GHOST_REPO": workDir.Dir,
	}

	dstDir, err := util.CloneWorkDir(srcDir)
	if err != nil {
		srcDir.Remove()
		return nil, nil, err
	}
	dstDir.Env = map[string]string{
		"GIT_GHOST_REPO": workDir.Dir,
	}
	return srcDir, dstDir, nil
}

func setupBasicGitRepo(wd *util.WorkDir) error {
	var err error
	_, _, err = wd.RunCommmand("bash", "-c", "echo a > sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "add", "sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "commit", "sample.txt", "-m", "initial commit")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("bash", "-c", "echo b > sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "commit", "sample.txt", "-m", "second commit")
	if err != nil {
		return err
	}
	return nil
}
