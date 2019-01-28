package git

import (
	"bytes"
	"errors"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	ORIGIN string = "origin"
)

// CreateTempGitDir creates a temporary directory for a specified git repo.
// It is the caller's responsibility to remove the directory when no longer needed.
// e.g. defer os.RemoveAll(dir)
func CreateTempGitDir(dir, repo, branch string) (string, error) {
	dirPath, err := ioutil.TempDir(dir, "git-ghost-")
	if err != nil {
		return "", err
	}

	args := []string{"clone", "-q", "-o", ORIGIN} // assuring default remote name to be ORIGIN
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repo, dirPath)
	cmd := exec.Command("git", args...)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(dirPath)
		s := stderr.String()
		if s != "" {
			return "", errors.New(s)
		}
		return "", err
	}

	return dirPath, nil
}

func CommitAndPush(dir, filename, message, refspec string) error {
	err := CommitFile(dir, filename, message)
	if err != nil {
		return err
	}
	err = Push(dir, refspec)
	if err != nil {
		return err
	}
	return nil
}

func CommitFile(dir, filename, message string) error {
	cmd := exec.Command("git", "-C", dir, "add", filename)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	cmd = exec.Command("git", "-C", dir, "commit", "-q", filename, "-m", message)
	stderr.Reset()
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}

func Push(dir string, refspecs ...string) error {
	args := []string{"-C", dir, "push", "origin"}
	args = append(args, refspecs...)
	cmd := exec.Command("git", args...)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}

func Pull(dir, refspec string) error {
	cmd := exec.Command("git", "-C", dir, "pull", "origin", refspec)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}

func CreateOrphanBranch(dir, branch string) error {
	cmd := exec.Command("git", "-C", dir, "checkout", "--orphan", branch)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}
	return nil
}

func ResetHardToBranch(dir, branch string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "reset", "--hard", branch),
	)
}
