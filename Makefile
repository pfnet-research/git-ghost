NAME        := git-ghost
PROJECTROOT := $(shell pwd)
VERSION     := $(shell cat ${PROJECTROOT}/VERSION)
REVISION    := $(shell git rev-parse --short HEAD)
IMAGE_PREFIX:=
IMAGE_TAG   ?= $(VERSION)
OUTDIR      ?= $(PROJECTROOT)/bin

LDFLAGS := -ldflags="-s -w -X \"git-ghost/cmd.Version=$(VERSION)\" -X \"git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: build
build: deps
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o $(OUTDIR)/$(NAME)

.PHONY: build-linux-amd64
build-linux-amd64:
	make build \
		GOOS=linux \
		GOARCH=amd64 \
		NAME=git-ghost-amd64

.PHONY: build-linux-ppc64le
build-linux-ppc64le:
	make build \
		GOOS=linux \
		GOARCH=ppc64le \
		NAME=git-ghost-ppc64le

.PHONY: build-linux-s390x
build-linux-s390x:
	make build \
		GOOS=linux \
		GOARCH=s390x \
		NAME=git-ghost-s390x

.PHONY: build-linux
build-linux: build-linux-amd64 build-linux-ppc64le build-linux-s390x

.PHONY: build-darwin
build-darwin:
	make build \
		GOOS=darwin \
		NAME=git-darwin-amd64

.PHONY: build-windows
build-windows:
	make build \
		GOARCH=amd64 \
		GOOS=windows \
		NAME=git-windows-amd64

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-all-in-docker
docker-build-all: build-image-dev
	docker run --rm -v $(PROJECTROOT)/bin:/tmp/git-ghost/bin $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) make build-all OUTDIR=/tmp/git-ghost/bin

.PHONY: lint
lint:
	gometalinter --config gometalinter.json ./...

.PHONY: deps
deps:
	dep ensure

.PHONY: build-image-dev
build-image-dev:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) --target git-ghost-dev $(PROJECTROOT)

.PHONY: build-image-test
build-image-test:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-test:$(IMAGE_TAG) --target git-ghost-test $(PROJECTROOT)

.PHONY: build-image-e2e
build-image-e2e:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-e2e:$(IMAGE_TAG) --target git-ghost-e2e $(PROJECTROOT)

.PHONY: build-image-cli
build-image-cli:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-cli:$(IMAGE_TAG) --target git-ghost-cli $(PROJECTROOT)

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
