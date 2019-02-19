package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"git-ghost/test/util"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}

func TestAll(t *testing.T) {
	ghostDir, err := util.CreateGitWorkDir()
	if err != nil {
		t.Fatal(err)
	}
	defer ghostDir.Remove()

	t.Run("TypeDefault", CreateTestTypeDefault(ghostDir))
	t.Run("TypeCommits", CreateTestTypeCommits(ghostDir))
	t.Run("TypeDiff", CreateTestTypeDiff(ghostDir))
	t.Run("TypeAll", CreateTestTypeAll(ghostDir))
	t.Run("IncludeFile", CreateTestIncludeFile(ghostDir))
	t.Run("IncludeLinkFile", CreateTestIncludeLinkFile(ghostDir))
}

func CreateTestTypeDefault(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
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

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "push")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := stdout
		assert.NotEqual(t, "", diffHash)

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "show", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, "-b\n+c\n")

		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "b\n", stdout)
		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "c\n", stdout)

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "delete", "--all")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))
	}
}

func CreateTestTypeCommits(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
		srcDir, dstDir, err := setupBasicEnv(ghostDir)
		if err != nil {
			t.Fatal(err)
		}
		defer srcDir.Remove()
		defer dstDir.Remove()

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push", "commits", "HEAD~1")
		if err != nil {
			t.Fatal(err)
		}
		hashes := strings.Split(stdout, " ")
		assert.Equal(t, 2, len(hashes))
		baseCommit := hashes[0]
		targetCommit := hashes[1]
		assert.NotEqual(t, "", baseCommit)
		assert.NotEqual(t, "", targetCommit)

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "show", "commits", baseCommit, targetCommit)
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
		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", "commits", baseCommit, targetCommit)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "b\n", stdout)

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "commits")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "delete", "commits", "--all")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "commits")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
	}
}

func CreateTestTypeDiff(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
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

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "push", "diff")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := stdout
		assert.NotEqual(t, "", diffHash)

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "show", "diff", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, "-b\n+c\n")

		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", "diff", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "c\n", stdout)

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "diff")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "delete", "diff", "--all")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "diff")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, diffHash))
	}
}

func CreateTestTypeAll(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
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

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push", "all", "HEAD~1")
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
		diffHash := lines[1]
		assert.NotEqual(t, "", diffHash)

		stdout, _, err = srcDir.RunCommmand("git", "ghost", "show", "all", baseCommit, targetCommit, diffHash)
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
		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", "all", baseCommit, targetCommit, diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "c\n", stdout)

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "all")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "delete", "all", "-v", "--all")
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
		assert.Contains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))

		stdout, _, err = dstDir.RunCommmand("git", "ghost", "list", "all")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, stdout, fmt.Sprintf("%s %s", baseCommit, targetCommit))
		assert.NotContains(t, stdout, fmt.Sprintf("%s %s", targetCommit, diffHash))
	}
}

func CreateTestIncludeFile(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
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

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push", "--include", "included_file")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := strings.TrimRight(stdout, "\n")
		assert.NotEqual(t, "", diffHash)

		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "included_file")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "this is an included file\n", stdout)
	}
}

func CreateTestIncludeLinkFile(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
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

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push", "--include", "included_link", "--follow-symlinks")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := strings.TrimRight(stdout, "\n")
		assert.NotEqual(t, "", diffHash)

		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", diffHash)
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
