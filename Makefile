RADIGODIR=/tmp/radigo
RADIGOPKG=$(shell go list ./... | grep -v "/vendor/")

.PHONY: all help init clean test

all: help

help:
	@echo "make init          #=> Run init scripts"
	@echo "make clean         #=> Remove cache and downloaded files"
	@echo "make test          #=> Run unit tests"

init:
	mkdir -p $(RADIGODIR) && mkdir -p $(RADIGODIR)/output && mkdir -p $(RADIGODIR)/.cache

clean: rm-cache
	cd $(RADIGODIR) && ls | grep -v 'output' | xargs rm -r

rm-cache:
	rm -rf $(RADIGODIR)/.cache/*


test:
	go test $(RADIGOPKG)

test-cover:
	@echo "" > coverage.txt; \
	for d in $(RADIGOPKG); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done

echo:
	@echo "$(HOGE)"
