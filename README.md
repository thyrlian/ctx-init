# ctx-init

Bootstrap a tool-agnostic AI context system for your project.

AI coding assistants are only as good as the context they have about your project.  `ctx-init` scaffolds a structured `.context/` directory in any project, giving AI agents a consistent, navigable source of truth they can load progressively on demand, not everything at once, regardless of which tool you use (OpenAI Codex, Claude Code, Google Antigravity, etc.).

Beyond context management, the bundled templates also serve as a lightweight design blueprint, prompting you to think through every dimension of your project from the start: architecture, product, conventions, workflows, and more.

## How It Works

1. A [`manifest.yml`](assets/manifest.yml) defines which context files to include and how to organize them
2. `ctx-init` copies the files into `.context/` in your target project
3. A [`_INDEX.md`](assets/context/_INDEX.md) is generated as an entry point for AI agents
4. Each `.md` file receives a unique `ctx-id` token, as a proof-of-read that agents must include in responses to confirm they actually loaded the file

## Quick Start

```bash
# Initialize with the standard preset (recommended)
go run ./cmd/ctx-init/ -out /path/to/your/project

# Preview what would happen without writing anything
go run ./cmd/ctx-init/ -out /path/to/your/project -dry-run

# Use a different preset
go run ./cmd/ctx-init/ -out /path/to/your/project -preset minimal

# Overwrite existing context files
go run ./cmd/ctx-init/ -out /path/to/your/project -force
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-out` | *(required)* | Target project directory |
| `-preset` | `standard` | Context preset: `minimal`, `standard`, `full` |
| `-dry-run` | `false` | Preview actions without writing any files |
| `-force` | `false` | Overwrite existing destination files |
| `-manifest` | `assets/manifest.yml` | Path to the context manifest file |

## Presets

| Preset | Description |
|--------|-------------|
| `minimal` | Core files only — just `ai_protocol.md` |
| `standard` | Full working set: product, standards, architecture, workflows |
| `full` | Everything, including prompts and ADR examples |

## Output Structure

Running `ctx-init` creates a `.context/` directory in your target project:

```
.context/
├── _INDEX.md                    ← auto-generated entry point for AI agents
├── ai_protocol.md               ← AI interaction rules (always loaded)
├── product/
│   ├── vision.md
│   ├── roadmap.md
│   ├── features.md
│   └── glossary.md
├── standards/
│   ├── project_conventions.md
│   ├── tech_stack.md
│   └── ...
├── architecture/
│   ├── system_overview.md
│   ├── adr/
│   └── ...
└── workflows/
    ├── feature_dev.md
    └── ...
```

### `_INDEX.md`

Regenerated on every run.  Lists all files in the active preset with their tags, so AI agents can decide which files to load on demand without opening them all:

```markdown
- [product/vision.md](product/vision.md) `product`, `global`
- [standards/observability.md](standards/observability.md) `standards`, `logging`, `monitoring`, `alerting`
```

Tags signal load priority:
- `core` — must always load (defines mandatory rules)
- `global` — always-load context relevant to every task
- everything else — load on demand based on the task at hand

### `ctx-id` — Proof of Read

Every `.md` file gets a unique token appended on copy:

```markdown
<!-- ctx-id: a3f8c2d1e4b09f7e -->
```

AI agents are expected to echo back the `ctx-id` of files they have read.  This makes it verifiable that context was actually loaded, not hallucinated.

## CI / Docker

A Docker-based script is provided for CI and local container runs:

```bash
# Run the CLI (default mode)
./scripts/ci.sh

# Or explicitly
./scripts/ci.sh run

# Run tests
./scripts/ci.sh test
```

## License

Copyright © 2026 [Jing Li](https://github.com/thyrlian)

Released under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

See the [LICENSE](./LICENSE) file for full details.
