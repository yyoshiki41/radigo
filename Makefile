.PHONY: all help installdeps build build-4-docker vet test test-cover docker-build

RADIGOPKG=$(shell go list ./...)

all: help

help:
	@echo "make build         #=> Build binary"
	@echo "make test          #=> Run unit tests"

installdeps:
	dep ensure

build: installdeps
	go build ./cmd/radigo/...

build-4-docker: installdeps
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
