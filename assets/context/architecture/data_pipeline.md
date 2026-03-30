---
ctx:
  tags: [architecture]
  purpose: Describe how data moves through the system and where transformations, retries, and failures occur.
  use_when:
    - changing ingestion, ETL, sync, or event processing logic
    - diagnosing latency, freshness, or data loss issues
  fill:
    - required: sources and destinations
    - required: pipeline stages
    - required: failure handling
    - optional: latency, throughput, and ownership
---

# Data Pipeline

> Project-specific. If the project has no meaningful pipeline, say so explicitly.
