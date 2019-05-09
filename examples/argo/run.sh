#!/bin/bash

set -eu
set -o pipefail

DIR=$(dirname "$0")

[[ -z "$GIT_GHOST_REPO" ]] && (echo "GIT_GHOST_REPO must be specified."; exit 1)
[[ -z "$GIT_GHOST_REGISTRY" ]] && (echo "GIT_GHOST_REGISTRY must be specified."; exit 1)

kubectl get secret git-ghost-git-cred > /dev/null || (cat << EOS
git-ghost-git-cred secret is needed.
You can create it by:

$ kubectl create secret generic git-ghost-git-cred --from-file=sshkey=\$HOME/.ssh/id_rsa

EOS
exit 1)

kubectl get secret git-ghost-docker-cred > /dev/null || (cat << EOS
git-ghost-docker-cred secret is needed.
You can create it by:

$ kubectl create secret docker-registry git-ghost-docker-cred --docker-username=YOUR_NAME --docker-password=YOUR_PASSWORD

EOS
exit 1)

IMAGE_PREFIX=
if [[ -z "$IMAGE_PREFIX" ]]; then
  IMAGE_PREFIX=dtaniwaki/
fi

IMAGE_TAG=
if [[ -z "$IMAGE_TAG" ]]; then
  IMAGE_TAG="v0.1.1"
fi

GIT_REPO=`git remote get-url origin`
git diff --exit-code HEAD > /dev/null && (echo "Make some local modification to try the example."; exit 1)
HASHES=`git ghost push origin/master`
GIT_COMMIT_HASH=`echo $HASHES | cut -f 1 -d " "`
DIFF_HASH=`echo $HASHES | cut -f 2 -d " "`

argo submit --watch -p image-prefix=$IMAGE_PREFIX -p image-tag=$IMAGE_TAG -p git-repo=$GIT_REPO -p git-ghost-repo=$GIT_GHOST_REPO -p git-ghost-registry=$GIT_GHOST_REGISTRY -p git-commit-hash=$GIT_COMMIT_HASH -p diff-hash=$DIFF_HASH $DIR/workflow.yaml --loglevel debug
