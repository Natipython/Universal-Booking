# Threat model and security controls

Status: **Draft — architecture baseline; independent review and implementation evidence required**  
Evidence checked: 2026-09-05 (Africa/Addis_Ababa)  
Depends on: actor/authority, workflows, platform, database, storage and recovery architecture  
Blocks: executable security requirements, production credentials/data, integrations and release approval

## Purpose

This document identifies who may attack or accidentally damage Ubi, what they can reach, the highest-impact abuse paths and the control/test evidence required before production. It covers Android, future desktop/customer clients, APIs, databases, objects, integrations, payments, messaging/calls, Ubi operations, CI/CD and vendors.

It is a living threat model, not a claim that a checklist makes Ubi secure.

## Reference baseline

Use these as verification inputs, adapted to Ubi rather than copied mechanically:

- [OWASP ASVS 5.0](https://owasp.org/www-project-application-security-verification-standard/) for application-control verification;
- [OWASP API Security Top 10: 2023](https://owasp.org/API-Security/editions/2023/en/0x03-introduction/) for object/function/property authorization, resource consumption, sensitive flows, SSRF, inventory and unsafe provider consumption;
- [OWASP MASVS](https://mas.owasp.org/MASVS/) for Android storage, cryptography, authentication, network, platform, code, resilience and privacy testing;
- [NIST SP 800-63-4](https://www.nist.gov/publications/nist-sp-800-63-4-digital-identity-guidelines) as risk-based identity/authentication guidance, not as an automatic compliance label;
- [OAuth 2.0 Security BCP, RFC 9700](https://www.rfc-editor.org/info/rfc9700/) for exact redirect matching, authorization code + PKCE, CSRF/mix-up controls, refresh-token rotation or sender constraint, audience/scope restriction and avoidance of deprecated insecure grants.

The pinned versions and applicability are reviewed at each major release.

## Security objectives

1. No account, merchant, employee, integration or support user can exceed explicit authority.
2. One organization cannot see or influence another organization’s data or capacity.
3. A retry, replay or race cannot create duplicate bookings, charges, refunds or messages.
4. “Reported,” “submitted,” “delivered,” “verified,” “settled” and “completed” remain distinct.
5. Account recovery cannot become silent business-ownership transfer.
6. Compromise of one device, tenant or provider credential has a bounded blast radius.
7. Security controls remain usable on low-cost devices and weak networks; users are not pushed into shared passwords or unsafe workarounds.
8. Powerful actions and sensitive access are attributable, reviewable and recoverable where possible.
9. Ubi fails closed for authority and money, and degrades honestly for optional providers.
10. Security/privacy claims match tested reality.

## Protected assets

### Critical owner truth

- organization ownership, membership, roles, approvals and restrictions;
- account authenticators, recovery evidence, sessions and devices;
- booking state, immutable snapshots, holds, reservations and capacity;
- payment obligations, provider evidence, allocations, refunds and reconciliation;
- privacy cases, legal holds, deletion/revocation tombstones and audit integrity;
- signing/encryption keys, service/provider credentials and release-signing identity.

### Confidential operational data

- customer contacts, addresses and booking details;
- staff calendars/assignments and business configuration/pricing;
- conversations, attachments, call metadata and abuse reports;
- provider identifiers/events and support case information;
- backups, exports, logs/traces and local device projections.

### Availability/reputation assets

- booking and manage-link availability;
- correct public business profiles/subdomains;
- notification/message/call reputation and provider accounts;
- app/update/build pipeline and domain/DNS;
- merchant trust, evidence and incident response capability.

## Threat actors and failure sources

- opportunistic internet attacker, bot or spammer;
- fraudulent customer, merchant owner or staff member;
- disgruntled/former employee or contractor;
- attacker with a lost/shared/rooted phone or stolen session;
- malicious or compromised Ubi workforce account;
- compromised merchant/provider/CI/cloud credential;
- vulnerable/malicious dependency, SDK, build action or app update;
- forged/replayed provider webhook or malicious provider payload;
- provider insider/outage/account takeover/policy change;
- accidental staff/support/developer mistake;
- automation/model prompt injection or unsafe autonomous decision;
- extortion/harassment/report abuse in messaging/calling;
- infrastructure failure, clock/network inconsistency or destructive migration.

## Trust boundaries

```text
untrusted Android/web/desktop clients
→ public edge/API boundary
→ authenticated application modules
→ tenant-bound PostgreSQL/object/search/cache
→ background queue/workers
→ external messaging/payment/call/notification providers
→ Ubi operations/JIT access
→ engineering CI/CD and cloud control planes
→ backup/recovery boundary
```

Every crossing authenticates origin where possible, validates schema/size/semantics, binds tenant and purpose, limits resources and records significant outcomes. TLS protects transport; it does not make a payload or caller trustworthy.

Diagram: `diagrams/50-security-trust-boundaries.mmd`.

## Risk method

Track threat ID, asset, actor, precondition, attack path, impact, likelihood, existing/proposed controls, detection, response, test evidence, residual risk, owner and review trigger.

Severity considers harm to customers/merchants, financial loss, cross-tenant scope, irreversibility, legal/reporting impact, exploitability, detectability and recovery—not only technical CVSS.

Critical/high risks cannot be accepted silently by engineering. Named business/security/privacy owners approve residual risk and release scope.

## Priority threat register

| ID | Threat and consequence | Required prevention/detection direction |
|---|---|---|
| T-001 | BOLA/IDOR exposes another tenant’s booking/customer/message | Server authorization per object/action; immutable tenant context; composite tenant FKs; forced RLS; negative/property tests; anomaly alerts |
| T-002 | Role/scope bypass grants refund, export, ownership or integration control | Capability + scope + state; step-up/approval; no owner superuser shortcut; field authorization; audit |
| T-003 | Recovery/social engineering transfers business ownership | Separate account recovery from ownership; verified owners; cool-down/alerts; multi-party evidence for high risk; redress |
| T-004 | Stolen phone/session exposes business/customer data | Short access life, refresh rotation/reuse detection, device/session list/revoke, local encryption/lock, risk-based reauth, bounded offline cache |
| T-005 | Shared staff password destroys accountability | Membership per person; invitations; no shared credentials; assisted setup; fast staff revoke; audit |
| T-006 | Offline/retry race double-books capacity | Idempotency request hash/result replay; database exclusion/locked capacity ledger; expected versions; concurrency tests |
| T-007 | Duplicate/forged payment event marks paid or duplicates refund | Webhook authenticity + replay/dedup; source namespacing; reported vs verified; unknown-state reconciliation; limits/approval |
| T-008 | Provider callback routes data to wrong merchant | Bind app/environment/connection/business asset; verify signature/token/timestamp; fetch/validate referenced owner; quarantine unknown |
| T-009 | OAuth interception/mix-up/token theft compromises merchant inbox | Authorization code + PKCE; exact redirect; state/nonce/issuer; server-side token exchange; minimum scopes; encrypted tokens; rotation/revocation |
| T-010 | Merchant gives Ubi social password or unofficial automation gets account banned | Never request/store password/session cookie; official API only; truthful capability eligibility; direct-link/manual fallback |
| T-011 | SSRF/media fetch reaches metadata/private network | No arbitrary server fetch; URL allowlist/resolution/re-resolution, DNS/IP/private-range blocks, size/time/redirect limits, isolated fetcher |
| T-012 | Malicious file/parser causes malware, decompression bomb or stored XSS | Direct quarantined upload, actual-type/checksum/size/dimension limits, isolated scanner/transformer, safe download headers, no active content by default |
| T-013 | API resource exhaustion blocks small merchants | Per principal/tenant/IP/device/action budgets; bounded pages/uploads/queries; queues/backpressure; cost alerts; fair-use protection |
| T-014 | Enumeration reveals phone/account/business/customer existence | Consistent responses/timing where practical; opaque IDs; rate/risk controls; no global contact/customer search |
| T-015 | Scoped booking link leaks or is replayed | High-entropy verifier, purpose/resource/actions/expiry/use limit, exchange for short session, referrer/log protection, revoke/supersede/step-up |
| T-016 | Cache/search/object key leaks cross-tenant data | Trusted tenant key/ACL on every projection; server filter; reauthorize hydration; private objects + short URLs; poison tests |
| T-017 | Ubi support browses or changes customer data invisibly | Case + purpose + scoped expiring JIT access; mask; approval/step-up; access audit and review; break-glass alerts |
| T-018 | Logs/analytics/crash SDK collect messages/tokens/contacts | Allowlisted structured telemetry; redaction; payload tests; sampling/short retention; vendor controls; secret/PII scanning |
| T-019 | CI/dependency/build compromise ships malicious app/server | Protected branches/reviews, pinned lockfiles/actions/images, SBOM/provenance/signing, isolated secrets, SAST/SCA/secret scan, reproducible promotion |
| T-020 | DNS/domain/subdomain takeover impersonates Ubi/business | Registrar MFA/lock, least-privilege DNS, automated certificates, reserved/tombstoned slugs, dangling-record checks, change audit |
| T-021 | Messaging/call feature enables harassment/scam/spam | Booking/context boundaries, block/report, rate/age/identity controls, staff takeover, safe notification privacy, moderation/redress |
| T-022 | Automation prompt injection or hallucination confirms money/time | Treat message/content as untrusted data; tool allowlist/schema; human confirmation for high risk; no secret exposure; provenance/evaluation/kill switch |
| T-023 | Destructive migration/admin/ransomware corrupts owner data/backups | JIT migration, lock-safe review, separate immutable backup control, PITR/restore drill, dual control, anomaly alerts |
| T-024 | Clock/timezone manipulation changes expiry/schedule | Trusted server time; stored timezone/calendar/version; bounded client timestamps; monotonic sequences; DST/calendar tests |
| T-025 | Insider/provider exports data beyond purpose | Least privilege, separation, encryption/key boundary, export approval/rate/manifest, DLP/anomaly, contract/audit/deletion |
| T-026 | Old Android client bypasses new rule | Server-enforced rules; contract/version compatibility; minimum safe version/kill switch; migration gates |

## Authentication architecture

### Account levels

- Browsing public information needs no Account.
- Basic customer/merchant account proves control of an authenticator/contact, not legal identity.
- Verified merchant ownership/representative authority is a separate evidence lifecycle.
- High-risk operation requires current assurance/step-up appropriate to impact.

Do not collect identity documents for everyone merely to appear secure.

### Authenticators

Preferred direction:

- passkeys/platform credentials where supported and comprehensible;
- phone OTP can bootstrap/recover in selected markets but is not treated as phishing-resistant and must consider SIM swap/recycled/shared numbers;
- passwords, if offered, use modern slow memory-hard hashing, breached/common-password screening, generous length/paste/password-manager support and no arbitrary periodic rotation;
- recovery codes or additional verified factor for owners where feasible;
- no security questions.

Exact assurance policy follows risk and local device/telecom research. A control people cannot use will create account sharing or lockout.

### OTP and recovery protections

- single purpose, short expiry, attempt/rate/resend bounds and hashed verifier;
- invalidate superseded codes; never log or reveal them to staff;
- generic enumeration-safe responses;
- risk alerts for number/device/recovery changes;
- high-risk ownership/payment/integration changes wait or use stronger independent evidence;
- supported human redress without bypassing controls through persuasion.

### Sessions

- short-lived audience-bound access token; rotating refresh-token family with reuse detection;
- server records device/session/family, assurance, issued/last-use, revocation and security events;
- revoke one/all sessions and notify on risky changes;
- cookies, if used, are Secure, HttpOnly, appropriate SameSite, scoped and CSRF-protected;
- tokens never in URLs, analytics, logs or browser local storage where safer server/session patterns apply;
- logout clears local protected data according to device/offline policy;
- risky actions re-check current authority/version and step-up, not just token age.

## Authorization controls

The complete model is in document 44. Security tests cover object, property and function authorization separately.

- No mass-assignment from client DTO to owner row.
- Response DTOs expose approved fields, not database objects.
- Batch/export/search/sync authorize every scope and prevent count/timing leaks.
- Membership role change, organization restriction and ownership transfer invalidate authority epochs.
- Service-to-service identity is authenticated, audience-restricted and cannot forge a user context without explicit delegation evidence.
- Database RLS is defense in depth, not the only authority layer.

## API and web controls

- TLS only; modern provider configuration; no certificate-validation bypass in released clients.
- Exact request schemas, content types, methods, size/depth/count bounds and normalized identifiers.
- Consistent safe errors; detailed internal cause linked by trace ID without secrets.
- State-changing browser requests use CSRF defenses; CORS is exact allowlist, never credentialed wildcard.
- Security headers include CSP/frame-ancestors, no-sniff, referrer policy and HSTS after deployment readiness.
- Redirect/deep-link allowlists prevent open redirect and app-link hijacking.
- Idempotency on retryable high-impact commands; optimistic/domain versions for stale writes.
- Rate limits are action/principal/tenant/risk aware, not one IP-only rule that punishes shared mobile networks.
- Query complexity, pagination and exports are bounded.
- Public docs/inventory expose only intended versions; old/debug/admin endpoints are removed or equally protected.

## Android controls

- Server remains authority; no privileged static API secret in APK.
- App-private storage and encrypted approved offline data; keys protected by Android Keystore where supported.
- Minimum permissions, requested at point of understandable use; no contact/SMS/call-log/location access unless separately approved.
- Verified Android App Links for Ubi domains; untrusted incoming intents/URIs strictly parsed and reauthorized.
- Disable cleartext traffic; pinning is a risk/operability decision with safe rotation if used.
- Exported components, WebViews, file providers, notifications, clipboard, screenshots/recents and backup flags reviewed.
- Release signing keys protected; Play/app integrity is a risk signal, not authority or a reason to block legitimate low-cost devices automatically.
- Root/debug/hook resistance can raise attacker cost but cannot make a compromised device trusted.
- Local purge/logout/revocation and offline pending-command behavior tested.

Android’s official guidance supports permission minimization and Keystore-protected, non-exportable key use: [Design for Safety](https://developer.android.com/quality/privacy-and-security) and [Android Keystore](https://developer.android.com/privacy-and-security/keystore).

## Integration/webhook controls

- Register exact HTTPS callback; authenticate provider using its current documented mechanism.
- Preserve raw bytes needed for signature verification before parsing/transformation.
- Verify algorithm/key/timestamp/replay window/event ID/environment/application/connection.
- Acknowledge within provider budget after durable receipt; process asynchronously/idempotently.
- Unknown tenant/asset/schema/state quarantines—never guesses success/routing.
- Deduplicate but preserve legitimate provider event revisions.
- Outbound calls use allowlisted provider hosts, TLS, strict timeouts, bounded retry with jitter and circuit/bulkhead controls.
- Provider response is untrusted: validate type/size/semantic mapping and prevent SSRF/injection.
- Secret rotation supports overlap and emergency revoke; disconnect revokes provider authorization and stops processing.

## Payment-specific controls

- Prefer licensed provider-hosted collection; Ubi avoids card/mobile-money credential custody and unapproved fund custody.
- Exact currency/amount/merchant/booking/purpose binds intent and callback.
- Customer screenshot/reference stays reported/unverified until authorized manual or authenticated provider verification.
- Unknown provider outcome is reconciled before retrying charge/refund.
- Refund/override/export/connection changes require capability, limits, step-up and approval by risk.
- Allocation/refund/reversal constraints and reconciliation protect accounting truth.
- Webhook/signature success alone is not settlement.
- Payment/security standards and regulatory scope are determined with the selected provider, payment model and launch-country counsel.

## Messaging, automation and call controls

- In-app messaging/calls remain future-gated but are threat-modeled now.
- User/business identity and automation status are visible; no deceptive human impersonation.
- Customer must initiate where provider rules require; templates/windows/opt-out and prohibited content follow current platform rules.
- Block/report/mute/takeover/escalation work across devices and prevent automation from continuing.
- Calls require consent/expectation, ring/spam limits, caller identity, safety reporting and network failure behavior.
- Recording is off by default and needs separate lawful purpose, notice/consent, storage/access/retention and jurisdiction approval.
- AI receives the minimum conversation context, never secrets; high-risk tools are schema/permission/state constrained and human-confirmed.

## Cryptography and secret management

- Use approved maintained libraries; no home-grown cryptographic protocol.
- TLS in transit and managed encryption at rest are baseline.
- KMS/secret manager protects service keys; envelope encryption may protect high-cardinality provider tokens.
- Keys are environment/purpose/version separated with owner, rotation, revoke and recovery procedure.
- Passwords use approved one-way hashing; high-entropy tokens stored as verifier/hash where replay is unnecessary.
- Authenticated encryption binds ciphertext to immutable environment/provider/tenant/credential context.
- Production secrets never enter repository, mobile binary, ticket/chat, log, analytics or backup outside the key design.
- Crypto agility retains algorithm/key version and safe rewrap/migration path.

## Operations, audit and insider controls

- Production human access is no standing privilege where practical: ticket/case, purpose, step-up, scoped JIT grant, expiry and audit.
- Support sees masked/minimum fields and cannot query raw database/object/log/search arbitrarily.
- Break-glass immediately alerts independent reviewers and triggers post-use review.
- High-risk changes use two-person approval where pilot staffing can sustain it; otherwise record compensating control/risk.
- Audit is append-controlled, access-controlled and integrity-checked; do not copy full sensitive content into it.
- Security alert review, access certification and dormant account/key removal have named cadence/owners.

## Supply-chain and CI/CD controls

- Protected review/merge; signed commits/tags where useful; least-privilege maintainers.
- Lock dependencies and CI actions/images by immutable version/digest; automated update PRs are reviewed/tested.
- SAST, dependency/vulnerability, license, secret and infrastructure scans; SBOM for releases.
- Build and deployment use short-lived federated identity, isolated protected environments and approval for production.
- Artefact provenance/signature and promotion of the same tested artefact; no rebuilding differently in production.
- Android signing/recovery, backend rollback/forward-fix and compromised-release response rehearsed.
- Vulnerability intake, triage, patch SLA and coordinated disclosure contact exist before public launch.

## Monitoring and response signals

Monitor without collecting unnecessary content:

- authentication/recovery/OTP failures and token reuse;
- membership/ownership/integration/payment/exports/JIT access changes;
- cross-tenant authorization denials and anomalous object patterns;
- idempotency conflicts, capacity collisions and provider duplicate/unknown outcomes;
- rate-limit/resource/cost anomalies and bulk scraping;
- webhook signature/schema/replay failures;
- secret/DNS/bucket/RLS/backup/CI policy changes;
- message/call block/report/spam patterns;
- audit/telemetry pipeline gaps and clock drift.

Signals create proportionate cases; they do not automatically punish customers based on one weak device/network/IP signal.

## Verification programme

Before production:

- map applicable ASVS/MASVS/API controls to tests/evidence and justified exclusions;
- abuse-case and data-flow threat modeling for every first-release slice;
- unit/property/concurrency/contract tests for authorization, idempotency, money and capacity;
- automated cross-tenant endpoint/database/search/cache/object/sync tests;
- Android static/dynamic/network/storage/deep-link/WebView/backup review on supported low-cost devices;
- API/web/integration/payment penetration test by a qualified independent tester;
- dependency/CI/cloud/IAM/infrastructure review;
- restore, credential leak, account recovery, support misuse and provider replay exercises;
- remediation verified, not merely ticketed; residual critical risk blocks launch.

Production uses continuous monitoring, vulnerability review and periodic penetration/threat-model refresh after major features/providers/incidents.

## Release gates

- Every route/command/event/provider callback has authentication, capability/scope/state, schema, idempotency, rate and audit rules.
- Account recovery and ownership transfer pass adversarial tests and have redress.
- Runtime roles cannot DDL/bypass RLS/read other tenants; pooled context leakage tests pass.
- No channel integration needs a merchant password, copied browser session or unofficial client automation.
- Payment release uses a licensed provider design and passes unknown/duplicate/reversal/reconciliation tests.
- All secrets/keys have owner/rotation/revoke/recovery; secret scanning is enforced.
- Support production access is case/JIT/audited before support receives access.
- Android and API meet the approved ASVS/MASVS-derived baseline on supported versions/devices.
- Incident runbooks, contacts and tested backups exist.
- Independent review finds no unresolved release-blocking issue.

## Open decisions

- Exact authentication methods, assurance levels and owner recovery evidence for launch country.
- Security testing provider/budget and vulnerability disclosure/bug bounty timing.
- Chosen cloud/identity/payment/message/call providers and their concrete controls.
- First-release Android OS/device support and safe minimum version policy.
- Whether certificate pinning/device attestation improves risk enough for its accessibility/operational cost.
- SOC/monitoring staffing and 24/7 escalation at pilot scale.
- Call recording, AI processing and high-risk automation remain separately gated.
