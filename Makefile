NAME        := git-ghost
PROJECTROOT := $(shell pwd)
REVISION    := $(shell git rev-parse --short HEAD)
OUTDIR      ?= $(PROJECTROOT)/dist
RELEASE_TAG ?=
GITHUB_USER := pfnet-research
GITHUB_REPO := git-ghost
GITHUB_REPO_URL := git@github.com:pfnet-research/git-ghost.git
GITHUB_TOKEN ?=

LDFLAGS := -ldflags="-s -w -X \"github.com/pfnet-research/git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: build
build:
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o $(OUTDIR)/$(NAME)

.PHONY: install
install:
	go install -tags netgo -installsuffix netgo $(LDFLAGS)

.PHONY: build-linux-amd64
build-linux-amd64:
	make build \
		GOOS=linux \
		GOARCH=amd64 \
		NAME=git-ghost-linux-amd64

.PHONY: build-linux
build-linux: build-linux-amd64

.PHONY: build-darwin
build-darwin:
	make build \
		GOOS=darwin \
		NAME=git-ghost-darwin-amd64

.PHONY: build-windows
build-windows:
	make build \
		GOARCH=amd64 \
		GOOS=windows \
		NAME=git-ghost-windows-amd64.exe

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: lint
lint:
	golangci-lint run --config golangci.yml

.PHONY: e2e
e2e:
	@go test -v $(PROJECTROOT)/test/e2e/e2e_test.go

.PHONY: update-license
update-license:
	@python3 ./scripts/license/add.py -v

.PHONY: check-license
check-license:
	@python3 ./scripts/license/check.py -v

.PHONY: coverage
coverage:
	@go test -tags no_e2e -covermode=count -coverprofile=profile.cov -coverpkg ./pkg/...,./cmd/... $(shell go list ./... | grep -v /vendor/)
	@go tool cover -func=profile.cov

.PHONY: clean
clean:
	rm -rf $(OUTDIR)/*
