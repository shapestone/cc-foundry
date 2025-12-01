---
name: open-source-project-setup
description: Assesses open source projects (new or existing) and recommends appropriate documentation and structure based on project characteristics, complexity, and maturity. Takes a diagnostic-first approach to guide projects from minimal setup through mature governance.
---

# Open Source Project Setup Skill

## Overview
This skill assesses open source projects (new or existing) and recommends appropriate documentation and structure based on project characteristics, complexity, and maturity. It takes a diagnostic-first approach, understanding context before prescribing solutions.

## When to Use This Skill
- Starting a new open source project
- Transitioning a private repository to open source
- Auditing existing open source project documentation
- Scaling up project structure as community grows
- Identifying documentation gaps or mismatches

## Assessment Framework

### For New Projects

The skill asks diagnostic questions to understand project context:

#### Project Type
- **Library:** Reusable code components
- **Framework:** Application scaffolding and patterns
- **CLI Tool:** Command-line application
- **Web Application:** Deployable web service
- **Language/Compiler:** Programming language implementation
- **Plugin/Extension:** Extends existing platform
- **Other:** Custom categorization

#### Technical Characteristics
- **Language/ecosystem:** Determines .gitignore, CI/CD, package files
- **Architectural complexity:** Single component vs. multi-layered system
- **API surface area:** Small focused API vs. extensive public interface
- **Performance requirements:** Latency-sensitive, throughput-critical, or standard
- **Security sensitivity:** Cryptography, authentication, data handling

#### Stability & Commitment
- **Experimental:** Exploring ideas, frequent breaking changes expected
- **Alpha/Beta:** Functional but evolving, limited stability guarantees
- **Production:** Committed to stability, semantic versioning, migration guides
- **Mature:** Long-term support, formal deprecation policy

#### Expected Community
- **Solo project:** Personal tool or learning project
- **Small team:** 2-5 core contributors
- **Open collaboration:** Accepting external contributions
- **Organization-backed:** Company or foundation support
- **Multi-stakeholder:** Multiple organizations with competing interests

### For Existing Projects

The skill analyzes current state to identify gaps and recommend improvements:

#### Documentation Audit
- **Present files:** What documentation already exists?
- **Missing standard files:** README, LICENSE, CONTRIBUTING, etc.
- **Outdated content:** Last update dates, version mismatches
- **Quality assessment:** Completeness, clarity, accuracy

#### Repository Analysis
- **Contributor metrics:** Count, frequency, diversity
- **Issue patterns:** Common questions, bug report quality, feature request volume
- **PR patterns:** Rejection rate, review bottlenecks, contribution friction
- **Commit history:** Breaking changes frequency, release cadence
- **Code structure:** Monorepo, multi-package, complexity indicators

#### Community Health Signals
- **Response times:** How quickly are issues/PRs addressed?
- **Conflict indicators:** Heated discussions, unresolved disputes
- **Onboarding friction:** First-time contributor experience
- **Documentation requests:** Frequency of "where do I start?" questions
- **Maintainer bandwidth:** Signs of burnout or capacity issues

#### Maturity Indicators
- **Age:** How long has project been public?
- **Adoption:** Stars, forks, dependents, downloads
- **Stability:** Breaking change frequency
- **Ecosystem impact:** Downstream dependencies affected by changes
- **External pressure:** Security audits, compliance requirements, industry scrutiny

## Recommendation Tiers

Based on assessment, the skill recommends one of three structural tiers:

### Tier 1: Minimal (New or Simple Projects)

**When to use:**
- New projects just going public
- Solo or small team projects
- Simple, single-purpose tools
- Experimental or learning projects
- Low external contribution expectations

**Required files:**
```
project-root/
├── README.md
├── LICENSE
├── .gitignore
└── [source code]
```

**README.md contents:**
- One-paragraph description
- Installation instructions
- Basic usage example
- License statement

**Next steps trigger:**
- First external contribution → Add CONTRIBUTING.md
- First bug report → Add issue template
- Security concern raised → Add SECURITY.md

### Tier 2: Standard (Growing Projects)

**When to use:**
- Active external contributions
- Multiple maintainers
- Production use by others
- API stability commitments
- Need for contribution guidelines

**Required files:**
```
project-root/
├── README.md
├── LICENSE
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── CHANGELOG.md
├── .gitignore
├── .github/
│   ├── ISSUE_TEMPLATE.md
│   └── PULL_REQUEST_TEMPLATE.md
└── [source code]
```

**Enhanced documentation:**
- Detailed README with examples
- Clear contribution workflow
- Issue and PR templates
- Version history tracking
- Community behavior standards

**Scale-up triggers:**
- 20+ regular contributors → Consider governance
- Complex architecture questions → Add ARCHITECTURE.md
- Feature request debates → Add ROADMAP.md
- Multiple organizations involved → Add GOVERNANCE.md

### Tier 3: Mature (Complex or High-Impact Projects)

**When to use:**
- Large contributor base
- Multiple organizations involved
- Critical infrastructure or dependencies
- Formal stability guarantees
- Complex technical architecture
- Security-sensitive domain
- Breaking changes require coordination

**Complete structure:**
```
project-root/
├── README.md
├── LICENSE
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── CHANGELOG.md
├── SECURITY.md
├── GOVERNANCE.md
├── ARCHITECTURE.md
├── ROADMAP.md
├── SUPPORT.md
├── AUTHORS
├── CONTRIBUTORS.md
├── .gitignore
├── .github/
│   ├── workflows/
│   │   ├── ci.yml
│   │   ├── release.yml
│   │   └── security.yml
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   ├── feature_request.md
│   │   └── security.md
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── FUNDING.yml
├── docs/
│   ├── index.md
│   ├── getting-started/
│   ├── guides/
│   ├── api/
│   ├── examples/
│   ├── contributing/
│   └── architecture/
└── [source code]
```

**Comprehensive governance:**
- Decision-making processes documented
- Maintainer roles and responsibilities
- Release and deprecation policies
- Security response procedures
- Conflict resolution mechanisms

## Domain-Specific Adaptations

The skill recognizes that different project types need different emphasis:

### Libraries
**Priority documentation:**
- API reference (critical)
- Usage examples (critical)
- Migration guides for breaking changes
- Performance characteristics
- Compatibility matrix

**Structure emphasis:**
- Clear semantic versioning
- Detailed CHANGELOG
- Deprecation policy

### CLI Tools
**Priority documentation:**
- Installation across platforms
- Command reference
- Configuration examples
- Common workflows

**Structure emphasis:**
- Simple, scannable README
- Man pages or help text
- Platform-specific .gitignore

### Frameworks
**Priority documentation:**
- Getting started tutorial
- Architecture concepts
- Extension/plugin guides
- Best practices

**Structure emphasis:**
- ARCHITECTURE.md early
- Examples directory
- Plugin contribution guidelines

### Languages/Compilers
**Priority documentation:**
- Language specification
- Compiler architecture
- Enhancement proposal process
- Standard library reference

**Structure emphasis:**
- Formal governance from start
- RFC or PEP process
- Breaking change policy
- Multiple documentation layers

### Security-Critical Projects
**Priority documentation:**
- SECURITY.md (mandatory)
- Threat model
- Security audit results
- Vulnerability disclosure timeline

**Structure emphasis:**
- Private security reporting
- Security review checklist
- Dependency audit automation

## Progressive Documentation Strategy

Rather than creating all files at once, the skill recommends a growth path:

### Phase 1: Foundation (Day 1)
- README.md with clear description
- LICENSE file
- .gitignore for your ecosystem
- Basic project structure

### Phase 2: Community Ready (First External Contributor)
- CONTRIBUTING.md with PR process
- CODE_OF_CONDUCT.md
- Issue template
- CHANGELOG.md started

### Phase 3: Scaling (10+ Contributors)
- ARCHITECTURE.md for complex projects
- docs/ directory with organized content
- CI/CD automation
- Multiple issue templates
- SUPPORT.md to direct questions

### Phase 4: Mature Governance (Multi-Org or High Impact)
- GOVERNANCE.md with decision process
- ROADMAP.md with feature planning
- SECURITY.md with disclosure process
- Release automation
- Formal RFC or proposal process

## File Descriptions

### Core Files (All Projects Eventually Need)

#### README.md
**Purpose:** Primary entry point for the project  
**Tier 1 contents:**
- Project name and one-sentence description
- Installation command
- Minimal usage example
- License mention

**Tier 2+ enhancements:**
- Badges (build status, version, coverage)
- Feature list
- Multiple usage examples
- Links to detailed docs
- Contributing quick-start
- Acknowledgments

**Anti-patterns:**
- Duplicate information from other docs (link instead)
- Outdated version information
- Broken links or examples
- No clear "what is this?" statement

#### LICENSE
**Purpose:** Legal framework for project usage  
**Common choices:**
- **MIT:** Simple, permissive
- **Apache 2.0:** Permissive with patent grant
- **GPL v3:** Copyleft, derivative works must be open
- **BSD:** Similar to MIT with additional clauses

**Selection criteria:**
- Personal projects: MIT or Apache 2.0
- Company projects: Check corporate policy
- Want derivatives to stay open: GPL
- Patent concerns: Apache 2.0

**Critical:** Must be present before accepting contributions

#### .gitignore
**Purpose:** Exclude unnecessary files from version control  
**Standard patterns for:**
- Build artifacts
- Dependency directories
- IDE configuration files
- Environment configuration
- OS-specific files (`.DS_Store`, `Thumbs.db`)
- Logs and temporary files

**Use templates from:**
- github.com/github/gitignore
- gitignore.io

### Community Files (Tier 2+)

#### CONTRIBUTING.md
**Purpose:** Lower barrier for contributions  
**Essential contents:**
- How to report bugs (template link)
- How to suggest features
- Pull request process
- Code style requirements
- Testing expectations
- Development setup

**Tier 2 additions:**
- Commit message format
- Branch naming conventions
- Review process timeline

**Tier 3 additions:**
- Architectural decision process
- Breaking change policy
- Release criteria
- Maintainer contact info

**Triggers for creation:**
- First external PR with issues
- Repeated questions about process
- Inconsistent contribution quality

#### CODE_OF_CONDUCT.md
**Purpose:** Establish community behavior standards  
**Standard template:** Contributor Covenant 2.1  
**Essential elements:**
- Expected behaviors
- Unacceptable behaviors
- Enforcement process
- Contact information for reports

**When critical:**
- Project has public discussions
- Diverse contributor base
- Organization-backed projects
- History of interpersonal conflict

**Enforcement requirement:**
- Must actually enforce it
- Have designated responders
- Document incident handling

#### CHANGELOG.md
**Purpose:** Track what changed between versions  
**Format:** Keep a Changelog (keepachangelog.com)  
**Structure:**
```markdown
# Changelog

## [Unreleased]
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [1.2.0] - 2025-01-15
### Added
- New feature X
### Fixed
- Bug Y
```

**Best practices:**
- Write for users, not developers
- Group by change type
- Link to issues/PRs
- Note breaking changes prominently
- Update before each release

**Triggers for creation:**
- First versioned release
- Users asking "what changed?"
- Difficulty tracking breaking changes

### Process Files (Tier 2+)

#### .github/ISSUE_TEMPLATE.md
**Purpose:** Improve issue quality and reduce back-and-forth  
**Basic template:**
```markdown
**Description:**
A clear description of the issue

**Steps to Reproduce:**
1. Step one
2. Step two

**Expected Behavior:**
What should happen

**Actual Behavior:**
What actually happens

**Environment:**
- OS:
- Version:
```

**Tier 3: Multiple templates:**
- `bug_report.md`
- `feature_request.md`
- `question.md`
- `security.md` (private reporting)

**Triggers for creation:**
- Vague or incomplete bug reports
- Missing environment information
- Feature requests without context

#### .github/PULL_REQUEST_TEMPLATE.md
**Purpose:** Ensure PRs include necessary information  
**Basic template:**
```markdown
**Description:**
What does this PR do?

**Related Issue:**
Fixes #123

**Changes:**
- Change 1
- Change 2

**Testing:**
- [ ] Added tests
- [ ] All tests pass
- [ ] Updated documentation

**Breaking Changes:**
- [ ] No breaking changes
- [ ] Breaking changes described below
```

**Tier 3 additions:**
- Performance impact assessment
- Security implications
- Upgrade path for breaking changes
- Reviewer checklist

### Architecture Files (Tier 2+ for Complex Projects)

#### ARCHITECTURE.md
**Purpose:** Explain system design and key decisions  
**When needed:**
- Multi-component systems
- Non-obvious design choices
- Frequent architecture questions
- New contributors need orientation

**Contents:**
- High-level component diagram
- Data flow
- Key abstractions
- Design decisions and rationale
- Extension points
- Performance considerations

**Anti-patterns:**
- Duplicating code comments
- Implementation details
- Out-of-sync with code

#### docs/ Directory Structure
**Purpose:** Organize detailed documentation  
**Recommended structure:**
```
docs/
├── index.md                    # Documentation home
├── getting-started/
│   ├── installation.md
│   ├── quick-start.md
│   └── configuration.md
├── guides/
│   ├── user-guide.md
│   ├── advanced-usage.md
│   └── troubleshooting.md
├── api/
│   └── reference.md            # Or auto-generated
├── examples/
│   ├── basic-example.md
│   └── advanced-example.md
├── contributing/
│   ├── development-setup.md
│   ├── testing-guide.md
│   └── release-process.md
└── architecture/
    ├── overview.md
    ├── design-decisions.md
    └── performance.md
```

**Documentation hosting:**
- **GitHub Pages:** Simple, free, works with Jekyll
- **Read the Docs:** Auto-builds, versioning, search
- **Docusaurus:** Modern React-based, good for large docs
- **MkDocs:** Python-based, Material theme popular

### Governance Files (Tier 3)

#### GOVERNANCE.md
**Purpose:** Define decision-making processes  
**When needed:**
- Multiple organizations involved
- Conflicting stakeholder interests
- Maintainer succession planning
- Formal release approval needed

**Contents:**
- Maintainer roles (core, committer, reviewer)
- Decision-making process (consensus, voting, BDFL)
- How maintainers are added/removed
- Conflict resolution
- Code ownership (CODEOWNERS file)
- Release authority

**Examples to study:**
- Node.js governance model
- Rust governance structure
- Python's steering council

#### SECURITY.md
**Purpose:** Responsible vulnerability disclosure  
**When mandatory:**
- Cryptography or authentication
- Network-facing applications
- Data handling
- Dependencies in security-critical contexts

**Essential contents:**
```markdown
# Security Policy

## Supported Versions
| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**DO NOT** open a public issue.

Email: security@project.org
Expected response: 48 hours
```

**Best practices:**
- Private reporting channel
- Response time commitment
- Disclosure timeline (typically 90 days)
- Recognition policy for reporters
- CVE assignment process

#### ROADMAP.md
**Purpose:** Communicate project direction  
**When valuable:**
- Active feature development
- Breaking changes planned
- Strategic direction debates
- User/contributor alignment needed

**Structure:**
- **Now:** Current focus (1-3 months)
- **Next:** Planned features (3-12 months)
- **Later:** Vision items (12+ months)
- **Not planned:** Explicitly out of scope

**Maintenance:**
- Review quarterly
- Link to tracking issues
- Note when items move between categories
- Archive completed items to CHANGELOG

### Support Files (Tier 2+)

#### SUPPORT.md
**Purpose:** Direct support questions away from issues  
**When needed:**
- Growing issue tracker noise
- Repeated "how do I?" questions
- Available community channels

**Contents:**
- Where to ask questions (Stack Overflow, Discord, Discussions)
- Where NOT to ask (GitHub issues for bugs only)
- How to search existing answers
- Expected response times
- Commercial support options (if any)

**Reduces:**
- Issue tracker clutter
- Maintainer support burden
- Duplicate questions

## Assessment Methodology

### For New Projects

The skill asks a series of diagnostic questions:

1. **What are you building?**
    - Library, framework, tool, language, application, other

2. **What language/ecosystem?**
    - Determines .gitignore, CI templates, package files

3. **How complex is the architecture?**
    - Single file → Simple
    - Single component → Moderate
    - Multi-component system → Complex
    - Distributed system → Very complex

4. **What's your stability commitment?**
    - Experimental → Tier 1
    - Beta → Tier 2
    - Production → Tier 2-3

5. **Is this security-sensitive?**
    - Yes → SECURITY.md mandatory
    - No → Optional

6. **Expected contribution model?**
    - Solo → Tier 1
    - Small team → Tier 1-2
    - Open collaboration → Tier 2
    - Multi-organization → Tier 3

7. **Impact if it breaks?**
    - Personal inconvenience → Tier 1
    - Team productivity → Tier 2
    - Production systems → Tier 3
    - Critical infrastructure → Tier 3 + enhanced

Based on answers, recommend initial structure and explain scaling triggers.

### For Existing Projects

The skill performs analysis:

1. **Scan repository structure**
    - List existing documentation files
    - Identify standard files present/missing
    - Check for outdated content (last modified dates)

2. **Analyze contributor metrics**
    - Total contributors
    - Active contributors (last 90 days)
    - First-time contributor rate
    - PR acceptance rate

3. **Review issue patterns**
    - Open vs closed ratio
    - Common question themes
    - Bug report quality
    - Response time distribution

4. **Assess project maturity**
    - Repository age
    - Release frequency
    - Semantic versioning usage
    - Breaking change frequency
    - Adoption metrics (stars, forks, dependents)

5. **Identify pain points**
    - Documentation gaps referenced in issues
    - Contribution friction (rejected PRs, process questions)
    - Governance conflicts
    - Security concerns

6. **Recommend improvements**
    - Missing files with priority ranking
    - Files to update/expand
    - Process improvements
    - Structure upgrades

### Gap Analysis Framework

When assessing existing projects, the skill identifies:

#### Critical Gaps (Fix Immediately)
- No LICENSE file
- No README or inadequate description
- Missing CONTRIBUTING.md with active external PRs
- No CODE_OF_CONDUCT with community issues
- Missing SECURITY.md for security-sensitive projects
- Broken links or examples in documentation

#### Important Gaps (Address Soon)
- Outdated CHANGELOG
- Missing issue/PR templates causing friction
- Inadequate architecture documentation for complex projects
- No clear governance with multiple organizations
- Missing .gitignore causing repository pollution

#### Optional Enhancements (Nice to Have)
- CONTRIBUTORS.md for recognition
- ROADMAP.md for transparency
- Enhanced documentation structure
- Additional CI/CD workflows
- FUNDING.yml for sponsorship

## Recommendations Output

The skill provides actionable recommendations structured as:

### Immediate Actions
Files to create now with brief justification

### Short-term Improvements (1-4 weeks)
Documentation to expand or update

### Long-term Enhancements
Structure to consider as project grows

### Maintenance Reminders
Regular updates needed (CHANGELOG, ROADMAP)

### Scaling Triggers
Conditions that indicate need for additional structure

## Best Practices

### Start Minimal, Grow Deliberately
- Don't over-document prematurely
- Add structure when pain points emerge
- Let community needs drive additions

### Keep Documentation Current
- Review docs with each release
- Update examples to match current API
- Archive outdated information
- Set reminders for periodic reviews

### Make Contributing Easy
- Clear, simple contribution process
- Good first issues labeled
- Welcoming tone in all docs
- Fast initial response to PRs

### Be Explicit About Project Status
- "Experimental - APIs will change"
- "Beta - suitable for testing"
- "Stable - production ready"
- "Mature - long-term support"

### Document Decisions, Not Just Code
- Why architectural choices were made
- What alternatives were considered
- Trade-offs accepted
- Future reconsideration triggers

### Build Community Thoughtfully
- Enforce CODE_OF_CONDUCT consistently
- Recognize contributions publicly
- Create pathways for growth (contributor → maintainer)
- Share decision-making gradually

## Common Anti-Patterns to Avoid

### Over-Engineering
- **Problem:** GOVERNANCE.md for solo project
- **Fix:** Add governance only when multi-party decisions need structure

### Documentation Rot
- **Problem:** README examples don't work anymore
- **Fix:** Test documentation as part of CI/CD

### Process Bureaucracy
- **Problem:** 10-step PR process for small fixes
- **Fix:** Scale process to contribution size

### Unclear Scope
- **Problem:** Project accepts everything, focus dilutes
- **Fix:** Document explicit scope and non-goals

### Maintainer Bottleneck
- **Problem:** One person reviews everything
- **Fix:** Distribute authority via CODEOWNERS

### Silent Decisions
- **Problem:** Major changes with no discussion
- **Fix:** Proposal process (RFC, ADR) for significant changes

### Toxic Permissiveness
- **Problem:** No CODE_OF_CONDUCT, community degrades
- **Fix:** Establish and enforce community standards

## Platform-Specific Considerations

### GitHub
- Use GitHub Discussions for Q&A
- GitHub Actions for CI/CD
- GitHub Projects for roadmap
- GitHub Sponsors for funding
- Branch protection rules
- CODEOWNERS file for distributed review

### GitLab
- `.gitlab-ci.yml` for CI/CD
- GitLab Pages for docs
- Merge request templates
- Issue boards for project management
- GitLab-specific issue templates

### Other Platforms
- Adapt directory conventions
- Use platform-native CI/CD
- Leverage platform community features

## Resources and Templates

### License Selection
- [choosealicense.com](https://choosealicense.com/)
- [opensource.org/licenses](https://opensource.org/licenses)

### Community Standards
- [Contributor Covenant](https://www.contributor-covenant.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)

### Learning Resources
- [opensource.guide](https://opensource.guide/)
- [The Architecture of Open Source Applications](http://aosabook.org/)

### Template Repositories
- GitHub's community templates
- Language-specific project templates
- Organization-specific templates

## Skill Usage Examples

### Example 1: New CLI Tool
**User:** "I'm creating a command-line tool for developers"

**Skill assesses:**
- Type: CLI tool
- Complexity: Simple
- Expected users: Developers
- Status: New

**Recommendation:**
```
Tier 1 structure:
├── README.md (installation, usage examples)
├── LICENSE (suggest MIT or Apache 2.0)
├── .gitignore (language-specific)
└── src/

Next steps:
- When first external PR → add CONTRIBUTING.md
- When 5+ stars → consider CODE_OF_CONDUCT.md
- At v1.0 release → start CHANGELOG.md
```

### Example 2: Existing Library with Growth
**User:** "I have a Python library with 15 contributors, getting 5-10 PRs/month"

**Skill analyzes:**
- Has: README, LICENSE, .gitignore
- Missing: CONTRIBUTING, CODE_OF_CONDUCT, issue templates
- Pain points: PRs often need rework, style inconsistencies

**Recommendation:**
```
Critical additions:
1. CONTRIBUTING.md - define PR process and code style
2. CODE_OF_CONDUCT.md - community is growing
3. .github/PULL_REQUEST_TEMPLATE.md - reduce review cycles

Consider soon:
- ARCHITECTURE.md - if design questions recurring
- docs/ directory - if README exceeds 500 lines
- CHANGELOG.md - track breaking changes

You're transitioning from Tier 1 to Tier 2
```

### Example 3: Security-Critical Project
**User:** "I'm open-sourcing our authentication library"

**Skill flags:**
- Security-sensitive domain
- Immediate SECURITY.md required
- Enhanced documentation needed
- Formal process recommended

**Recommendation:**
```
Mandatory immediately:
├── SECURITY.md (private disclosure process)
├── docs/threat-model.md
├── CONTRIBUTING.md (security review requirements)

Strongly recommended:
├── .github/workflows/security-scan.yml
├── docs/security-best-practices.md
├── Automated dependency auditing

Start at Tier 2+, plan for Tier 3 governance
```

## Continuous Improvement

The skill should periodically suggest reassessment:

- **Every 6 months:** Review documentation accuracy
- **Major version releases:** Update all version-specific docs
- **Every 25 contributors:** Reassess governance needs
- **When conflicts arise:** Strengthen CODE_OF_CONDUCT or add GOVERNANCE.md
- **Security incidents:** Review and enhance SECURITY.md

The goal is documentation that grows with the project, adding structure when it provides value, not before.
