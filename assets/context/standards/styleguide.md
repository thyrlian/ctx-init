---
ctx:
  tags: [gateway]
  purpose: Point AI systems to language, framework, or repository style guides maintained elsewhere in the project.
  use_when:
    - formatting or refactoring code in a language-specific area
    - locating the authoritative style rules for a code area
  points_to:
    - # e.g. docs/style/swift.md
    - # e.g. docs/style/kotlin.md
  include:
    - # e.g. "**/styleguide-*.md"
    - # e.g. "**/.editorconfig"
---

# Style Guides

> **Gateway** — fill in `points_to` with known style guide files, and `include` with
> glob patterns to discover others. Add one entry per language or framework. Remove placeholder comments when done.

Style guides enforce consistency across the codebase. This project may have guides for
multiple languages or frameworks — list them all above so AI agents can locate them.
