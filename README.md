# claude-code-foundry

A community tool for managing and standardizing Claude Code files (commands, agents, and skills) across development teams.

## ⚠️ Trademark Notice

**claude-code-foundry** is an independent, community-created tool and is **NOT** officially affiliated with, endorsed by, or sponsored by Anthropic PBC.

- "Claude" and "Claude Code" are trademarks of Anthropic PBC
- This tool is designed to work with Claude Code but is developed and maintained independently
- Use of these names is purely descriptive to indicate compatibility

For official Claude Code documentation and support, visit [Anthropic's official documentation](https://docs.claude.com/).

---

## Overview

claude-code-foundry helps teams maintain consistent Claude Code configurations by providing a centralized repository and CLI tool for managing:

- **Commands** - Custom slash commands for repeated workflows
- **Agents** - Specialized AI agents for specific tasks
- **Skills** - Reusable capabilities and knowledge bases

### Problem Solved

- ✅ **Consistency** - Everyone uses the same proven commands, agents, and skills
- ✅ **Discoverability** - Easily find and install useful configurations from your team
- ✅ **Version Control** - Track changes and updates to shared configurations
- ✅ **Onboarding** - New team members get up and running quickly

---

## Features

### Core Functionality
- 📦 **Category-based organization** - Group related commands/agents/skills by purpose (e.g., development, deployment, testing)
- 🔄 **Install/Upgrade/Remove** - Manage files at category or type level
- 🔍 **Smart conflict detection** - Prevents overwriting user files with unique naming
- 🛡️ **Transactional operations** - Automatic rollback on failure with backup/restore
- 📊 **List and discover** - Browse available categories and their contents

### Experimental Features
- 🩺 **Doctor command** - Clean bloated `~/.claude.json` files and verify installation integrity
  - Remove stale project data that causes Claude Code slowdowns
  - Trim excessive conversation history
  - Verify ccf-managed file integrity
  - Detect and fix corrupted configurations

> ⚠️ **Warning**: The `doctor` command modifies system files. Always review changes before confirming. Backups are created automatically.

---

## Installation

### Prerequisites
- Go 1.21 or higher
- Claude Code installed on your system

### Build from Source

```bash
# Clone the repository
git clone https://github.com/shapestone/claude-code-foundry.git
cd claude-code-foundry

# Build and install
make install
```

This installs the `claude-code-foundry` binary to your system.

### Verify Installation

```bash
claude-code-foundry --version
```

---

## Usage

### Getting Started

#### 1. List Available Categories

```bash
# Show all categories with their contents
claude-code-foundry list all

# Show specific category
claude-code-foundry list development
```

Example output:
```
Available Categories:

📁 development/
  Agents:
    - oss-auditor.md
  Skills:
    - oss-project-setup.md
```

#### 2. Install Files

```bash
# Install everything
claude-code-foundry install all

# Install specific category
claude-code-foundry install development

# Install specific type from category
claude-code-foundry install development commands
claude-code-foundry install development agents
claude-code-foundry install development skills
```

Files are installed to:
- `~/.claudecode/commands/` (or `~/.config/claude/commands/` on Linux)
- `~/.claudecode/agents/`
- `~/.claudecode/skills/`

**Naming Convention**: All installed files use the `ccf-[category]-[filename]` format to prevent conflicts:
- `development/commands/deploy.md` → `~/.claudecode/commands/ccf-development-deploy.md`
- `deployment/commands/deploy.md` → `~/.claudecode/commands/ccf-deployment-deploy.md`

#### 3. Upgrade Files

```bash
# Upgrade everything installed by foundry
claude-code-foundry upgrade all

# Upgrade specific category
claude-code-foundry upgrade development

# Upgrade specific type
claude-code-foundry upgrade development commands
```

The tool will:
- Compare file hashes to detect changes
- Prompt if you've modified files locally
- Show diffs before overwriting

#### 4. Remove Files

```bash
# Remove everything foundry installed
claude-code-foundry remove all

# Remove specific category
claude-code-foundry remove development

# Remove specific type
claude-code-foundry remove development commands
```

Only removes files that foundry installed (tracked in state file).

#### 5. Health Check (Experimental)

```bash
# Run all checks
claude-code-foundry doctor

# Check only ~/.claude.json
claude-code-foundry doctor config

# Check only ccf-managed files
claude-code-foundry doctor files

# Preview without making changes
claude-code-foundry doctor --dry-run
```

The doctor command analyzes:
- **Config issues**: Bloated `~/.claude.json` files causing performance problems
- **File integrity**: Missing, orphaned, or modified ccf-managed files
- **Smart cleanup**: Removes stale projects and excessive conversation history

Example output:
```
Analysis of ~/.claude.json:
- File size: 65.2 MB
- Total projects: 147
- Conversation history: 58.3 MB (89%)

Detected issues:
✓ 89 projects not accessed in 120+ days (48.2 MB)
✓ 23 projects with excessive history (10.1 MB)
✓ 5 orphaned projects (directories deleted) (2.3 MB)

Checking foundry-managed files...
✓ 23 files installed and verified
⚠ 2 files modified locally
⚠ 1 orphaned file found

Suggested cleanup: 60.6 MB → 4.6 MB (93% reduction)

What would you like to do?
1. Clean ~/.claude.json (recommended)
2. Fix foundry issues
3. Both
4. Exit
```

### Interactive Help

```bash
# Launch interactive guide
claude-code-foundry help

# Quick reference
claude-code-foundry --help
claude-code-foundry install --help
```

---

## Available Agents & Skills

### Development Category

#### 🤖 oss-auditor Agent
**Purpose**: Comprehensive open source project repository auditing

Analyzes repositories to assess documentation completeness, identify gaps, evaluate maturity, and provide actionable recommendations.

**Key Features**:
- Scans for standard OSS files (LICENSE, README, CONTRIBUTING, etc.)
- Evaluates documentation quality and freshness
- Analyzes git history and contributor patterns
- Performs gap analysis with prioritized recommendations
- Determines appropriate project tier (Minimal/Standard/Mature)
- Generates detailed audit reports

**Use Cases**:
- "Audit this repository's open source structure"
- "Check if our security documentation is adequate for an auth library"
- "We're at 50 contributors now - have we outgrown our current structure?"
- "We're about to open source this - what do we need to add?"

**Installation**:
```bash
claude-code-foundry install development agents
```

After installation, invoke the agent in conversation to analyze any repository.

---

#### 📚 oss-project-setup Skill
**Purpose**: Knowledge framework for open source project setup and growth

Provides comprehensive guidance on structuring open source projects based on characteristics, complexity, and maturity. Takes a diagnostic-first approach.

**Key Concepts**:
- **Tier 1 (Minimal)**: New/simple projects, solo/small team
- **Tier 2 (Standard)**: Active contributions, production use
- **Tier 3 (Mature)**: Large community, multi-org, critical infrastructure

**Coverage**:
- Assessment framework for new and existing projects
- Progressive documentation strategy
- Domain-specific adaptations (libraries, CLI tools, frameworks, etc.)
- File descriptions and best practices
- Common anti-patterns to avoid
- Scaling triggers for growth

**Installation**:
```bash
claude-code-foundry install development skills
```

The skill is automatically available as context for conversations and is used by the oss-auditor agent.

---

### Using Together

The **oss-auditor agent** and **oss-project-setup skill** work together:

1. The **skill** provides the knowledge framework:
   - Tier definitions
   - Best practices
   - File templates
   - Assessment criteria

2. The **agent** applies that framework:
   - Analyzes actual repositories
   - Generates audit reports
   - Prioritizes recommendations
   - Provides actionable next steps

**Recommended workflow**:
```bash
# Install both
claude-code-foundry install development

# In conversation, ask the oss-auditor agent to audit your project
# The agent will use the skill's framework to provide recommendations
```

---

## Repository Structure

```
claude-code-foundry/
├── categories/              # Managed files organized by purpose
│   └── development/
│       ├── agents/
│       │   └── oss-auditor.md
│       └── skills/
│           └── oss-project-setup.md
├── cmd/                     # CLI application code
├── pkg/                     # Core library packages
├── Makefile                 # Build and install
├── go.mod
├── go.sum
├── LICENSE                  # Apache 2.0
└── README.md               # This file
```

### Adding New Categories

To contribute new categories or files:

1. Create category folder: `categories/[category-name]/`
2. Add subdirectories: `commands/`, `agents/`, `skills/`
3. Add markdown files with descriptive names
4. Submit pull request

Files will automatically be embedded in the next binary build.

---

## How It Works

### State Management

The CLI tracks installations in `~/.claude-code-foundry.json`:

```json
{
  "version": "1.0.0",
  "installations": [
    {
      "category": "development",
      "type": "command",
      "file": "deploy-to-production.md",
      "installed_path": "~/.claudecode/commands/ccf-development-deploy-to-production.md",
      "hash": "abc123...",
      "installed_at": "2025-12-01T10:00:00Z"
    }
  ]
}
```

### Backup & Rollback

Every operation creates a timestamped backup:
- Location: `~/.claude-code-foundry-backups/[timestamp]/`
- Contains: Complete snapshot of affected Claude Code directories
- Lifecycle: Automatically deleted after successful operation or rollback
- Rollback: Automatic on any failure during transactional operations

### File Embedding

All category files are embedded directly in the binary using Go's `embed.FS`:
- No network requests needed
- Works offline
- Single binary distribution
- Version controlled with the CLI code

---

## Troubleshooting

### Claude Code Directory Not Found

If Claude Code directories don't exist, the tool will prompt you:

```
Claude Code directories not found at:
  ~/.claudecode/
  ~/.config/claude/

Would you like to create them? (y/n)
```

### File Conflicts

The `ccf-[category]-[filename]` naming convention prevents conflicts with:
- User's existing files
- Files from different categories
- Manual installations

### Performance Issues

If Claude Code feels slow:
1. Run `claude-code-foundry doctor` to analyze your setup
2. The most common issue is a bloated `~/.claude.json` file (50MB+)
3. Follow the doctor's recommendations to clean up stale data

See: [Claude Code Performance Issues](https://github.com/anthropics/claude-code/issues/5024)

---

## Contributing

We welcome contributions! Here's how you can help:

### Adding Commands/Agents/Skills

1. Fork the repository
2. Create a new category or add to existing one
3. Follow naming conventions (descriptive, human-friendly names)
4. Submit a pull request with:
   - Clear description of what the files do
   - Use cases and examples
   - Any dependencies or requirements

### Reporting Issues

- Use GitHub Issues for bugs or feature requests
- Include your OS, Go version, and Claude Code version
- Provide steps to reproduce for bugs

### Development

```bash
# Run tests
make test

# Build locally
make build

# Install development version
make install
```

---

## Roadmap

### v1.0.0 (Current)
- ✅ Core install/upgrade/remove functionality
- ✅ Category-based organization
- ✅ Transactional operations with rollback
- ✅ Doctor command (experimental)

### Potential Future Releases
- [ ] Binary distributions (no build required)
- [ ] Package manager support (Homebrew, apt, etc.)
- [ ] Team sync features (shared remote repositories)
- [ ] Enhanced doctor diagnostics
- [ ] Template system for creating new files

---

## Known Issues & Limitations

### Experimental Features
The `doctor` command is experimental and should be used with caution:
- Always review the proposed changes before confirming
- Backups are created automatically, but verify critical data separately
- Report any issues or unexpected behavior on GitHub

### Claude Code Compatibility
- Tested with Claude Code v1.0.x - v2.0.x
- Directory structure may vary between versions
- The tool attempts to detect the correct location automatically

### Performance
- First install may prompt for directory creation
- Large categories (100+ files) may take a few seconds to install
- Doctor analysis on very large `~/.claude.json` files (100MB+) may take 10-30 seconds

---

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

### What This Means
- ✅ Free to use, modify, and distribute
- ✅ Can be used in commercial projects
- ✅ Includes patent protection
- ⚠️ Must include copyright notice and license
- ⚠️ Must state significant changes made

---

## Acknowledgments

- Built for the Claude Code community
- Inspired by common pain points in team collaboration
- Special thanks to all contributors

---

## Legal Disclaimer

**IMPORTANT - PLEASE READ:**

This tool modifies system files and Claude Code configurations. While we implement safety measures (backups, rollback, validation), you use this tool at your own risk.

**The developers and contributors:**
- Provide this software "AS IS" without warranty of any kind
- Are not responsible for data loss, corruption, or system issues
- Recommend testing in non-production environments first
- Suggest regular backups of your Claude Code configurations

**Specific Warnings:**
- The `doctor` command modifies `~/.claude.json` which is critical for Claude Code operation
- File operations are transactional but system crashes during operations could cause issues
- Always ensure you have backups before running destructive operations
- Test with `--dry-run` flag first when available

**Not Affiliated with Anthropic:**
This is an independent, community tool. For official support:
- Visit [Anthropic's Documentation](https://docs.claude.com/)
- Contact Anthropic Support for Claude Code issues
- Check [Anthropic's Status Page](https://status.anthropic.com/) for service issues

---

## Support

### Getting Help
- 📖 Documentation: This README and `claude-code-foundry help`
- 🐛 Bug Reports: [GitHub Issues](https://github.com/shapestone/claude-code-foundry/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/shapestone/claude-code-foundry/discussions)

### Useful Resources
- [Claude Code Official Docs](https://docs.claude.com/)
- [Claude Code GitHub](https://github.com/anthropics/claude-code)
- [Claude Code Best Practices](https://www.anthropic.com/engineering/claude-code-best-practices)

---

**Made with ❤️ by the community, for the community**
