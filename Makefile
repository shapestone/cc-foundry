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

# Code signing and notarization settings (macOS only)
# Set these as environment variables or pass to make:
# export CODESIGN_IDENTITY="Developer ID Application: Your Name (TEAM_ID)"
# export NOTARIZE_APPLE_ID="your-apple-id@example.com"
# export NOTARIZE_PASSWORD="xxxx-xxxx-xxxx-xxxx"  # App-specific password from appleid.apple.com
# export NOTARIZE_TEAM_ID="YOUR_TEAM_ID"
CODESIGN_IDENTITY ?=
NOTARIZE_APPLE_ID ?=
NOTARIZE_PASSWORD ?=
NOTARIZE_TEAM_ID ?=

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
build: ## Build the application (development, no signing for speed)
	@echo "Building $(PROJECT_NAME) for $(DETECTED_OS)..."
	@$(MKDIR) $(BIN_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./$(CMD_DIR)
	@echo "✓ Built to $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
	@echo "Note: For signed/notarized build, use 'make build-signed' or 'make release'"

.PHONY: build-signed
build-signed: ## Build and sign with Developer ID (no notarization)
	@echo "Building $(PROJECT_NAME) for $(DETECTED_OS)..."
	@$(MKDIR) $(BIN_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./$(CMD_DIR)
ifeq ($(DETECTED_OS),macOS)
	@if [ -z "$(CODESIGN_IDENTITY)" ]; then \
		echo "Error: CODESIGN_IDENTITY must be set"; \
		echo ""; \
		echo "Find your identity with:"; \
		echo "  security find-identity -v -p codesigning"; \
		echo ""; \
		echo "Then set:"; \
		echo "  export CODESIGN_IDENTITY='Developer ID Application: Your Name (TEAM_ID)'"; \
		exit 1; \
	fi
	@xattr -cr $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || true
	@codesign --sign "$(CODESIGN_IDENTITY)" --timestamp --options=runtime --force $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || echo "Warning: Could not sign binary (continuing...)"
	@echo "✓ Built and signed with Developer ID: $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
	@echo "Note: Binary is signed but NOT notarized. Use 'make release' for full notarization."
else
	@echo "✓ Built to $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
endif

.PHONY: build-prod
build-prod: ## Build optimized production binary
	@echo "Building $(PROJECT_NAME) (production) for $(DETECTED_OS)..."
	@$(MKDIR) $(BIN_DIR)
	$(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./$(CMD_DIR)
ifeq ($(DETECTED_OS),macOS)
	@xattr -cr $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || true
	@codesign --sign "$(CODESIGN_IDENTITY)" --timestamp --options=runtime --force $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) 2>/dev/null || echo "Warning: Could not sign binary (continuing...)"
	@echo "✓ Production build signed with Developer ID: $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)"
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
	@xattr -cr $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 2>/dev/null || true
	@codesign --sign "$(CODESIGN_IDENTITY)" --timestamp --options=runtime --force $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 2>/dev/null || echo "Warning: Could not sign binary"
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 (signed with Developer ID)"

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS Apple Silicon
	@echo "Building for macOS Apple Silicon..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 ./$(CMD_DIR)
	@xattr -cr $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 2>/dev/null || true
	@codesign --sign "$(CODESIGN_IDENTITY)" --timestamp --options=runtime --force $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 2>/dev/null || echo "Warning: Could not sign binary"
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 (signed with Developer ID)"

.PHONY: build-windows
build-windows: ## Build for Windows x64
	@echo "Building for Windows x64..."
	@$(MKDIR) $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(PROD_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "✓ $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe"

# Notarization (macOS only)
.PHONY: notarize
notarize: ## Notarize macOS binary with Apple (requires NOTARIZE_APPLE_ID and NOTARIZE_PASSWORD)
	@if [ "$(DETECTED_OS)" != "macOS" ]; then \
		echo "Notarization is only available on macOS"; \
		exit 1; \
	fi
	@if [ -z "$(NOTARIZE_APPLE_ID)" ] || [ -z "$(NOTARIZE_PASSWORD)" ] || [ -z "$(NOTARIZE_TEAM_ID)" ]; then \
		echo "Error: NOTARIZE_APPLE_ID, NOTARIZE_PASSWORD, and NOTARIZE_TEAM_ID must be set"; \
		echo ""; \
		echo "Get app-specific password from: https://appleid.apple.com/account/manage"; \
		echo "Find your Team ID with:"; \
		echo "  security find-identity -v -p codesigning"; \
		echo "  (Look for the 10-character code in parentheses)"; \
		echo ""; \
		echo "Then run:"; \
		echo "  export NOTARIZE_APPLE_ID='your-email@example.com'"; \
		echo "  export NOTARIZE_PASSWORD='xxxx-xxxx-xxxx-xxxx'"; \
		echo "  export NOTARIZE_TEAM_ID='YOUR_TEAM_ID'"; \
		echo "  make notarize"; \
		exit 1; \
	fi
	@echo "Notarizing $(PROJECT_NAME) with Apple..."
	@echo "This may take 5-15 minutes..."
	@$(RM) $(BIN_DIR)/$(PROJECT_NAME).zip 2>/dev/null || true
	@cd $(BIN_DIR) && zip -q $(PROJECT_NAME).zip $(PROJECT_NAME)
	@xcrun notarytool submit $(BIN_DIR)/$(PROJECT_NAME).zip \
		--apple-id "$(NOTARIZE_APPLE_ID)" \
		--password "$(NOTARIZE_PASSWORD)" \
		--team-id "$(NOTARIZE_TEAM_ID)" \
		--wait
	@echo "Stapling notarization ticket..."
	@xcrun stapler staple $(BIN_DIR)/$(PROJECT_NAME) 2>/dev/null || echo "Warning: Could not staple ticket"
	@$(RM) $(BIN_DIR)/$(PROJECT_NAME).zip
	@echo "✓ Notarization complete"
	@echo ""
	@echo "Verify with:"
	@echo "  spctl --assess --verbose=4 --type execute $(BIN_DIR)/$(PROJECT_NAME)"

.PHONY: notarize-check
notarize-check: ## Check if binary is notarized
	@if [ "$(DETECTED_OS)" != "macOS" ]; then \
		echo "Notarization check only available on macOS"; \
		exit 1; \
	fi
	@echo "Checking notarization status..."
	@spctl --assess --verbose=4 --type execute $(BIN_DIR)/$(PROJECT_NAME) 2>&1 || true
	@echo ""
	@codesign -dv $(BIN_DIR)/$(PROJECT_NAME) 2>&1 | grep -E "(Authority|TeamIdentifier|Timestamp)" || true

.PHONY: build-notarize
build-notarize: build notarize ## Build and notarize in one step

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
release: clean build-signed notarize ## Clean, build, sign, and notarize for macOS release
	@echo ""
	@echo "✓ Release build complete and notarized!"
	@echo ""
	@echo "Verify with:"
	@echo "  spctl --assess --verbose=4 --type execute $(BIN_DIR)/$(PROJECT_NAME)"
	@echo ""
	@echo "Install with:"
	@echo "  make install"

.PHONY: release-all
release-all: clean build-all ## Clean and build for all platforms (cross-compile only)
	@echo ""
	@echo "✓ Release builds complete:"
	@ls -lh $(BIN_DIR)/
	@echo ""
	@echo "Note: Only macOS builds are signed. Use 'make release' for notarized macOS binary."
