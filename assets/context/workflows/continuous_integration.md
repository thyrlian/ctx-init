---
ctx:
  tags: [workflow, ci]
  purpose: Record the CI pipelines, merge gates, and failure handling rules that actually apply in this repository.
  use_when:
    - modifying CI pipelines or branch protections
    - deciding which checks should gate merges
  fill:
    - required: pipeline stages and required checks
    - required: trigger strategy for branches, PRs, main, and scheduled runs
    - required: which automated checks block approval or merge
    - required: failure ownership and retry policy
    - optional: caching, parallelization, environment parity, and quarantine rules
---

# Continuous Integration

## Baseline

CI should provide fast, trustworthy feedback. Required checks should be stable, actionable, and difficult to bypass.

## Project CI Rules

Document:

- the stages that run in CI
- which checks are merge-blocking
- what runs on feature branches, pull requests, main, and scheduled workflows
- which failures may be retried and which must be fixed before rerun
- who owns broken pipelines and flaky jobs

## CI Environment and Exceptions

Record:

- runtime versions, container images, and critical environment assumptions
- approved quarantine or flaky-job handling
- any checks that are advisory rather than blocking
