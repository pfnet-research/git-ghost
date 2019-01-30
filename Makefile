NAME        := git-ghost
PROJECTROOT := $(shell pwd)
VERSION     := $(shell cat ${PROJECTROOT}/VERSION)
REVISION    := $(shell git rev-parse --short HEAD)
IMAGE_PREFIX:=
IMAGE_TAG   ?= $(VERSION)

LDFLAGS := -ldflags="-s -w -X \"git-ghost/cmd.Version=$(VERSION)\" -X \"git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: lint
lint:
	gometalinter --config gometalinter.json ./...

.PHONY: deps
deps:
	dep ensure

.PHONY: build
build: deps
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) $(PROJECTROOT)

.PHONY: docker-build-dev
docker-build-dev:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-dev:$(IMAGE_TAG) --target git-ghost-dev $(PROJECTROOT)

.PHONY: docker-build-test
docker-build-test:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-test:$(IMAGE_TAG) --target git-ghost-test $(PROJECTROOT)

.PHONY: docker-build
docker-build:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost:$(IMAGE_TAG) --target git-ghost $(PROJECTROOT)

test: deps
	@go test -v $(PROJECTROOT)/...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*
