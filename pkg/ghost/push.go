package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Push(options PushOptions) (string, error) {
	repo := options.GhostRepo
	prefix := options.GhostPrefix
	srcDir := options.SrcDir
	dstDir, err := git.CreateTempGitDir(options.DstDir, options.GhostRepo, "")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dstDir)

	err = git.ValidateRefspec(srcDir, options.RemoteBase)
	if err != nil {
		return "", err
	}
	err = git.ValidateRefspec(srcDir, options.LocalBase)
	if err != nil {
		return "", err
	}

	remoteBase, err := git.ResolveRefspec(srcDir, options.RemoteBase)
	if err != nil {
		return "", err
	}
	localBase, err := git.ResolveRefspec(srcDir, options.LocalBase)
	if err != nil {
		return "", err
	}

	branches := []GhostBranch{}
	if remoteBase != localBase {
		branch := LocalBaseBranch{
			RemoteBaseCommit: remoteBase,
			LocalBaseCommit:  localBase,
			Prefix:           prefix,
		}
		existence, err := git.ValidateRemoteBranchExistence(repo, branch.BranchName())
		if err != nil {
			return "", err
		}
		if existence {
			fmt.Printf("Skipped an existing local base branch '%s' in %s\n", branch.BranchName(), repo)
		} else {
			diffFilePath := filepath.Join(dstDir, branch.FileName())
			err := git.CreateDiffBundleFile(srcDir, diffFilePath, branch.RemoteBaseCommit, branch.LocalBaseCommit)
			if err != nil {
				return "", err
			}
			err = git.CreateOrphanBranch(dstDir, branch.BranchName())
			if err != nil {
				return "", err
			}
			err = git.CommitFile(dstDir, branch.FileName(), "this is a test")
			if err != nil {
				return "", err
			}
			branches = append(branches, branch)
			fmt.Printf("Created a local base branch '%s' in %s\n", branch.BranchName(), repo)
		}
	}

	tmpFile, err := ioutil.TempFile("", "git-ghost-local-mod")
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	err = git.CreateDiffPatchFile(srcDir, tmpFile.Name(), localBase)
	if err != nil {
		return "", err
	}
	size, err := util.FileSize(tmpFile.Name())
	if err != nil {
		return "", err
	}

	hash := ""
	if size > 0 {
		hash, err = util.GenerateFileContentHash(tmpFile.Name())
		if err != nil {
			return "", err
		}
		branch := LocalModBranch{
			LocalBaseCommit: localBase,
			LocalModHash:    hash,
			Prefix:          prefix,
		}
		diffFilePath := filepath.Join(dstDir, branch.FileName())
		err := os.Rename(tmpFile.Name(), diffFilePath)
		if err != nil {
			return "", err
		}
		existence, err := git.ValidateRemoteBranchExistence(repo, branch.BranchName())
		if err != nil {
			return "", err
		}
		if existence {
			fmt.Printf("Skipped an existing local mod branch '%s' in %s\n", branch.BranchName(), repo)
		} else {
			err = git.CreateOrphanBranch(dstDir, branch.BranchName())
			if err != nil {
				return "", err
			}
			err = git.CommitFile(dstDir, branch.FileName(), "this is a test")
			if err != nil {
				return "", err
			}
			branches = append(branches, branch)
			fmt.Printf("Created a local mod branch '%s' in %s\n", branch.BranchName(), repo)
		}
	}

	if len(branches) > 0 {
		fmt.Printf("Pushing branches %v\n", branches)
		branchNames := []string{}
		for _, name := range branches {
			branchNames = append(branchNames, name.BranchName())
		}
		err := git.Push(dstDir, branchNames...)
		if err != nil {
			return "", err
		}
	}

	return hash, nil
}
