RADIGODIR=/tmp/radigo

.PHONY: all help init

all: help

help:
	@echo "make init          #=> Run init scripts"

init:
	mkdir $(RADIGODIR) && mkdir $(RADIGODIR)/.cache

cleanup:
	rm -rf $(RADIGODIR)/* && rm -rf $(RADIGODIR)/.cache/*
