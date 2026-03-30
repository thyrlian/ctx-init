---
ctx:
  tags: [workflow, template]
  purpose: Provide a reusable task shape that makes work items specific, bounded, and verifiable.
  use_when:
    - writing a ticket, issue, task, or implementation brief
    - refining an ambiguous request before development starts
---

# Task Template

Use this structure when writing tasks, tickets, or issues. A well-formed task reduces ambiguity before work starts and makes completion verifiable.

---

## Template

```
## Context

[Why does this task exist? What situation or problem prompted it?
Link to related issues, PRs, or decisions if relevant.]

## Goal

[What is the desired outcome? State the end state, not the steps to get there.]

## Constraints

[What must not change? What are the boundaries of this task?
Examples: must not break existing API, must stay under X ms p99, must not require a migration.]

## Acceptance Criteria

- [ ] [Specific, verifiable condition 1]
- [ ] [Specific, verifiable condition 2]
- [ ] [Tests added or updated to cover the change]

## Out of Scope

[What could reasonably be assumed to be included but is not?]
```

---

## Guidance

- Acceptance criteria must be independently checkable.
- Goal should describe the desired end state, not the implementation approach.
- Out of scope prevents adjacent work from quietly expanding the task.

---

## Example

```
## Context

The `/api/export` endpoint currently loads the entire dataset into memory before
streaming it to the client. For large accounts (>500k rows), this causes OOM
crashes under load. Tracked in #1203.

## Goal

The export endpoint can handle datasets of 1M+ rows without exceeding 512MB
memory per request, measured under sustained load.

## Constraints

- The response format (NDJSON) must not change.
- Existing export API contract must be preserved.
- No changes to the export scheduling logic.

## Acceptance Criteria

- [ ] Load test with 1M row dataset stays under 512MB resident memory
- [ ] P95 latency for export requests does not regress vs. baseline
- [ ] Integration tests cover partial failure mid-stream
- [ ] Existing export tests still pass

## Out of Scope

- Adding new export formats
- Changing export scheduling or triggering logic
- UI changes to the export flow
```
