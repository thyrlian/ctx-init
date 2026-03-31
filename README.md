# ctx-init

Bootstrap a tool-agnostic AI context system for your project.

AI coding assistants are only as good as the context they have about your project.  `ctx-init` scaffolds a structured `.context/` directory in any project, giving AI agents a consistent, navigable source of truth they can load progressively on demand, not everything at once, regardless of which tool you use (OpenAI Codex, Claude Code, Google Antigravity, etc.).

Beyond context management, the bundled templates also serve as a lightweight design blueprint, prompting you to think through every dimension of your project from the start: architecture, product, conventions, workflows, and more.

## Design Philosophy

### Bootstrap, Not Babysit

`ctx-init` is a bootstrapper, not a daemon.  Run it once to scaffold the `.context/` directory and its context files in your project, everything after is yours.  No re-runs required, though following the conventions is expected.

### Modular by Default

Context should be broken into small, focused files instead of one monolithic document.  A modular structure makes the system easier to navigate, easier to evolve, and easier to adapt to different AI tools without rewriting the source of truth.  It also enables progressive loading, since agents can read only the files relevant to the task at hand.

### Single Responsibility

Each context file should do one job well.  Product files describe product facts, standards define engineering rules, architecture files explain structure and boundaries, and workflows capture execution patterns.  This separation keeps context clearer for both humans and AI agents.

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

Together they form a modular, tool-agnostic source of truth that requires no external tooling to interpret.

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
go run ./cmd/ctx-init/ -project-root /path/to/your/project

# Preview what would happen without writing anything
go run ./cmd/ctx-init/ -project-root /path/to/your/project -dry-run

# Use a different preset
go run ./cmd/ctx-init/ -project-root /path/to/your/project -preset minimal

# Overwrite existing generated context files
go run ./cmd/ctx-init/ -project-root /path/to/your/project -force

# Use a custom manifest file
go run ./cmd/ctx-init/ -project-root /path/to/your/project -manifest path/to/manifest.yml

# Generate a specific tool adapter entrypoint (e.g.: Claude Code)
go run ./cmd/ctx-init/ -project-root /path/to/your/project -adapter claude
```

When using a custom manifest file, any relative paths inside that manifest are resolved relative to the manifest file's location, not the current working directory.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-project-root` | *(required)* | Target project root directory |
| `-preset` | `standard` | Context preset: `minimal`, `standard`, `full` |
| `-dry-run` | `false` | Preview actions without writing any files |
| `-force` | `false` | Overwrite existing generated files; for adapters, only replaces existing `*.ctx-init.md` fallback files |
| `-manifest` | `assets/manifest.yml` | Path to the context manifest file |
| `-adapter` | *(optional)* | Generate an adapter entrypoint, currently: `claude`, `codex` |

## Presets

| Preset | Description |
|--------|-------------|
| `minimal` | Core files only, just `ai_protocol.md` |
| `standard` | Full working set: product, standards, architecture, workflows |
| `full` | Everything defined in the manifest, including ADR templates and optional sections |

## Adapters

Adapters generate tool-specific entrypoints that bridge a tool to the `.context/` system.

Common behavior:
- Each adapter checks only its official candidate locations within the project root.
- If the selected primary file already exists, `ctx-init` generates a `*.ctx-init.md` fallback file instead and tells you to append or merge that content manually.
- `-force` does not overwrite an existing primary tool file; it only allows `ctx-init` to replace an existing `*.ctx-init.md` fallback file.
- The relative path to `.context/ai_protocol.md` inside the generated AI agent instructions file is adjusted automatically based on the final output location.

### Claude Code

The `claude` adapter checks [`./.claude/CLAUDE.md`, `./CLAUDE.md`] and generates `./.claude/CLAUDE.md` if neither file exists.

**TODO**: Claude slash-command generation is intentionally left for a future adapter version.

[Official Reference](https://code.claude.com/docs/en/memory)

### OpenAI Codex

The `codex` adapter checks [`./AGENTS.md`] and generates `./AGENTS.md` if it does not exist.

**TODO**: Support scoped `AGENTS.md` generation for subdirectories in a future adapter version.

[Official Reference](https://developers.openai.com/codex/guides/agents-md)

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
