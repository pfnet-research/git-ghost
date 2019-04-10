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
