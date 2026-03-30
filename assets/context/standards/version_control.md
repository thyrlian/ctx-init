---
ctx:
  tags: [standards, vcs, git]
  purpose: Record the repository's actual branch, commit, review, merge, and release rules so collaboration follows local conventions.
  use_when:
    - creating branches, commits, pull requests, tags, or releases
    - checking how changes are expected to move into the main branch
  fill:
    - required: branching model and branch naming rules
    - required: commit message convention
    - required: merge strategy and required approvals or checks
    - optional: release tagging and release note policy
---

# Version Control

## Branches

Document the branch model this repository actually uses:

- default branch
- whether work happens on short-lived branches, release branches, hotfix branches, or trunk only
- branch naming rules
- expectations for rebasing, syncing, and branch lifetime

## Commits and Pull Requests

Record:

- commit message convention
- expected PR size or scope
- required PR description fields
- whether self-review is expected before requesting review
- who must approve before merge

## Merge Policy

Record:

- allowed merge strategies
- whether history should stay linear
- when merge commits are allowed or prohibited
- whether branches must be up to date before merge
- whether branches are deleted after merge

## Releases and Protected Refs

If applicable, document:

- tagging scheme
- release branch or release cut process
- release notes format
- protected branches and protected tags
