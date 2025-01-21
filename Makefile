# Variables
BINARY_NAME=kthcloud
BUILD_DIR=bin
MAIN_FILE=main.go
BUILDTIMESTAMP=$(shell date -u +%Y%m%d%H%M%S)
EXT=$(if $(filter windows,$(GOOS)),.exe,)

# Targets
.PHONY: all build run test release install all-platforms clean lint

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) .
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT)

test:
	@go test ./...

release:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) .
	@echo "Build complete."

install: release
	@echo "installing"
	@mkdir -p ~/.local/kthcloud/bin
	@cp ./$(BUILD_DIR)/$(BINARY_NAME)$(EXT) ~/.local/kthcloud/bin/$(BINARY_NAME)$(EXT)
	@echo "add to PATH"

all-platforms:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_amd64_linux . &
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_arm64_linux . &
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_amd64_windows.exe . &
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_arm64_windows.exe . &
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_amd64_macos . &
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s -X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/kthcloud_arm64_macos . &
	@wait
	@echo "All builds complete."

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

lint:
	@./scripts/util/check-lint.sh
