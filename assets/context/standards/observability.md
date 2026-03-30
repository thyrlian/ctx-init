---
ctx:
  tags: [standards, logging, monitoring, alerting]
  purpose: Define the baseline logging, metrics, tracing, and alerting expectations for the project.
  use_when:
    - instrumenting code or debugging production behavior
    - designing monitoring and alerting for a service or workflow
---

# Observability Standards

## Philosophy

Observability is the ability to understand what a system is doing from its external outputs — without modifying it. The three pillars are logs, metrics, and traces. A system that cannot be observed cannot be debugged in production.

**Core principles:**

- Instrument code as you write it, not after something breaks.
- Observability data is for humans and machines. Structure it so both can consume it.
- High-cardinality context (user ID, request ID, trace ID) is what makes the difference between "something broke" and "this broke for this user on this request."
- Alert on symptoms that affect users, not on internal signals that may or may not matter.

---

## Logging

### Structure

Use structured logging. Every log entry should be a machine-parseable key-value record (JSON or equivalent), not a free-form string.

**Required fields on every log entry:**

| Field | Description |
|-------|-------------|
| `timestamp` | ISO 8601, UTC |
| `level` | `debug`, `info`, `warn`, `error` |
| `message` | Human-readable summary of what happened |
| `service` | Name of the service or component emitting the log |
| `trace_id` | Distributed trace ID (if applicable) |
| `request_id` | Per-request correlation ID (if applicable) |

Add domain-specific fields as needed (e.g. `user_id`, `order_id`, `job_id`). Fields should be consistent across the codebase — agree on names and stick to them.

### Log Levels

| Level | When to use |
|-------|-------------|
| `debug` | Detailed internal state, only useful during development. Must be disabled in production by default. |
| `info` | Normal operational events: service started, request received, job completed. |
| `warn` | Unexpected conditions that were handled and did not cause failure. Worth investigating. |
| `error` | Failures that require attention. Include enough context to diagnose: what was attempted, what failed, relevant IDs. |

Do not use `error` for expected user errors (e.g. validation failures). Those are `warn` or `info`.

### What to Log

**Log:**
- Service startup and shutdown with configuration summary (redact secrets)
- Incoming requests and their outcomes (status code, latency)
- Significant state transitions (job started, order placed, user promoted)
- External dependency calls and their outcomes (especially failures)
- Errors and unexpected conditions, with full context

**Do not log:**
- Passwords, tokens, API keys, or any credentials
- PII (names, emails, phone numbers, payment details) unless explicitly required and compliant with data policy
- High-volume noise that adds no diagnostic value (e.g. every cache hit)
- Stack traces to stdout in normal operation — capture them in error tracking instead

---

## Metrics

Metrics capture the quantitative state of the system over time. Use them to track trends, set SLOs, and drive alerts.

### The Four Golden Signals

For any service, instrument at minimum:

| Signal | What it measures |
|--------|-----------------|
| **Latency** | How long requests take. Track p50, p95, p99 — not just averages. |
| **Traffic** | Request rate. How much demand is the system handling? |
| **Errors** | Error rate. What fraction of requests are failing? |
| **Saturation** | How full is the system? CPU, memory, queue depth, connection pool usage. |

### Naming Conventions

Use a consistent naming scheme across the codebase. A common convention:

```
<service>.<subsystem>.<metric_name>
```

Examples: `api.auth.login_duration_seconds`, `worker.queue.depth`, `db.pool.connections_active`

- Use units in metric names (`_seconds`, `_bytes`, `_total`).
- Use `_total` suffix for counters.
- Keep names lowercase and underscore-separated.

---

## Distributed Tracing

For systems that span multiple services or processes, use distributed tracing to follow a request end-to-end.

**Rules:**

- Propagate trace context (trace ID, span ID) across all service boundaries via headers (e.g. W3C `traceparent`, B3).
- Create spans for meaningful units of work: incoming requests, outbound calls, database queries, background jobs.
- Annotate spans with relevant attributes: user ID, resource ID, operation name.
- Sample traces at a rate appropriate for your traffic volume. 100% sampling is usually only practical in development.

---

## Alerting

Alerts must be actionable. An alert that fires and requires no action is noise — and noise trains engineers to ignore alerts.

**Rules:**

- Alert on user-facing symptoms: elevated error rate, high latency, service unavailability.
- Do not alert on internal signals unless they reliably predict a user-facing symptom.
- Every alert must have a clear owner and a runbook or documented response procedure.
- Review and prune alerts regularly. An alert that has never fired — or always fires — needs attention.

**Alert severity levels:**

| Level | Meaning | Response |
|-------|---------|----------|
| Critical | User-facing impact now | Page on-call immediately |
| Warning | Degraded or trending toward impact | Investigate within business hours |
| Info | Informational only | No action required; review periodically |

---

## Health Checks

Every service must expose a health check endpoint.

**Minimum:**

- `GET /healthz` (or equivalent) — liveness check. Returns `200` if the process is running and able to serve requests. Should be fast and have no external dependencies.
- `GET /readyz` (or equivalent) — readiness check. Returns `200` only if the service is ready to accept traffic (dependencies connected, warmup complete). Used by load balancers and orchestrators.

Health check endpoints must not require authentication.
