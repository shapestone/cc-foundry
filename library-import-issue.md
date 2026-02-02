# Library Import Issue: Module Path Mismatch and Local Replace Directives

This document describes a pattern of issues that prevent Go libraries from being consumed as external dependencies. While discovered in `shape-yaml`, this applies to any library that follows the same approach.

## Problem Summary

Two common practices in multi-module Go development break external consumption when modules are published to the Go module proxy:

1. **Local `replace` directives left in published `go.mod`**
2. **Module path mismatches between tagged versions and the `require` path**

Any library that uses local `replace` directives during development and tags releases without removing them will exhibit these issues.

## Issue 1: `replace` Directive in Published go.mod

### What happens

A library's `go.mod` contains a local path replacement:

```
replace github.com/shapestone/shape-core => ../shape-core
```

`replace` directives with local paths are **ignored by the Go module proxy** when the module is consumed as a dependency. The proxy resolves `shape-core` from the registry, not the local path. This means the `replace` only works during local development in the author's workspace.

### Impact

Any external consumer running `go get` will bypass the `replace` and pull `shape-core` directly from the proxy, potentially hitting version or path mismatches that the `replace` was masking.

### Fix

Remove `replace` directives before tagging releases. For local multi-module development, use Go workspaces (`go.work`) instead:

```bash
# go.work (not committed, local development only)
go 1.23

use (
    ./shape-yaml
    ./shape-core
)
```

This keeps published `go.mod` files clean while still allowing cross-module local development.

## Issue 2: Module Path Mismatch on Tagged Version

### What happens

When Go fetches a dependency at a specific version (e.g., `github.com/shapestone/shape-core@v0.9.2`), the `go.mod` inside that tagged version declares a different module path:

```
module github.com/shapestone/shape
```

But the consuming library requires it as:

```
require github.com/shapestone/shape-core v0.9.2
```

Go enforces that the declared module path in `go.mod` must match the path used in `require`. This mismatch causes a hard failure:

```
github.com/shapestone/shape-core@v0.9.2: parsing go.mod:
    module declares its path as: github.com/shapestone/shape
            but was required as: github.com/shapestone/shape-core
```

### Impact

The module is completely unusable as a dependency at the affected version. Even if the `main` branch has the correct module path, the tagged version is immutable on the proxy.

### Fix

Tag a new release from the corrected source. Go module versions are immutable once published — the only path forward is a new version.

## Recommended Release Sequence

Using `shape-yaml` / `shape-core` as the concrete example:

1. Verify `shape-core` `main` branch has the correct `module github.com/shapestone/shape-core` in `go.mod` (confirmed)
2. Tag `shape-core` v0.9.3 from `main`
3. In `shape-yaml`:
   - Update `require` to `github.com/shapestone/shape-core v0.9.3`
   - Remove the `replace github.com/shapestone/shape-core => ../shape-core` directive
   - Run `go mod tidy`
4. Tag `shape-yaml` v0.9.3
5. Consumers can then run `go get github.com/shapestone/shape-yaml@v0.9.3`

## General Prevention Guidelines

These practices apply to any Go library intended for external consumption:

- **Never publish tags with local `replace` directives** — use `go.work` for local development instead
- **Verify module path consistency** before tagging: the `module` declaration in `go.mod` must match the repository path
- **Test external consumption** before releasing: create a throwaway module and `go get` the new version to confirm it resolves correctly
- **Treat tags as immutable** — if a bad version is published, the only fix is a new version number
