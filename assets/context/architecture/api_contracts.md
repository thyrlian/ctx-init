---
ctx:
  tags: [architecture]
  purpose: Document the interfaces this system exposes and consumes so changes do not break contract boundaries.
  use_when:
    - changing request or response shapes
    - integrating with internal or external services
  fill:
    - required: APIs exposed by this system
    - required: APIs consumed by this system
    - optional: schema source of truth and versioning rules
---

# API Contracts

> Project-specific. Use this as a concise map. Link to OpenAPI, protobuf, GraphQL, or other canonical specs when they exist.
