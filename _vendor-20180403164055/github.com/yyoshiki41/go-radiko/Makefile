PKGS=$(shell go list ./... | grep -v examples)
BASE_FOLDERS=$(shell ls -d */ | grep -v vendor | grep -v testdata)

.PHONY: all help init test test-out

all: help

help:
	@echo "make init          #=> Run init scripts"
	@echo "make get-deps      #=> Install dependencies"
	@echo "make verify        #=> Verify tests"
	@echo "make lint          #=> Run golint"
	@echo "make vet           #=> Run go vet"
	@echo "make test          #=> Run tests"
	@echo "make test-out      #=> Run tests from outside Japan"

init: get-deps

test:
	go test $(PKGS)

test-out:
	GO_RADIKO_OUTSIDE_JP=true go test $(PKGS)

test-ci:
	echo "GO_RADIKO_OUTSIDE_JP=true go test"
	echo "" > coverage.txt
	$(foreach pkg, $(PKGS), \
		GO_RADIKO_OUTSIDE_JP=true \
		go test -race -coverprofile=profile.txt -covermode=atomic $(pkg); \
		[ ! -f profile.txt ] && continue; \
		cat profile.txt >> coverage.txt; \
		rm profile.txt;)

verify: lint vet

lint:
	golint ./...

vet:
	go tool vet -all -structtags -shadow $(BASE_FOLDERS)

get-deps:
	@echo "go get go-radiko dependencies"
	@go get -v $(PKGS)
