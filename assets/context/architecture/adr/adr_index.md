---
ctx:
  tags: [decision, gateway]
  purpose: Point AI systems to the project's ADR sources so architectural decisions can be discovered progressively.
  use_when:
    - tracing the rationale behind an architecture choice
    - locating ADRs that live outside .context
  points_to:
    - # e.g. docs/adr/
    - # e.g. legacy/decisions/
  include:
    - # e.g. "**/ADR-*.md"
    - # e.g. "**/decisions/*.md"
---

# Architecture Decision Records

> **Gateway** — fill in `points_to` with known ADR directories/files, and `include` with
> glob patterns for scattered ADRs. Remove placeholder comments when done.

ADRs capture significant architectural decisions: the context, options considered, and rationale.
Refer to `adr_template.md` for the standard format.
