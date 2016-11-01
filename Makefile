RADIGODIR=/tmp/radigo
RADIGOPKG=$(shell go list ./... | grep -v "/vendor/")

.PHONY: all help init

all: help

help:
	@echo "make init          #=> Run init scripts"
	@echo "make cleanup       #=> Remove cache and downloaded files"

init:
	mkdir $(RADIGODIR) && mkdir $(RADIGODIR)/.cache

cleanup:
	rm -rf $(RADIGODIR)/* && rm -rf $(RADIGODIR)/.cache/*

test:
	go test $(RADIGOPKG)
