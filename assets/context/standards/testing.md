---
ctx:
  tags: [standards]
  purpose: Record the project's actual testing strategy, quality gates, and ownership so new work follows the local test model instead of generic advice.
  use_when:
    - adding or reviewing tests
    - deciding what kinds of tests are required for a change
  fill:
    - required: supported test layers and what belongs in each
    - required: mandatory local and CI test gates
    - required: mocking, fixture, and external dependency policy
    - optional: coverage thresholds, flaky test handling, and performance test policy
---

# Testing Standards

## Baseline

Tests should be deterministic, readable, and meaningful enough that failures provide trustworthy signal. Prefer testing behavior over internal implementation details.

## Project Testing Strategy

Document the local testing model for this project:

- Test layers: which layers exist and when each is required.
- Scope boundaries: what counts as unit, integration, end-to-end, contract, smoke, or performance testing here.
- Dependency policy: when mocks are expected, when real infrastructure is required, and how third-party systems are handled.
- Test data: fixtures, builders, seeded environments, and cleanup expectations.
- Required commands: what must pass locally and what is enforced in CI.
- Quality gates: merge-blocking suites, coverage expectations if any, and rules for skipped or quarantined tests.

## Ownership and Exceptions

Record:

- who owns test infrastructure and flaky test triage
- where test environments are defined
- any approved exceptions to normal coverage or test gate expectations
