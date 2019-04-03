####################################################################################################
# builder image
####################################################################################################
FROM golang:1.11.4 as builder

RUN apt-get update && apt-get install -y \
    git \
    make \
    wget \
    gcc \
    zip \
    bzip2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp

# Install dep
ENV DEP_VERSION=0.5.0
RUN wget https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

# Install gometalinter
ENV GOMETALINTER_VERSION=2.0.12
RUN curl -sLo- https://github.com/alecthomas/gometalinter/releases/download/v${GOMETALINTER_VERSION}/gometalinter-${GOMETALINTER_VERSION}-linux-amd64.tar.gz | \
    tar -xzC "$GOPATH/bin" --exclude COPYING --exclude README.md --strip-components 1 -f- && \
    ln -s $GOPATH/bin/gometalinter $GOPATH/bin/gometalinter.v2

# Install github-release
ENV GITHUB_RELEASE_VERSION=0.7.2
RUN curl -sLo- https://github.com/aktau/github-release/releases/download/v${GITHUB_RELEASE_VERSION}/linux-amd64-github-release.tar.bz2 | \
    tar -xjC "$GOPATH/bin" --strip-components 3 -f-


####################################################################################################
# git-ghost-dev
####################################################################################################
FROM builder as git-ghost-dev

# A dummy directory is created under $GOPATH/src/dummy so we are able to use dep
# to install all the packages of our dep lock file
COPY Gopkg.toml ${GOPATH}/src/dummy/Gopkg.toml
COPY Gopkg.lock ${GOPATH}/src/dummy/Gopkg.lock

RUN cd ${GOPATH}/src/dummy && \
    dep ensure -vendor-only && \
    mv vendor/* ${GOPATH}/src/ && \
    rmdir vendor

# Perform the build
WORKDIR /go/src/git-ghost
COPY . .
RUN make build

####################################################################################################
# git-ghost-test
####################################################################################################
FROM ubuntu:16.04 as git-ghost-test

RUN apt-get update && apt-get install -y \
    vim \
    git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=git-ghost-dev /go/src/git-ghost/dist/git-ghost /usr/local/bin/

COPY scripts/create-test-repo.sh /work/create-test-repo.sh
RUN mkdir -p /work/local /work/remote /work/ghost-repo
RUN /work/create-test-repo.sh /work/local /work/remote /work/ghost-repo
ENV GIT_GHOST_REPO=/work/ghost-repo

WORKDIR /work/local

####################################################################################################
# git-ghost-e2e
####################################################################################################
FROM git-ghost-dev as git-ghost-e2e

COPY --from=git-ghost-dev /go/src/git-ghost/dist/git-ghost /usr/local/bin/
WORKDIR /go/src/git-ghost
RUN git config --global user.email you@example.com \
    && git config --global user.name "Your Name"

####################################################################################################
# git-ghost-cli
####################################################################################################
FROM ubuntu:16.04 as git-ghost-cli

RUN apt-get update && apt-get install -y \
    git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=git-ghost-dev /go/src/git-ghost/dist/git-ghost /usr/local/bin/
