---
name: makefile-skills-guide
description: Comprehensive guide for creating and managing Makefiles for build automation, full-stack projects, cross-platform compilation, and development workflows. Includes production-proven patterns for Go, npm, testing, and process management.
---

# Makefile Skills Guide

## Usage

### When to Use This Skill

Claude should consult this skill when the user's request involves:

**Primary Triggers:**
- Creating, modifying, or debugging Makefiles
- Setting up build automation for projects
- Questions about Make syntax, patterns, or best practices
- Orchestrating multi-language builds (e.g., Go + npm, Rust + frontend)
- Managing development workflows (start/stop servers, testing, deployment)
- Cross-platform compilation and builds
- Process management and PID tracking

**Specific Keywords/Phrases:**
- "Makefile", "make command", "make target"
- "build automation", "build system"
- "cross-compile", "multi-platform build"
- "full-stack build", "monorepo build"
- "start/stop server", "background process", "daemon"

**Context Clues:**
- User has a project with multiple languages/components
- User mentions wanting reproducible builds
- User asks about automating development tasks
- User is debugging build or process management issues

### When NOT to Use This Skill

Do not consult this skill for:
- Simple single-language projects with good native tooling (e.g., pure npm/yarn projects)
- CI/CD pipeline configuration (GitHub Actions, Jenkins, etc.) unless Make is part of it
- General shell scripting unrelated to builds
- Docker/containerization (unless Make is orchestrating it)

### How to Apply This Skill

1. **For Makefile Creation:**
    - Read the "Example: Creating a New Makefile" section
    - Adapt patterns from relevant sections (e.g., "Full-Stack Integration" for Go+npm projects)
    - Include self-documenting help system from the start
    - Use appropriate patterns based on project complexity

2. **For Debugging:**
    - Reference "Debugging Makefiles" section
    - Check "Common Issues" for known problems
    - Suggest using `make -n` for dry runs

3. **For Enhancement:**
    - Review "Common Patterns from NomadAI" for advanced features
    - Suggest relevant best practices from "Best Practices" section
    - Consider user's tech stack when recommending patterns

4. **For Explanation:**
    - Start with "Basic Makefile Concepts" for beginners
    - Use concrete examples from the guide
    - Reference the specific pattern or section that applies

### Output Guidelines

When using this skill:
- **Provide working, tested patterns** - all examples in this guide are production-proven
- **Explain the "why"** - Makefiles can be cryptic; help users understand the patterns
- **Use proper formatting** - Remember TAB characters for recipes, not spaces
- **Include help documentation** - Always add `## comments` for the help system
- **Consider cross-platform** - Detect OS and provide portable solutions when possible
- **Suggest testing** - Recommend `make -n` before running destructive targets
- **Prioritize portability** - Use language-native tools (Go, npm) over shell commands when possible

---

## What is Make?

Make is a build automation tool that uses a `Makefile` to define targets (tasks) and their dependencies. It's language-agnostic and excels at orchestrating complex build processes, especially for full-stack applications.

**Cross-Platform Note:** Make is available on macOS, Linux/Unix, and Windows (via WSL, Git Bash, or tools like `make` from GnuWin32 or Chocolatey). This guide uses POSIX-compatible patterns where possible, with platform-specific notes where differences exist.

### When to Use Make

**Use Make when:**
- Building multi-language projects (e.g., Go backend + npm frontend)
- Need reproducible builds across environments
- Managing development workflows (start/stop servers, run tests)
- Cross-platform compilation required
- Want self-documenting build commands

**Consider alternatives when:**
- Single-language project with good native tooling (e.g., just npm for Node projects)
- Need complex conditional logic (shell scripts might be better)
- Team unfamiliar with Make syntax

## Basic Makefile Concepts

### Target Structure

```makefile
target: dependency1 dependency2
	recipe command 1
	recipe command 2
```

- **Target**: The name of the task (e.g., `build`, `test`)
- **Dependencies**: Other targets that must run first
- **Recipe**: Shell commands (MUST be indented with TAB, not spaces)

### Phony Targets

Declare targets that don't create files:

```makefile
.PHONY: clean test run

clean:
	rm -rf build/

test:
	go test ./...
```

Without `.PHONY`, Make would check for a file named "clean" or "test".

### Variables

```makefile
# Simple variable
VERSION := 1.0.0

# Shell command variable
COMMIT := $(shell git rev-parse HEAD)

# Using variables
build:
	go build -ldflags="-X main.Version=$(VERSION)"
```

**Variable assignment types:**
- `:=` - Simple expansion (evaluated once)
- `=` - Recursive expansion (evaluated each use)
- `?=` - Set only if not already set
- `+=` - Append

### Cross-Platform Shell Commands

Different platforms require different shell syntax:

```makefile
# Detect OS
ifeq ($(OS),Windows_NT)
    # Windows-specific (use cmd.exe syntax or PowerShell)
    RM := del /Q
    MKDIR := mkdir
    SEP := \\
else
    # Unix-like (macOS, Linux)
    RM := rm -rf
    MKDIR := mkdir -p
    SEP := /
endif

# Use platform-aware commands
clean:
	$(RM) build$(SEP)*
```

**Best Practice:** Use cross-platform tools where possible (Go, Node.js, Python scripts) instead of shell-specific commands.

### Silent Commands

Prefix with `@` to suppress command echo:

```makefile
status:
	@echo "Checking status..."  # Only output shows
	@ps aux | grep myapp        # Command not shown
```

## Common Patterns for Production Makefiles

### 1. Self-Documenting Help System

```makefile
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
```

Usage in targets:
```makefile
.PHONY: build
build: ## Build frontend and backend
	@echo "Building..."
```

Running `make help` automatically generates formatted help from `##` comments.

### 2. Conditional Tool Detection

```makefile
.PHONY: lint-backend
lint-backend: ## Run Go linters
	@echo "Running Go linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(BACKEND_DIR) && golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		cd $(BACKEND_DIR) && go vet ./...; \
	fi
```

**Pattern**: Prefer optimal tool, fallback to available alternative.

### 3. Multi-Stage Build Dependencies

```makefile
.PHONY: build
build: build-frontend build-backend ## Build frontend and backend

.PHONY: build-frontend
build-frontend: ## Build frontend
	cd $(FRONTEND_DIR) && npm ci && npm run build

.PHONY: build-backend
build-backend: ## Build backend
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BIN_DIR)/$(BINARY_NAME) cmd/main.go
```

Running `make build` executes both sub-targets in order.

### 4. Cross-Platform Compilation

```makefile
.PHONY: build-all
build-all: build-frontend build-linux build-darwin-arm64 build-windows ## Build for all platforms

.PHONY: build-linux
build-linux: ## Build for Linux x64
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME)-linux-amd64 cmd/main.go

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS Apple Silicon
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/main.go
```

### 5. Version Embedding

```makefile
# Extract version info from git
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Embed into binary
VERSION_FLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
BUILD_FLAGS := -ldflags="$(VERSION_FLAGS)"

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o bin/app cmd/main.go
```

Access in Go code:
```go
var (
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"
)
```

## Advanced Features

### Process & Port Management

#### The Problem: PID/Port Staleness

A naive approach stores PIDs and ports in files. When a process dies unexpectedly, those files become stale. If another service (or external process) grabs the same port, the original service may incorrectly kill it on restart. **This is a real coordination bug in multi-service development environments.**

#### The Solution: PID + Port Ownership Verification

**Never trust a stored PID/port at face value.** Before acting on any recorded state, verify:

1. **Is the PID still alive?** (`kill -0 $PID`)
2. **Does that PID actually own the recorded port?** (Check with `lsof`)

Only if BOTH conditions are true should you treat the process as "ours." If the PID is dead, or alive but not bound to the expected port, the entry is stale — clean it up, don't kill anything.

#### State File Format

Use a flat text file at `.run/state` with one entry per line:

```
role:pid:port
```

Examples:
```
http:1234:8080
grpc:1235:9090
worker:1236:0
```

- **role**: A label for the process (e.g., `http`, `grpc`, `worker`, `scheduler`)
- **pid**: The OS process ID
- **port**: The port number, or `0` if the process doesn't listen on a port

This format is trivially parseable with `awk`, `grep`, and `cut` — ideal for Makefile shell recipes.

#### Cross-Platform Port Ownership Check

macOS and Linux have different `lsof` behavior. Use this pattern:

```bash
# Check if a specific PID owns a specific port.
# Returns 0 (true) if the PID is listening on the port, 1 (false) otherwise.
pid_owns_port() {
    local pid="$1" port="$2"
    if [ "$port" = "0" ]; then
        # No port to check — just verify PID is alive
        kill -0 "$pid" 2>/dev/null
        return $?
    fi
    # lsof works on both macOS and Linux
    # -iTCP:PORT matches the port, -sTCP:LISTEN filters to listeners,
    # then we grep for the PID
    lsof -iTCP:"$port" -sTCP:LISTEN -nP 2>/dev/null | awk '{print $2}' | grep -q "^${pid}$"
}
```

**Platform notes:**
- `lsof` is available on macOS (preinstalled), Linux (install `lsof` package), and WSL (Linux behavior)
- `-nP` avoids DNS/port name resolution for speed
- On Linux, `ss -tlnp` is an alternative but output format differs; `lsof` is more portable
- If `lsof` is not available, fall back to: `ss -tlnp 2>/dev/null | grep ":$port " | grep "pid=$pid"` (Linux/WSL only)

#### Finding a Free Port

```bash
# Find an available port. Checks that nothing is listening on it.
find_free_port() {
    local port
    for port in $(seq 8000 9000); do
        if ! lsof -iTCP:"$port" -sTCP:LISTEN -nP >/dev/null 2>&1; then
            echo "$port"
            return 0
        fi
    done
    echo "ERROR: No free port found in range 8000-9000" >&2
    return 1
}
```

Adjust the port range as needed. Some services may want a specific range (e.g., 3000-3999 for frontend, 8000-8999 for backend).

#### Verification Logic: The Core Pattern

This function validates all entries in the state file and returns what's still actually ours:

```bash
# Verify all entries in the state file.
# Prints verified entries to stdout, warnings to stderr.
# Returns the number of stale entries found.
verify_state() {
    local state_file="$1"
    local stale=0
    if [ ! -f "$state_file" ]; then
        return 0
    fi
    while IFS=: read -r role pid port; do
        if kill -0 "$pid" 2>/dev/null; then
            if pid_owns_port "$pid" "$port"; then
                echo "$role:$pid:$port"  # Verified — still ours
            else
                echo "WARNING: PID $pid alive but NOT on port $port (role: $role) — stale entry" >&2
                stale=$((stale + 1))
            fi
        else
            echo "WARNING: PID $pid dead (role: $role, port: $port) — stale entry" >&2
            stale=$((stale + 1))
        fi
    done < "$state_file"
    return $stale
}
```

#### Recommended Script: `scripts/process-manager.sh`

Bundle the above functions into a single script with subcommands. The Makefile calls it like:

```bash
./scripts/process-manager.sh start <role> <command> [--port-range 8000-9000]
./scripts/process-manager.sh stop [role]          # Stop one or all
./scripts/process-manager.sh status               # Show verified state
./scripts/process-manager.sh cleanup              # Remove stale entries
./scripts/process-manager.sh find-port [range]    # Just find and print a free port
```

Claude Code should implement this script with these behaviors:

**`start <role> <command>`:**
1. Run `cleanup` first (verify existing state, remove stale entries, warn)
2. Check if this role is already running (verified) — if so, error or restart
3. Find a free port (unless the process doesn't need one)
4. Launch the process with `nohup`, capture PID
5. Brief sleep + verify the process actually started and bound the port
6. Append `role:pid:port` to `.run/state`

**`stop [role]`:**
1. Read state file, verify each entry with `pid_owns_port`
2. For verified entries matching the role (or all if no role given): send `SIGTERM`, wait, `SIGKILL` as fallback
3. Remove stopped entries from state file
4. Warn about any stale entries (don't try to kill them — they're not ours)

**`status`:**
1. Read and verify all state entries
2. Print a table: role, PID, port, status (running/stale/dead)

**`cleanup`:**
1. Verify all entries, rewrite state file with only verified entries
2. Print warnings for anything removed

#### Makefile Integration

```makefile
RUN_DIR := .run
STATE_FILE := $(RUN_DIR)/state
LOG_DIR := $(RUN_DIR)/logs
PM := ./scripts/process-manager.sh

.PHONY: start
start: build ## Start service(s) in background
	@mkdir -p $(RUN_DIR) $(LOG_DIR)
	@$(PM) start http "$(BIN_DIR)/$(BINARY_NAME) -addr :PORT" --port-range 8000-9000 --log-dir $(LOG_DIR)

.PHONY: stop
stop: ## Stop all service processes
	@$(PM) stop

.PHONY: restart
restart: stop build start ## Restart service(s)

.PHONY: status
status: ## Show status of service processes
	@$(PM) status

.PHONY: force-stop
force-stop: ## Emergency: kill all verified processes immediately
	@$(PM) stop --force
```

**Notes for Claude Code:**
- The `PORT` placeholder in the command string should be replaced by process-manager.sh with the actual allocated port
- For services with multiple processes, call `start` multiple times with different roles:
  ```makefile
  start: build
  	@$(PM) start http "$(BIN_DIR)/$(BINARY_NAME) serve -addr :PORT" --port-range 8000-8999
  	@$(PM) start worker "$(BIN_DIR)/$(BINARY_NAME) worker" --no-port
  ```
- The `--no-port` flag indicates a process that doesn't listen on any port (stored as port `0`)

#### Key Principles for Claude Code

1. **NEVER kill a process based solely on a stored PID.** Always verify PID + port ownership first.
2. **Stale entries are informational, not actionable.** If PID is dead or doesn't own the port, just clean up the state file. Don't try to kill whatever is now on that port.
3. **Always allocate fresh ports on start.** Don't reuse a port from the state file — scan for a free one.
4. **Warn loudly about stale state.** The user should know when entries didn't match reality.
5. **Graceful shutdown first.** `SIGTERM` → wait (2-5 seconds) → `SIGKILL` as last resort.
6. **Verify after launch.** After starting a process, briefly sleep and confirm it's alive and bound to the expected port before recording state.
7. **Clean up state file on stop.** Remove entries for processes that were successfully stopped.
8. **Handle the "no state file" case.** First run, missing `.run/` dir — these should all work cleanly.
9. **Cross-platform**: Use `lsof -iTCP:PORT -sTCP:LISTEN -nP` for port checks — works on macOS, Linux, and WSL. Fall back to `ss` on Linux if `lsof` is unavailable.

### Development vs Production Builds

```makefile
# Development: Fast builds with debug info
BUILD_FLAGS := -ldflags="$(VERSION_FLAGS)"

# Production: Optimized, stripped binaries
PROD_BUILD_FLAGS := -ldflags="$(VERSION_FLAGS) -s -w"

.PHONY: build
build: build-frontend build-backend ## Development build

.PHONY: build-prod
build-prod: build-frontend build-backend-prod ## Production build

.PHONY: build-backend-prod
build-backend-prod:
	cd $(BACKEND_DIR) && go build $(PROD_BUILD_FLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME) cmd/main.go
```

**Flags:**
- `-s` - Strip symbol table
- `-w` - Strip DWARF debugging info
- Result: Smaller binary, faster startup, no debugging

## Full-Stack Integration (Go + npm)

### Coordinated Build Process

```makefile
FRONTEND_DIR := frontend
BACKEND_DIR := backend
DIST_DIR := $(FRONTEND_DIR)/dist

.PHONY: build
build: build-frontend build-backend ## Build full stack

.PHONY: build-frontend
build-frontend: ## Build frontend with npm
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm ci && npm run build
	# Vite outputs to backend/internal/server/web for embedding

.PHONY: build-backend
build-backend: ## Build Go backend (embeds frontend)
	@echo "Building backend..."
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BIN_DIR)/$(BINARY_NAME) cmd/main.go
```

### Development Workflows

```makefile
.PHONY: run
run: build-frontend ## Run with hot reload (go run)
	cd $(BACKEND_DIR) && go run cmd/main.go

.PHONY: dev-frontend
dev-frontend: ## Vite dev server (port 5173)
	cd $(FRONTEND_DIR) && npm run dev

.PHONY: watch
watch: ## Auto-restart on file changes (requires entr)
	@if command -v entr >/dev/null 2>&1; then \
		find $(BACKEND_DIR) -name "*.go" -o -path "./$(FRONTEND_DIR)/src" | entr -r make run; \
	else \
		echo "entr not installed. Install: brew install entr"; \
		make run; \
	fi
```

**Development modes:**
1. `make run` - Foreground, `go run` auto-recompiles
2. `make dev-frontend` - Vite dev server with HMR (hot module replacement)
3. `make watch` - Auto-restart on any file change
4. `make start` - Background production-like mode

## Testing Integration

### Unit Tests

```makefile
.PHONY: test
test: test-backend test-frontend ## Run all unit tests

.PHONY: test-backend
test-backend: ## Run Go tests
	cd $(BACKEND_DIR) && go test -v ./...

.PHONY: test-frontend
test-frontend: ## Run frontend tests
	cd $(FRONTEND_DIR) && npm test

.PHONY: test-coverage
test-coverage: ## Go tests with coverage report
	cd $(BACKEND_DIR) && go test -v -coverprofile=../coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
```

### E2E Tests (Playwright)

```makefile
.PHONY: test-e2e-setup
test-e2e-setup: ## Install Playwright browsers (one-time)
	cd $(FRONTEND_DIR) && npx playwright install --with-deps

.PHONY: test-e2e
test-e2e: ## Run E2E tests (server must be running)
	cd $(FRONTEND_DIR) && npx playwright test

.PHONY: test-e2e-ui
test-e2e-ui: ## Run E2E tests with UI
	cd $(FRONTEND_DIR) && npx playwright test --ui

.PHONY: test-e2e-debug
test-e2e-debug: ## Run E2E tests with inspector
	cd $(FRONTEND_DIR) && npx playwright test --debug

.PHONY: test-all
test-all: test test-e2e ## Run all tests
```

## Linting and Formatting

```makefile
.PHONY: lint
lint: lint-backend lint-frontend ## Run all linters

.PHONY: lint-backend
lint-backend:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(BACKEND_DIR) && golangci-lint run; \
	else \
		cd $(BACKEND_DIR) && go vet ./...; \
	fi

.PHONY: lint-frontend
lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

.PHONY: format
format: format-backend format-frontend ## Format all code

.PHONY: format-backend
format-backend:
	cd $(BACKEND_DIR) && go fmt ./...

.PHONY: format-frontend
format-frontend:
	cd $(FRONTEND_DIR) && npm run format

.PHONY: validate
validate: deps lint test ## Run all validation checks
```

## Clean Targets

```makefile
# Cross-platform clean targets
ifeq ($(OS),Windows_NT)
    RM := del /Q /F
    RMDIR := rmdir /S /Q
else
    RM := rm -f
    RMDIR := rm -rf
endif

.PHONY: clean
clean: clean-backend clean-frontend clean-runtime ## Clean all

.PHONY: clean-backend
clean-backend:
	$(RMDIR) $(BUILD_DIR)
	$(RM) coverage.out coverage.html

.PHONY: clean-frontend
clean-frontend:
	$(RMDIR) $(DIST_DIR)

.PHONY: clean-runtime
clean-runtime:
	@if [ -f $(RUN_DIR)/state ]; then \
		echo "Verifying running processes before cleanup..."; \
		if ./scripts/process-manager.sh status 2>/dev/null | grep -q "running"; then \
			echo "Warning: Services still running. Use 'make stop' first."; \
		exit 1; \
		fi; \
	fi
	$(RMDIR) $(RUN_DIR)

.PHONY: clean-deps
clean-deps:
	cd $(BACKEND_DIR) && go clean -modcache
	cd $(FRONTEND_DIR) && $(RMDIR) node_modules && $(RM) package-lock.json
```

**Windows Note:** The `clean-runtime` target uses Unix-style process checking. On Windows, consider using PowerShell commands or checking for the PID file existence only.

## Dependency Management

```makefile
.PHONY: deps
deps: deps-backend deps-frontend ## Install all dependencies

.PHONY: deps-backend
deps-backend:
	cd $(BACKEND_DIR) && go mod tidy
	cd $(BACKEND_DIR) && go mod download

.PHONY: deps-frontend
deps-frontend:
	cd $(FRONTEND_DIR) && npm ci  # Uses package-lock.json for reproducibility

.PHONY: deps-update
deps-update:
	cd $(BACKEND_DIR) && go get -u ./...
	cd $(BACKEND_DIR) && go mod tidy
	cd $(FRONTEND_DIR) && npm update
```

**Best practices:**
- Use `npm ci` (not `npm install`) for reproducible builds
- Run `go mod tidy` to clean up unused dependencies
- Separate update target for intentional upgrades

## Best Practices

### 1. Organization

Group related targets with comments:

```makefile
# Build targets
.PHONY: build build-frontend build-backend

# Testing targets
.PHONY: test test-backend test-frontend test-e2e

# Development targets
.PHONY: run dev-frontend dev-backend
```

### 2. Default Target

```makefile
.DEFAULT_GOAL := build

# Or explicitly:
.PHONY: default
default: build
```

### 3. Error Handling

```makefile
.PHONY: build-backend
build-backend:
	@mkdir -p $(BIN_DIR) || { echo "Failed to create bin dir"; exit 1; }
	cd $(BACKEND_DIR) && go build -o ../$(BIN_DIR)/app cmd/main.go
```

### 4. Parallel Execution

Allow parallel builds:

```makefile
# Enable parallel execution
MAKEFLAGS += -j4

# Or: make -j4 build
```

### 5. Directory Variables

```makefile
# Centralize paths
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
SRC_DIR := src

# Use consistently
build:
	@mkdir -p $(BIN_DIR)
	gcc $(SRC_DIR)/main.c -o $(BIN_DIR)/app
```

### 6. Cross-Platform Detection and Compilation

```makefile
# Detect OS for platform-specific behavior
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    EXE_EXT := .exe
    RM := del /Q
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
    RM := rm -rf
    MKDIR := mkdir -p
endif

build:
	@echo "Building for $(DETECTED_OS)..."
	go build -o bin/app$(EXE_EXT) cmd/main.go

# Cross-compilation targets for Go
.PHONY: build-all
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app-linux-amd64 cmd/main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/app-darwin-amd64 cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/app-darwin-arm64 cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/app-windows-amd64.exe cmd/main.go
```

**Cross-Platform Tips:**
- **Windows:** Use WSL, Git Bash, or install Make via Chocolatey (`choco install make`)
- **macOS:** Make is pre-installed via Xcode Command Line Tools
- **Linux:** Usually pre-installed; if not: `apt-get install build-essential` or `yum install make`
- **Path separators:** Use forward slashes `/` in Makefiles (works on all platforms with modern Make)
- **Executables:** Add `.exe` extension for Windows binaries automatically with `$(EXE_EXT)`

### 7. Writing Portable Makefiles

**DO:**
- Use `$(shell ...)` for cross-platform command execution
- Detect OS and set variables accordingly
- Use language-specific tools (Go, Node.js) that are inherently cross-platform
- Use forward slashes `/` for paths (works everywhere)
- Test on target platforms

**DON'T:**
- Hard-code Unix commands like `rm`, `cp`, `mv` without OS detection
- Assume bash-specific syntax (use POSIX shell)
- Use platform-specific tools without fallbacks
- Hard-code path separators

**Example of portable clean target:**
```makefile
# Bad - Unix only
clean:
	rm -rf build/

# Good - Cross-platform
ifeq ($(OS),Windows_NT)
    RM_CMD := if exist build rmdir /S /Q build
else
    RM_CMD := rm -rf build
endif

clean:
	$(RM_CMD)

# Better - Use Go's cross-platform capabilities
clean:
	go clean -cache -testcache
	@echo "Cleaned build artifacts"
```

## Quick Reference Guide

### Setup & Installation
```bash
make check-tools    # Verify Go, Node, npm installed
make setup          # Complete dev setup (deps + build)
make install        # Install binary to /usr/local/bin
```

### Development Workflows
```bash
make run                    # Foreground dev server (go run)
make dev-frontend           # Vite dev server only (HMR)
make watch-production-like  # Auto-restart on changes (requires entr)
```

### Server Management
```bash
make start         # Background server with PID+port verification
make stop          # Graceful shutdown (verified processes only)
make restart       # Stop + build + start
make restart-clean # Restart + clear frontend cache
make force-stop    # Emergency kill all verified processes
make status        # Show verified status (role, PID, port, state)
make logs          # Show last 50 lines
make logs-follow   # Tail logs (Ctrl+C to stop)
```

### Building
```bash
make build           # Dev build (current platform)
make build-prod      # Production build (optimized)
make build-all       # All platforms (Linux, macOS, Windows)
make build-frontend  # Frontend only (npm ci + build)
make build-backend   # Backend only (go build)
```

### Testing
```bash
make test              # Unit tests (Go + npm)
make test-coverage     # Go tests with HTML coverage
make test-e2e-setup    # Install Playwright (one-time)
make test-e2e          # E2E tests (server must be running)
make test-e2e-ui       # E2E with Playwright UI
make test-e2e-debug    # E2E with inspector
make test-all          # All tests (unit + E2E)
```

### Code Quality
```bash
make lint         # Run all linters
make format       # Format all code
make validate     # deps + lint + test
```

### Cleanup
```bash
make clean           # All build artifacts
make clean-deps      # Remove all dependencies
make clean-runtime   # PID/log files (checks if server running)
```

### Dependencies
```bash
make deps         # Install all dependencies
make deps-update  # Update all dependencies
```

### Other
```bash
make help     # Show all targets with descriptions
make version  # Show version info (git-based)
make release  # Clean + build all platforms
```

## Debugging Makefiles

### See What Make Will Do

```bash
make -n build      # Dry run (show commands, don't execute)
make -d build      # Debug mode (verbose output)
```

### Print Variables

```makefile
.PHONY: debug
debug:
	@echo "VERSION: $(VERSION)"
	@echo "BUILD_FLAGS: $(BUILD_FLAGS)"
	@echo "BIN_DIR: $(BIN_DIR)"
```

### Check Syntax

```bash
make --version           # Verify Make is installed
make -f Makefile build   # Explicitly specify Makefile
```

### Common Issues

**Issue**: "missing separator" error  
**Fix**: Ensure recipes use TAB (not spaces) for indentation

**Issue**: Variable not expanding  
**Fix**: Use `:=` instead of `=` for immediate expansion

**Issue**: Target always runs (even when up-to-date)  
**Fix**: Add `.PHONY` declaration if target doesn't create a file

**Issue**: Command not found in PATH  
**Fix**: Use absolute paths or check with `command -v <tool>`

**Issue**: Make not installed on Windows  
**Fix**: Install via WSL (`apt-get install build-essential`), Git Bash (comes with Git for Windows), Chocolatey (`choco install make`), or use nmake (Microsoft's Make)

**Issue**: Shell syntax errors on Windows  
**Fix**: Detect OS with `ifeq ($(OS),Windows_NT)` and provide Windows-specific commands, or use WSL/Git Bash for Unix-like shell

## Example: Creating a New Makefile

```makefile
# Project variables
PROJECT_NAME := myapp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin

# Detect OS for cross-platform compatibility
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    EXE_EXT := .exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
    else
        DETECTED_OS := Linux
    endif
    EXE_EXT :=
endif

# Build flags
BUILD_FLAGS := -ldflags="-X main.Version=$(VERSION)"

# Default target
.DEFAULT_GOAL := help

# Help system
.PHONY: help
help: ## Show this help
	@echo "Building for $(DETECTED_OS)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build
.PHONY: build
build: ## Build the application
	@echo "Building $(PROJECT_NAME) for $(DETECTED_OS)..."
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT) ./cmd

# Test
.PHONY: test
test: ## Run tests
	go test -v ./...

# Clean
.PHONY: clean
clean: ## Remove build artifacts
	go clean -cache
	@echo "Cleaned build artifacts"

# Run
.PHONY: run
run: build ## Build and run
	$(BIN_DIR)/$(PROJECT_NAME)$(EXE_EXT)
```

## Additional Resources

- **GNU Make Manual**: https://www.gnu.org/software/make/manual/
- **Make Tutorial**: https://makefiletutorial.com/
- **Go Build**: https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
- **npm Scripts**: https://docs.npmjs.com/cli/v9/using-npm/scripts

## Summary

Makefiles excel at:
- Orchestrating multi-language builds
- Creating reproducible, documented workflows
- Managing complex dependencies
- Providing consistent interfaces across projects

Key takeaways:
1. Use `.PHONY` for non-file targets
2. Implement self-documenting help with AWK
3. Leverage conditional logic for tool detection
4. Separate dev and production builds
5. Organize targets logically with comments
6. Use variables for paths and configuration
7. Test thoroughly with `-n` (dry run) first

These patterns demonstrate advanced techniques for full-stack applications, including dynamic port management, graceful process handling, cross-platform builds, and comprehensive testing integration.