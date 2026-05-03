---
ctx:
  tags: [workflow]
  purpose: Define the review conventions that changes must satisfy before merge.
  use_when:
    - preparing a PR for review
    - reviewing another engineer's change
  fill:
    - required: review comment taxonomy and blocking criteria
    - required: expected PR scope and self-review requirements
    - optional: reviewer SLA, escalation path, and special reviewer roles
---

# Code Review

## Baseline

Reviews should optimize for correctness, clarity, risk reduction, and shared understanding. Feedback should be specific, scoped to the change, and explicit about whether it blocks merge.

## Project Review Rules

Document:

- the prefixes or conventions used for blocking, non-blocking, and clarifying feedback
- what must be checked before approval
- expected PR size or scope limits
- whether self-review is required before requesting review
- whether approval is allowed before required CI, lint, format, and test checks pass
- who must review security-sensitive, data-sensitive, or architecture-significant changes

## Resolution and Escalation

Record:

- how disagreements are resolved
- when to bring in an additional reviewer or owner
- whether unresolved comments block merge by default
