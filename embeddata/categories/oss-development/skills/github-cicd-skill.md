---
name: github-cicd
description: Comprehensive CI/CD pipeline setup and management for open source GitHub projects. Use this skill when Claude needs to: (1) Set up or configure GitHub Actions workflows, (2) Create CI/CD pipelines for testing, building, and deployment, (3) Configure automated checks for pull requests, (4) Set up release automation and semantic versioning, (5) Implement code quality checks (linting, formatting, security scanning), (6) Configure multi-platform builds and deployments, (7) Troubleshoot or debug existing CI/CD workflows, or (8) Optimize build performance and caching strategies
---

# GitHub CI/CD Management Skill

This skill provides comprehensive guidance for setting up and managing CI/CD processes for open source projects on GitHub using GitHub Actions.

## Core Principles

1. **Start Simple, Scale Gradually**: Begin with basic workflows and add complexity as needed
2. **Fail Fast**: Configure checks to identify issues early in the pipeline
3. **Security First**: Always use minimal permissions and secure secret handling
4. **Developer Experience**: Optimize for fast feedback and clear error messages
5. **Cost Awareness**: Use caching and conditional execution to minimize CI minutes

## Workflow Organization

### Directory Structure

```
.github/
├── workflows/
│   ├── ci.yml              # Main CI pipeline (tests, linting)
│   ├── release.yml         # Release automation
│   ├── pr-checks.yml       # PR-specific validation
│   ├── deploy.yml          # Deployment pipeline
│   └── security.yml        # Security scanning
├── actions/                # Reusable composite actions
│   └── setup-env/
│       └── action.yml
└── dependabot.yml          # Dependency updates
```

### Naming Conventions

- Use descriptive, kebab-case names: `ci.yml`, `deploy-production.yml`
- Prefix specialized workflows: `pr-`, `release-`, `deploy-`
- Name jobs clearly: `test`, `lint`, `build`, `deploy-staging`

## Essential Workflows

### 1. Continuous Integration (ci.yml)

**Triggers**: Push to main/develop branches, all PRs
**Purpose**: Run tests, linting, and basic validation

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [18, 20, 22]
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}
          cache: 'npm'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        if: matrix.node-version == 20
        with:
          token: ${{ secrets.CODECOV_TOKEN }}

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'npm'
      - run: npm ci
      - run: npm run lint
      - run: npm run format:check

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'npm'
      - run: npm ci
      - run: npm run build
      - name: Upload build artifacts
        uses: actions/upload-artifact@v4
        with:
          name: build
          path: dist/
```

### 2. Pull Request Checks (pr-checks.yml)

**Triggers**: PR opened, synchronized, reopened
**Purpose**: Validate PR quality before merging

```yaml
name: PR Checks

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  validate-pr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Validate PR title
        uses: amannn/action-semantic-pull-request@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Check for breaking changes
        run: |
          if git diff origin/main --name-only | grep -q "BREAKING"; then
            echo "::warning::This PR may contain breaking changes"
          fi
      
      - name: Size label
        uses: pascalgn/size-label-action@v0.5.3
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  dependency-review:
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/dependency-review-action@v4
```

### 3. Release Automation (release.yml)

**Triggers**: Push to main with version tags
**Purpose**: Automate releases and changelog generation

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'npm'
          registry-url: 'https://registry.npmjs.org'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Build
        run: npm run build
      
      - name: Run tests
        run: npm test
      
      - name: Publish to npm
        run: npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
      
      - name: Generate changelog
        id: changelog
        uses: mikepenz/release-changelog-builder-action@v4
        with:
          configuration: ".github/changelog-config.json"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          body: ${{ steps.changelog.outputs.changelog }}
          files: |
            dist/**/*
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 4. Security Scanning (security.yml)

**Triggers**: Weekly schedule, push to main, PRs
**Purpose**: Scan for vulnerabilities and security issues

```yaml
name: Security

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sundays

jobs:
  codeql:
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      actions: read
      contents: read
    steps:
      - uses: actions/checkout@v4
      
      - name: Initialize CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: javascript
      
      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3

  dependency-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Snyk
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
        with:
          args: --severity-threshold=high
      
      - name: Audit dependencies
        run: npm audit --audit-level=moderate

  secret-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: TruffleHog scan
        uses: trufflesecurity/trufflehog@main
        with:
          path: ./
          base: ${{ github.event.repository.default_branch }}
          head: HEAD
```

## Language-Specific Configurations

### Python Projects

```yaml
- uses: actions/setup-python@v5
  with:
    python-version: '3.11'
    cache: 'pip'

- name: Install dependencies
  run: |
    pip install -r requirements.txt
    pip install -r requirements-dev.txt

- name: Run tests with pytest
  run: pytest --cov=src tests/

- name: Lint with ruff
  run: ruff check .

- name: Type check with mypy
  run: mypy src/
```

### Go Projects

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.21'
    cache: true

- name: Run tests
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Lint
  uses: golangci/golangci-lint-action@v3
  with:
    version: latest

- name: Build
  run: go build -v ./...
```

### Rust Projects

```yaml
- uses: dtolnay/rust-toolchain@stable
  with:
    components: rustfmt, clippy

- uses: Swatinem/rust-cache@v2

- name: Format check
  run: cargo fmt -- --check

- name: Lint
  run: cargo clippy -- -D warnings

- name: Test
  run: cargo test --all-features

- name: Build
  run: cargo build --release
```

### Docker Projects

```yaml
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3

- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}

- name: Build and push
  uses: docker/build-push-action@v5
  with:
    context: .
    push: true
    tags: user/app:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

## Performance Optimization

### Caching Strategies

```yaml
# Node.js npm cache
- uses: actions/setup-node@v4
  with:
    cache: 'npm'

# Python pip cache
- uses: actions/setup-python@v5
  with:
    cache: 'pip'

# Custom cache for build artifacts
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache
      dist/
    key: ${{ runner.os }}-build-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-build-
```

### Parallel Execution

```yaml
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        node: [18, 20, 22]
      fail-fast: false  # Continue testing other combinations on failure
    runs-on: ${{ matrix.os }}
```

### Conditional Execution

```yaml
# Skip CI for documentation-only changes
on:
  push:
    paths-ignore:
      - '**.md'
      - 'docs/**'

# Run only on specific file changes
on:
  push:
    paths:
      - 'src/**'
      - 'tests/**'
      - 'package.json'

# Conditional steps
- name: Deploy
  if: github.ref == 'refs/heads/main' && github.event_name == 'push'
  run: npm run deploy
```

## Advanced Patterns

### Reusable Workflows

Create `.github/workflows/reusable-test.yml`:

```yaml
name: Reusable Test Workflow

on:
  workflow_call:
    inputs:
      node-version:
        required: true
        type: string
    secrets:
      NPM_TOKEN:
        required: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}
      - run: npm ci
      - run: npm test
```

Use in another workflow:

```yaml
jobs:
  test-node-18:
    uses: ./.github/workflows/reusable-test.yml
    with:
      node-version: '18'
    secrets:
      NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
```

### Composite Actions

Create `.github/actions/setup-env/action.yml`:

```yaml
name: 'Setup Environment'
description: 'Setup Node.js and install dependencies'
inputs:
  node-version:
    description: 'Node.js version'
    required: false
    default: '20'
runs:
  using: 'composite'
  steps:
    - uses: actions/setup-node@v4
      with:
        node-version: ${{ inputs.node-version }}
        cache: 'npm'
    - run: npm ci
      shell: bash
```

Use in workflows:

```yaml
- uses: ./.github/actions/setup-env
  with:
    node-version: '20'
```

### Matrix with Exclusions

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    node: [18, 20, 22]
    exclude:
      - os: macos-latest
        node: 18  # Skip macOS + Node 18
    include:
      - os: ubuntu-latest
        node: 22
        experimental: true  # Add extra data
```

## Deployment Workflows

### Deploy to GitHub Pages

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'npm'
      - run: npm ci
      - run: npm run build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: dist/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### Deploy to Cloud Platforms

```yaml
# AWS S3
- name: Deploy to S3
  run: |
    aws s3 sync dist/ s3://${{ secrets.S3_BUCKET }} --delete
  env:
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    AWS_REGION: us-east-1

# Vercel
- name: Deploy to Vercel
  uses: amondnet/vercel-action@v25
  with:
    vercel-token: ${{ secrets.VERCEL_TOKEN }}
    vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
    vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}

# Netlify
- name: Deploy to Netlify
  uses: nwtgck/actions-netlify@v2
  with:
    publish-dir: './dist'
    production-branch: main
  env:
    NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
    NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
```

## Dependency Management

### Dependabot Configuration

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    reviewers:
      - "maintainer-username"
    labels:
      - "dependencies"
    commit-message:
      prefix: "chore"
      include: "scope"
    groups:
      development-dependencies:
        dependency-type: "development"
      production-dependencies:
        dependency-type: "production"
        update-types:
          - "minor"
          - "patch"

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "monthly"
    labels:
      - "ci"
      - "dependencies"
```

### Auto-merge Dependabot PRs

```yaml
name: Auto-merge Dependabot

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  auto-merge:
    runs-on: ubuntu-latest
    if: github.actor == 'dependabot[bot]'
    steps:
      - name: Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@v1
        with:
          github-token: "${{ secrets.GITHUB_TOKEN }}"
      
      - name: Auto-merge minor and patch updates
        if: steps.metadata.outputs.update-type == 'version-update:semver-minor' || steps.metadata.outputs.update-type == 'version-update:semver-patch'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Secrets Management

### Required Secrets Setup

For npm publishing:
- `NPM_TOKEN`: npm authentication token

For Docker Hub:
- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password or access token

For cloud deployments:
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- `VERCEL_TOKEN`, `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`
- `NETLIFY_AUTH_TOKEN`, `NETLIFY_SITE_ID`

For security scanning:
- `CODECOV_TOKEN`: Code coverage reporting
- `SNYK_TOKEN`: Vulnerability scanning

### Environment-specific Secrets

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production  # Use environment-specific secrets
    steps:
      - name: Deploy
        run: ./deploy.sh
        env:
          API_KEY: ${{ secrets.API_KEY }}  # Scoped to production environment
```

## Monitoring and Notifications

### Slack Notifications

```yaml
- name: Notify Slack on failure
  if: failure()
  uses: slackapi/slack-github-action@v1
  with:
    payload: |
      {
        "text": "CI Failed: ${{ github.workflow }}",
        "blocks": [
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": "Workflow *${{ github.workflow }}* failed\n*Repository:* ${{ github.repository }}\n*Branch:* ${{ github.ref }}\n*Author:* ${{ github.actor }}"
            }
          }
        ]
      }
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### GitHub Status Checks

```yaml
# Set as required in branch protection rules
jobs:
  test:
    name: Required Tests  # This name appears in GitHub UI
    runs-on: ubuntu-latest
    steps:
      - run: npm test
```

## Troubleshooting

### Debug Mode

Enable debug logging:
1. Set repository secret: `ACTIONS_STEP_DEBUG` = `true`
2. Or add to workflow:

```yaml
- name: Debug info
  run: |
    echo "Event: ${{ github.event_name }}"
    echo "Ref: ${{ github.ref }}"
    echo "SHA: ${{ github.sha }}"
    echo "Actor: ${{ github.actor }}"
```

### Common Issues and Solutions

**Issue**: Workflow not triggering
- Check trigger conditions and branch names
- Verify workflow file is in `.github/workflows/`
- Ensure YAML is valid (no syntax errors)

**Issue**: Permission denied errors
- Add required permissions to workflow:
```yaml
permissions:
  contents: write
  packages: write
```

**Issue**: Slow builds
- Implement caching for dependencies
- Use `concurrency` to cancel redundant runs
- Optimize matrix strategy

**Issue**: Secrets not accessible
- Verify secret names match exactly (case-sensitive)
- Check if using environment-specific secrets correctly
- Ensure secrets are set at repository or organization level

**Issue**: Artifacts not persisting
- Check artifact upload/download action versions match
- Verify artifact retention settings
- Use unique artifact names to avoid conflicts

## Best Practices Checklist

- [ ] Use latest stable versions of actions (e.g., `@v4`, not `@master`)
- [ ] Pin action versions for reproducibility
- [ ] Implement proper caching strategies
- [ ] Use matrix builds for multi-platform support
- [ ] Set appropriate permissions (principle of least privilege)
- [ ] Add status badges to README
- [ ] Configure branch protection rules
- [ ] Set up required status checks
- [ ] Use environments for deployment workflows
- [ ] Document CI/CD setup in repository
- [ ] Monitor workflow execution times
- [ ] Regularly update dependencies and actions
- [ ] Test workflows on fork before merging
- [ ] Use reusable workflows to reduce duplication

## Quick Reference Commands

```bash
# Validate workflow syntax locally (requires act)
act -n

# List available workflows
gh workflow list

# View workflow runs
gh run list --workflow=ci.yml

# View specific run logs
gh run view <run-id> --log

# Re-run failed workflow
gh run rerun <run-id>

# Download artifacts
gh run download <run-id>

# Trigger workflow manually (if configured)
gh workflow run release.yml
```

## Additional Resources

When implementing or troubleshooting CI/CD:
- GitHub Actions documentation: https://docs.github.com/actions
- Actions marketplace: https://github.com/marketplace?type=actions
- Community forums: https://github.community/c/actions
- Status page: https://www.githubstatus.com/
