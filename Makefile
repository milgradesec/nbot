
VERSION       :=  $(shell git describe --tags --always --abbrev=0)
SYSTEM        := 
BUILDFLAGS    := -v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH   := github.com/milgradesec/nbot

all: build

build:
	$(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/nbot

clean:
	go clean
	rm nbot.exe

test:
	go test ./...

lint:
	golangci-lint run

release:
	docker buildx build . -f build.Dockerfile \
		--platform linux/arm64 \
		--tag ghcr.io/milgradesec/nbot:$(VERSION) \
		--tag ghcr.io/milgradesec/nbot:latest \
		--push