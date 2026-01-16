BINARY_NAME=skills
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w \
    -X github.com/takoeight0821/skills/jig/pkg/version.Version=$(VERSION) \
    -X github.com/takoeight0821/skills/jig/pkg/version.Commit=$(COMMIT) \
    -X github.com/takoeight0821/skills/jig/pkg/version.Date=$(DATE)"

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/skills

.PHONY: install
install:
	go install $(LDFLAGS) ./cmd/skills

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf bin/ dist/

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --clean

.PHONY: release
release:
	goreleaser release --clean

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.DEFAULT_GOAL := build
