#!/bin/bash

# Copyright 2019 Preferred Networks, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
  IMAGE_PREFIX=ghcr.io/pfnet-research/
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
