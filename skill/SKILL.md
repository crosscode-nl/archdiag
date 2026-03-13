---
name: archdiag-yaml
description: Use when the user asks to create, generate, or write an architecture diagram, infrastructure diagram, system diagram, service map, or archdiag YAML file. Triggers on requests like "diagram our VPC", "create an architecture diagram", "generate archdiag YAML for...", "visualize our architecture", "make a service map", "draw a system diagram". Also use when the user wants to modify or refine an existing archdiag YAML file.
---

# Generating archdiag YAML

Generate valid archdiag YAML files from natural-language descriptions of infrastructure. Produce a complete first draft immediately, then offer to refine.

## Workflow

1. Check that `archdiag` is available by running `which archdiag`. If not found, install it:
   - Requires Go 1.24+. Verify with `go version`.
   - Run `go install github.com/crosscode-nl/archdiag/cmd/archdiag@latest`.
   - If `go` is not installed either, tell the user they need Go first (https://go.dev/dl/) and stop.
2. Parse the user's description for: components, relationships, layers/boundaries, data flows.
2. Choose a palette of 4-6 semantic colors based on the domains mentioned.
3. Structure top-down: title/subtitle → palette → connectors between layers → sections for boundaries → cards/flows/steps inside.
4. Write the complete YAML file with the `Write` tool.
5. Run `archdiag validate <file>` to catch errors. If validation fails, fix and re-validate.
6. Run `archdiag render <file>` to produce HTML (outputs `.html` next to the `.yaml` by default; use `-o <dir>` for a different output directory; use `--light` or `--dark` to override the theme).
7. Present the result and offer refinement.

**Do NOT ask clarifying questions before generating.** Get something on screen fast. After the first draft, ask:

> "Here's the first draft. Want me to adjust anything? Common refinements: add/remove components, change detail level, split into multiple diagrams, adjust layout."

**Refinement rules:**
- Edit the existing YAML — do not regenerate from scratch.
- Re-validate and re-render after each change.
- Offer `archdiag watch <file> --open` for iterative refinement if multiple rounds are expected (also supports `--port <N>`, `--light`, `--dark`).

**Error handling:**
- Validation failures: fix the YAML automatically, re-validate.
- Vague descriptions: generate a minimal skeleton and explain what information would improve it.

---

## Schema Reference

### Root Structure

```yaml
diagram:
  title: string        # required
  subtitle: string     # optional
  theme: dark | light  # default: dark
  palette:
    <name>: "<hex>"    # e.g. cloud: "#4a9eff"
  elements:
    - <primitive>: { ... }
```

### Primitives

There are 9 primitives. Each element in `elements` is a single-key mapping: `- <type>: { properties }`.

#### section
Bordered container with a label. Nestable.

| Property    | Required | Default | Notes                                |
|-------------|----------|---------|--------------------------------------|
| name        | yes      |         | Section label                        |
| color       | no       |         | Palette name or hex                  |
| border      | no       | solid   | `solid` or `dashed`                  |
| tag         | no       |         | Right-aligned badge in label         |
| placeholder | no       |         | Shown instead of children            |
| span        | no       | full    | `full` or `half`                     |
| children    | no*      |         | Required unless `placeholder` is set |

#### cards
Flex-wrapped row of cards. Value is a YAML sequence (not a mapping):

```yaml
- cards:
    - name: "Card A"
      color: blue
    - name: "Card B"
```

Each item has the same properties as `card` (except `span`).

#### card
Single card component.

| Property | Required | Default | Notes                                                   |
|----------|----------|---------|---------------------------------------------------------|
| name     | yes      |         | Card heading                                            |
| color    | no       |         | Palette name or hex                                     |
| details  | no       |         | List of descriptive text lines                          |
| badges   | no       |         | `[{ text, color, style: filled|outlined }]`             |
| version  | no       |         | Version tag                                             |
| subtitle | no       |         | Monospace sub-label (service accounts, identifiers)     |
| footer   | no       |         | Muted footer with separator (dependencies)              |
| groups   | no       |         | `[{ label, mono: bool, items: [strings] }]`             |
| span     | no       |         | `full` to span grid columns (only inside a `grid`)      |

#### flow
Horizontal step sequence with connectors.

| Property  | Required | Default | Notes                          |
|-----------|----------|---------|--------------------------------|
| name      | yes      |         | Flow label                     |
| color     | no       |         | Label color                    |
| steps     | no       |         | Single-row shorthand           |
| connector | no       | arrow   | `arrow`, `plus`, `none`        |
| suffix    | no       |         | Text after last step           |
| rows      | no       |         | Multi-row: `[{ steps, connector, suffix }]` |

**FlowStep:** `"text"` (string) or `{ text, color, details: [strings] }`.

Use `steps` for single-row flows, `rows` for multi-row. Each row carries its own `steps`, `connector`, and `suffix`.

#### grid
Multi-column layout.

| Property | Required | Default | Notes               |
|----------|----------|---------|---------------------|
| columns  | yes      |         | Number of columns   |
| children | yes      |         | Primitives to lay out |

#### steps
Numbered step list.

| Property | Required | Default | Notes                  |
|----------|----------|---------|------------------------|
| name     | yes      |         | Phase label            |
| color    | no       |         | Label and number color |
| items    | yes      |         | List of step strings   |
| start    | no       | 1       | Starting number        |

#### info
Key-value metadata grid.

| Property | Required | Default | Notes                              |
|----------|----------|---------|-------------------------------------|
| name     | yes      |         | Section label                      |
| color    | no       |         | Label color                        |
| items    | yes      |         | `[{ key: string, value: string }]` |

#### connector
Directional separator between sections.

| Property  | Required | Default | Notes                         |
|-----------|----------|---------|-------------------------------|
| direction | yes      |         | `down`, `right`, `up`, `left` |
| text      | no       |         | Pill label                    |
| color     | no       |         | Text and border color         |
| style     | no       | arrow   | `arrow` or `plain`            |

#### note
Footer callout.

| Property | Required | Default | Notes                  |
|----------|----------|---------|------------------------|
| text     | yes      |         | Note content           |
| style    | no       | muted   | `highlight` or `muted` |

---

## Composition Patterns

### Structure
- **Top-down layering:** Use connectors (`direction: down`) between major layers — internet, cloud, VPC, workloads.
- **Sections for boundaries:** Each region, VPC, subnet, namespace, or logical group gets a section.
- **Nesting depth:** Max 2-3 levels. Deeper nesting hurts readability.
- **Grids for peers:** 2-3 column grids for equal-weight items (IAM roles, AZs, security layers).

### Color Semantics
- Define 4-6 palette colors with **semantic** names — not "blue"/"red" but "cloud", "network", "workload", "security".
- One color per logical domain, used consistently throughout.
- Section border color = label color = child accent color.
- Recommended hex values:
  - `#4a9eff` — cloud/infrastructure
  - `#22c55e` — networking/public
  - `#a855f7` — platform/operators
  - `#f59e0b` — critical infrastructure
  - `#ef4444` — security/private

### Primitive Selection

| Concept                        | Use                   |
|--------------------------------|-----------------------|
| Service / resource / component | Card (inside Cards)   |
| Request path / data flow       | Flow (arrow)          |
| Dependency chain               | Flow (plus)           |
| Deployment phases              | Steps                 |
| Configuration / specs          | Info                  |
| Layer separation               | Connector (down)      |
| Important warning              | Note (highlight)      |
| General footnote               | Note (muted)          |
| Logical boundary               | Section               |
| Equal-weight layout            | Grid (2-3 columns)    |

### Card Detail Levels
- **Minimal:** name + 1-2 details — supporting components
- **Standard:** name + details + badges — primary components
- **Rich:** name + subtitle + groups + footer — components with permissions, dependencies

### Diagram Sizing
- **Small** (5-10 elements): single section, flat
- **Medium** (10-20 elements): 2-3 top-level sections with children
- **Large** (20+): nested sections, connectors between layers, grids for width

---

## Example

A well-structured diagram using multiple primitives:

```yaml
diagram:
  title: "Cloud Infrastructure Overview"
  subtitle: "VPC, compute, and networking"
  palette:
    cloud: "#4a9eff"
    network: "#22c55e"
    compute: "#a855f7"
    security: "#ef4444"

  elements:
    - connector:
        direction: down
        text: "Internet"
        color: network
        style: arrow

    - section:
        name: "Cloud Region"
        color: cloud
        tag: "eu-west-2"
        children:
          - section:
              name: "Public Subnet"
              color: network
              border: dashed
              children:
                - cards:
                    - name: "Load Balancer"
                      color: network
                      details:
                        - "TLS termination"
                        - "Health checks"
                      badges:
                        - { text: "NLB", color: network, style: outlined }

          - section:
              name: "Private Subnet"
              color: compute
              children:
                - cards:
                    - name: "App Server"
                      color: compute
                      details:
                        - "2x instances"
                      subtitle: "app-service (default)"
                      groups:
                        - label: "Permissions"
                          mono: true
                          items:
                            - "s3:GetObject"
                            - "s3:PutObject"
                      footer: "Depends on: Database"
                    - name: "Database"
                      color: security
                      details:
                        - "Encrypted at rest"
                      badges:
                        - { text: "RDS", color: security, style: filled }

    - flow:
        name: "Request Flow"
        color: cloud
        steps:
          - { text: "Client", color: network }
          - { text: "LB", color: network }
          - { text: "App", color: compute }
          - { text: "DB", color: security }

    - flow:
        name: "Platform Dependencies"
        color: compute
        rows:
          - steps:
              - { text: "Ingress", color: network }
              - { text: "Service Mesh", color: compute }
              - "Backend"
            connector: arrow
            suffix: "mTLS"
          - steps:
              - { text: "Secrets" }
              - { text: "Config Maps" }
              - { text: "CRDs" }
            connector: plus

    - note:
        text: "All traffic encrypted in transit via TLS"
        style: highlight
```
