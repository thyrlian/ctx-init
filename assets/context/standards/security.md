---
ctx:
  tags: [standards]
  purpose: Record the project's concrete security controls, owners, and decisions so security-sensitive work follows local policy instead of generic assumptions.
  use_when:
    - changing code that touches auth, secrets, input handling, external exposure, or sensitive data
    - reviewing a change for security risk
  fill:
    - required: secrets manager, auth provider, and credential rotation approach
    - required: sensitive data classification and encryption expectations
    - required: security scanning tools, gates, and review ownership
    - optional: API abuse controls, webhook verification, and infrastructure security controls
    - optional: incident response contacts, escalation path, and emergency procedures
---

# Security Standards

## Baseline

Use established security primitives and libraries. Apply least privilege, validate untrusted input at system boundaries, protect secrets, and treat sensitive data handling as a deliberate design concern.

## Project Security Decisions

Document the concrete controls used by this project:

- Identity and access: auth provider, protocol, session model, MFA requirements, and privileged access rules.
- Secrets: secrets manager, local development secret handling, rotation policy, and break-glass access.
- Sensitive data: what counts as PII, financial data, regulated data, or internal-only data in this project.
- Encryption: what must be encrypted at rest, in transit, or application-level.
- Scanning and review: SAST, dependency scanning, container or image scanning, and who must review security-sensitive changes.
- API and edge protections: rate limiting, CORS policy, API key lifecycle, bot or abuse controls, and public endpoint restrictions.
- Webhooks and external inputs: signature verification, replay protection, source allowlists, and timeout handling.
- Infrastructure controls: network boundaries, runtime isolation, image provenance, and platform-level policy enforcement.

## Ownership and Escalation

Record:

- security reviewer or owning team
- incident escalation path
- where audit logs and security alerts are reviewed
- how emergency revocation or hotfix deployment is performed

## Exceptions

Document any approved security exceptions, legacy constraints, or temporary compensating controls.
