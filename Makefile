NAME := library
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
		   -X 'main.revision=$(REVISION)'

# 必要なツール類をセットアップする
## Setup
setup:
	go get github.com/golang/dep/...
	go get github.com/Songmu/make2help/cmd/make2help

# depを使って依存パッケージをインストールする
## Install dependencies
deps: setup
	dep ensure

## build binaries ex. make bin/library
build: main.go deps
	go build -ldflags "$(LDFLAGS)" -o "$(NAME)" $<

## show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps help
