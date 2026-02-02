---
name: project-layout-go
description: Standard Go project layout guidelines following community best practices. Use when creating or organizing Go projects to ensure proper structure, maintainability, and adherence to Go conventions.
version: 4
last_updated: 2026-02-02
trigger_patterns:
  - Creating a new Go/Golang project
  - Organizing Go code structure
  - Setting up Go application architecture
  - Questions about Go project directories
  - "where should I put" + Go code
---

# Golang Project Layout Skill

## Core Principles

1. **Don't create directories until needed** - Start simple, add structure as the project grows
2. **Follow Go conventions** - The Go compiler enforces some patterns (like `internal/`)
3. **Separate concerns** - Keep application logic, libraries, and configuration organized
4. **Think about reusability** - Code in `/pkg` is public API, code in `/internal` is private

## Standard Project Structure

```
myapp/
├── cmd/                    # Main applications
│   └── myapp/
│       └── main.go
├── internal/               # Private application code (compiler-enforced)
│   ├── app/               # Application-specific code
│   │   └── myapp/
│   └── pkg/               # Shared internal libraries
│       └── myprivlib/
├── pkg/                    # Public library code (importable by others)
│   └── mypublib/
├── api/                    # API definitions
├── web/                    # Web assets
├── scripts/                # Build/deployment scripts
├── build/                  # Packaging & CI
│   ├── package/
│   └── ci/
├── docs/                   # Documentation
├── test/                   # External test files & data
├── examples/               # Usage examples
├── tools/                  # Supporting tools
├── assets/                 # Images, logos, etc.
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Build automation (optional)
└── README.md
```

## Directory Detailed Guide

### `/cmd` - Main Applications
**When to use:** Every executable your project produces

- Directory name should match executable name: `/cmd/myapp/`
- Keep `main.go` minimal - it should mostly call code from `/internal` or `/pkg`
- Multiple applications? Create separate subdirectories: `/cmd/api-server/`, `/cmd/worker/`
- **Anti-pattern:** Putting business logic directly in `main.go`

```go
// cmd/myapp/main.go - GOOD
package main

import "myapp/internal/app/myapp"

func main() {
    app.Run()  // Actual logic lives in internal/
}
```

### `/internal` - Private Code
**When to use:** Code you don't want others importing (compiler-enforced since Go 1.4)

- `/internal/app/myapp/` - Application-specific code
- `/internal/pkg/mylib/` - Shared code between your internal apps
- Go compiler prevents imports like `import "github.com/user/repo/internal/pkg/secret"`
- **Best for:** Business logic, application-specific utilities, implementation details

### `/pkg` - Public Libraries
**When to use:** Code that's safe for external projects to import

- Think carefully before adding here - this is a public contract
- Other projects may depend on this code
- Consider API stability and backwards compatibility
- **Good candidates:** Reusable utilities, well-tested libraries, stable interfaces
- **Alternative:** Many modern Go projects skip `/pkg` and put libraries directly in root

### `/api` - API Definitions
**Contents:**
- OpenAPI/Swagger specs (`.yaml`, `.json`)
- Protocol Buffer definitions (`.proto`)
- JSON Schema files
- GraphQL schemas (`.graphql`)
- API documentation

### `/web` - Web Application Assets
**Contents:**
- Static files: CSS, JavaScript, images
- Server-side templates
- Single-page application (SPA) files
- **Note:** For pure API services, this directory isn't needed

### `/scripts` - Build & Automation Scripts
**When to use:** Build, install, analysis, deployment automation

- Keeps root Makefile clean
- Examples: `scripts/build.sh`, `scripts/test.sh`, `scripts/deploy.sh`
- Can be in any language (bash, Python, etc.)

### `/build` - Packaging & CI
**Structure:**
- `/build/package/` - Container configs (Dockerfile), package configs (deb, rpm)
- `/build/ci/` - CI configuration (GitHub Actions, CircleCI, etc.)

**Note:** Some CI tools require config in root (like `.github/workflows/`). That's fine - use what works.

### `/test` - Additional Test Data
**When to use:** Integration tests, test fixtures, external test apps

- `/test/data/` or `/test/testdata/` for test fixtures
- Go ignores directories starting with `.` or `_`
- **Note:** Unit tests should live alongside code in `*_test.go` files

### `/docs` - Documentation
**Contents:**
- Design documents
- Architecture diagrams
- User guides
- Changelog, contributing guidelines
- **Note:** GoDoc comments in code are still primary documentation

### `/examples` - Usage Examples
**When to use:** Demonstrating how to use your library or application

- Runnable example code
- Sample configurations
- Tutorial projects

### `/tools` - Supporting Tools
**When to use:** Project-specific development tools

- Code generators
- Build tools specific to this project
- Can import from `/pkg` and `/internal`

### `/assets` - Non-code Assets
**Contents:**
- Images, logos
- Fonts
- Other binary assets for documentation or distribution

## Common Patterns & Best Practices

### Small Projects (Single Application)
```
myapp/
├── main.go              # Simple projects can start here
├── handlers.go
├── models.go
├── go.mod
└── README.md
```
Grow the structure as needed.

### Medium Projects (Multiple Packages)
```
myapp/
├── cmd/myapp/
│   └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   └── database/
├── go.mod
└── README.md
```

### Large Projects (Multiple Applications + Libraries)
Use the full structure shown at the top of this document.

## What NOT to Do

❌ **Don't create `/src`** - This is a Java/Node.js pattern, not Go
❌ **Don't put everything in `/pkg`** - Reserve for truly reusable code
❌ **Don't make deep nested hierarchies** - Go favors flat structures
❌ **Don't create directories you don't need yet** - Start simple

## Modern Variations

Some modern Go projects use simpler layouts:
- No `/pkg` directory (libraries in root or internal)
- `/app` instead of `/cmd` for applications
- Flat structure for smaller projects

**The key:** Be consistent and follow Go community conventions.

## Quick Decision Tree

**Where should my code go?**

1. Is it a `main` package (executable)?
   → `/cmd/appname/`

2. Is it reusable by external projects?
   → `/pkg/libname/` (or root if skipping `/pkg`)

3. Is it private to this project?
   → `/internal/pkg/libname/`

4. Is it specific to one application?
   → `/internal/app/appname/`

5. Is it a test fixture or external test?
   → `/test/`

6. Is it an API definition?
   → `/api/`

7. Is it a build script or tool?
   → `/scripts/` or `/tools/`

## References

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) - Community reference
- [Go 1.4 Release Notes](https://golang.org/doc/go1.4#internalpackages) - Internal packages documentation
- [Effective Go](https://golang.org/doc/effective_go) - Official Go guidelines
