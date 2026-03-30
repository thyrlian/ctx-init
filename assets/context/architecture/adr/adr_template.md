---
ctx:
  tags: [architecture, adr, template]
  purpose: Provide the standard shape for recording a significant architectural decision and its rationale.
  use_when:
    - writing a new ADR
    - reviewing whether an architectural decision is documented clearly enough
---

# ADR-NNNN: [Title — short, imperative phrase describing the decision]

**Date:** YYYY-MM-DD
**Status:** Draft | Proposed | Accepted | Rejected | Deprecated | Superseded by [ADR-NNNN]
**Participants:** [Names or roles of contributors to the discussion]
**Approvers:** [Names or roles of final approvers]

---

## Context

[What situation prompted this decision? Describe the forces at play: technical constraints,
business requirements, team constraints, prior decisions. Be specific enough that a new team
member can understand why this decision was necessary without additional context.

Do not describe the decision here — only the situation that required one.]

## Decision

[What was decided? State it clearly and directly in one or two sentences, then elaborate
as needed. This section should be unambiguous — a reader should know exactly what was
chosen and what was rejected.]

## Alternatives Considered

[What other options were evaluated? For each:
- What was the option
- Why it was considered
- Why it was not chosen

Documenting rejected alternatives is as important as documenting the chosen one.
It prevents relitigating decisions when someone rediscovers the same alternatives later.]

| Option | Pros | Cons | Reason not chosen |
|--------|------|------|-------------------|
| [Option A — the chosen one] | ... | ... | — |
| [Option B] | ... | ... | ... |
| [Option C] | ... | ... | ... |

## Consequences

[What are the outcomes of this decision — both positive and negative?

- What becomes easier or possible?
- What becomes harder or constrained?
- What follow-up actions or decisions does this create?
- What technical debt, if any, is being accepted?]

## References

- [Link to relevant issue, PR, RFC, or discussion]
- [Link to related ADRs]
- [Link to external documentation or prior art]
