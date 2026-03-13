# Examples

Architecture diagrams of the archdiag project itself, progressing from simple to complex.

| # | File | Primitives Used |
|---|------|-----------------|
| 1 | [01-hello-world.yaml](01-hello-world.yaml) | `note` |
| 2 | [02-project-overview.yaml](02-project-overview.yaml) | `cards`, `note` |
| 3 | [03-data-pipeline.yaml](03-data-pipeline.yaml) | `flow`, `note` |
| 4 | [04-cli-commands.yaml](04-cli-commands.yaml) | `grid`, `card`, `info` |
| 5 | [05-nine-primitives.yaml](05-nine-primitives.yaml) | `section`, `cards` (with badges) |
| 6 | [06-template-system.yaml](06-template-system.yaml) | `section`, `grid`, `cards`, `connector` |
| 7 | [07-watch-and-live-reload.yaml](07-watch-and-live-reload.yaml) | `flow`, `section`, `grid`, `cards`, `info`, `connector` |
| 8 | [08-color-system.yaml](08-color-system.yaml) | `section`, `flow`, `grid`, `card`, `info`, `connector`, `note` |
| 9 | [09-validation-rules.yaml](09-validation-rules.yaml) | `section`, `cards`, `card`, `connector`, `note` (with groups) |
| 10 | [10-full-architecture.yaml](10-full-architecture.yaml) | All 9 primitives |

## Render all examples

```bash
archdiag render examples/
```

This produces a `.html` file next to each `.yaml` file in the `examples/` directory.

## Render to a different directory

```bash
archdiag render examples/ -o out/
```

## Render a single example

```bash
archdiag render examples/03-data-pipeline.yaml
```

## Watch with live reload

Watch all examples and open in your browser:

```bash
archdiag watch examples/ --open
```

Watch a single file:

```bash
archdiag watch examples/10-full-architecture.yaml --open
```

## Validate without rendering

```bash
archdiag validate examples/
```
