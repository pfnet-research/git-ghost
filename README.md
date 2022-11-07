# git-ghost

[![GoDoc][godoc-image]][godoc-link]
[![Build Status][build-image]][build-link]
[![Coverage Status][cov-image]][cov-link]

Git Ghost is a command line tool for synchronizing your working directory efficiently to a remote place without commiting changes.

## Concept

Git Ghost creates 2 types of branches to synchronize your working directory.

### Commits Branch

This type of branch contains commits between 2 commits which can exist only in your working directory.

### Diff Branch

This type of branch contains modifications from a specific commit in your working directory.

## Installing

### From Source

Install the binary from source: execute,

```bash
$ git clone https://github.com/pfnet-research/git-ghost
$ cd git-ghost
$ make install
```

Compiled binary is located in `dist` folder.

### Releases

The binaries of each releases are available in [Releases](../../releases).

## Getting Started

First, create an empty repository which can be accessible from a remote place. Set the URL as `GIT_GHOST_REPO` env.

Assume you have a local working directory `DIR_L` and a remote directory to be synchronized `DIR_R`.

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
$ git clone git@github.com:pfnet-research/git-ghost.git
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

## Release

Release flow is fully automated. All the maintainer need to do is just to approve and merge [the next release PR](https://github.com/pfnet-research/git-ghost/pulls?q=is%3Apr+is%3Aopen+label%3Atagpr+)

## Copyright

Copyright (c) 2019 Preferred Networks. See [LICENSE](LICENSE) for details.

[build-image]: https://github.com/pfnet-research/git-ghost/actions/workflows/ci.yml/badge.svg
[build-link]:  https://github.com/pfnet-research/git-ghost/actions/workflows/ci.yml
[cov-image]:   https://coveralls.io/repos/github/pfnet-research/git-ghost/badge.svg?branch=master
[cov-link]:    https://coveralls.io/github/pfnet-research/git-ghost?branch=master
[godoc-image]: https://godoc.org/github.com/pfnet-research/git-ghost?status.svg
[godoc-link]:  https://godoc.org/github.com/pfnet-research/git-ghost
