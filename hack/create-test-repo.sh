#!/bin/bash

if [ "$#" -ne 3 ]; then
  echo "Illegal number of arguments"
  exit
fi

set -eux

LOCAL=$(realpath $1)
REMOTE=$(realpath $2)
GHOST_REPO=$(realpath $3)

(
cd $REMOTE
git init
git config --global user.email "you@example.com"
git config --global user.name "Your Name"
echo a > sample.txt
git add sample.txt
git commit sample.txt -m 'Initial commit'
echo b > sample.txt
git commit sample.txt -m 'Second commit'
echo c > sample.txt
)

(
git clone $REMOTE $LOCAL
)

(
cd $GHOST_REPO
git init
)
