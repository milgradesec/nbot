
VERSION:=$(shell git describe --tags --always --abbrev=0 --dirty="-dev")
SYSTEM:=
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/nbot
CGO_ENABLED:=0

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/nbot

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run