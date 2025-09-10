# Makefile for Neuron CLI

# --- Variables ---
BINARY_NAME=neuron
BUILD_DIR=build
GOPATH=$(shell go env GOPATH)
INSTALL_PATH=$(GOPATH)/bin

.PHONY: all build clean install uninstall run test help

all: build

## build: Compiles the project into a binary in the ./build directory.
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

## install: Compiles and installs the binary to your Go bin path with the correct name.
install:
	@echo "Building and installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	go build -o $(INSTALL_PATH)/$(BINARY_NAME) .
	@echo "$(BINARY_NAME) installed successfully. Run 'neuron --help' to get started."

## uninstall: Removes the installed binary from your Go bin path.
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f "$(INSTALL_PATH)/$(BINARY_NAME)"
	@rm -f "$(INSTALL_PATH)/neuron-cli"
	@echo "$(BINARY_NAME) has been removed."

## run: A shortcut to run the program from source.
run:
	go run . $(filter-out $@,$(MAKECMDGOALS))

## clean: Removes the build directory and artifacts.
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

## test: Runs all tests in the project.
test:
	@echo "Running tests..."
	go test -v ./...

## help: Shows this help message.
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
