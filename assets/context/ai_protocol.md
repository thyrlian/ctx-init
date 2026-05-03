---
ctx:
  tags: [core]
  purpose: Define the mandatory protocol for how AI systems must discover, read, and cite project context.
  use_when:
    - starting any AI-assisted work in the project
    - deciding which context files to load and how to declare them
---

# AI Collaboration Protocol

This file defines the rules for AI-assisted development in this project.
Follow these rules in every AI coding session for this project.

---

## Proof-of-Read Protocol

Every `.context/` file deployed by `ctx-init` contains a unique proof-of-read token
as a Markdown comment: `<!-- ctx-id: <16-hex-char token> -->`

**Rules for AI systems:**

1. When you read any `.context/` file, locate its `ctx-id` comment line.
2. At the start of your response, declare **only the files you actually opened and used**.
   Use exactly one `CTX-READ:` line per file, with paths relative to `.context/`:
   `CTX-READ: <relative-path> = <ctx-id>`
3. If you cannot find the `ctx-id` in a file, use `NOT_FOUND` as the value.
4. **Never fabricate a `ctx-id`.** If you did not read the file, do not claim you did.
5. **Never store or cache a `ctx-id`**. **Never reuse one from memory**, prior sessions, previous responses, or cached context.
   Only declare a `ctx-id` after freshly reading the corresponding file in the current task.

---

## Context Navigation

Start every session by reading `_INDEX.md` in this directory.
It lists the files available for the current preset. Use it as a map, not a reading list.
Always load files marked `core` or `global`:
- `core` files define mandatory rules.
- `global` files provide project-wide context relevant to every task.

For all others, use ancestor directory names within `.context/`, file name, and tags
together to infer relevance and load them on demand.

---

## File Frontmatter

Some files carry a `ctx:` YAML frontmatter block at the top. It is the authoritative
metadata for that file and takes the following format:

```yaml
---
ctx:
  tags: [decision, gateway]             # load-priority tags
  points_to:                            # explicit paths to follow (dirs or files)
    - docs/adr/
  include:                              # glob patterns for scattered files
    - "**/ADR-*.md"
---
```

**Rules for AI systems:**

- `tags` use the same semantics as `_INDEX.md` and inform relevance decisions.
- Effective tags are the union of:
  - tags inherited from `manifest.yml` sections and file entries
  - tags declared in the file's own `ctx.tags`
- Repeating a tag in both `manifest.yml` and frontmatter is allowed. Treat duplicates as a single tag after merging.
- `points_to` entries are paths relative to the project root. Read the target directory or file as part of this file's context.
- `include` entries are glob patterns. Scan the project for matches and read them as part of this file's context.
- Files tagged `gateway` in `_INDEX.md` always carry `points_to` and/or `include`. Open them and follow those fields; do not treat them as self-contained content.
- If `points_to` or `include` entries are still placeholder comments, skip them and note that the user has not configured this gateway yet.
- For files not listed in `_INDEX.md`: any file whose frontmatter contains at least one non-placeholder entry in `points_to` or `include` is treated as a gateway.
  A `gateway` tag without either field has no effect.

---

## Collaboration Rules

- Prefer editing existing files over creating new ones.
- Do not make changes outside the scope of the current task.
- If requirements are unclear, ask before implementing.
- All key decisions must be traceable to a document in this `.context/` directory.
- When making decisions, cite the relevant `.context/` file paths that informed them.

---

## Context Evolution

Context files should evolve with the project. AI systems should help turn repeated decisions,
conventions, and lessons learned into written project context.

**Rules for AI systems:**

- When you detect a recurring convention, clarified requirement, or stable team preference, propose updating the relevant `.context/` file instead of leaving that knowledge only in the current thread.
- Prefer updating an existing context file over creating a new one.
- Do not silently rewrite project conventions. Ask for user approval before making substantive `.context/` changes that codify new norms, rules, or decisions.
- If the user approves, update the relevant `.context/` file so the knowledge becomes durable and reusable.
- If the user wants the documented context preserved in version control, stage and commit those `.context/` changes with a clear commit message.
- Do not commit context changes without explicit user approval.
