
VERSION       :=  $(shell git describe --tags --always --abbrev=0)
SYSTEM        := 
BUILDFLAGS    := -v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH   := github.com/milgradesec/nbot

.PHONY: all
all: build

.PHONY: build
build:
	$(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/nbot

.PHONY: clean
clean:
	go clean
	rm nbot.exe

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run