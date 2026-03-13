# GitHub Actions CI/CD Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add GitHub Actions workflows for running tests on push/PR and building cross-platform release binaries on tag push, plus README badges.

**Architecture:** Two independent workflow files (test + release), a GoReleaser v2 config for cross-compilation, and badge additions to README. All files are new except README which gets a one-line edit.

**Tech Stack:** GitHub Actions, GoReleaser v2, actions/checkout, actions/setup-go@v5, goreleaser/goreleaser-action

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `.github/workflows/test.yml` | Create | Run `go test ./...` on push/PR |
| `.github/workflows/release.yml` | Create | Run GoReleaser on tag push |
| `.goreleaser.yaml` | Create | Cross-compilation config for 5 targets |
| `README.md` | Edit | Add test status + release badges |

---

## Chunk 1: CI/CD Workflows and Badges

### Task 1: Create test workflow

**Files:**
- Create: `.github/workflows/test.yml`

- [ ] **Step 1: Create the workflow file**

```yaml
name: Tests

on:
  push:
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - run: go test ./...
```

- [ ] **Step 2: Validate the workflow syntax**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && cat .github/workflows/test.yml`
Verify: YAML is valid, triggers and steps match spec.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/test.yml
git commit -m "ci: add test workflow for push and PR"
```

---

### Task 2: Create GoReleaser configuration

**Files:**
- Create: `.goreleaser.yaml`

- [ ] **Step 1: Create the GoReleaser v2 config**

```yaml
version: 2

builds:
  - main: ./cmd/archdiag
    binary: archdiag
    env:
      - CGO_ENABLED=0
    ldflags:
      - -s -w
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64

archives:
  - name_template: "archdiag_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
```

- [ ] **Step 2: Validate the config locally (if goreleaser is installed)**

Run: `goreleaser check` (optional — skip if not installed)

- [ ] **Step 3: Commit**

```bash
git add .goreleaser.yaml
git commit -m "ci: add GoReleaser v2 configuration"
```

---

### Task 3: Create release workflow

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create the workflow file**

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 2: Validate the workflow syntax**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && cat .github/workflows/release.yml`
Verify: YAML is valid, `fetch-depth: 0` present, permissions set, GITHUB_TOKEN passed.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add release workflow with GoReleaser"
```

---

### Task 4: Add README badges

**Files:**
- Modify: `README.md:3-4` (existing badge lines)

- [ ] **Step 1: Add badges to existing badge block**

Add these two lines after line 4 (`[![License: MIT]...`) in `README.md`:

```markdown
[![Tests](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/crosscode-nl/archdiag)](https://github.com/crosscode-nl/archdiag/releases/latest)
```

The badge block should look like:

```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/crosscode-nl/archdiag.svg)](https://pkg.go.dev/github.com/crosscode-nl/archdiag)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Tests](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/crosscode-nl/archdiag)](https://github.com/crosscode-nl/archdiag/releases/latest)
```

- [ ] **Step 2: Verify the README renders correctly**

Run: `head -6 README.md`
Verify: 4 badge lines between the title and the blank line.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add test status and release badges to README"
```

---

### Task 5: Run tests to verify nothing is broken

- [ ] **Step 1: Run the full test suite**

Run: `go test ./...`
Expected: All tests pass. No workflow files affect test behavior.

---

### Verification Checklist

- [ ] `.github/workflows/test.yml` exists with push + PR triggers
- [ ] `.github/workflows/release.yml` exists with `v*` tag trigger, `fetch-depth: 0`, `contents: write` permission
- [ ] `.goreleaser.yaml` exists with `version: 2`, 5 build targets, zip override for Windows
- [ ] `README.md` has 4 badge lines (Go Reference, License, Tests, Release)
- [ ] `go test ./...` passes
