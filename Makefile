NAME        := git-ghost
PROJECTROOT := $(shell pwd)
VERSION     := $(shell cat ${PROJECTROOT}/VERSION)
REVISION    := $(shell git rev-parse --short HEAD)
IMAGE_PREFIX ?=
IMAGE_TAG   ?= $(VERSION)
OUTDIR      ?= $(PROJECTROOT)/bin
RELEASE_TAG ?=
GITHUB_API  ?=
GITHUB_USER ?=
GITHUB_REPO ?=
GITHUB_TOKEN ?=
DOCKER_GITHUB_ENV_FLAGS := -e GITHUB_API=$(GITHUB_API) -e GITHUB_USER=$(GITHUB_USER) -e GITHUB_REPO=$(GITHUB_REPO) -e GITHUB_TOKEN=$(GITHUB_TOKEN)

LDFLAGS := -ldflags="-s -w -X \"git-ghost/cmd.Version=$(VERSION)\" -X \"git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

guard-%:
	@ if [ "${${*}}" = "" ]; then \
    echo "Environment variable $* is not set"; \
		exit 1; \
	fi

.PHONY: build
build: deps
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o $(OUTDIR)/$(NAME)

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
		NAME=git-ghost-windows-amd64

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-all-in-docker
docker-build-all: build-image-dev
	docker run --rm -v $(OUTDIR):/tmp/git-ghost/bin $(IMAGE_PREFIX)git-ghost-dev:$(RELEASE_TAG) make build-all OUTDIR=/tmp/git-ghost/bin

.PHONY: release
release: guard-RELEASE_TAG
	make build-image-dev IMAGE_TAG=$(RELEASE_TAG)
	docker run --rm $(DOCKER_GITHUB_ENV_FLAGS) $(IMAGE_PREFIX)/git-ghost-dev:$(RELEASE_TAG) github-release release --tag $(RELEASE_TAG)
	make release-assets
	make release-image IMAGE_TAG=$(RELEASE_TAG)

.PHONY: release-assets
release-assets: guard-RELEASE_TAG
	make build-image-dev IMAGE_TAG=$(RELEASE_TAG)
	docker run --rm $(DOCKER_GITHUB_ENV_FLAGS) $(IMAGE_PREFIX)/git-ghost-dev:$(RELEASE_TAG) /bin/bash -c "\
	  set -eux; \
		make build-all OUTDIR=/tmp/git-ghost/dist; \
		for target in linux-amd64 darwin-amd64 windows-amd64; do \
			github-release upload \
				--tag $(RELEASE_TAG) \
				--name git-ghost-\$$target \
				--file /tmp/git-ghost/dist/git-ghost-\$$target; \
		done"

.PHONY: release-image
release-image: guard-RELEASE_TAG
	make build-image-cli IMAGE_TAG=$(RELEASE_TAG)
	docker push $(IMAGE_PREFIX)/git-ghost-cli:$(RELEASE_TAG)

.PHONY: lint
lint: deps
	gometalinter --config gometalinter.json ./...

.PHONY: docker-lint
docker-lint: build-image-dev
	@docker run $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) make lint

.PHONY: deps
deps:
	dep ensure

.PHONY: build-image-dev
build-image-dev:
	docker build -t $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) --target git-ghost-dev $(PROJECTROOT)

.PHONY: build-image-test
build-image-test:
	docker build -t $(IMAGE_PREFIX)git-ghost-test:$(IMAGE_TAG) --target git-ghost-test $(PROJECTROOT)

.PHONY: build-image-e2e
build-image-e2e:
	docker build -t $(IMAGE_PREFIX)git-ghost-e2e:$(IMAGE_TAG) --target git-ghost-e2e $(PROJECTROOT)

.PHONY: build-image-cli
build-image-cli:
	docker build -t $(IMAGE_PREFIX)git-ghost-cli:$(IMAGE_TAG) --target git-ghost-cli $(PROJECTROOT)

.PHONY: build-image-all
build-image-all: build-image-test build-image-e2e build-image-cli

test: deps
	@go test -v $(PROJECTROOT)/...

.PHONY: e2e
e2e:
	@go test -v $(PROJECTROOT)/test/e2e/e2e_test.go

.PHONY: docker-e2e
docker-e2e: build-image-e2e
	@docker run $(IMAGE_PREFIX)git-ghost-e2e:$(IMAGE_TAG) make e2e

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*
