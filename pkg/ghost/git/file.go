package git

import (
	"fmt"
	"git-ghost/pkg/util"
	"git-ghost/pkg/util/errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
)

// CreateDiffBundleFile creates patches for fromComittish..toComittish and save it to filepath
func CreateDiffBundleFile(dir, filepath, fromComittish, toComittish string) errors.GitGhostError {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.LogDeferredError(f.Close)

	return util.JustStreamOutputCmd(
		exec.Command("git", "-C", dir,
			"log", "-p", "--reverse", "--pretty=email", "--stat", "-m", "--first-parent", "--binary", fmt.Sprintf("%s..%s", fromComittish, toComittish),
		),
		f,
	)
}

// ApplyDiffBundleFile apply a patch file created in CreateDiffBundleFile
func ApplyDiffBundleFile(dir, filepath string) errors.GitGhostError {
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
	return errors.WithStack(errs)
}

// CreateDiffPatchFile creates a diff from comittish to current working state of `dir` and save it to filepath
func CreateDiffPatchFile(dir, filepath, comittish string) errors.GitGhostError {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.LogDeferredError(f.Close)

	return util.JustStreamOutputCmd(
		exec.Command("git", "-C", dir, "diff", "--patience", "--binary", comittish),
		f,
	)
	return nil
}

// AppendNonIndexedDiffFiles appends non-indexed diff files
func AppendNonIndexedDiffFiles(dir, filepath string, nonIndexedFilepaths []string) errors.GitGhostError {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.LogDeferredError(f.Close)

	var errs error
	for _, p := range nonIndexedFilepaths {
		cmd := exec.Command("git", "-C", dir, "diff", "--patience", "--binary", "--no-index", os.DevNull, p)
		cmd.Stdout = f
		ggerr := util.JustRunCmd(cmd)
		if ggerr != nil {
			if exiterr, ok := ggerr.Cause().(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					// exit 1 is valid for git diff
					if status.ExitStatus() == 1 {
						continue
					}
				}
			}
			errs = multierror.Append(errs, ggerr)
		}
	}
	return errors.WithStack(errs)
}

// ApplyDiffPatchFile apply a diff file created by CreateDiffPatchFile
func ApplyDiffPatchFile(dir, filepath string) errors.GitGhostError {
	return util.JustRunCmd(
		exec.Command("git", "-C", dir, "apply", filepath),
	)
}
