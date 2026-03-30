---
ctx:
  tags: [workflow]
  purpose: Define the expected path for delivering a feature from task definition through merge and release.
  use_when:
    - starting a new feature or substantial change
    - checking the expected sequence of development steps
  fill:
    - required: feature workflow from task to merge
    - optional: feature flag tooling and rollout rules for this project
    - optional: definition of done adjustments for this project
---

# Feature Development

## Baseline

Feature work should start from a clear task, move in small reviewable increments, and leave the main branch in a releasable state.

## Project Feature Workflow

Document the normal sequence for this project:

1. how a task is defined and approved
2. how work is branched or otherwise started
3. what must happen before review is requested
4. what must happen before merge
5. how rollout or feature-flag enablement is handled after merge

## Definition of Done

Record any project-specific conditions that must be true before feature work is considered complete in this project.

If this project has no additional definition-of-done rules, default to the gates documented in `testing.md`, `code_review.md`, `continuous_integration.md`, and `deployment.md`.
