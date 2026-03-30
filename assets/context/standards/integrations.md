---
ctx:
  tags: [standards, third-party]
  purpose: Record third-party and external system integrations so changes account for their constraints and ownership.
  use_when:
    - modifying code that depends on external services
    - reviewing auth, rate limits, or data flow for an integration
  fill:
    - required: service and purpose
    - required: integration method
    - required: authentication and secret handling
    - optional: fallback handling, retries, rate-limit handling, and degradation behavior
    - optional: constraints, ownership, and runbooks
---

# Integrations

> Project-specific. One concise entry per external system is usually enough.
