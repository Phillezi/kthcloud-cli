# Variables
BINARY_NAME=kthcloud
BUILD_DIR=bin
CMD_DIR=cmd/kthcloud-cli
MAIN_FILE=$(CMD_DIR)/main.go
BUILDTIMESTAMP=$(shell date -u +%Y%m%d%H%M%S)

# Targets
.PHONY: all clean build run

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/*
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."