RADIGOPKG=$(shell go list ./... | grep -v "/vendor/")

.PHONY: all help build installdeps vet test test-cover

all: help

help:
	@echo "make build         #=> Build binary"
	@echo "make test          #=> Run unit tests"

installdeps:
	glide i

build: installdeps
	go build cmd/radigo/...

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
