---
ctx:
  tags: [standards, global]
  purpose: Record baseline engineering conventions that shape how changes are scoped, reasoned about, and documented.
  use_when:
    - deciding how to structure a change
    - checking project-specific conventions that are broader than language style
  fill:
    - optional: project-specific deviations from the baseline below
---

# Project Conventions

## Change Discipline

Keep changes focused and minimal. Avoid mixing unrelated refactoring with functional changes. If both are necessary, separate them into distinct commits or PRs.

A change should do one thing. If you find yourself writing "and also" in a change summary, split the work.

## Documentation Discipline

When behavior or architecture changes, update the relevant context files. Documentation and implementation should evolve together. Avoid leaving outdated descriptions.

If a decision was made during implementation that is not obvious from the code, record it. Future maintainers, including yourself, will need to understand why, not just what.

## Design Discipline

Avoid speculative abstractions for anticipated future use. At the same time, do not settle for code that merely works. Structure code clearly, extract meaningful functions, and keep related logic organized. Generalize only when multiple concrete cases appear.

Three similar code paths are a signal to consider abstraction. One is not.
