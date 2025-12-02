# Project variables
PROJECT_NAME := claude-code-foundry
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Directory structure
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
CMD_DIR := cmd/$(PROJECT_NAME)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOCLEAN := $(GOCMD) clean
GOVET := $(GOCMD) vet
GOFMT := $(GOCMD) fmt

# Detect OS for cross-platform compatibility
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    EXE_EXT := .exe
    RM := del /Q /F
    RMDIR := rmdir /S /Q
    MKDIR := mkdir
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
    else ifeq ($(UNAME_S),Linux)
        DETECTED_OS := Linux
    else
        DETECTED_OS := Unix
    endif
    EXE_EXT :=
    RM := rm -f
    RMDIR := rm -rf
    MKDIR := mkdir -p
endif

# Build flags
VERSION_FLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)
BUILD_FLAGS := -ldflags="$(VERSION_FLAGS)"
PROD_BUILD_FLAGS := -ldflags="$(VERSION_FLAGS) -s -w"

# Install location
INSTALL_DIR := /usr/local/bin

# Default target
.DEFAULT_GOAL := help

# Help system - self-documenting Makefile
.PHONY: help
help: ## Show this help message
	@echo "$(PROJECT_NAME) - Build and install"
	@echo ""
	@echo "Building for $(DETECTED_OS)"
	@echo "Version: $(VERSION)"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the application (development)
	@echo "Building $(PROJECT_NAME) for $(DETECTED_OS)..."
	@$(MKDIR) $(BIN_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./$(CMD_DIR)
ifeq ($(DETECTED_OS),macOS)
	@codesign -s - $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || echo "Warning: Could not sign binary (continuing...)"
	@echo "✓ Built and signed: $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
else
	@echo "✓ Built to $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
endif

.PHONY: build-prod
build-prod: ## Build optimized production binary
	@echo "Building $(PROJECT_NAME) (production) for $(DETECTED_OS)..."
	@$(MKDIR) $(BIN_DIR)
	$(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./$(CMD_DIR)
ifeq ($(DETECTED_OS),macOS)
	@codesign -s - $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || echo "Warning: Could not sign binary (continuing...)"
	@echo "✓ Production build signed: $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
else
	@echo "✓ Production build: $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
endif

# Cross-platform builds
.PHONY: build-all
build-all: build-linux build-darwin-amd64 build-darwin-arm64 build-windows ## Build for all platforms

.PHONY: build-linux
build-linux: ## Build for Linux x64
	@echo "Building for Linux x64..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-linux-amd64 ./$(CMD_DIR)
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-linux-amd64"

.PHONY: build-darwin-amd64
build-darwin-amd64: ## Build for macOS Intel
	@echo "Building for macOS Intel..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 ./$(CMD_DIR)
	@codesign -s - $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 2>/dev/null || echo "Warning: Could not sign binary"
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 (signed)"

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS Apple Silicon
	@echo "Building for macOS Apple Silicon..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 ./$(CMD_DIR)
	@codesign -s - $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 2>/dev/null || echo "Warning: Could not sign binary"
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 (signed)"

.PHONY: build-windows
build-windows: ## Build for Windows x64
	@echo "Building for Windows x64..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe"

# Dependency management
.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies ready"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	$(GOCMD) get -u ./...
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Testing
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Code quality
.PHONY: lint
lint: ## Run Go linters
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		$(GOVET) ./...; \
	fi

.PHONY: format
format: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: validate
validate: deps lint test ## Run all validation checks

# Installation
.PHONY: install
install: build ## Install to /usr/local/bin
	@echo "Installing $(PROJECT_NAME) to $(INSTALL_DIR)..."
	@if [ "$(DETECTED_OS)" = "Windows" ]; then \
		echo "Please manually copy $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) to your PATH"; \
	else \
		sudo cp $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) $(INSTALL_DIR)/; \
		sudo chmod +x $(INSTALL_DIR)/$(PROJECT_NAME)$(EXE_EXT); \
		echo "✓ Installed to $(INSTALL_DIR)/$(PROJECT_NAME)"; \
		echo ""; \
		echo "Run: $(PROJECT_NAME) --help"; \
	fi

.PHONY: uninstall
uninstall: ## Remove from /usr/local/bin
	@echo "Removing $(PROJECT_NAME) from $(INSTALL_DIR)..."
	@if [ "$(DETECTED_OS)" = "Windows" ]; then \
		echo "Please manually remove $(PROJECT_NAME)$(EXE_EXT) from your PATH"; \
	else \
		sudo rm -f $(INSTALL_DIR)/$(PROJECT_NAME); \
		echo "✓ Uninstalled"; \
	fi

# Cleanup
.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@$(RMDIR) $(BUILD_DIR)
	$(GOCLEAN) -cache -testcache
	@$(RM) coverage.out coverage.html 2>/dev/null || true
	@echo "✓ Cleaned"

.PHONY: clean-deps
clean-deps: ## Remove dependency cache
	@echo "Cleaning dependency cache..."
	$(GOCLEAN) -modcache
	@echo "✓ Dependency cache cleaned"

# Development
.PHONY: run
run: build ## Build and run
	@echo "Running $(PROJECT_NAME)..."
	$(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)

.PHONY: run-help
run-help: build ## Build and show help
	$(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) --help

# Information
.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Built:   $(BUILD_TIME)"
	@echo "OS:      $(DETECTED_OS)"

.PHONY: debug
debug: ## Show Makefile variables
	@echo "PROJECT_NAME: $(PROJECT_NAME)"
	@echo "VERSION: $(VERSION)"
	@echo "BUILD_FLAGS: $(BUILD_FLAGS)"
	@echo "BIN_DIR: $(BIN_DIR)"
	@echo "DETECTED_OS: $(DETECTED_OS)"
	@echo "EXE_EXT: $(EXE_EXT)"
	@echo "INSTALL_DIR: $(INSTALL_DIR)"

# Release
.PHONY: release
release: clean build-all ## Clean and build for all platforms
	@echo ""
	@echo "✓ Release builds complete:"
	@ls -lh $(BIN_DIR)/
