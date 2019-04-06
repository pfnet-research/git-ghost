#!/bin/bash

set -eu
set -o pipefail

DIR=$(dirname "$0")

[[ -z "$GIT_GHOST_REPO" ]] && (echo "GIT_GHOST_REPO must be specified."; exit 1)

kubectl get secret git-ghost-creds > /dev/null || (cat << EOS
git-ghost-creds secret is needed.
You can create it by:

$ kubectl create secret generic git-ghost-creds --from-file=sshkey=\$HOME/.ssh/id_rsa

EOS
exit 1)

IMAGE_PREFIX=
if [[ -z "$IMAGE_PREFIX" ]]; then
  IMAGE_PREFIX=dtaniwaki
fi

IMAGE_TAG=
if [[ -z "$IMAGE_TAG" ]]; then
  IMAGE_TAG="v0.1.0"
fi

GIT_REPO=`git remote get-url origin`
GIT_COMMIT_HASH=`git rev-parse HEAD`
git diff --exit-code HEAD > /dev/null && (echo "Make some local modification to try the example."; exit 1)
DIFF_HASH=`git ghost push $GIT_COMMIT_HASH | tail -n 1 | cut -f 2 -d " "`

cat $DIR/pod.yaml | \
sed "s|{{IMAGE_PREFIX}}|$IMAGE_PREFIX|g" | \
sed "s|{{IMAGE_TAG}}|$IMAGE_TAG|g" | \
sed "s|{{GIT_REPO}}|$GIT_REPO|g" | \
sed "s|{{GIT_GHOST_REPO}}|$GIT_GHOST_REPO|g" | \
sed "s|{{GIT_COMMIT_HASH}}|$GIT_COMMIT_HASH|g" | \
sed "s|{{DIFF_HASH}}|$DIFF_HASH|g" | kubectl create -f -

POD_NAME=`kubectl get pod -l app=git-ghost-example --sort-by=.metadata.creationTimestamp --no-headers -o custom-columns=:metadata.name | tail -n 1`

sleep 10

echo "Print modifications synchronized from the local."
echo

kubectl logs -f $POD_NAME
