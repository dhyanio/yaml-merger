# Variables
BINARY_NAME := yaml-merger
BINARY_DIR := bin
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
GO_FILES := $(shell find . -type f -name '*.go')

# Targets
.PHONY: all build clean run runfollower test lint fmt

all: build

build: $(GO_FILES)
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_PATH)

test: build
	@go test -v $(BINARY_PATH)/...

lint:
	@golangci-lint run ./...

fmt:
	@go fmt ./...

clean:
	@rm -rf $(BINARY_DIR)