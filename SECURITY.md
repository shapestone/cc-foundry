# Security Policy

## Supported Versions

We actively support the following versions of claude-code-foundry with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Security Context

claude-code-foundry modifies critical system files and directories:

- `~/.claude.json` - Claude Code's main configuration file
- `~/.claude/` - User-level Claude Code file directories
- `.claude/` - Project-level Claude Code file directories (if using `--project`)
- `~/.claude-code-foundry.json` - State tracking file
- `~/.claude-code-foundry-backups/` - Backup directory

**Potential Security Concerns:**
- File system operations (install/remove/backup)
- Path traversal vulnerabilities
- Unsafe file content handling
- State file tampering
- Backup file exposure

## Security Measures

### Current Protections

1. **Atomic Operations**:
   - All file operations are transactional with automatic rollback
   - Backups created before any destructive operations

2. **Path Validation**:
   - All file paths are validated to prevent directory traversal
   - Operations restricted to designated Claude Code directories

3. **Safe File Handling**:
   - No execution of embedded file content
   - Markdown files are written as-is without interpretation

4. **State Integrity**:
   - State file uses JSON for structure validation
   - Hash verification for installed files

### Known Limitations

- State file (`~/.claude-code-foundry.json`) is user-readable/writable
- Backup files stored in user's home directory
- No encryption of backup files
- Doctor command modifies `~/.claude.json` directly (experimental)

## Reporting a Vulnerability

**DO NOT** open a public GitHub issue for security vulnerabilities.

### How to Report

Please report security vulnerabilities by emailing:

**security@shapestone.io**

Include in your report:

1. **Description** of the vulnerability
2. **Steps to reproduce** the issue
3. **Potential impact** (what could an attacker do?)
4. **Suggested fix** (if you have one)
5. **Your contact information** for follow-up

### What to Expect

- **Initial Response**: Within 48 hours (business days)
- **Status Update**: Within 5 business days
- **Fix Timeline**: Varies by severity
  - Critical: 7-14 days
  - High: 14-30 days
  - Medium: 30-60 days
  - Low: 60-90 days or next release

### Our Process

1. **Acknowledgment**: We confirm receipt of your report
2. **Assessment**: We validate and assess severity
3. **Development**: We develop and test a fix
4. **Disclosure**: We coordinate disclosure with you
5. **Release**: We release the fix and publish a security advisory
6. **Recognition**: We credit you in the advisory (unless you prefer anonymity)

## Security Disclosure Policy

### Coordinated Disclosure

We practice coordinated disclosure:

- **Timeline**: 90 days from initial report to public disclosure
- **Extension**: May request extension if fix is complex
- **Early Disclosure**: May disclose earlier if actively exploited

### Public Disclosure

After a fix is released, we will:

1. Publish a GitHub Security Advisory
2. Update CHANGELOG.md with security fix details
3. Tag a new release with the fix
4. Credit the reporter (unless anonymity requested)

## Security Best Practices for Users

### When Using claude-code-foundry

1. **Verify Source**:
   ```bash
   # Build from source or verify binary signature
   git clone https://github.com/shapestone/claude-code-foundry.git
   cd claude-code-foundry
   make build
   ```

2. **Review Before Installing**:
   ```bash
   # List what will be installed
   claude-code-foundry list all

   # Review specific category files before installation
   # (Files are embedded in the binary)
   ```

3. **Use Dry Run (when available)**:
   ```bash
   # For doctor command
   claude-code-foundry doctor --dry-run
   ```

4. **Regular Backups**:
   - claude-code-foundry creates automatic backups
   - Consider additional backups of `~/.claude.json`:
     ```bash
     cp ~/.claude.json ~/.claude.json.backup-$(date +%Y%m%d)
     ```

5. **Review State File**:
   ```bash
   # Check what's installed
   cat ~/.claude-code-foundry.json | jq .
   ```

### What to Avoid

- **Don't run with elevated privileges** - Normal user permissions are sufficient
- **Don't modify state file manually** - Use the CLI commands
- **Don't install from untrusted sources** - Build from official repository
- **Don't ignore warnings** - The doctor command warns before making changes

## Scope

### In Scope

Security issues in:
- File operations (install/remove/backup/restore)
- State management
- Path validation and sanitization
- Doctor command operations
- CLI input handling

### Out of Scope

- Issues in Claude Code itself (report to Anthropic)
- Issues in embedded content (commands/agents/skills) - These are user-executed in Claude Code context
- Third-party dependencies (report to upstream)
- Social engineering or phishing attacks
- Physical access scenarios

## Dependencies

We monitor our Go dependencies for known vulnerabilities:

```bash
# Check dependencies for vulnerabilities
go list -json -deps | nancy sleuth
```

Current direct dependencies are minimal (standard library focused).

## Security Updates

Subscribe to security updates:

- **GitHub Watch**: Click "Watch" → "Custom" → "Security alerts"
- **Releases**: Watch release tags for security fixes
- **RSS Feed**: Subscribe to GitHub releases feed

## Questions?

For general security questions (not vulnerability reports):
- Open a [GitHub Discussion](https://github.com/shapestone/claude-code-foundry/discussions)
- Tag with "security" label

For vulnerability reports:
- Email: security@shapestone.io

---

## Acknowledgments

We appreciate security researchers who help keep claude-code-foundry safe:

- *Your name could be here!*

Thank you for helping keep claude-code-foundry secure!
