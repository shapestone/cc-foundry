---
name: github-badges-skill
description: Expert guidance on selecting and implementing GitHub badges that build trust, credibility, and professionalism in open-source projects
---

# GitHub Badges Skill

## Description
Expert guidance on selecting and implementing GitHub badges that build trust, credibility, and professionalism in open-source projects. This skill helps you choose the right badges to signal code quality, maintenance status, and community health.

## Usage
Use this skill when:
- Creating a new README.md for an open-source project
- Improving the visibility and credibility of existing repositories
- Deciding which badges add value vs clutter
- Setting up CI/CD pipelines that generate status badges
- Establishing community health indicators

## Best Practices

### Badge Selection Philosophy
1. **Quality over Quantity**: Too many badges create visual noise. Focus on 5-10 meaningful ones.
2. **Action-Oriented**: Only include badges that reflect actual project practices (don't add a test coverage badge if you have no tests).
3. **User-Focused**: Prioritize badges that help users make decisions (license, stability, security).
4. **Keep Current**: Remove badges for deprecated services or outdated workflows.

### Badge Categories

#### ✅ Essential Badges (Must-Have)
These build immediate trust and answer basic user questions:

**1. Build/CI Status**
- **Purpose**: Shows the code is functional and tests pass
- **Trust Signal**: Active maintenance and stability
- **Implementation**:
```markdown
![Build](https://github.com/OWNER/REPO/actions/workflows/ci.yml/badge.svg)
```
- **Alternatives**: Travis CI, CircleCI, Jenkins

**2. Test Coverage**
- **Purpose**: Indicates code testing thoroughness
- **Trust Signal**: Reliability and fewer bugs
- **Providers**: Codecov, Coveralls
- **Implementation**:
```markdown
[![Coverage](https://codecov.io/gh/OWNER/REPO/branch/main/graph/badge.svg)](https://codecov.io/gh/OWNER/REPO)
```

**3. License**
- **Purpose**: Legal clarity for usage
- **Trust Signal**: Transparency about rights
- **Implementation**:
```markdown
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
```
- **Common Licenses**: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause

**4. Latest Release/Version**
- **Purpose**: Shows release cadence and current version
- **Trust Signal**: Active development
- **Implementation**:
```markdown
![GitHub release](https://img.shields.io/github/v/release/OWNER/REPO)
```

**5. Downloads/Installs**
- **Purpose**: Social proof of usage
- **Trust Signal**: Real-world adoption
- **Implementation** (npm example):
```markdown
![npm](https://img.shields.io/npm/dm/package-name)
```
- **Alternatives**: PyPI downloads, Docker pulls, gem downloads

#### ⭐ Highly Recommended Badges (Community Health)

**6. Open Issues/PRs**
- **Purpose**: Shows maintenance responsiveness
- **Implementation**:
```markdown
![GitHub issues](https://img.shields.io/github/issues/OWNER/REPO)
![GitHub PRs](https://img.shields.io/github/issues-pr/OWNER/REPO)
```

**7. Contributors**
- **Purpose**: Indicates community strength
- **Trust Signal**: Sustainable project
- **Implementation**:
```markdown
![GitHub contributors](https://img.shields.io/github/contributors/OWNER/REPO)
```

**8. Code Style/Linting**
- **Purpose**: Shows quality standards
- **Examples**: Prettier, ESLint, Black, Ruff
- **Implementation** (Black example):
```markdown
![code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)
```

**9. Documentation Status**
- **Purpose**: Indicates documentation availability
- **Providers**: ReadTheDocs, Docusaurus, GitHub Pages
- **Implementation**:
```markdown
![docs](https://readthedocs.org/projects/yourproject/badge/?version=latest)
```

**10. Community Links**
- **Purpose**: Shows where to get help
- **Examples**: Discord, Slack, Gitter, Code of Conduct
- **Implementation**:
```markdown
![Discord](https://img.shields.io/discord/SERVER_ID)
![Code of Conduct](https://img.shields.io/badge/Code%20of%20Conduct-Active-blue)
```

#### 🔐 Security & Compliance Badges (Strong Trust Signals)

**11. Security Policy**
- **Purpose**: Indicates vulnerability reporting process
- **Implementation**:
```markdown
![Security Policy](https://img.shields.io/badge/Security-Policy-brightgreen)
```

**12. Dependency Scanning**
- **Purpose**: Shows dependency health monitoring
- **Tools**: Snyk, Dependabot, Renovate
- **Implementation**:
```markdown
![Dependabot](https://img.shields.io/badge/Dependabot-enabled-blueviolet)
```

#### 🎯 Optional Badges (Social Proof)

**13. Stars/Forks**
- **Use Case**: When popularity is relevant (libraries, frameworks)
- **Implementation**:
```markdown
![GitHub stars](https://img.shields.io/github/stars/OWNER/REPO?style=social)
```

**14. Code Quality Scores**
- **Providers**: Codacy, LGTM, SonarCloud
- **Implementation**:
```markdown
![Codacy grade](https://img.shields.io/codacy/grade/PROJECT_ID)
```

**15. Sponsorship**
- **Purpose**: Shows funding options
- **Implementation**:
```markdown
![GitHub Sponsors](https://img.shields.io/github/sponsors/USERNAME)
```

## Implementation Guide

### Step 1: Identify Your Project's Needs
Ask these questions:
- Is this a library, application, or tool?
- Who is the primary audience (developers, end-users, enterprises)?
- What assurances do users need (stability, security, support)?
- What processes do we actually have in place?

### Step 2: Set Up Supporting Infrastructure
Before adding badges, ensure you have:
- CI/CD pipeline configured
- Test coverage reporting (if applicable)
- License file in repository
- Security policy (SECURITY.md)
- Contributing guidelines (CONTRIBUTING.md)

### Step 3: Badge Placement
**Optimal Layout**:
```markdown
# Project Name

[Brief one-line description]

![Build](badge1) ![Coverage](badge2) ![License](badge3)
![Version](badge4) ![Downloads](badge5)

[Longer description...]

## Features
...
```

**Alternative for Many Badges**:
```markdown
# Project Name

## Status & Quality
![Build](badge1) ![Coverage](badge2) ![Quality](badge3)

## Community
![Contributors](badge4) ![Issues](badge5) ![Discord](badge6)

## Installation
...
```

### Step 4: Customization with Shields.io
For custom badges, use [shields.io](https://shields.io):
```markdown
![Custom Badge](https://img.shields.io/badge/LABEL-MESSAGE-COLOR)
```

**Colors**: brightgreen, green, yellowgreen, yellow, orange, red, blue, lightgrey

## Common Patterns by Project Type

### Library/Package
Priority badges:
1. Build status
2. Test coverage
3. License
4. Version/Release
5. Downloads
6. Documentation

### Application/Tool
Priority badges:
1. Build status
2. License
3. Latest release
4. Security policy
5. Contributors

### Framework
Priority badges:
1. Build status
2. Test coverage
3. Documentation
4. Version
5. Contributors
6. Community (Discord/Slack)
7. Stars (social proof)

## Ready-to-Use Templates

### Minimal Set (5 badges)
```markdown
![Build](https://github.com/OWNER/REPO/actions/workflows/ci.yml/badge.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Release](https://img.shields.io/github/v/release/OWNER/REPO)
![Issues](https://img.shields.io/github/issues/OWNER/REPO)
![Contributors](https://img.shields.io/github/contributors/OWNER/REPO)
```

### Standard Set (8 badges)
```markdown
![Build](https://github.com/OWNER/REPO/actions/workflows/ci.yml/badge.svg)
[![Coverage](https://codecov.io/gh/OWNER/REPO/branch/main/graph/badge.svg)](https://codecov.io/gh/OWNER/REPO)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Release](https://img.shields.io/github/v/release/OWNER/REPO)
![Issues](https://img.shields.io/github/issues/OWNER/REPO)
![PRs](https://img.shields.io/github/issues-pr/OWNER/REPO)
![Contributors](https://img.shields.io/github/contributors/OWNER/REPO)
![Dependabot](https://img.shields.io/badge/Dependabot-enabled-blueviolet)
```

### Comprehensive Set (10+ badges)
```markdown
![Build](https://github.com/OWNER/REPO/actions/workflows/ci.yml/badge.svg)
[![Coverage](https://codecov.io/gh/OWNER/REPO/branch/main/graph/badge.svg)](https://codecov.io/gh/OWNER/REPO)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Release](https://img.shields.io/github/v/release/OWNER/REPO)
![npm](https://img.shields.io/npm/dm/package-name)
![Issues](https://img.shields.io/github/issues/OWNER/REPO)
![PRs](https://img.shields.io/github/issues-pr/OWNER/REPO)
![Contributors](https://img.shields.io/github/contributors/OWNER/REPO)
![Docs](https://readthedocs.org/projects/yourproject/badge/?version=latest)
![Code of Conduct](https://img.shields.io/badge/Code%20of%20Conduct-Active-blue)
![Dependabot](https://img.shields.io/badge/Dependabot-enabled-blueviolet)
```

## Anti-Patterns to Avoid

1. **Badge Walls**: More than 15 badges is overwhelming
2. **Broken Badges**: Regularly verify all badges work
3. **Vanity Metrics**: Don't add badges just to look impressive
4. **Outdated Services**: Remove badges for defunct CI services
5. **Fake Status**: Don't add a "passing" badge if tests don't exist

## Maintenance Tips

- **Quarterly Review**: Check all badges still work
- **Update URLs**: When changing CI providers, update badges
- **Remove Noise**: Delete badges that no longer add value
- **Monitor Performance**: Ensure badge loading doesn't slow README rendering

## Additional Resources

- [Shields.io](https://shields.io) - Badge creation service
- [GitHub Badges](https://github.com/badges/shields) - Open-source badge system
- [Awesome Badges](https://github.com/Naereen/badges) - Badge collection
- [README Best Practices](https://github.com/jehna/readme-best-practices)

## Placeholders to Replace

When using templates, replace these:
- `OWNER` - GitHub username or organization
- `REPO` - Repository name
- `package-name` - npm/PyPI package name
- `ci.yml` - Your actual workflow filename
- `main` - Your default branch name (could be `master` or `develop`)
- `PROJECT_ID` - Your project ID from quality services
- `SERVER_ID` - Discord server ID
- `USERNAME` - Your GitHub username
