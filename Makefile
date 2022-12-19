
VERSION     :=  $(shell git describe --tags --always --abbrev=0)
BUILDFLAGS  := -v -trimpath -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH := github.com/milgradesec/nbot

GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

.PHONY: build
build:
	go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/nbot

.PHONY: lint
lint: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: docker
docker:
	docker build . -f build.Dockerfile

.PHONY: release
release:
	docker --log-level=debug buildx build . \
		-f build.Dockerfile \
		--platform linux/amd64,linux/arm64 \
		--tag ghcr.io/milgradesec/nbot:$(VERSION) \
		--tag ghcr.io/milgradesec/nbot:latest \
		--push
