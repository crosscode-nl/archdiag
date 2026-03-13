# archdiag

[![Go Reference](https://pkg.go.dev/badge/github.com/crosscode-nl/archdiag.svg)](https://pkg.go.dev/github.com/crosscode-nl/archdiag)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Tests](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/crosscode-nl/archdiag/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/crosscode-nl/archdiag)](https://github.com/crosscode-nl/archdiag/releases/latest)

A YAML-to-HTML architecture diagram generator. Define your infrastructure and system architecture in YAML, render it to self-contained HTML with live reload support.

**Features:**
- 9 composable diagram primitives (sections, cards, flows, grids, steps, info, connectors, notes)
- Dark and light themes
- Watch mode with live reload (SSE-based)
- Single self-contained HTML output — no external dependencies
- Validate diagrams before rendering
- AI skill for generating diagrams from natural language descriptions (Claude Code, Gemini CLI)

## Installation

Requires [Go](https://go.dev/doc/install) 1.24 or later.

```bash
go install github.com/crosscode-nl/archdiag/cmd/archdiag@latest
```

## Quick Start

Create a file called `diagram.yaml`:

```yaml
diagram:
  title: "My Service"
  subtitle: "Architecture Overview"
  elements:
    - section:
        name: "Frontend"
        color: blue
        children:
          - cards:
              - name: "Web App"
                color: blue
              - name: "Mobile App"
                color: blue
    - connector:
        direction: down
        text: "REST API"
    - section:
        name: "Backend"
        color: green
        children:
          - card:
              name: "API Gateway"
              color: green
```

Validate and render:

```bash
# Validate your diagram
archdiag validate diagram.yaml

# Render to HTML
archdiag render diagram.yaml

# Or watch for changes with live reload
archdiag watch diagram.yaml --open
```

## Examples

The [`examples/`](examples/) directory contains 10 progressively complex diagrams — from a single `note` to a full architecture using all 9 primitives. See [`examples/README.md`](examples/README.md) for the full list.

## CLI Commands

### `archdiag render`

Render YAML file(s) to self-contained HTML.

```bash
archdiag render <path> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output directory | Same as input |
| `--light` | Use light theme | |
| `--dark` | Use dark theme | |

`<path>` can be a single `.yaml` file or a directory of `.yaml` files.

### `archdiag validate`

Validate YAML file(s) without rendering. Reports errors to stderr and exits with code 1 on failure.

```bash
archdiag validate <path>
```

### `archdiag watch`

Watch YAML file(s) for changes and serve with live reload via Server-Sent Events.

```bash
archdiag watch <path> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-p, --port` | HTTP port | `3210` |
| `--open` | Open browser automatically | |
| `--light` | Use light theme | |
| `--dark` | Use dark theme | |

## YAML Format Reference

### Root Structure

```yaml
diagram:
  title: string          # required
  subtitle: string       # optional
  theme: dark | light    # default: dark
  palette:               # optional custom colors
    colorName: "#RRGGBB"
  elements:              # sequence of primitives
    - primitive:
        property: value
```

Colors can be referenced by palette name or as inline hex codes (`#RRGGBB`) throughout the diagram.

### Primitives

#### `section`

Bordered container for grouping elements. Can be nested.

```yaml
- section:
    name: "Kubernetes Cluster"    # required
    color: blue                   # optional
    border: solid                 # solid | dashed (default: solid)
    tag: "production"             # optional label
    placeholder: "..."            # optional placeholder text
    span: full                    # full | half (default: full)
    children:                     # nested primitives
      - card:
          name: "Pod"
```

#### `cards`

Flex-wrapped row of multiple cards.

```yaml
- cards:
    - name: "Service A"           # required
      color: green
      subtitle: "v2.1"
      version: "2.1.0"
      details: "Handles auth"
      badges:
        - "REST"
        - "gRPC"
      footer: "Port 8080"
      groups:
        - name: "Endpoints"
          items:
            - "/login"
            - "/logout"
    - name: "Service B"
      color: green
```

#### `card`

Single card component.

```yaml
- card:
    name: "API Gateway"           # required
    color: green
    span: full                    # full | half (default: half)
    subtitle: "nginx"
    version: "1.25"
    details: "Reverse proxy"
    badges:
      - "HTTPS"
    footer: "Port 443"
    groups:
      - name: "Routes"
        items:
          - "/api/*"
          - "/health"
```

#### `flow`

Horizontal step sequence with connectors.

```yaml
- flow:
    name: "Deploy Pipeline"       # required
    color: purple
    connector: arrow              # arrow | plus | none (default: arrow)
    suffix: "Done"                # optional end label
    steps:                        # single row
      - "Build"
      - "Test"
      - text: "Deploy"            # rich step with details
        color: green
        details: "to production"
    rows:                         # multi-row alternative
      - ["Build", "Test"]
      - ["Stage", "Deploy"]
```

#### `grid`

Multi-column layout for arranging primitives.

```yaml
- grid:
    columns: 3                    # required
    children:                     # sequence of any primitives
      - card:
          name: "Cell 1"
      - card:
          name: "Cell 2"
      - card:
          name: "Cell 3"
```

#### `steps`

Numbered step list.

```yaml
- steps:
    name: "Setup Guide"           # required
    color: blue
    start: 1                      # starting number (default: 1)
    items:
      - "Install dependencies"
      - "Configure environment"
      - "Run migrations"
```

#### `info`

Key-value metadata grid.

```yaml
- info:
    name: "Service Details"       # required
    color: blue
    items:
      - key: "Region"
        value: "eu-west-1"
      - key: "Runtime"
        value: "Go 1.24"
```

#### `connector`

Directional separator between sections.

```yaml
- connector:
    direction: down               # required: down | right | up | left
    text: "HTTPS"                 # optional label
    color: gray
    style: arrow                  # arrow | plain (default: arrow)
```

#### `note`

Footer callout or annotation.

```yaml
- note:
    text: "Diagram generated by archdiag"   # required
    style: muted                             # highlight | muted (default: muted)
```

## Themes

archdiag includes dark and light themes. The theme can be set in the YAML file or overridden via CLI flags.

```yaml
diagram:
  title: "My Diagram"
  theme: light    # dark (default) | light
```

```bash
# Override theme via CLI
archdiag render diagram.yaml --light
archdiag render diagram.yaml --dark
```

## AI Skill

archdiag ships with an AI skill (`skill/SKILL.md`) that enables AI assistants to generate architecture diagrams from natural language descriptions.

When working in the archdiag project directory, Claude Code and Gemini CLI automatically discover the skill. To install it globally:

```bash
npx skills add crosscode-nl/archdiag@archdiag-yaml -g -y
```

Then ask your AI assistant to create a diagram:

> "Create an architecture diagram for a Kubernetes deployment with an nginx ingress, two microservices, and a PostgreSQL database"

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
