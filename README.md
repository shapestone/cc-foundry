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

### macOS Gatekeeper Security

On macOS, you may see this warning when running the binary:

> "claude-code-foundry" cannot be opened because it is from an unidentified developer.

This is expected for locally-built binaries. The `make build` command automatically ad-hoc signs the binary to prevent this issue.

**If you still see the warning:**

**Option 1: Quick Fix (Recommended)**
```bash
# Ad-hoc sign the binary
codesign -s - build/bin/claude-code-foundry
```

**Option 2: Manual Override**
1. Right-click the binary and select **Open**
2. Click **Open** in the security dialog
3. OR go to **System Settings** → **Privacy & Security** → Click **"Open Anyway"**

**What is ad-hoc signing?**
- Local-only code signing without Apple Developer certificate
- Marks the binary as trusted by you
- Prevents Gatekeeper from blocking it
- Safe for binaries you built yourself
- **Note:** Does not allow distribution to other users

The Makefile automatically signs all macOS builds (`make build`, `make install`, `make build-darwin-*`).

---

## Usage

### Interactive Mode

Simply run the command without arguments to launch the interactive menu:

```bash
claude-code-foundry
```

You'll see a main menu with arrow-key navigation:

```
🔧 claude-code-foundry - Manage Claude Code files

What would you like to do?
❯ Show directory structure
  List available files
  Install files
  Remove files
  Doctor (verify & repair)
  Version information
  Help
  Exit
```

Navigate with **↑/↓ arrows**, select with **Enter**, cancel with **Ctrl+C**.

---

### Main Menu Options

#### 1. Show Directory Structure

Displays your Claude Code directory structure and installed files:

```
📁 Claude Code Directory Structure

User-level (~/.claude/):
  commands/  (5 files)
  agents/    (3 files)
  skills/    (2 skills)

Project-level (.claude/):
  ✗ Directory does not exist

📦 Installed Files (managed by foundry)

  development: 2 commands, 1 agent, 2 skills

  Total: 5 files installed
```

#### 2. List Available Files

Browse available categories and their contents:

1. Select category from list:
   ```
   Select category to list
   ❯ development (2 commands, 1 agent, 2 skills)
     deployment (1 command, 1 agent, 0 skills)
     ← Back to main menu
   ```

2. View files in category:
   ```
   Category: development

   Commands:
     - implement.md
     - review.md

   Agents:
     - oss-auditor.md

   Skills:
     - makefile-skills-guide.md
     - oss-project-setup.md
   ```

#### 3. Install Files

Install commands, agents, and skills to your system:

1. **Select category**:
   ```
   Select category to install
   ❯ All categories
     development (2 commands, 1 agent, 2 skills)
     deployment (1 command, 1 agent, 0 skills)
     ← Back to main menu
   ```

2. **Choose location**:
   ```
   Choose location
   ❯ 1. Project (.claude/)
     2. Personal (~/.claude/)
   ```

3. **Preview changes**:
   ```
   Preview: development [user-level (~/.claude/)]

     + command: ccf-development-implement.md → ~/.claude/commands/ccf-development-implement.md
     ↻ agent: ccf-development-oss-auditor.md → ~/.claude/agents/ccf-development-oss-auditor.md (will update)
     · skill: ccf-development-makefile-skills-guide/SKILL.md → ~/.claude/skills/... (unchanged)

   Summary: 1 to install, 1 to update, 1 unchanged
   ```

4. **Confirm installation**:
   ```
   Proceed with installation?
   ❯ Yes, proceed
     No, cancel
   ```

**Symbols:**
- `+` New installation
- `↻` Update (content changed)
- `·` Skip (unchanged)

#### 4. Remove Files

Remove installed files:

1. **Select category**:
   ```
   Select category to remove
   ❯ All categories
     development (2 commands, 1 agent, 2 skills)
     ← Back to main menu
   ```

2. **Choose location**:
   ```
   Choose location
   ❯ 1. Project (.claude/)
     2. Personal (~/.claude/)
   ```

3. **Preview removal**:
   ```
   Preview: Remove category development [user-level (~/.claude/)]

     - command: ~/.claude/commands/ccf-development-implement.md
     - agent: ~/.claude/agents/ccf-development-oss-auditor.md
     - skill: ~/.claude/skills/ccf-development-makefile-skills-guide

   Summary: 3 files will be removed
   ```

4. **Confirm removal**:
   ```
   Proceed with removal?
   ❯ Yes, remove
     No, cancel
   ```

#### 5. Doctor (Verify & Repair)

Run diagnostics and fix issues:

```
🏥 Running doctor diagnostics...

Checking Claude Code configuration...
Checking installed file integrity...
Detecting conflicts...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Health Report
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Files checked: 5
⚠️  Warnings: 2
Modified files: 1
Orphaned files: 1

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Issues Found:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  [development] Modified file detected: ~/.claude/commands/ccf-development-implement.md (hash mismatch)

⚠️  [orphaned] Orphaned foundry file: ~/.claude/commands/ccf-old-command.md (not tracked in state)
   (can be fixed)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1 issue(s) can be automatically fixed.

Would you like to fix these issues?
❯ Yes, fix all issues
  No, leave as is
```

The doctor command checks:
- **~/.claude.json validity**: Verifies config file exists and is valid JSON
- **File integrity**: Compares installed file hashes to detect modifications
- **Conflict detection**: Finds orphaned ccf- files not tracked in state
- **Auto-repair**: Offers to fix detected issues

#### 6. Version Information

Shows the current version:

```
claude-code-foundry v1.0.0
```

#### 7. Help

Displays usage information and file structure details.

---

### Installation Locations

**Personal (User-level)** - `~/.claude/`:
- **Scope**: Available across all projects
- **Use when**: You want these files available everywhere you use Claude Code
- **Example**: General-purpose agents/skills you use daily

**Project (Project-level)** - `.claude/`:
- **Scope**: Only available in this specific project
- **Version control**: Can be committed to git and shared with team
- **Use when**: Project-specific configurations or team-shared commands

**Directory Structure:**

```
~/.claude/                          OR    .claude/
├── commands/                             ├── commands/
│   └── ccf-[category]-[name].md          │   └── ccf-[category]-[name].md
├── agents/                               ├── agents/
│   └── ccf-[category]-[name].md          │   └── ccf-[category]-[name].md
└── skills/                               └── skills/
    └── ccf-[category]-[name]/                └── ccf-[category]-[name]/
        └── SKILL.md                              └── SKILL.md
```

**Naming Convention**: All files use `ccf-[category]-[name]` format to prevent conflicts:
- `development/commands/implement.md` → `ccf-development-implement.md`
- `development/agents/oss-auditor.md` → `ccf-development-oss-auditor.md`
- `development/skills/oss-project-setup.md` → `ccf-development-oss-project-setup/SKILL.md`

---

### Tips

- **Navigate anywhere**: Use ← Back option or Ctrl+C to return to previous menu
- **Preview before changes**: All operations show preview before making changes
- **Safe operations**: Automatic backups created, rollback on failure
- **Update files**: Re-run install to update to latest versions (unchanged files skipped)
- **Check health regularly**: Run doctor to verify installation integrity

---

### Future: Non-Interactive Mode

Command-line arguments for scripting will be added in a future release, allowing:

```bash
# Not yet available - coming in future version
claude-code-foundry install --category development --location user --yes
```

For now, the interactive menu provides the full feature set.

---

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

#### 📚 makefile-skills-guide Skill
**Purpose**: Comprehensive reference for Makefile build automation

Production-proven patterns and best practices for creating Makefiles for multi-language projects, cross-platform builds, and development workflows.

**Key Features**:
- Full-stack integration (Go + npm/Node.js)
- Cross-platform compilation (Linux, macOS, Windows)
- Process management (start/stop/restart with PID tracking)
- Testing integration (unit, E2E with Playwright)
- Development workflows (hot reload, watch mode)
- Self-documenting help system

**Coverage**:
- Basic and advanced Makefile concepts
- Common production patterns
- Full-stack build orchestration
- Dependency management
- Linting, formatting, and validation
- Cross-platform compatibility
- Debugging techniques

**Use Cases**:
- "Create a Makefile for my Go + React project"
- "How do I manage background processes with Make?"
- "Set up cross-platform builds"
- "Integrate Playwright E2E tests into Make"

**Installation**:
```bash
claude-code-foundry install development skills
```

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
│           ├── makefile-skills-guide.md
│           └── oss-project-setup.md
├── cmd/                     # CLI application code
├── pkg/                     # Core library packages
├── Makefile                 # Build and install
├── go.mod
├── go.sum
├── LICENSE                  # Apache 2.0
└── README.md               # This file
```

### File Format Conventions

All commands, agents, and skills must follow Claude Code's file format conventions:

#### Commands (`commands/*.md`)
```markdown
---
name: command-name
description: Brief description of what this command does
---

# Command Name

[Command content and instructions]
```

#### Agents (`agents/*.md`)
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

**Required fields:**
- `name`: Kebab-case identifier matching the filename (e.g., `oss-auditor` for `oss-auditor.md`)
- `description`: One-sentence summary (used in listings)
- `tools`: Array of tools the agent can use (agents only)

#### Skills (`skills/*.md`)
```markdown
---
name: skill-name
description: Brief description of skill's knowledge domain
---

# Skill Name

[Knowledge content and guidance]
```

### Adding New Categories

To contribute new categories or files:

1. Create category folder: `categories/[category-name]/`
2. Add subdirectories: `commands/`, `agents/`, `skills/`
3. Create markdown files following format conventions above:
   - Use kebab-case for filenames (e.g., `my-command.md`)
   - Ensure `name` field in frontmatter matches filename
   - Provide clear, concise descriptions
   - For agents, specify required tools
4. Verify formatting:
   - Frontmatter uses YAML format with `---` delimiters
   - Required fields are present
   - Content is well-structured and documented
5. Submit pull request with:
   - Clear description of what the files do
   - Use cases and examples
   - Any dependencies or requirements

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
      "installed_path": "~/.claude/commands/ccf-development-deploy-to-production.md",
      "hash": "abc123...",
      "installed_at": "2025-12-01T10:00:00Z"
    },
    {
      "category": "development",
      "type": "skill",
      "file": "oss-project-setup.md",
      "installed_path": "~/.claude/skills/ccf-development-oss-project-setup/SKILL.md",
      "hash": "def456...",
      "installed_at": "2025-12-01T10:01:00Z"
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

If Claude Code directories don't exist, the tool will create them automatically:

**User-level**:
```
~/.claude/commands/
~/.claude/agents/
~/.claude/skills/
```

**Project-level**:
```
.claude/commands/
.claude/agents/
.claude/skills/
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
- ✅ Core install/remove functionality (install handles updates automatically)
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
