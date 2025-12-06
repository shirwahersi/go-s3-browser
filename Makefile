# Go S3 Browser Makefile
# Supports cross-compilation for macOS, Linux, and Windows

BINARY_NAME=go-s3-browser
DIST_DIR=dist

# Default target
.PHONY: build
build:
	go build -o $(BINARY_NAME)

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Linux builds
.PHONY: build-linux
build-linux: build-linux-amd64 build-linux-arm64

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64

# macOS builds
.PHONY: build-darwin
build-darwin: build-darwin-amd64 build-darwin-arm64

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64

# Windows builds
.PHONY: build-windows
build-windows: build-windows-amd64 build-windows-arm64

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe

.PHONY: build-windows-arm64
build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build for current OS/arch"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run all tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  tidy           - Run go mod tidy"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-linux    - Build for Linux (amd64, arm64)"
	@echo "  build-darwin   - Build for macOS (amd64, arm64)"
	@echo "  build-windows  - Build for Windows (amd64, arm64)"
	@echo ""
	@echo "Individual platform builds:"
	@echo "  build-linux-amd64, build-linux-arm64"
	@echo "  build-darwin-amd64, build-darwin-arm64"
	@echo "  build-windows-amd64, build-windows-arm64"
