####################################################################################################
# git-ghost-dev
####################################################################################################
FROM golang:1.11.4 as git-ghost-dev

RUN apt-get update -q && apt-get install -yq --no-install-recommends \
    git \
    make \
    wget \
    gcc \
    zip \
    bzip2 \
    lsb-release \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp

# Install docker client
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
RUN apt-get update -q && apt-get install -yq --no-install-recommends docker-ce-cli && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install dep
ENV DEP_VERSION=0.5.0
RUN wget https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

# Install golangci-lint
ENV GOLANGCI_LINT_VERSION=1.16.0
RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/v$GOLANGCI_LINT_VERSION/install.sh| sh -s -- -b $(go env GOPATH)/bin v$GOLANGCI_LINT_VERSION

# Install github-release
ENV GITHUB_RELEASE_VERSION=0.7.2
RUN curl -sLo- https://github.com/aktau/github-release/releases/download/v${GITHUB_RELEASE_VERSION}/linux-amd64-github-release.tar.bz2 | \
    tar -xjC "$GOPATH/bin" --strip-components 3 -f-

# A dummy directory is created under $GOPATH/src/dummy so we are able to use dep
# to install all the packages of our dep lock file
COPY Gopkg.toml ${GOPATH}/src/dummy/Gopkg.toml
COPY Gopkg.lock ${GOPATH}/src/dummy/Gopkg.lock

RUN cd ${GOPATH}/src/dummy && \
    dep ensure -vendor-only && \
    mv vendor/* ${GOPATH}/src/ && \
    rmdir vendor

WORKDIR $GOPATH/src/git-ghost
COPY . .

####################################################################################################
# builder
####################################################################################################
FROM git-ghost-dev as builder

# Perform the build
RUN make build

####################################################################################################
# git-ghost-cli
####################################################################################################
FROM ubuntu:16.04 as git-ghost-cli

COPY --from=builder /go/src/git-ghost/dist/git-ghost /usr/local/bin/

RUN apt-get update -q && apt-get install -yq --no-install-recommends git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

