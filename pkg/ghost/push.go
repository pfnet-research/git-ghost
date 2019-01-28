package ghost

import (
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"
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

	branches := []string{}
	if remoteBase != localBase {
		filename := "commits.patch"
		filepath := filepath.Join(dstDir, filename)
		branchName := fmt.Sprintf("%s/%s-%s", prefix, remoteBase, localBase)
		existence, err := git.ValidateRemoteBranchExistence(repo, branchName)
		if err != nil {
			return "", err
		}
		if existence {
			fmt.Printf("Skipped an existing local base branch '%s' in %s\n", branchName, repo)
		} else {
			err := git.CreateDiffBundleFile(srcDir, filepath, remoteBase, localBase)
			if err != nil {
				return "", err
			}
			err = git.CreateOrphanBranch(dstDir, branchName)
			if err != nil {
				return "", err
			}
			err = git.CommitFile(dstDir, filename, "this is a test")
			if err != nil {
				return "", err
			}
			branches = append(branches, branchName)
			fmt.Printf("Created a local base branch '%s' in %s\n", branchName, repo)
		}
	}

	filename := "local-mod.patch"
	filepath := filepath.Join(dstDir, filename)
	err = git.CreateDiffPatchFile(srcDir, filepath, options.LocalBase)
	if err != nil {
		return "", err
	}
	size, err := util.FileSize(filepath)
	if err != nil {
		return "", err
	}

	hash := ""
	if size > 0 {
		hash, err = util.GenerateFileContentHash(filepath)
		if err != nil {
			return "", err
		}
		branchName := fmt.Sprintf("%s/%s/%s", prefix, localBase, hash)
		existence, err := git.ValidateRemoteBranchExistence(repo, branchName)
		if err != nil {
			return "", err
		}
		if existence {
			fmt.Printf("Skipped an existing local mod branch '%s' in %s\n", branchName, repo)
		} else {
			err = git.CreateOrphanBranch(dstDir, branchName)
			if err != nil {
				return "", err
			}
			err = git.CommitFile(dstDir, filename, "this is a test")
			if err != nil {
				return "", err
			}
			branches = append(branches, branchName)
			fmt.Printf("Created a local mod branch '%s' in %s\n", branchName, repo)
		}
	}

	if len(branches) > 0 {
		fmt.Printf("Pushing branches %v\n", branches)
		err := git.Push(dstDir, branches...)
		if err != nil {
			return "", err
		}
	}

	return hash, nil
}
