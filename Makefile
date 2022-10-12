NAME        := git-ghost
PROJECTROOT := $(shell pwd)
VERSION     := $(if $(VERSION),$(VERSION),$(shell cat ${PROJECTROOT}/VERSION)-dev)
REVISION    := $(shell git rev-parse --short HEAD)
IMAGE_PREFIX ?=
IMAGE_TAG   ?= $(VERSION)
OUTDIR      ?= $(PROJECTROOT)/dist
RELEASE_TAG ?=
GITHUB_USER := pfnet-research
GITHUB_REPO := git-ghost
GITHUB_REPO_URL := git@github.com:pfnet-research/git-ghost.git
GITHUB_TOKEN ?=

LDFLAGS := -ldflags="-s -w -X \"github.com/pfnet-research/git-ghost/cmd.Version=$(VERSION)\" -X \"github.com/pfnet-research/git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

guard-%:
	@ if [ "${${*}}" = "" ]; then \
    echo "Environment variable $* is not set"; \
		exit 1; \
	fi

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

.PHONY: release
release: release-code release-assets release-image

.PHONY: release-code
release-code: guard-RELEASE_TAG guard-RELEASE_COMMIT guard-GITHUB_USER guard-GITHUB_REPO guard-GITHUB_REPO_URL guard-GITHUB_TOKEN
	@GITHUB_TOKEN=$(GITHUB_TOKEN)
	git tag $(RELEASE_TAG) $(RELEASE_COMMIT)
	git push $(GITHUB_REPO_URL) $(RELEASE_TAG)
	github-release release \
	  --user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(RELEASE_TAG)

.PHONY: release-assets
release-assets: guard-RELEASE_TAG guard-RELEASE_COMMIT guard-GITHUB_USER guard-GITHUB_REPO guard-GITHUB_REPO_URL guard-GITHUB_TOKEN
	@GITHUB_TOKEN=$(GITHUB_TOKEN)
	git diff --quiet HEAD || (echo "your current branch is dirty" && exit 1)
	git checkout $(RELEASE_COMMIT)
	make clean build-all VERSION=$(shell cat ${PROJECTROOT}/VERSION)
	for target in linux-amd64 darwin-amd64 windows-amd64.exe; do \
		github-release upload \
		  --user $(GITHUB_USER) \
			--repo $(GITHUB_REPO) \
			--tag $(RELEASE_TAG) \
			--name git-ghost-$$target \
			--file $(OUTDIR)/git-ghost-$$target; \
	done
	git checkout -

.PHONY: release-image
release-image: guard-RELEASE_TAG
	git diff --quiet HEAD || (echo "your current branch is dirty" && exit 1)
	git checkout $(RELEASE_COMMIT)
	make build-image-cli VERSION=$(shell cat ${PROJECTROOT}/VERSION)
	docker push $(IMAGE_PREFIX)git-ghost-cli:$(RELEASE_TAG)
	git checkout -

.PHONY: lint
lint:
	golangci-lint run --config golangci.yml

.PHONY: build-image-dev
build-image-dev:
	docker build -t $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) --target git-ghost-dev $(PROJECTROOT)

.PHONY: build-image-cli
build-image-cli:
	docker build -t $(IMAGE_PREFIX)git-ghost-cli:$(IMAGE_TAG) --build-arg VERSION=$(VERSION) --target git-ghost-cli $(PROJECTROOT)

.PHONY: build-image-all
build-image-all: build-image-dev build-image-cli

test:
	@go test -v -race -short -tags no_e2e ./...

.PHONY: shell
shell: build-image-cli
	docker run -it $(IMAGE_PREFIX)git-ghost-cli:$(IMAGE_TAG) bash

.PHONY: dev-shell
dev-shell: build-image-dev
	docker run -it $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) bash

.PHONY: e2e
e2e:
	@go test -v $(PROJECTROOT)/test/e2e/e2e_test.go

.PHONY: docker-e2e
docker-e2e: build-image-dev
	@docker run $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) make install e2e DEBUG=$(DEBUG)

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
