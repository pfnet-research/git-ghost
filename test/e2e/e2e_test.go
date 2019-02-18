package e2e

import (
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

	t.Run("BasicScenario", CreateTestBasicScenario(ghostDir))
	t.Run("IncludeFile", CreateTestIncludeFile(ghostDir))
	t.Run("IncludeLinkFile", CreateTestIncludeLinkFile(ghostDir))
}

func CreateTestBasicScenario(ghostDir *util.WorkDir) func(t *testing.T) {
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

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := strings.TrimRight(stdout, "\n")
		assert.NotEqual(t, "", diffHash)

		_, _, err = srcDir.RunCommmand("git", "ghost", "show", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: Do some assertion

		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "c\n", stdout)

		_, _, err = dstDir.RunCommmand("git", "ghost", "list")
		if err != nil {
			t.Fatal(err)
		}
		// TODO: Do some assertion

		// TODO: delete the ghost branches and do some assertion
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

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "-vvv", "push", "--include", "included_file")
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
