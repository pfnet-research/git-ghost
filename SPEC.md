# Spec
 ## Ghost Repo
 The ghost creates a patch for your uncommitted modifications with the following spec into a ghost repo.
 __Word Definitions in this section__
 | name | description | example |
|--------|--------|--------|
| `GHOST_BRANCH_PREFIX` | A prefix to identify branches for ghost. | `git-ghost` |
| `REMOTE_BASE_COMMIT` | A full base commit hash of a source repo. It is supposed to exist in a remote source repo when the ghost apply patches. | `f2e8fbf0f2c1527ad208faa5d08d4b377ce962a3` |
| `LOCAL_BASE_COMMIT` | A full base commit hash of a source repo to make a local modifications diff. | `1fd368bbe90a2e9079d86fa63398a7cd06a79577` |
| `LOCAL_MOD_HASH` | A content hash of local modifications diff. | `35d5d6474c92a780c6be4679cee3a72cd1bdfe99` |
 ### Branch Structure
 The ghost creates 2 kinds of branches.
 #### Local Base Branch
 __Format__: `$GHOST_BRANCH_PREFIX/$REMOTE_BASE_COMMIT..$LOCAL_BASE_COMMIT`
 __Directory Structure__
 ```
/
└─ commits.bundle
```
 `commits.patch` is a diff bundle from a remote base commit to a local base commit. It is not created if a remote base commit equals to a local base commit.
 The file is created by the following command.
 ```
$ git format-patch --binary --stdout $REMOTE_BASE_COMMIT..$LOCAL_BASE_COMMIT > commits.patch
```
 And it can be applied by the following command.
 ```
$ git pull --ff-only --no-tags commits.patch $GHOST_BRANCH_PREFIX/$REMOTE_BASE_COMMIT..$LOCAL_BASE_COMMIT
```
 #### Local Mod Branch
 __Format__: `$GHOST_BRANCH_PREFIX/$LOCAL_BASE_COMMIT/$LOCAL_MOD_HASH`
 __Directory Structure__
 ```
/
└─ patch.diff
```
 `patch.diff` is a patch file of uncommitted local modifications from a local base commit.
 The file is created by the following command.
 ```
$ git diff --binary $LOCAL_BASE_COMMIT > local-mod.patch
```
 And it can be applied by the following command.
 ```
$ git apply local-mod.patch
```
