# Contributing to claude-code-foundry

Thank you for your interest in contributing to claude-code-foundry! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md) (coming soon). Please be respectful and constructive in all interactions.

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details**:
  - OS (macOS, Linux, Windows/WSL)
  - Go version (`go version`)
  - Claude Code version
  - claude-code-foundry version
- **Relevant logs or error messages**

**Security Issues**: If you discover a security vulnerability, please see [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

### Suggesting Features or Enhancements

We welcome feature requests! Please:

1. **Check existing issues** to avoid duplicates
2. **Describe the problem** your feature would solve
3. **Propose a solution** with examples if possible
4. **Consider alternatives** you've thought about
5. **Explain the impact** (who benefits, how much effort)

### Contributing Code

We love pull requests! Here's how to contribute code:

#### Development Setup

1. **Prerequisites**:
   ```bash
   # Go 1.21 or higher
   go version

   # Make (for build automation)
   make --version

   # Git
   git --version
   ```

2. **Fork and clone**:
   ```bash
   # Fork the repository on GitHub, then:
   git clone https://github.com/YOUR-USERNAME/claude-code-foundry.git
   cd claude-code-foundry
   ```

3. **Install dependencies**:
   ```bash
   make deps
   ```

4. **Build the project**:
   ```bash
   make build
   ```

5. **Verify installation**:
   ```bash
   ./build/bin/claude-code-foundry --version
   ```

#### Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Make your changes**:
   - Write clean, readable code
   - Follow Go conventions (`go fmt`, `go vet`)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   # Format code
   make format

   # Run linters
   make lint

   # Run tests
   make test

   # Build to verify compilation
   make build
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Brief description of changes"
   ```

   **Commit message format**:
   - Use present tense ("Add feature" not "Added feature")
   - Be concise but descriptive
   - Reference issues when applicable ("Fix #123: ...")

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**:
   - Describe what changed and why
   - Reference related issues
   - Include screenshots for UI changes
   - Mark as draft if work is in progress

#### Code Style Guidelines

**Go Code**:
- Follow standard Go conventions
- Run `go fmt` before committing
- Use `golangci-lint` if available (fallback to `go vet`)
- Keep functions focused and testable
- Add comments for exported functions and complex logic

**Project Structure**:
```
claude-code-foundry/
├── cmd/                  # CLI application entry points
├── pkg/                  # Public library packages
├── internal/             # Private application code
├── embeddata/            # Embedded category files
│   └── categories/
│       └── [category]/
│           ├── commands/
│           ├── agents/
│           └── skills/
├── Makefile             # Build automation
└── go.mod
```

#### Testing

- Add unit tests for new functionality
- Ensure existing tests pass: `make test`
- Test cross-platform compatibility when possible
- For system file operations, test backup/rollback behavior

**Running tests**:
```bash
# All tests
make test

# Specific package
cd pkg/installer && go test -v

# With coverage
make test-coverage
```

### Contributing Commands, Agents, or Skills

To add new category files (commands/agents/skills):

#### 1. Choose or Create a Category

Categories group related files by purpose:
- `development/` - Development workflow tools
- `deployment/` - Deployment and release tools
- `testing/` - Testing and QA tools
- Or create a new category if appropriate

#### 2. File Format Requirements

All files must follow Claude Code's format conventions:

**Commands** (`commands/*.md`):
```markdown
---
name: command-name
description: Brief description of what this command does
---

# Command Name

[Command content and instructions]
```

**Agents** (`agents/*.md`):
```markdown
---
name: agent-name
description: Brief description of agent's purpose and capabilities
tools: [Bash, Glob, Grep, Read, Write, Edit, WebFetch, TodoWrite]
---

# Agent Name

## Purpose
[Detailed agent description]

[Agent instructions and workflow]
```

**Skills** (`skills/*.md`):
```markdown
---
name: skill-name
description: Brief description of skill's knowledge domain
---

# Skill Name

[Knowledge content and guidance]
```

**Requirements**:
- `name`: Kebab-case identifier matching filename (e.g., `my-skill` for `my-skill.md`)
- `description`: One-sentence summary for listings
- `tools`: Array of tools for agents only
- Frontmatter uses YAML format with `---` delimiters

#### 3. File Placement

```bash
# Place file in appropriate category directory
embeddata/categories/[category-name]/[type]/[filename].md

# Examples:
embeddata/categories/development/agents/code-reviewer.md
embeddata/categories/deployment/commands/deploy-prod.md
embeddata/categories/testing/skills/playwright-guide.md
```

#### 4. Testing Your Addition

```bash
# Rebuild with embedded files
make build

# List to verify your file appears
./build/bin/claude-code-foundry list all

# Test installation
./build/bin/claude-code-foundry install [category]
```

#### 5. Documentation

In your pull request, include:
- **Description**: What the file does
- **Use cases**: When to use it
- **Examples**: Sample invocations or workflows
- **Dependencies**: Any required tools or setup

### Documentation Improvements

Documentation is just as important as code! You can help by:

- Fixing typos or unclear explanations
- Adding examples or use cases
- Improving README organization
- Creating guides or tutorials
- Updating outdated information

## Pull Request Process

1. **Ensure all tests pass** (`make test`)
2. **Update documentation** if you changed functionality
3. **Add yourself to contributors** if this is your first PR
4. **Keep PRs focused** - one feature or fix per PR
5. **Respond to feedback** - maintainers may request changes
6. **Be patient** - reviews may take a few days

### What to Expect

- **Initial response**: Within 3-5 business days
- **Code review**: Maintainers will review for:
  - Code quality and style
  - Test coverage
  - Documentation completeness
  - Backward compatibility
- **Merge**: After approval and passing CI checks

## Development Tips

### Useful Make Targets

```bash
make help           # Show all available targets
make build          # Build the binary
make test           # Run tests
make lint           # Run linters
make format         # Format code
make clean          # Clean build artifacts
make install        # Install binary to $GOPATH/bin (use GLOBAL=1 for /usr/local/bin)
```

### Debugging

```bash
# Build with debug info
go build -gcflags="all=-N -l" -o build/bin/claude-code-foundry cmd/claude-code-foundry/main.go

# Run with verbose output
./build/bin/claude-code-foundry --debug [command]

# Check state file
cat ~/.claude-code-foundry.json | jq .
```

### Testing File Operations

When testing install/remove operations:

```bash
# Use a test directory instead of ~/.claude
export CCF_TEST_MODE=true

# Or manually specify paths for testing
# (Feature request: Add test mode flag)
```

## Project Governance

Currently, this project is maintained by Shapestone. Major decisions are made by core maintainers. As the project grows, we may adopt a more formal governance model.

## Questions?

- **General questions**: Open a [GitHub Discussion](https://github.com/shapestone/claude-code-foundry/discussions)
- **Bug reports**: Open a [GitHub Issue](https://github.com/shapestone/claude-code-foundry/issues)
- **Security issues**: See [SECURITY.md](SECURITY.md)

## Recognition

Contributors will be recognized in:
- Git commit history
- Release notes (for significant contributions)
- GitHub contributors page

Thank you for contributing to claude-code-foundry!

---

## License

By contributing to claude-code-foundry, you agree that your contributions will be licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
