# Variables
BINARY_NAME=kthcloud
BUILD_DIR=bin
MAIN_FILE=main.go
BUILDTIMESTAMP=$(shell date -u +%Y%m%d%H%M%S)

# Targets
.PHONY: all clean build run

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

install: build
	@echo "installing"
	@mkdir -p ~/.local/kthcloud/bin
	@cp ./$(BUILD_DIR)/$(BINARY_NAME) ~/.local/kthcloud/bin/$(BINARY_NAME)
	@echo "add to PATH"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."