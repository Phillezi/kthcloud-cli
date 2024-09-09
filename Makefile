# Variables
BINARY_NAME=kthcloud
BUILD_DIR=bin
CMD_DIR=cmd/kthcloud-cli
MAIN_FILE=$(CMD_DIR)/main.go

# Targets
.PHONY: all clean build run

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGOENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/*
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."