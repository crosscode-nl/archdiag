# README and Module Rename Design

## Summary

Generate a public-facing open-source README.md for archdiag and update the Go module path from `go.domain.example/archdiag` to `github.com/crosscode-nl/archdiag`. Add an MIT LICENSE file.

## Changes

### 1. Update Go module path

- Change `go.mod` module to `github.com/crosscode-nl/archdiag`
- Update all import paths across all `.go` files (13 files)

### 2. Generate MIT LICENSE file

- Copyright holder: Patrick Vollebregt
- Year: 2026

### 3. Generate README.md

Structure (Reference-First approach):

1. **Header** — Title, Go Reference badge, MIT License badge
2. **Introduction** — One-liner description, feature bullet list (including AI skill)
3. **Installation** — `go install` only
4. **Quick Start** — Example YAML + validate/render/watch commands
5. **CLI Commands** — `render`, `validate`, `watch` with flags tables
6. **YAML Format Reference** — Root structure + all 9 primitives with full property docs:
   - section, cards, card, flow, grid, steps, info, connector, note
7. **Themes** — YAML config and CLI flag overrides
8. **AI Skill** — Description, auto-discovery, `npx skills add` global install
9. **License** — MIT, link to LICENSE file

## Design Decisions

- **Single README**: Everything in one file for discoverability. Standard for Go CLI tools.
- **MIT License**: Most permissive, standard for developer tools.
- **go install only**: Simplest installation method, no release infrastructure needed yet.
- **Comprehensive YAML reference**: All 9 primitives documented inline with examples.
