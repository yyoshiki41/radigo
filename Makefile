.PHONY: all help installdeps build build-4-docker vet test test-cover docker-build

RADIGOPKG=$(shell go list ./...)

REVISION := $(shell git rev-parse --short HEAD)

all: help

help:
	@echo "Useful targets:"
	@echo "  make installdeps   => Install dependencies"
	@echo "  make build         => Build a binary"
	@echo "  make test          => Run unit tests"
	@echo "  make vet           => Run go vet"
	@echo "  make docker-build  => Build a docker image"

installdeps:
	go mod download

build:
	go build ./cmd/radigo/...

build-4-docker:
	CGO_ENABLED=0 GOOS=linux go build -o /bin/radigo cmd/radigo/main.go

vet:
	@go vet -v $(RADIGOPKG)

test: vet
	@go test $(RADIGOPKG)

test-cover:
	@echo "" > coverage.txt; \
	for d in $(RADIGOPKG); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done

docker-build:
	docker build -t yyoshiki41/radigo .
