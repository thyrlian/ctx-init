---
ctx:
  tags: [standards]
  purpose: Define the universal compliance baseline and record any project-specific regulatory obligations.
  use_when:
    - handling user data or operational controls with legal or audit implications
    - checking whether a change introduces compliance-sensitive behavior
  fill:
    - required: applicable regulations and jurisdictions
    - optional: project-specific control owners, evidence locations, and exceptions
---

# Compliance

## Universal Baseline

These principles apply regardless of industry, jurisdiction, or regulatory framework. They represent the minimum responsible behavior for any software system that handles user data or operates in a production environment.

### Data Minimization

Collect only the data you need to deliver the product. Do not log, store, or transmit data that serves no defined purpose. Every field in your data model should have an owner and a reason to exist.

When a data retention period expires, data must be deleted — not just hidden or marked as inactive.

Minimize not only the amount of data collected, but also the scope of access to it.
Each service, job, or component should access only the subset of data required for
its responsibility. Prefer ownership boundaries that reduce unnecessary exposure
and limit the blast radius of a bug, misuse, or compromise.

In distributed systems, this often means separating data access by bounded context
or service responsibility rather than allowing broad shared access across the
system. The goal is not microservices for their own sake; the goal is to keep
data access narrow, auditable, and proportionate to each component's role.

### Access Control

Access to production systems, databases, and sensitive data must be:
- Granted on a need-to-know basis
- Reviewed and revoked when roles change or team members leave
- Logged and auditable

No individual should have standing, unrestricted access to production data. Privileged access should require approval and generate an audit trail.

### Audit Logging

Security-relevant events must be logged and logs must be tamper-evident:
- Authentication events (login success, failure, logout)
- Authorization failures
- Changes to user permissions or roles
- Access to sensitive data
- Administrative actions

Audit logs must be retained for a defined period and stored outside the direct write path of the application.

### Incident Handling

Any suspected breach, unauthorized access, or data exposure must be escalated immediately according to the incident response process. Do not attempt to contain or resolve a suspected breach silently.

### Third-Party Dependencies

Third-party services that process or store user data must be reviewed for their own compliance posture before integration. Vendor agreements must reflect data processing responsibilities.

---

## Industry-Specific Requirements

> **Project-specific:** Add your applicable regulatory frameworks and their requirements here.
>
> Common examples:
>
> **GDPR (EU personal data)**
> - Legal basis for each data processing activity must be documented
> - Data subject rights: access, rectification, erasure, portability
> - Data Protection Impact Assessment (DPIA) required for high-risk processing
> - Breach notification to supervisory authority within 72 hours
>
> **HIPAA (US healthcare)**
> - Protected Health Information (PHI) must be encrypted at rest and in transit
> - Access to PHI must be logged and auditable
> - Business Associate Agreements required with vendors handling PHI
> - Minimum necessary standard for data access
>
> **PCI DSS (payment card data)**
> - Cardholder data must never be stored unless absolutely required
> - All cardholder data transmission must use TLS 1.2+
> - Regular vulnerability scans and penetration testing required
> - Access to cardholder data must be restricted and logged
>
> **DSA (EU Digital Services Act)**
> - Applicability depends on the type of service provided, especially for covered online intermediary services offered in the EU
> - Illegal content, goods, or services must have reporting and handling processes where the regulation applies
> - Users may need transparency and redress mechanisms for moderation or restriction decisions, depending on the service type
> - Online marketplaces may need trader traceability and related disclosures
> - Advertising, recommender systems, and platform operations may require additional transparency obligations
> - Very large platforms and search engines have additional risk assessment, mitigation, audit, and oversight duties
>
> **SOC 2**
> - Controls mapped to Trust Service Criteria (Security, Availability, Confidentiality, etc.)
> - Evidence collection for annual audit
> - Vendor management and risk assessment processes
>
> Replace this block with the specific requirements applicable to your product and markets.
