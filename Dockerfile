####################################################################################################
# builder image
####################################################################################################
FROM golang:1.11.4 as builder

RUN apt-get update && apt-get install -y \
    git \
    make \
    wget \
    gcc \
    zip && \
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


####################################################################################################
# git-ghost-dev
####################################################################################################
FROM builder as git-ghost-build

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

COPY --from=git-ghost-build /go/src/git-ghost/bin/git-ghost /usr/local/bin/

COPY hack/create-test-repo.sh /work/create-test-repo.sh
RUN mkdir -p /work/local /work/remote /work/ghost-repo
RUN /work/create-test-repo.sh /work/local /work/remote /work/ghost-repo
ENV GHOST_REPO=/work/ghost-repo

WORKDIR /work/local

####################################################################################################
# git-ghost
####################################################################################################
FROM ubuntu:16.04 as git-ghost

RUN apt-get update && apt-get install -y \
    git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=git-ghost-build /go/src/git-ghost/bin/git-ghost /usr/local/bin/
