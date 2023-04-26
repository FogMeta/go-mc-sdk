# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOBIN=$(shell pwd)/build

PROJECT_NAME=client
BINARY_NAME=$(PROJECT_NAME)

PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all ffi build clean help

all: build


ffi:
	git submodule update --init --recursive
	./extern/filecoin-ffi/install-filcrypto
.PHONY: ffi

build: ## Build the binary file
	@go mod download
	@go mod tidy
	@go build -o $(GOBIN)/$(BINARY_NAME)  ./demo/main.go
	@echo "Done building."
	@echo "Go to build folder and run \"$(GOBIN)/$(BINARY_NAME)\" to launch swan client."
.PHONY: build

clean: ## Remove previous build
	@go clean
	@rm -rf $(shell pwd)/build
	@echo "Done cleaning."
.PHONY: clean

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: clean
