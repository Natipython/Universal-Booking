# Records, audit, evidence and lawful-request process

Status: **Draft — operating/control architecture; country-specific procedure requires licensed counsel**  
Depends on: data classification/retention, identity and permissions, legal-document register, incident response, launch country/entity and support model  
Blocks: production privileged support, real-money operations, formal disputes, sensitive integrations and external data disclosure

## Purpose

Ubi needs trustworthy operational records without turning routine logs into permanent surveillance. This document defines which records support booking, payment, security, accountability and lawful process; how their integrity and access are controlled; and how Ubi responds to government, court, emergency, civil and merchant requests.

An audit log does not prove that every underlying statement is true. It proves, to the designed assurance level, that a recorded actor/system performed or observed a specific event. Provider verification, user reports and staff conclusions remain distinct.

## Record families and authoritative ownership

| Family | Examples | Authoritative owner/source | Key caution |
|---|---|---|---|
| Transaction records | Booking/request/quote/policy snapshots, allocation and status transitions | Booking/domain system | Preserve history; do not overwrite disputes away |
| Financial records | Payment intent, reported/processing/verified status, provider references, refund/reversal/reconciliation | Payment ledger plus provider evidence | Never treat screenshots as provider verification |
| Contract/choice records | Terms/policy version, presentation, acceptance, consent/withdrawal | Legal-evidence service/domain event | Minimize metadata; retain the actual approved version |
| Identity/access records | Account, membership, role, recovery, MFA/step-up and session security events | Identity/authorization system | Exclude secrets and raw credentials |
| Privileged/support records | Approval, purpose, scope, accessed objects, action and revocation | Privileged-access system | Just-in-time, monitored and reviewed |
| Integration records | Connection/grant state, webhook receipt, validation, delivery/retry and provider outcome | Integration adapter/event journal | Tokens and message bodies are not general log fields |
| Safety/moderation records | Report, evidence pointer, triage, action, appeal and outcome | Trust-and-safety case system | Restrict access; preserve fairness/context |
| Security/incident records | Alerts, investigation timeline, affected scope, containment and notification decisions | Security incident system | Legal privilege/confidentiality rules require counsel |
| Technical telemetry | Performance, errors, traces, crash data and capacity | Observability platform | Sample/redact; short retention unless justified |
| Corporate/compliance records | Vendor review, training, decisions, approvals, tax/accounting and regulator correspondence | Designated business owner | Country/entity schedule controls retention |

Records must link by stable identifiers and evidence pointers, not by duplicating complete personal payloads into every system.

## Audit-event minimum contract

Risk-relevant events should include:

- unique event ID and occurred/recorded timestamps;
- tenant/organization scope where applicable;
- actor type and stable actor ID, or explicit unauthenticated/system/provider actor;
- authenticated session/client/service and assurance level where relevant;
- action, target type/ID and result;
- reason/purpose and approval/case/change reference when required;
- before/after field names or protected change reference for material mutations;
- source event/provider reference and idempotency/correlation/trace IDs;
- policy/permission/version evaluated;
- integrity-chain or protected-store metadata appropriate to risk.

Audit events must not contain passwords, OTP values, full tokens, cryptographic private material, complete card/mobile-money credentials, arbitrary message bodies, call media or unnecessary personal data.

## Events that require security-grade audit

- organization ownership, membership, role and permission changes;
- login recovery, MFA/step-up, phone/email/identity changes and session revocation;
- booking allocation/confirmation/cancellation/no-show and policy override;
- quote/price/deposit/refund/reversal/payment-verification/reconciliation changes;
- provider connection, OAuth grant/scope, token rotation, webhook and disconnect;
- export, deletion, legal hold, privacy request and retention exception;
- staff/support impersonation or privileged data access;
- moderation action, appeal and restoration;
- configuration/schema/migration/deployment/secret/production-access changes;
- incident declaration, evidence access and external disclosure.

Ordinary read auditing should be risk-based. Logging every customer screen view forever would increase privacy, cost and breach risk.

## Integrity, availability and access controls

- Application roles cannot update/delete their own finalized security audit events.
- Events are appended through a narrow service/path with server-generated identity context.
- High-risk logs use restricted storage, encryption, integrity checks/chaining/signing or equivalent tamper-evidence, and independently monitored export.
- “Tamper-evident” means unauthorized modification can reasonably be detected; it must not be marketed as mathematically immutable unless proven.
- Clock synchronization and original provider timestamps are retained with source and uncertainty.
- Duplicate/out-of-order/replayed events are tolerated through event IDs, idempotency and reconciliation.
- Backups/restores include audit continuity tests and do not silently reopen expired access.
- Access is least-privilege, purpose-limited and logged; investigators see only the case scope.
- Exports are encrypted, access-controlled, checksummed and time-limited.

## Record lifecycle

1. **Create/classify:** identify purpose, owner, legal/business basis, sensitivity and source.
2. **Validate:** perform schema, completeness, source, sequence and reconciliation checks.
3. **Protect/use:** authorize by role/purpose/case; monitor exceptional access.
4. **Retain:** apply event- or relationship-triggered approved schedule by country/record class.
5. **Hold:** suspend normal deletion only for a documented, scoped legal/regulatory/dispute need.
6. **Review:** periodically narrow or release holds; never use “possible future need” as an indefinite hold.
7. **Dispose:** delete/anonymize across operational, analytical, search, object, cache and vendor copies; expire backups under documented lifecycle.
8. **Prove:** retain minimal deletion/hold/approval evidence where justified without recreating deleted content.

Concrete periods remain unapproved until DEC-001, DEC-005 and DEC-051 are resolved. Product code must consume versioned retention rules rather than scattering hard-coded periods.

## Evidence export and chain of custody

For a dispute, incident or lawful disclosure:

- create a unique case and name an authorized owner;
- define allegation/request, lawful authority, exact scope and time window;
- preserve original records and record every collection/handling action;
- export the minimum relevant fields/media with source metadata;
- generate inventory, counts and cryptographic checksums;
- store working copies separately from originals;
- record redactions, transformations, timezone/calendar conversions and limitations;
- require appropriate review before external release;
- use a secure, expiring delivery method and verify recipient authority;
- record delivery, receipt, later correction and final disposition.

An export must identify whether an item is customer-reported, merchant-entered, system-derived or provider-verified.

## Request classification

Incoming requests must be classified before action:

| Request | Governing path |
|---|---|
| Data subject/authorized agent | Privacy-rights process and applicable exceptions |
| Merchant seeking its operational records | Contract/role/tenant authorization and export policy |
| Merchant seeking customer data beyond service purpose | Reject or require verified lawful basis/instruction; privacy review |
| Customer-versus-merchant dispute | Dispute process; Ubi does not invent adjudicatory authority |
| Police/administrative/government request | Lawful-request process below |
| Court order/subpoena/civil discovery | Counsel validates authority, jurisdiction, service and scope |
| Emergency threat request | Emergency process; documented good-faith necessity and later review |
| Security researcher report | Coordinated vulnerability intake; never send production data casually |

## Government and law-enforcement request process

Ubi will cooperate with valid law while protecting users against informal, excessive or misdirected disclosure.

### Intake and preservation

- Use a published controlled channel; support/chat/social-media requests are redirected and preserved, not fulfilled.
- Record requester identity, agency, official contact, authority, jurisdiction, case, signature, service method, data sought, period, deadline and nondisclosure demand.
- Authenticate through independently obtained official contact details where proportionate.
- A preservation demand preserves only the valid, available, specified scope; it is not authorization to disclose.

### Validation and decision

Licensed counsel or an approved legal decision-maker verifies:

1. the requesting body and document are authentic;
2. the request is validly served on the correct Ubi entity;
3. jurisdiction and legal power reach the requested data/entity;
4. the legal basis and procedural safeguards are sufficient;
5. scope, accounts, data categories and dates are specific and proportionate;
6. conflicts with other applicable law, localization, privilege or third-party rights are addressed;
7. user/merchant notice is allowed, required, delayed or prohibited;
8. challenge, narrowing or transfer to another entity/provider is appropriate.

Ubi should seek clarification or narrowing for vague, bulk, technically impossible, wrong-entity or disproportionate requests. No employee may create a permanent government “back door,” bypass encryption/access control, or disclose merely because a requester sounds authoritative.

### Collection, review and disclosure

- Query by verified stable identifiers and narrow dates; usernames/phone numbers alone may be ambiguous or reassigned.
- Use two-person review for sensitive/bulk/content/emergency disclosure when circumstances allow.
- Separate responsive from non-responsive/other-tenant data and redact accordingly.
- Disclose only data Ubi actually holds in intelligible form; do not reconstruct deleted data or promise provider-held data.
- Record legal approval, exact exported fields, counts, redactions, checksum, delivery and recipient.
- Preserve a defensible case record and delete working copies when their purpose ends.

### User notice and transparency

Unless prohibited or likely to create a legally recognized serious risk, Ubi should provide affected users/merchants meaningful notice with enough time to seek advice. Delayed notice restrictions are calendared and reviewed for expiry. Country-specific wording and exceptions require counsel.

Ubi should publish aggregate transparency reporting when volume and legal safety make it practical, including request categories, countries, accounts specified, outcomes and emergency requests. Counts may need bands or delayed publication to protect investigations and privacy.

## Emergency requests

An emergency assertion does not erase verification. The trained on-call decision-maker documents:

- credible imminent risk of death or serious physical harm;
- why the requested Ubi-held data is necessary and time-critical;
- requester authentication and emergency contact;
- minimum data disclosed and legal/good-faith basis;
- approval, delivery and post-event counsel review.

If time permits, require formal legal process. Emergency disclosure cannot be used for routine property, debt, commercial, immigration or reputational disputes without the applicable lawful authority.

## Cross-border and provider-held data

- A foreign request is not automatically valid against every Ubi entity or data location.
- Counsel determines the appropriate domestic, treaty, mutual-assistance, provider or court route.
- Ubi must not tell authorities or merchants that it controls connected-platform data it cannot lawfully/technically retrieve.
- Requests sent to a subprocessor are governed by contract notification/cooperation terms and the applicable law; the vendor register records the path.
- Conflicting legal obligations are escalated; staff do not improvise cross-border transfers.

## Merchant, customer and staff protections

- Tenant membership alone does not permit exporting another tenant or unrelated customer history.
- Business owners cannot silently read staff credentials/private data outside the approved model.
- Support personnel cannot use production search for curiosity, personal relationships or informal favors.
- Sensitive searches and bulk exports require reason, case and approval; anomalous access is reviewed.
- Parties receive correction mechanisms where an evidence package contains a demonstrable factual or attribution error; original history remains traceable.
- Retaliation against reporters, appellants or data-rights requesters is prohibited where applicable.

## Operational roles

| Role | Responsibility |
|---|---|
| Request intake owner | Secure intake, case completeness and deadline tracking |
| Licensed counsel/legal approver | Validity, jurisdiction, challenge, notice and disclosure authority |
| Privacy lead | Rights, minimization, conflicts, notice and transfer assessment |
| Security/evidence custodian | Preservation, scoped collection, integrity and secure delivery |
| Product/data owner | Source meaning, completeness, reconciliation and technical limitations |
| Executive approver | Exceptional high-risk/bulk/emergency decisions within law/policy |
| Independent reviewer/audit | Sample control testing, access review and corrective action |

The same person should not request, collect, approve and release a high-risk disclosure alone where staffing allows separation.

## Control testing and metrics

- percentage of required high-risk actions producing complete audit events;
- unauthorized cross-tenant access tests and privileged-access review findings;
- audit ingestion delay, loss/duplicate rate and integrity verification failures;
- open legal holds past review date and deletion jobs blocked by undocumented holds;
- requests received, authenticated, narrowed/rejected/disclosed and completed by deadline;
- disclosure packages with two-person review, inventory and checksum;
- user notices sent/delayed/prohibited and restrictions reviewed on time;
- rights/export evidence discrepancies and reconciliation time;
- working-copy deletion and recipient-delivery confirmation completion.

Metrics should expose control failure, not create incentives for automatic disclosure or denial.

## Release gates

- Record classes, authoritative sources, owners, classifications and versioned retention/hold rules are approved.
- SEC-004 events exist for first-release booking, permission, payment, integration, support and disclosure actions.
- Audit schema excludes secrets/content by default and tenant access is independently tested.
- Integrity, restricted access, clock/source handling, export/checksum and restore continuity are tested.
- Launch-country counsel approves request intake, service, preservation, emergency, notice, challenge and disclosure procedures.
- Only named trained roles can search/export/disclose; high-risk cases use separation/two-person review.
- Data-rights, merchant export, dispute, incident and lawful-request paths cannot bypass one another.
- Legal holds are scoped, reviewable and reversible; deletion propagates after release.
- Vendor contracts support preservation, incident notice, scoped export, request notification where lawful and deletion evidence.
- Public claims accurately distinguish what Ubi stores from data held only by merchants or connected providers.

## Open decisions

- Launch entity/country, service address, legal intake channel and counsel/on-call coverage.
- Approved retention schedule and legal-hold authorities by record family.
- Audit integrity/storage design and independent export cadence.
- Which read events justify audit at pilot scale.
- Government/user notice policy and transparency-report threshold by country.
- Emergency approval coverage outside business hours.
- Merchant/customer evidence access, dispute role and fees for exceptional exports where lawful.
- Call content/metadata and recording scope remain deferred with DEC-056.

