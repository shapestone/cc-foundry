# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Binary distributions (Homebrew, apt, etc.)
- Package manager support
- Team sync features (shared remote repositories)
- Enhanced doctor diagnostics
- Template system for creating new files

## [1.0.0] - 2025-12-01

### Added
- Core install/remove/upgrade functionality for Claude Code files
- Category-based organization system (commands, agents, skills)
- Smart conflict detection with `ccf-[category]-[filename]` naming convention
- Transactional operations with automatic rollback on failure
- Backup and restore system for safe file operations
- State tracking in `~/.claude-code-foundry.json`
- File embedding system using Go's `embed.FS`
- List command to browse available categories and contents
- Self-documenting help system
- Cross-platform support (macOS, Linux, Windows via WSL)

#### Development Category
- **oss-auditor agent**: Comprehensive OSS project repository auditing
- **oss-project-setup skill**: Knowledge framework for OSS project structure and growth
- **makefile-skills-guide skill**: Production-proven Makefile patterns and best practices

### Experimental
- Doctor command for cleaning bloated `~/.claude.json` files
  - Remove stale project data
  - Trim excessive conversation history
  - Verify ccf-managed file integrity
  - Detect and fix corrupted configurations

### Changed
- Install command automatically updates existing files (replaces separate upgrade command)
- Improved output formatting with file type and path display

### Security
- All file operations use atomic writes with rollback
- Backup creation before any destructive operations
- Validation of file paths to prevent directory traversal
- Safe handling of `~/.claude.json` modifications

## [0.1.0] - 2025-11-15 (Pre-release)

### Added
- Initial project structure
- Basic install/remove commands
- Go module setup
- Makefile build system
- Embedded file system for categories

---

## Version History Notes

### Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

### Breaking Changes

Starting with v1.0.0, any breaking changes will be clearly documented here with migration guides.

### Experimental Features

Features marked as "Experimental" may change or be removed without following semantic versioning. Use them with caution and expect potential changes.

[unreleased]: https://github.com/shapestone/claude-code-foundry/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/shapestone/claude-code-foundry/releases/tag/v1.0.0
[0.1.0]: https://github.com/shapestone/claude-code-foundry/releases/tag/v0.1.0
