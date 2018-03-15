RADIGOPKG=$(shell go list ./... | grep -v "/vendor/")

.PHONY: all help build test

all: help

help:
	@echo "make build         #=> Build binary"
	@echo "make test          #=> Run unit tests"

build:
	go build ./cmd/radigo/...

test:
	go test $(RADIGOPKG)

test-cover:
	@echo "" > coverage.txt; \
	for d in $(RADIGOPKG); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done
