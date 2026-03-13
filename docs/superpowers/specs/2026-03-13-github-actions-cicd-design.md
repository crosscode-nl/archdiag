# GitHub Actions CI/CD for archdiag

## Overview

Add GitHub Actions workflows to the archdiag project for two purposes:
1. Run tests on every push and pull request
2. Build and attach cross-platform binaries to GitHub releases when a version tag is pushed

Add README badges showing test status and latest release version.

## Test Workflow

**File:** `.github/workflows/test.yml`

**Triggers:**
- Push to any branch
- Pull requests targeting `main`

**Job:** Single job on `ubuntu-latest`

**Steps:**
1. Checkout repository (`actions/checkout`)
2. Set up Go via `actions/setup-go@v5` with `go-version-file: go.mod` (reads version from `go.mod`, includes automatic module caching)
3. Run `go test ./...`

## Release Workflow

**File:** `.github/workflows/release.yml`

**Triggers:**
- Push of a tag matching `v*` (e.g. `v1.0.0`, `v0.2.3`)

**Job:** Single job on `ubuntu-latest`

**Steps:**
1. Checkout repository with `fetch-depth: 0` (full git history required for GoReleaser changelog generation)
2. Set up Go via `actions/setup-go@v5` with `go-version-file: go.mod`
3. Run GoReleaser via `goreleaser/goreleaser-action` (uses GoReleaser v2)

**Permissions:** `contents: write` (required to create/upload release assets)

**Environment:** `GITHUB_TOKEN` provided via `secrets.GITHUB_TOKEN`

## GoReleaser Configuration

**File:** `.goreleaser.yaml`

**Version:** GoReleaser v2 config format (requires `version: 2` at top of file).

**Build targets:**

| OS      | Architecture |
|---------|-------------|
| Linux   | amd64       |
| Linux   | arm64       |
| macOS   | amd64       |
| macOS   | arm64       |
| Windows | amd64       |

Windows arm64 is omitted due to low demand for CLI tools on that platform.

**Build settings:**
- Binary name: `archdiag`
- Main package: `./cmd/archdiag`
- CGO disabled (pure Go, required for cross-compilation)
- ldflags: `-s -w` (strip debug info for smaller binaries)

**Archives:**
- Format: `.tar.gz` for Linux and macOS, `.zip` for Windows
- Name template: `archdiag_{{ .Os }}_{{ .Arch }}`

**Checksum:**
- SHA256 checksums file generated for all archives

**Changelog:**
- Auto-generated from git commit history between tags

## README Badges

Add two badges to the existing badge line in `README.md` (line 3), alongside the existing Go Reference and License badges. Use the linked badge format (`[![alt](image)](url)`) to match the existing style:

```markdown
[![Tests](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/crosscode-nl/archdiag)](https://github.com/crosscode-nl/archdiag/releases/latest)
```

## Files Changed

| File | Action |
|------|--------|
| `.github/workflows/test.yml` | Create |
| `.github/workflows/release.yml` | Create |
| `.goreleaser.yaml` | Create |
| `README.md` | Edit (add badges to existing badge line) |

## Release Process

After these changes, the release workflow is:

1. Develop on feature branches, tests run automatically
2. When ready to release, create and push a tag: `git tag v1.0.0 && git push origin v1.0.0`
3. The release workflow builds binaries for all 5 targets and creates a GitHub release with attached archives and checksums
4. Users can download pre-built binaries from the GitHub releases page
