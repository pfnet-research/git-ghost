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

.PHONY: build-image-test
build-image-test:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-test:$(IMAGE_TAG) --target git-ghost-test $(PROJECTROOT)

.PHONY: build-image-e2e
build-image-e2e:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-e2e:$(IMAGE_TAG) --target git-ghost-e2e $(PROJECTROOT)

.PHONY: build-image-cli
build-image-cli:
	docker build --build-arg NAME=$(NAME) --build-arg VERSION=$(VERSION) --build-arg REVISION=$(REVISION) -t $(IMAGE_PREFIX)git-ghost-cli:$(IMAGE_TAG) --target git-ghost-cli $(PROJECTROOT)

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
