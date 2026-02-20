# バージョン情報の動的取得
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# ldflags でバージョン情報を埋め込み
LDFLAGS := -ldflags "-s -w \
    -X github.com/tacokumo/tacoctl/internal/version.Version=$(VERSION) \
    -X github.com/tacokumo/tacoctl/internal/version.GitCommit=$(GIT_COMMIT) \
    -X github.com/tacokumo/tacoctl/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: all
all: format test build lint

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/tacoctl .

.PHONY: install
install:
	go install $(LDFLAGS) .

.PHONY: lint
lint:
	golangci-lint run