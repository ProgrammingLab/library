NAME := library.out
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
		   -X 'main.revision=$(REVISION)'
SRCS := $(wildcard *.go)

# 必要なツール類をセットアップする
## Setup
.PHONY: setup
setup:
	go get github.com/golang/dep/...
	go get github.com/Songmu/make2help/cmd/make2help

# depを使って依存パッケージをインストールする
## Install dependencies
.PHONY: deps
deps: setup
	dep ensure

## build binaries ex. make bin/library
.PHONY: build
build: deps
	go build -ldflags "$(LDFLAGS)" -o "$(NAME)" $(SRCS)

## show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps help
