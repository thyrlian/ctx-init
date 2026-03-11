# AI Collaboration Protocol

This file defines the rules for AI-assisted development in this project.
Follow these rules in every AI coding session for this project.

---

## Proof-of-Read Protocol

Every `.context/` file deployed by `ctx-init` contains a unique proof-of-read token
as a Markdown comment: `<!-- ctx-id: <16-hex-char token> -->`

**Rules for AI systems:**

1. When you read any `.context/` file, locate its `ctx-id` comment line.

2. At the start of your response, declare **only the files you actually opened
   and used**. Use exactly one `CTX-READ:` line per file, paths relative to `.context/`:
   `CTX-READ: <relative-path> = <ctx-id>`

3. If you cannot find the `ctx-id` in a file, use `NOT_FOUND` as the value.

4. **Never fabricate a `ctx-id`.** If you did not read the file, do not claim you did.

---

## Context Navigation

Start every session by reading `_INDEX.md` in this directory.
It lists the files available for the current preset — use it as a map, not a reading list.
Always load files marked `core` or `global`:
* `core` files define mandatory rules.
* `global` files provide project-wide context relevant to every task.

For all others, use ancestor directory names within `.context/`, file name, and tags
together to infer relevance and load them on demand.

---

## Collaboration Rules

- Prefer editing existing files over creating new ones.
- Do not make changes outside the scope of the current task.
- If requirements are unclear, ask before implementing.
- All key decisions must be traceable to a document in this `.context/` directory.
- When making decisions, cite the relevant `.context/` file paths that informed them.
