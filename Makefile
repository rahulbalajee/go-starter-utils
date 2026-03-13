MODULE   := github.com/rahulbalajee/go-starter-utils
BINARY   := go-starter-utils
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE     ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  := -s -w \
	-X '$(MODULE)/cmd.Version=$(VERSION)' \
	-X '$(MODULE)/cmd.Commit=$(COMMIT)' \
	-X '$(MODULE)/cmd.Date=$(DATE)'

.PHONY: build run test lint fmt vet clean install

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -f $(BINARY) coverage.out coverage.html

install:
	go install -ldflags "$(LDFLAGS)" .

coverage: test
	go tool cover -html=coverage.out -o coverage.html
