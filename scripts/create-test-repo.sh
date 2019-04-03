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
