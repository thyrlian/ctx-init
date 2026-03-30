---
ctx:
  tags: [workflow, release]
  purpose: Record the deployment, rollback, and release-control rules that actually apply to this project.
  use_when:
    - preparing a production or staging release
    - changing rollout, migration, or rollback procedures
  fill:
    - required: deployment strategy and environments
    - required: rollback policy and post-deploy verification steps
    - required: release approvals and deployment authorization
    - required: migration and schema change safety rules
    - optional: feature flag rollout, freeze windows, and emergency exceptions
---

# Deployment

## Baseline

Deployments should be repeatable, observable, and reversible. Release decisions should optimize for controlled change, not speed at any cost.

## Project Deployment Rules

Document:

- environments and how code moves between them
- rollout strategy used by this project
- who can approve or trigger deployments
- rollback triggers, rollback owner, and rollback mechanics
- post-deploy verification steps and observation window
- database or schema migration safety rules

## Exceptions and Freeze Rules

Record:

- change freeze windows
- emergency deployment path
- any releases that follow a different approval or rollout model
