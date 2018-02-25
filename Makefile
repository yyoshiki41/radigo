RADIGODIR=/tmp/radigo
RADIGOPKG=$(shell go list ./... | grep -v "/vendor/")

.PHONY: all help init clean test

all: help

help:
	@echo "make init          #=> Run init scripts"
	@echo "make clean         #=> Remove downloaded files"
	@echo "make test          #=> Run unit tests"

init:
	mkdir -p $(RADIGODIR)/output

clean:
	rm -rf $(RADIGODIR)/output/*

test:
	go test $(RADIGOPKG)

test-cover:
	@echo "" > coverage.txt; \
	for d in $(RADIGOPKG); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done
