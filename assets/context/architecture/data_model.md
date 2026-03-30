---
ctx:
  tags: [architecture]
  purpose: Describe the core domain entities, relationships, and invariants of the system.
  use_when:
    - designing schema changes
    - writing queries or business logic tied to domain entities
  fill:
    - required: core entities
    - required: relationships
    - required: key invariants
    - optional: schema source of truth
---

# Data Model

> Project-specific. Focus on stable domain concepts and constraints, not every column in every table.
