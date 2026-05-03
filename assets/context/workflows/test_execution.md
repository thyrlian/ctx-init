---
ctx:
  tags: [workflow]
  purpose: Record the project's actual expectations for running tests locally and interpreting test failures in CI.
  use_when:
    - deciding which tests to run before pushing or opening a PR
    - interpreting CI failures or skipped tests
  fill:
    - required: mandatory local test commands and environments
    - required: CI failure and retry policy
    - optional: skipped test policy, performance test policy, and staging validation rules
---

# Test Execution

## Baseline

Before asking others to trust a change, run the tests that this project considers mandatory for that type of change. CI failures should be treated as real signal until proven otherwise.

## Project Test Execution Rules

Document:

- which commands must be run locally before opening a PR
- when integration or end-to-end tests are required
- what test environments exist and how they are provisioned
- what counts as a legitimate retry versus a failure that must be fixed
- how skipped, quarantined, or flaky tests are tracked

## Exceptions

Record any approved shortcuts, non-blocking suites, or exceptional paths used for large migrations, incidents, or release work.
