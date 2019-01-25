NAME        := git-ghost
VERSION     := v0.0.1
REVISION    := $(shell git rev-parse --short HEAD)
PROJECTROOT := "./"

LDFLAGS := -ldflags="-s -w -X \"git-ghost/cmd.Version=$(VERSION)\" -X \"git-ghost/cmd.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: deps
deps:
	dep ensure

.PHONY: build
build: deps
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) $(PROJECTROOT)

test: deps
	@go test -v $(PROJECTROOT)/...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*
