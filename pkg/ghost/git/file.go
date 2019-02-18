package git

import (
	"bytes"
	"errors"
	"fmt"
	"git-ghost/pkg/util"
	"io"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
)

// CreateDiffBundleFile creates patches for fromComittish..toComittish and save it to filepath
func CreateDiffBundleFile(dir, filepath, fromComittish, toComittish string) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer util.LogDeferredError(f.Close)

	cmd := exec.Command("git", "-C", dir,
		"log", "-p", "--reverse", "--pretty=email", "--stat", "-m", "--first-parent", "--binary",
		fmt.Sprintf("%s..%s", fromComittish, toComittish),
	)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer util.LogDeferredError(reader.Close)
	err = cmd.Start()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}

	total := 0
	buf := make([]byte, 1024, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			d := buf[:n]
			_, err = f.Write(d)
			if err != nil {
				return err
			}
			total += n
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

// ApplyDiffBundleFile apply a patch file created in CreateDiffBundleFile
func ApplyDiffBundleFile(dir, filepath string) error {
	var errs error
	err := util.JustRunCmd(
		exec.Command("git", "-C", dir, "am", filepath),
	)
	if err != nil {
		errs = multierror.Append(errs, err)
		log.WithFields(util.MergeFields(
			log.Fields{
				"srcDir":   dir,
				"filepath": filepath,
				"error":    err.Error(),
			})).Info("apply('git am') failed. aborting.")
		resetErr := util.JustRunCmd(
			exec.Command("git", "-C", dir, "am", "--abort"),
		)
		if resetErr != nil {
			errs = multierror.Append(errs, resetErr)
		}
	}
	return errs
}

// CreateDiffPatchFile creates a diff from comittish to current working state of `dir` and save it to filepath
func CreateDiffPatchFile(dir, filepath, comittish string) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer util.LogDeferredError(f.Close)

	cmd := exec.Command("git", "-C", dir, "diff", "--patience", "--binary", comittish)
	stderr := bytes.NewBufferString("")
	cmd.Stderr = stderr
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer util.LogDeferredError(reader.Close)
	err = cmd.Start()
	if err != nil {
		s := stderr.String()
		if s != "" {
			return errors.New(s)
		}
		return err
	}

	total := 0
	buf := make([]byte, 1024, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			d := buf[:n]
			_, err = f.Write(d)
			if err != nil {
				return err
			}
			total += n
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

// AppendNonIndexedDiffFiles appends non-indexed diff files
func AppendNonIndexedDiffFiles(dir, filepath string, nonIndexedFilepaths []string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer util.LogDeferredError(f.Close)

	var errs error
	for _, p := range nonIndexedFilepaths {
		cmd := exec.Command("git", "-C", dir, "diff", "--patience", "--binary", "--no-index", os.DevNull, p)
		cmd.Stdout = f
		err = util.JustRunCmd(cmd)
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					// exit 1 is valid for git diff
					if status.ExitStatus() == 1 {
						continue
					}
				}
			}
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// ApplyDiffPatchFile apply a diff file created by CreateDiffPatchFile
func ApplyDiffPatchFile(dir, filepath string) error {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "apply", filepath),
	)
}
