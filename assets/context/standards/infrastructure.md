---
ctx:
  tags: [standards, infra, platform]
  purpose: Describe the runtime environments and infrastructure shape that the project depends on.
  use_when:
    - reasoning about deployment targets or environment differences
    - making infra-dependent implementation or debugging decisions
  fill:
    - required: environments
    - required: compute model
    - required: storage and networking basics
    - optional: infrastructure as code locations and ownership
---

# Infrastructure

> Project-specific. Capture the parts of the platform that affect how the code is built, deployed, or operated.
