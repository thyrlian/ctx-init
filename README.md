# ctx-init

Bootstrap a tool-agnostic AI context system for your project.

AI coding assistants are only as good as the context they have about your project.  `ctx-init` scaffolds a structured `.context/` directory in any project, giving AI agents a consistent, navigable source of truth they can load progressively on demand, not everything at once, regardless of which tool you use (OpenAI Codex, Claude Code, Google Antigravity, etc.).

Beyond context management, the bundled templates also serve as a lightweight design blueprint, prompting you to think through every dimension of your project from the start: architecture, product, conventions, workflows, and more.

## Design Philosophy

### Bootstrap, Not Babysit

`ctx-init` is a bootstrapper, not a daemon.  Run it once to scaffold the `.context/` directory and its context files in your project, everything after is yours.  No re-runs required, though following the conventions is expected.

### Convention + Configuration

The system works through the following complementary mechanisms:

- [**`ai_protocol.md`**](assets/context/ai_protocol.md): the rulebook.  Defines how AI agents must behave, what to load, and how to navigate the context system.  Always read first.
- [**`_INDEX.md`**](assets/context/_INDEX.md): the map.  Auto-generated on every run, it lists every file in the selected preset with its tags.  AI agents use it to navigate what's available and decide relevance based on tags.
- **Frontmatter**: the self-contained metadata.  Key files carry a `ctx:` YAML block at the top of the file, readable by humans, AI, and scripts alike:

  ```yaml
  ---
  ctx:
    tags: [architecture, adr, gateway]
  ---
  ```

Together they form a single source of truth that requires no external tooling to interpret.

### Gateway Files

Some files in `.context/` are **entry points**, not content.  Rather than containing information directly, they carry `points_to` (explicit paths) and/or `include` (glob patterns) in their frontmatter, directing AI agents to related files that live elsewhere in the project, keeping `.context/` lean while the full project remains navigable.

  ```yaml
  ---
  ctx:
    points_to:
      - docs/adr/
      - docs/architecture/
    include:
      - "**/ADR-*.md"
      - "**/architecture/*.md"
  ---
  ```

Gateway files are signaled by the `gateway` tag in `_INDEX.md` or in their frontmatter.  Any file with `points_to` or `include` fields is treated as a gateway regardless of tag.

### Progressive Loading

Not every file is needed for every task.  Tags in `_INDEX.md` signal load priority:

- `core` / `global`: always load
- everything else: load on demand, guided by path, filename, and tags

This keeps token usage efficient and responses focused.

## How It Works

1. A [`manifest.yml`](assets/manifest.yml) defines which context files to include and how to organize them
2. `ctx-init` copies the files into `.context/` in your target project
3. A [`_INDEX.md`](assets/context/_INDEX.md) is generated as an entry point for AI agents
4. Each copied `.md` file receives a unique `ctx-id` token, as a proof-of-read that agents must include in responses to confirm they actually loaded the file

## Quick Start

```bash
# Initialize with the standard preset (recommended)
go run ./cmd/ctx-init/ -out /path/to/your/project

# Preview what would happen without writing anything
go run ./cmd/ctx-init/ -out /path/to/your/project -dry-run

# Use a different preset
go run ./cmd/ctx-init/ -out /path/to/your/project -preset minimal

# Overwrite existing generated context files
go run ./cmd/ctx-init/ -out /path/to/your/project -force

# Use a custom manifest file
go run ./cmd/ctx-init/ -out /path/to/your/project -manifest path/to/manifest.yml

# Generate a specific tool adapter entrypoint (e.g.: Claude Code)
go run ./cmd/ctx-init/ -out /path/to/your/project -adapter claude
```

When using a custom manifest file, any relative paths inside that manifest are resolved relative to the manifest file's location, not the current working directory.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-out` | *(required)* | Target project directory |
| `-preset` | `standard` | Context preset: `minimal`, `standard`, `full` |
| `-dry-run` | `false` | Preview actions without writing any files |
| `-force` | `false` | Overwrite existing generated files; for adapters, only replaces existing `*.ctx-init.md` fallback files |
| `-manifest` | `assets/manifest.yml` | Path to the context manifest file |
| `-adapter` | *(optional)* | Generate an adapter entrypoint, currently: `claude` |

## Presets

| Preset | Description |
|--------|-------------|
| `minimal` | Core files only, just `ai_protocol.md` |
| `standard` | Full working set: product, standards, architecture, workflows |
| `full` | Everything defined in the manifest, including ADR templates and optional sections |

## Adapters

Adapters generate tool-specific entrypoints that bridge a tool to the `.context/` system.

### Claude Code

Use the Claude adapter to generate a root-level `CLAUDE.md` that points Claude Code at `.context/ai_protocol.md`. Claude slash-command generation is intentionally left for a future adapter version.

Behavior:
- If `CLAUDE.md` does not exist, `ctx-init` generates it in the project root.
- If `CLAUDE.md` already exists, `ctx-init` generates `CLAUDE.ctx-init.md` instead and tells you to append or merge that content manually.
- `-force` does not overwrite an existing `CLAUDE.md`; it only allows `ctx-init` to replace an existing `CLAUDE.ctx-init.md` fallback file.

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
- `core` -> must always load (defines mandatory rules)
- `global` -> always-load context relevant to every task
- everything else -> load on demand based on the task at hand

### `ctx-id` -> Proof of Read

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
