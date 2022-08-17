
VERSION     :=  $(shell git describe --tags --always --dirty='-dev')
SYSTEM      := 
BUILDFLAGS  := -v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH := github.com/milgradesec/nbot

all: build

build:
	$(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/nbot

clean:
	go clean

test:
	go test ./...

lint:
	golangci-lint run

docker:
	docker build . -f build.Dockerfile

release:
	docker --log-level=debug buildx build . \
		-f build.Dockerfile \
		--platform linux/amd64,linux/arm64 \
		--tag ghcr.io/milgradesec/nbot:$(VERSION) \
		--tag ghcr.io/milgradesec/nbot:latest \
		--push