---
ctx:
  tags: [standards, tracking, experimentation]
  purpose: Record this project's analytics model so instrumentation changes remain consistent with existing tracking.
  use_when:
    - adding or modifying analytics events
    - adding or modifying page views or screen views
    - reviewing experiment instrumentation or reporting assumptions
  fill:
    - required: analytics platform and ownership
    - required: taxonomy and naming rules for events, page views, and screen views
    - required: important tracked entities and properties
    - required: continuity rules for renames, semantic changes, and versioning
    - optional: experimentation and warehouse conventions
---

# Analytics

> Project-specific. Do not assume generic event names or metrics. Document the conventions this project actually uses.
>
> Analytics is not only about discrete events. For many products, page views,
> screen views, user actions, funnels, experiments, and derived business metrics
> all need to remain interpretable over time.
>
> Consistency matters more than local convenience. When naming or redefining
> tracking, preserve semantic continuity unless there is a clear reason to break it.
