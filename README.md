# git-ghost

Git Ghost is a command line tool for synchronizing your working directory efficiently to a remote place without commiting changes.

## Concept

Git Ghost creates 2 types of branches to synchronize your working directory.

### Commits Branch

This type of branch contains commits between 2 commits which can exist only in your working directory.

### Diff Branch

This type of branch contains modifications from a specific commit in your working directory.

## Installing

Install the binary from source: execute,

```bash
$ git clone <REPO_URL> git-ghost
$ cd git-ghost
$ make install
```

The binaries of each versions are available in [Releases](/releases).

## Getting Started

First, create an empty repository which can be accessible from a remote place. Set the URL as `GIT_GHOST_REPO` env.

Assume your have a local working directory `DIR_L` and a remote directory to be synchronized `DIR_R`.

## Case 1 (`DIR_L` HEAD == `DIR_R` HEAD)

You can synchoronize local modifications.

```bash
$ cd <DIR_L>
$ git-ghost push
<HASH>
$ git-ghost show <HASH>
...
$ cd <DIR_R>
$ git-ghost pull <HASH>
```

## Case 2 (`DIR_L` HEAD \> `DIR_R` HEAD)

You can synchronize local commits and modifications.

Assume `DIR_R`'s HEAD is HASH\_R.

```bash
$ cd <DIR_L>
$ git-ghost push all <HASH_R>
<HASH_1> <HASH_2>
<HASH_3>
$ git-ghost show all <HASH_2> <HASH_3>
...
$ cd <DIR_R>
$ git-ghost pull all <HASH_2> <HASH_3>
```

## Development

```
# checkout this repo to $GOPATH/src/git-ghost
$ cd $GOPATH/src
$ git clone git@github.pfidev.jp:Cluster/git-ghost.git
$ cd git-ghost

# build
$ make build

# see godoc
$ go get golang.org/x/tools/cmd/godoc
$ godoc -http=:6060  # access http://localhost:6060 in browser
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new [Pull Request](../../pull/new/master)

## Copyright

Copyright (c) 2019 Preferred Networks. See [LICENSE](LICENSE) for details.
