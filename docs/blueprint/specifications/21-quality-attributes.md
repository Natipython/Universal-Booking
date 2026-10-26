# Quality attributes and measurable service requirements

Status: **Draft — proposed targets, not a public SLA**  
Depends on: product principles and functional requirements  
Blocks: architecture, client design, infrastructure, testing, capacity planning and release approval

## Purpose

This document turns words such as *fast*, *smooth*, *reliable*, *secure* and *offline-capable* into testable engineering requirements. The figures are deliberately demanding but must be validated on real launch-market devices, networks and workflows before approval.

Unless stated otherwise:

- client measurements use production-like release builds, not debug builds;
- service measurements exclude time spent inside a clearly identified third-party provider, but Ubi's own adapter time is included;
- percentile targets are measured over a representative rolling period and segmented by device, network, operation and release;
- the user must receive an immediate truthful state even when authoritative work continues;
- a target does not authorize hiding a failure, dropping work, weakening validation or displaying a false success.

## Reference test conditions

The final matrix must be derived from field research and analytics that users knowingly permit. These profiles are the provisional baseline.

### Devices

| Profile | Provisional specification | Main purpose |
|---|---|---|
| D1 — minimum supported business Android | Android 11 or later, at least 4 GB RAM (founder DEC-022); provisional 32 GB storage, 720p-class screen, low-cost four-to-eight-core device | Release gate at the Android 11/4 GB floor, plus newer supported versions |
| D2 — common Android | Android 12 or later, 4 GB RAM, 64 GB storage, 1080p-class screen | Normal product and performance testing |
| D3 — larger business device | Android tablet or large-screen device with 4 GB RAM | Responsive/adaptive layout validation after tablet scope is approved |
| W1 — constrained customer web | Low-cost phone with a supported browser and an uncached first visit; exact OS/RAM matrix remains open independently of the business-app floor | Public booking/link release gate; do not automatically exclude guests on older phones |
| W2 — business desktop | Supported current browser on a dual-core/4 GB-class Windows device | Desktop business-workspace baseline |

Exact locally available D1 and D2 models must be purchased and recorded before the first external pilot. Emulator results cannot replace physical-device results.

Founder DEC-022 (2026-10-04) replaces the earlier provisional Android 10/2 GB business baseline. Android 11 maps to API 30; the prototype sets `minSdk = 30`, following [Android's SDK-version documentation](https://developer.android.com/guide/topics/manifest/uses-sdk-element). This is distinct from compile/target SDK selection and store-release policy. The 4 GB figure describes installed physical RAM, not current free memory; do not reject a user because other apps reduce available RAM. Actual performance, memory pressure and support coverage still require physical-device evidence. No speed guarantee follows from this minimum.

### Networks

| Profile | Provisional conditions | Expected behaviour |
|---|---|---|
| N0 — offline | No usable route | Cached schedule and eligible queued work remain available |
| N1 — severely constrained | 400 kbit/s down, 150 kbit/s up, 400 ms round-trip latency, 2% packet loss | Critical text workflows remain usable; media defers |
| N2 — basic mobile | 1.5 Mbit/s down, 500 kbit/s up, 200 ms round-trip latency, 1% packet loss | Complete critical workflows without confusing stalls |
| N3 — typical good connection | 10 Mbit/s down, 2 Mbit/s up, 80 ms round-trip latency | Full supported experience |

The matrix is a reproducible laboratory baseline, not a claim that it represents every East African network. Field measurements in the selected city must confirm or adjust it.

## Performance and responsiveness

### Android application

| ID | Proposed requirement | Verification |
|---|---|---|
| NFR-PERF-001 | On D1 with cached data, cold-start time to initial display (TTID) MUST be at or below 2.0 seconds at p75. | Instrumented physical-device test |
| NFR-PERF-002 | On D1 with cached data, cold-start time to fully usable Today screen (TTFD) MUST be at or below 3.0 seconds at p75; warm TTFD MUST be at or below 1.0 second. | Instrumented physical-device test |
| NFR-PERF-003 | A tap on an enabled local control MUST show visual acknowledgement within 100 ms, including when remote work is pending. | UI instrumentation/high-speed recording |
| NFR-PERF-004 | A cached primary screen SHOULD present useful content within 300 ms of navigation on D1; it MUST NOT wait for a network request before showing valid cached content. | N0/N1 navigation test |
| NFR-PERF-005 | Common scrolling and transitions MUST have no frozen frame longer than 700 ms; the missed-frame ratio SHOULD remain below 5% and p99 frame time below 33.3 ms on D1. | Android frame metrics in release build |
| NFR-PERF-006 | A normal incremental start synchronization SHOULD transfer no more than 200 KiB when there are no new attachments; a typical booking command SHOULD be no more than 25 KiB excluding media. | Network capture and contract test |
| NFR-PERF-007 | Delivered Android binary SHOULD be at or below 35 MB and MUST trigger architecture review above 50 MB. Installed base data SHOULD remain below 100 MB before optional cache/media. | Release artifact and clean-install measurement |
| NFR-PERF-008 | Local database, media and logs MUST use explicit quotas and safe eviction. Optional media MUST load on demand and MUST NOT prevent a text workflow. | Storage pressure and eviction test |
| NFR-PERF-009 | Background work MUST respect operating-system constraints, battery and metered-network policy; continuous polling is prohibited. | Battery profiler and background-policy test |

These figures are engineering budgets. They are not permission to remove necessary accessibility, integrity or security controls.

### Customer web and business desktop

| ID | Proposed requirement | Verification |
|---|---|---|
| NFR-WEB-001 | At the 75th percentile of real visits, customer-facing pages MUST target Largest Contentful Paint at or below 2.5 seconds, Interaction to Next Paint at or below 200 ms and Cumulative Layout Shift at or below 0.1. | Field data plus controlled W1 tests |
| NFR-WEB-002 | The critical customer booking path MUST remain usable under N1 without downloading promotional media. | Throttled end-to-end test |
| NFR-WEB-003 | The first critical customer route SHOULD transfer no more than 250 KiB compressed, excluding optional media, and SHOULD require no more than 150 KiB compressed initial JavaScript. | Bundle and network budget gate |
| NFR-WEB-004 | Server-rendered content MUST communicate the business, service, price method and next action before client enhancement completes. | JavaScript-delayed/degraded test |
| NFR-WEB-005 | Responsive images MUST declare dimensions, use appropriate modern formats where supported and MUST NOT block selection, review or confirmation. | Automated and visual performance test |
| NFR-WEB-006 | Desktop cached navigation SHOULD acknowledge input within 100 ms and present useful primary content within 500 ms on W2. | Browser performance test |

### Service and API latency

| ID | Proposed requirement |
|---|---|
| NFR-API-001 | Under normal regional conditions and forecast load, first-party p95 read latency SHOULD be at or below 500 ms and p99 at or below 1.5 seconds. |
| NFR-API-002 | First-party p95 command acceptance SHOULD be at or below 800 ms and p99 at or below 2 seconds when the command does not synchronously depend on a provider. |
| NFR-API-003 | Availability queries SHOULD complete at or below 1 second p95 for the launch archetype and bounded search window. |
| NFR-API-004 | Provider-dependent timing MUST be measured separately by provider and operation; the client MUST show processing/pending/failure rather than extending an indefinite spinner. |
| NFR-API-005 | Every outbound request MUST have a bounded timeout, retry classification and correlation identifier. Interactive retries MUST not multiply a logical command. |

Latency objectives start at the Ubi regional edge/API and end when Ubi returns a response. Client-perceived latency is measured separately end to end.

## Reliability, integrity and availability

| ID | Proposed requirement |
|---|---|
| NFR-REL-001 | The authoritative booking path SHOULD achieve 99.9% monthly availability after general launch; the pilot measures this objective but does not promise it contractually. |
| NFR-REL-002 | Authenticated webhook ingestion SHOULD achieve 99.95% monthly availability because provider retries may be time-limited. |
| NFR-REL-003 | Non-critical operations/reporting tools SHOULD achieve 99.5% monthly availability and MUST NOT share a failure mode that disables booking confirmation. |
| NFR-REL-004 | Every logical mutation MUST be idempotent across client retry, queue retry, worker retry and provider webhook replay. |
| NFR-REL-005 | Ubi-induced double booking MUST be zero in the approved concurrency, retry and partition test suite. Capacity constraints MUST be enforced by the authoritative transaction boundary. |
| NFR-REL-006 | Money MUST use exact currency-aware amounts; floating-point arithmetic MUST NOT determine charges, balances or refunds. |
| NFR-REL-007 | A notification, analytics or search-index failure MUST NOT roll back or falsify a committed booking. Recoverable downstream work MUST be durably recorded. |
| NFR-REL-008 | Every externally visible command MUST end as committed, rejected, expired, cancelled or action-required; an unbounded unknown state is a defect. |
| NFR-REL-009 | Crash-free Android sessions SHOULD be at least 99.8% and user-perceived ANR rate SHOULD remain below 0.1% in the pilot cohort. No known crash/ANR in a critical flow is releaseable. |
| NFR-REL-010 | Planned maintenance exclusions, if any, MUST be defined before a public SLA; silent reclassification of incidents as maintenance is prohibited. |

Availability is not correctness. A fast response that exposes another tenant, loses a command or confirms unavailable capacity is a failed request.

## Offline and synchronization quality

| ID | Proposed requirement |
|---|---|
| NFR-OFF-001 | On N0, an authenticated business user with valid local access MUST be able to open the cached Today schedule in no more than 1 second after the app is usable. |
| NFR-OFF-002 | An eligible offline action MUST be durably stored before the app labels it “Saved on this phone”; killing the process immediately afterward MUST not lose it. |
| NFR-OFF-003 | Reconnection MUST trigger bounded resumable synchronization without requiring the user to reopen each pending item. |
| NFR-OFF-004 | Duplicate transmission MUST produce one logical server effect and converge to the same local identifier/result. |
| NFR-OFF-005 | Permission, capacity, price and customer-visible term conflicts MUST stop automatic commitment and enter an understandable resolution flow. |
| NFR-OFF-006 | Device clock drift MUST not establish authorization, hold expiry, payment truth or final event order. Server authority and stable causal/version information are required. |
| NFR-OFF-007 | A user MUST always be able to see whether an operational item is local, pending, synchronized, rejected or conflicted. Connectivity color/icon alone is insufficient. |

## Recovery, continuity and durability

| ID | Proposed requirement |
|---|---|
| NFR-DR-001 | Authoritative operational data MUST target a recovery point objective of 5 minutes or less. Rebuildable analytics may target up to 24 hours only when its source remains durable. |
| NFR-DR-002 | Booking core MUST target a four-hour recovery time and supporting administration an eight-hour recovery time for the initial single-region design. |
| NFR-DR-003 | Backups MUST be encrypted, access-controlled, monitored for successful completion and protected from routine application credentials. |
| NFR-DR-004 | A backup is not accepted until a restore test proves schema, constraints, tenant isolation, application usability and reconciliation checkpoints. |
| NFR-DR-005 | Restore exercises MUST occur before production launch and at least quarterly thereafter; findings need owners and due dates. |
| NFR-DR-006 | The architecture MUST document degraded operation and reconciliation for provider, notification, object-storage, search and analytics outages. |
| NFR-DR-007 | Regional disaster recovery beyond the initial launch region requires a separately costed and tested design; marketing MUST NOT imply multi-region resilience before it exists. |

RPO and RTO are proposed engineering objectives until hosting topology, budget and business-impact analysis are approved.

## Security and abuse resistance

| ID | Proposed requirement |
|---|---|
| NFR-SEC-001 | Internet-facing application controls SHOULD be designed and assessed against OWASP ASVS 5.0 Level 2, with additional risk controls for privileged administration, identity, payments and integrations. No compliance claim may be made without assessment evidence. |
| NFR-SEC-002 | Threat modelling MUST cover account takeover, tenant escape, broken object authorization, mass assignment, replay, double spend/booking, token leakage, webhook forgery, insider misuse, abuse and denial of service. |
| NFR-SEC-003 | Every tenant-scoped endpoint and asynchronous handler MUST have positive and negative authorization tests, including changed/revoked membership. |
| NFR-SEC-004 | Supported network traffic MUST use modern transport encryption; TLS 1.2 is the minimum compatibility floor and TLS 1.3 is preferred where available. Sensitive stored data and backups MUST use provider-supported encryption with controlled key access. |
| NFR-SEC-005 | Passwords, OTPs, access/refresh tokens, provider tokens, payment credentials and raw secrets MUST never enter analytics or ordinary logs. |
| NFR-SEC-006 | Release builds MUST pass secret scanning, dependency/SBOM generation, static analysis and risk-based dynamic tests. Critical/high exploitable findings have zero release tolerance unless formally risk-accepted by accountable security leadership. |
| NFR-SEC-007 | Public and authentication endpoints MUST implement rate limits, enumeration resistance, bot/abuse controls and observable throttling without excluding legitimate constrained-network users. |
| NFR-SEC-008 | Security-sensitive events MUST create tamper-evident audit records and alerts appropriate to risk. Audit readers and retention are separately controlled. |
| NFR-SEC-009 | Provider credentials MUST be isolated by environment and tenant/connection, encrypted, scope-minimized, rotatable and revocable. |
| NFR-SEC-010 | Production security testing MUST include a pre-launch independent review and a defined vulnerability intake, triage and disclosure process. |

The OWASP Top 10 and API Security Top 10 are threat-awareness inputs, not complete verification standards.

## Privacy and data protection quality

| ID | Proposed requirement |
|---|---|
| NFR-PRV-001 | Every personal-data field and analytical event MUST have a recorded purpose, classification, owner, access scope, retention rule and deletion behaviour before production collection. |
| NFR-PRV-002 | Production observability MUST default to metadata and identifiers; viewing message/customer content requires an approved operational purpose and access path. |
| NFR-PRV-003 | Data-subject and merchant requests MUST use identity verification proportionate to the disclosure or action risk and retain case/audit evidence. |
| NFR-PRV-004 | Deletion and restriction jobs MUST be idempotent, report partial failures and verify propagation across operational, object, search, cache and analytical stores. |
| NFR-PRV-005 | Backup expiry, legal hold and records that Ubi must retain MUST be explained truthfully; a user-facing completion state cannot mean more than the system has performed. |
| NFR-PRV-006 | No country-specific response/deletion period is approved until the jurisdiction matrix and licensed counsel review establish it. Internal targets must be no slower than applicable law. |

## Accessibility and usability quality

Detailed requirements are controlled in `22-accessibility-localization-and-calendar.md`. The system-level quality gates are:

| ID | Proposed requirement |
|---|---|
| NFR-UX-001 | A first-time target user MUST be able to complete the pilot’s primary booking tasks from plain-language prompts without prior software training. Training may improve speed but cannot be required for correctness. |
| NFR-UX-002 | Critical flows MUST pass moderated usability testing with locally recruited owners/staff representing limited digital experience; internal team testing is insufficient. |
| NFR-UX-003 | Web experiences target WCAG 2.2 Level AA; native experiences MUST provide equivalent perceivable, operable, understandable and robust behaviour. |
| NFR-UX-004 | Destructive, financial and customer-visible actions MUST show scope and consequences, prevent accidental double action and provide a receipt/recovery path. |
| NFR-UX-005 | Error messages MUST state what happened, what is safely retained and what the user can do next without exposing security-sensitive detail. |
| NFR-UX-006 | No required meaning or state may depend on color, motion, sound, technical code or an icon without a text/semantic equivalent. |

Pilot usability thresholds must be set after task protocols exist. A proposed starting gate is at least 90% unassisted completion for each critical task, zero unrecoverable data-loss error and median System Usability Scale at or above 80, reported by user subgroup rather than only as a blended average.

## Scalability and capacity

The first architecture must scale economically without pretending to solve global scale.

| ID | Proposed requirement |
|---|---|
| NFR-CAP-001 | Before each release, finance/product MUST provide a 12-month forecast for organizations, operations, staff, customers, bookings, messages, media, notifications, webhook bursts and peak concurrency. |
| NFR-CAP-002 | Load tests MUST sustain the forecast peak plus a minimum 3× headroom and a short 10× burst for idempotent ingress without tenant leakage or lost accepted work. Exact duration depends on workload. |
| NFR-CAP-003 | One unusually active tenant MUST have quotas/backpressure that protect other tenants while preserving or clearly rejecting accepted work. |
| NFR-CAP-004 | Capacity dimensions and bottlenecks MUST be observable; scaling decisions cannot rely only on aggregate CPU. |
| NFR-CAP-005 | Partitioning, sharding and microservice extraction require measured limits and an approved migration path; they are not launch prerequisites by default. |

Illustrative pilot planning numbers—such as 100 organizations or 1,000 bookings per day—are scenarios, not hard-coded architectural limits and not approved forecasts.

## Maintainability, change safety and interoperability

| ID | Proposed requirement |
|---|---|
| NFR-MNT-001 | The initial server SHOULD be a modular monolith with enforceable domain boundaries; deployment separation follows demonstrated scaling, security or ownership needs. |
| NFR-MNT-002 | Public/client/provider contracts MUST be schema-defined, backward-compatible for a stated support window and tested with consumers. |
| NFR-MNT-003 | Booking, availability, money, authorization and synchronization invariants MUST have deterministic unit/property/concurrency tests; a percentage coverage number cannot substitute for invariant coverage. |
| NFR-MNT-004 | Database and event changes MUST be forward/backward compatible across the deployment window and include tested rollback or roll-forward procedures. |
| NFR-MNT-005 | A pull-request fast gate SHOULD complete in 15 minutes or less; slower suites may run separately but all release gates must finish before promotion. |
| NFR-MNT-006 | Feature flags MUST have owner, scope, expiry/removal condition and safe default. Flags MUST NOT bypass authorization or data migration requirements. |
| NFR-MNT-007 | Country, provider and channel differences MUST be isolated behind explicit policy/configuration/adapter contracts rather than scattered conditional code. |
| NFR-MNT-008 | All production components, jobs, schemas and runbooks MUST have a named accountable owner and escalation path. |

## Observability and supportability

| ID | Proposed requirement |
|---|---|
| NFR-OBS-001 | Every user command MUST be traceable across client, API, transaction, outbox/worker and provider by privacy-safe correlation and idempotency identifiers. |
| NFR-OBS-002 | Service-level indicators MUST include availability, latency, error outcome, queue age, synchronization lag, conflict rate, notification outcome, provider health and reconciliation backlog. |
| NFR-OBS-003 | Alerts MUST map to a user/business impact, owner and runbook. Alerting on infrastructure without an actionable threshold is insufficient. |
| NFR-OBS-004 | Logs, metrics and traces MUST have bounded retention/cardinality/cost and preserve tenant/privacy boundaries. |
| NFR-OBS-005 | A support user MUST be able to explain a booking/payment/message state from authorized audit and delivery evidence without direct database editing. |
| NFR-OBS-006 | Health dashboards MUST distinguish Ubi failure, customer network/device failure and named provider failure where evidence permits. |

## Cost, data and energy efficiency

| ID | Proposed requirement |
|---|---|
| NFR-EFF-001 | Unit economics MUST measure first-party infrastructure plus provider cost per active organization, booking, notification, conversation and media unit. |
| NFR-EFF-002 | Each expansion feature MUST identify its recurring provider, storage, egress, support and compliance cost before release approval. |
| NFR-EFF-003 | Synchronization MUST be incremental; polling, duplicate payloads and automatic full-resolution media downloads are prohibited unless justified. |
| NFR-EFF-004 | Retention and cache policies MUST bound storage growth without deleting authoritative or legally required records. |
| NFR-EFF-005 | Performance optimizations MUST be evaluated on D1 battery, thermal behaviour and mobile-data consumption, not only server latency. |

## Compatibility and graceful degradation

| ID | Proposed requirement |
|---|---|
| NFR-COMP-001 | The minimum Android version and supported browser matrix MUST be based on target-market evidence and published before pilot enrollment. |
| NFR-COMP-002 | Unsupported devices/browsers MUST receive a clear, safe fallback or explanation; they must not proceed into a flow known to lose work. |
| NFR-COMP-003 | Provider/channel capability differences MUST be displayed truthfully; Ubi must not promise read history, delivery, calls or interactive controls that a provider does not expose. |
| NFR-COMP-004 | Android process death, rotation/configuration change, background restriction and operating-system upgrade MUST not silently discard a committed local command. |

## Measurement governance

Each quality metric must record:

- definition, unit and measurement boundary;
- eligible and excluded events;
- device/network/region/release segmentation;
- target percentile/window and minimum sample size;
- instrumentation owner and privacy basis;
- dashboard, alert and response runbook;
- last validation date and known limitations.

Client telemetry is sampled and minimized. Synthetic tests detect regressions; real-user measurements confirm field performance; neither replaces task-based usability research.

## Release gates

A first production release cannot be approved until:

1. The D1/D2/W1 device and N0–N3 network profiles are field-validated.
2. Critical journeys have traceable performance, accessibility, security, integrity and recovery tests.
3. Offline process-death, retry, duplication and conflict tests pass.
4. Capacity tests meet the approved forecast and headroom.
5. A production-like backup restore completes within the approved objectives.
6. Security review has no unaccepted release-blocking findings.
7. Observability can distinguish committed, queued, rejected and provider-dependent states.
8. Pilot users representing the target skill level meet the approved task-completion gate.

## Authoritative external baselines

- [Android app startup time](https://developer.android.com/topic/performance/vitals/launch-time) defines TTID/TTFD concepts and startup measurement.
- [Android slow rendering guidance](https://developer.android.com/topic/performance/vitals/render) describes slow and frozen frame signals.
- [Android offline-first architecture](https://developer.android.com/topic/architecture/data-layer/offline-first) supports local reads, durable queues and explicit write strategies.
- [Core Web Vitals thresholds](https://web.dev/articles/defining-core-web-vitals-thresholds) define the proposed LCP, INP and CLS field targets.
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/) is the application verification baseline.
- [OWASP Top 10](https://owasp.org/www-project-top-ten/) and [OWASP API Security Top 10](https://owasp.org/www-project-api-security/) inform threat coverage.

External sources are reference baselines. Ubi remains responsible for its target-market validation, architecture and controls.

## Open decisions

- Launch country/city and locally representative D1/D2 hardware.
- Pilot vertical and critical task set.
- Commercially affordable availability, RPO and RTO after impact/cost analysis.
- Supported Android/browser versions and upgrade policy.
- Approved client telemetry, consent/notice and retention.
- Provider-specific latency/failure budgets.
- Contractual SLA, service credits and exclusions, if offered later.
