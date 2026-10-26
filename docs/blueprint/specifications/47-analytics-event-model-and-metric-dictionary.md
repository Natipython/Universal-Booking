# Analytics event model and metric dictionary

Status: **Draft — measurement architecture proposed; product baselines and consent/notice unresolved**  
Depends on: workflows, requirements, event architecture, data inventory and retention framework  
Blocks: analytics schemas/SDKs, KPI dashboard, experiments, pilot scorecard and data contracts

## Purpose

Ubi needs evidence that small businesses can set up, accept and manage bookings reliably on affordable devices and imperfect networks. Analytics must answer those questions without copying private bookings, contacts or conversations into an uncontrolled event warehouse.

This document separates operational truth, domain integration events, audit evidence, product interaction events and derived metrics. One event must not silently serve every purpose.

## Measurement principles

- Start from a decision/question, not from “track everything.”
- Authoritative domain tables remain the source for commercial counts and money.
- Events are versioned facts with bounded properties, not arbitrary JSON dumps.
- No raw message, note, attachment, address, token, OTP or payment credential in analytics.
- Tenant, account and device identifiers are minimized/pseudonymized according to purpose.
- Offline event time and server receipt time are both retained; neither is silently substituted.
- At-least-once delivery means every aggregation deduplicates by event ID/business key.
- Metrics publish formula, denominator, filters, timezone, freshness, owner and caveats.
- Every percentage exposes its denominator and sample size.
- A target is not a fact until pilot data establishes a baseline and leadership approves it.
- Analytics outages never block booking, sync, payment or messaging work.

## Five distinct evidence types

| Type | Purpose | Example | Owner/source |
|---|---|---|---|
| Owner/domain record | Current and historical business truth | Booking transitioned to confirmed | Booking PostgreSQL aggregate/history |
| Domain event | Notify other modules of committed fact | `booking.confirmed` | Transactional outbox after commit |
| Audit/security event | Prove significant access/action/decision | Staff exported customer list | Protected audit interface |
| Product interaction event | Understand UX flow/performance | User opened calendar or abandoned setup step | Approved client/server analytics contract |
| Metric/projection | Governed calculation for a decision | Confirmation success rate by app/network cohort | Reporting pipeline + definition version |

An analytics click cannot confirm a booking. An audit log is not a behavioral analytics feed. A domain event payload is not automatically safe for third-party analytics.

## Initial pipeline decision

For pilot:

1. authoritative facts commit to PostgreSQL;
2. outbox emits bounded versioned domain events;
3. a reporting consumer deduplicates and builds daily/hourly projections in PostgreSQL;
4. Android/web product events enter a first-party ingestion endpoint, are validated/limited, then queued/stored under the analytics retention policy;
5. dashboards query governed definitions/projections;
6. external analytics warehouse/vendor is deferred until volume, analysis need and privacy review justify it.

Crash/performance tooling may be selected separately after vendor/data-location review, with payload scrubbing and sampling. Do not install multiple overlapping SDKs in the Android pilot.

## Canonical event envelope

Every event contract declares required/optional fields and property allow-list. Conceptual envelope:

```text
event_id                 globally unique, dedup key
event_name               controlled namespace, past-tense fact/action
schema_version           explicit positive version
occurred_at              device/domain occurrence instant
recorded_at              trusted ingestion/commit instant
producer                 app/module/service and version
environment              production/staging/etc.
organization_id          only where tenant fact needs it
aggregate_type/id/version for domain facts
actor_type               account/member/customer/service/anonymous/workforce
actor_pseudonym          only when measurement needs continuity
session/installation pseudonym where approved
correlation_id / causation_id
source_channel           Android, customer web, walk-in, provider, API, etc.
locale/timezone/country cohort where justified
network/device/app cohort bounded values where approved
properties               schema-validated allow-listed non-secret fields
privacy_class/purpose_id
```

Do not assume events from different devices/services have a global order. Aggregate version orders one aggregate; reporting handles late/out-of-order events through event ID, source version and bounded correction windows.

## Naming and versioning

Examples:

```text
domain.booking.created.v1
domain.booking.confirmed.v1
domain.hold.expired.v1
product.onboarding.step_completed.v1
product.booking_form.validation_failed.v1
quality.sync.command_resolved.v1
reliability.api.request_completed.v1
```

Rules:

- Event name describes one fact, not a UI label such as `green_button_click`.
- Breaking semantic/property changes create a new schema version.
- Producers and consumers declare supported versions in the schema registry.
- Additive optional changes still receive review and documentation.
- Removed properties stop collection and are deleted/expired according to inventory—not merely ignored.
- Unknown events/properties quarantine or reject; they do not flow silently.

## Product event safety contract

Allowed properties are enums, bounded numbers, booleans, durations, coarse cohorts and opaque/pseudonymous IDs needed for an approved metric.

Disallowed by default:

- names, phone numbers, emails and exact addresses;
- service/customer free text, form answers, messages and call content;
- file names/URLs and scoped-link/token values;
- raw exception requests/responses/SQL/provider payloads;
- fine GPS or sensitive device advertising identifiers;
- arbitrary screen text or accessibility labels that may contain customer data;
- payment reference/credentials or exact private financial narrative.

Server rejects extra keys, oversized payloads, future timestamps beyond tolerance and invalid enums. Rate limits and batch limits protect ingestion. Debug sampling uses a separate short-lived controlled path.

## Offline and low-connectivity events

Android keeps a bounded encrypted analytics queue only when collection is allowed. Product work has priority over analytics traffic.

- Batch, compress and upload on suitable connectivity/battery policy.
- Bound event count, bytes and age; discard lowest-value analytics safely when full.
- Keep `occurred_at`, monotonic sequence where available, boot/session ID and `recorded_at`.
- Mark offline/cached/pending context with bounded flags, not network contents.
- Deduplicate retries by event ID.
- Respect opt-out/purpose changes before upload and purge queued disallowed events.
- Never make a user repeat a booking because analytics upload failed.

## Metric definition template

Every metric in `reporting.metric_definitions` records:

```text
metric_id/name/version and plain-language question
owner and decision/action when it changes
numerator, denominator and unit
authoritative sources/event versions
population, eligibility, exclusions and dedup rule
time window, timezone/calendar and late-arrival policy
dimensions allowed and minimum-group privacy rule
freshness/completeness/accuracy objective
baseline/target/guardrail status
known biases and change history
```

## Initial north-star and guardrails

The provisional north-star is:

> Weekly completed legitimate booking work managed through Ubi by retained pilot organizations.

It is deliberately stricter than app opens, registered businesses or booking-created counts. “Legitimate,” “completed,” “retained” and pilot eligibility need exact approved definitions.

North-star guardrails prevent optimization that harms people:

- confirmed double-booking incident rate;
- cancellation/no-show and customer complaint rate interpreted by vertical;
- cross-tenant/privacy/security incidents;
- payment reconciliation exceptions/loss;
- crash-free use and command success under weak networks;
- staff task time/error/help dependency;
- data use and operating cost per completed booking;
- notification/message abuse and provider block rate.

## Pilot metric dictionary

Values/targets are `TBD-BASELINE`; the definitions come first.

### Merchant value and adoption

| Metric | Definition direction | Why/caveat |
|---|---|---|
| Setup completion rate | Eligible pilot organizations completing required publishable config / organizations that started assisted or self setup | Segment assisted versus self setup; neither proves value |
| Time to first value | Organization activation to first legitimate booking reaching the chosen value state | Exclude staff/demo/test data; report distribution |
| Weekly active value organizations | Eligible organizations with at least one legitimate managed/completed booking action in local week | Do not call app open “active value” |
| Four-week retained value organizations | Activated cohort producing approved value in week 4 window | Small pilot cohorts need counts and confidence caveat |
| Booking work per active organization | Legitimate completed/managing actions / active value organizations | Median and distribution prevent large merchant distortion |
| Manual fallback rate | Booking tasks leaving Ubi for paper/call/spreadsheet due to missing/failed workflow / observed eligible tasks | Requires lightweight research/support tagging, not surveillance |

### Booking/customer funnel

| Metric | Definition direction | Why/caveat |
|---|---|---|
| Booking flow start-to-submit | Unique eligible flows submitted / unique eligible starts | Define inactivity timeout and cross-device behavior |
| Submit success rate | Accepted idempotent booking/request commands / valid submit attempts | Separate validation, conflict, network and server failure |
| Confirmation conversion | Confirmed bookings / eligible submitted standard-booking attempts | Instant/manual-confirm verticals are separate cohorts |
| Custom request quote rate | Requests receiving an issued quote / eligible complete requests | Exclude spam/test/withdrawn-before-review explicitly |
| Quote acceptance rate | Accepted issued quote options / eligible viewed/issued quote versions | Define denominator and expiry carefully |
| Completion rate | Fulfilled bookings / eligible confirmed bookings whose service window matured | Lagged cohort; cancellation/no-show separate |
| Customer correction/help rate | Flows requiring staff/support correction / completed customer flows | Useful proxy for “idiot-proof”; reason codes needed |

### Speed and usability

| Metric | Definition direction | Why/caveat |
|---|---|---|
| Time to create simple booking | Start of staff/customer task to server-accepted booking, median/p75/p90 | Break out data entry versus network/server wait |
| Time to next action | Work becomes actionable to first valid staff action | Business hours/timezone and staffing affect it |
| Form validation recovery | Validation failure followed by successful correction in same flow / validation failures | Do not collect entered value |
| Assisted help dependency | Eligible tasks needing trainer/support intervention / observed tasks | Capture during pilot research with clear method |
| Accessibility task success | Completed critical tasks for tested accessibility/literacy cohorts | Requires moderated usability tests, not analytics alone |
| Language-switch completion | Completion after language switch / eligible switched flows | Small sample; language quality research also needed |

### Offline, performance and reliability

| Metric | Definition direction | Why/caveat |
|---|---|---|
| Cold/warm screen-ready time | Time until meaningful locally available content is usable | Measure low-end device cohorts and p50/p90/p99 |
| Server latency | Trusted server duration by endpoint/result, p50/p90/p99 | Separate network round trip/device perceived latency |
| Command terminal success | Commands reaching accepted/rejected-known result / commands submitted | A truthful rejection is not an infrastructure failure |
| Sync convergence time | Server acceptance to client projection visible, distribution | Segment foreground/background/network cohort |
| Offline command conflict rate | Commands needing user resolution / offline commands replayed | Also count auto-resolved/rejected categories |
| Forced full-resync rate | Full bootstrap syncs / active installations/time | High rate signals protocol/storage defects |
| Crash-free critical task rate | Critical task attempts without crash/ANR / attempts | Session-level crash-free alone can hide task failure |
| Payload/data use | Compressed bytes uploaded/downloaded per critical task and active org | Important cost/accessibility guardrail |

### Correctness, safety and operations

| Metric | Definition direction | Why/caveat |
|---|---|---|
| Confirmed double-booking incidents | Verified same-resource conflicting active commitments | Target zero; distinguish configured capacity/authorized override |
| Capacity conflict rejection | Capacity-changed rejects / valid commit attempts | Expected under concurrency; trend and UX handling matter |
| Idempotent replay rate | Commands returning prior result / accepted commands | Helps find network/client retry conditions, not an error itself |
| Reconciliation exception rate | Unresolved payment/provider discrepancies / governed items/value | Report count and exact currency value separately |
| Notification delivery outcome | Authenticated provider delivered/failed/unknown by purpose/channel | “Sent” is not “read”; provider semantics differ |
| Isolation/security/privacy incident | Confirmed incidents by severity and affected scope | Never normalize nonzero high-severity events as acceptable growth |
| Support cases per value action | Product/support cases / completed value actions | Separate onboarding, defect, education, provider outage |

## Financial metrics

Money metrics come from exact currency-qualified payment/billing owner records, not client analytics.

- Never sum different currencies without an explicit dated FX source/version and separate native totals.
- Separate merchant customer funds from Ubi subscription revenue.
- Distinguish obligation, reported/unverified, verified, allocated, refunded, reversed, disputed and settled values.
- Gross booking value is not Ubi revenue, cash received or merchant settlement.
- Payment provider fees/taxes/refunds have exact metric definitions and reconciliation.

## Dimensions and privacy

Approved useful dimensions may include launch market, booking archetype, operation, source channel, assisted/self setup, app version, OS class, coarse device capability, connectivity cohort, locale and new/retained cohort.

Rules:

- avoid unique/small group combinations that expose a person/business;
- organization-level operational dashboards are visible only to that authorized business;
- platform analysis uses minimized IDs and minimum group thresholds where appropriate;
- protected/sensitive traits are not inferred for segmentation without a separately approved need;
- no employee ranking/surveillance dashboard by default;
- IP/location/device identifiers follow inventory/retention and are not casually retained.

## Metric computation and corrections

- Domain state transitions provide authoritative lifecycle facts.
- Daily metrics use the organization/operation’s stored local timezone; UTC remains available for platform operations.
- Late events enter a bounded correction window and change metric revision/checkpoint.
- Rebuilds are deterministic from versioned sources and record code/definition version.
- Backfills never mix old/new definitions under one metric version.
- Dashboard displays freshness, last successful projection and partial-data warning.
- Reconciliation samples metrics back to owner rows and investigates unexplained drift.

## Experiment governance

Experiments require a hypothesis, eligible population, assignment unit, primary metric, guardrails, minimum duration/sample reasoning, stop rule, exclusion, privacy review and owner.

- Assignment is deterministic and recorded by version.
- Never experiment on security/privacy disclosures, unauthorized dark patterns or irreversible high-risk payment/booking behavior without explicit governance.
- Small East African pilot samples favor qualitative usability plus simple evidence over false statistical certainty.
- Report absolute counts/effect ranges and operational context, not only a percentage uplift.
- Ended experiments remove flags/instrumentation on schedule.

## Data quality and monitoring

Pipeline monitors:

- accepted/rejected event count by schema/producer/app version;
- duplicate, late, future-time and out-of-order rates;
- outbox/ingestion/projection lag;
- missing required/cohort property rate;
- domain-record versus event/projection reconciliation;
- sudden cardinality/property-volume/cost increase;
- forbidden-field/redaction scanning;
- dashboard freshness and backfill state.

Bad analytics data quarantines and alerts; it never blocks the owner workflow.

## Release gates

- Each KPI links to a decision, owner, exact versioned definition and authoritative source.
- Pilot test/demo organizations and events are reliably excluded.
- Product event schemas reject unknown/high-risk properties and pass privacy/security review.
- Android batching works under offline, metered, low-storage and consent/choice changes.
- Metrics reconcile against synthetic owner-domain fixtures and deliberate duplicate/late events.
- Dashboards show denominator, counts, timezone, freshness and definition version.
- Crash/performance vendor, data location, retention and payload scrubbing are approved.
- Baselines are collected before targets are set; no fabricated adoption/performance claim enters planning.
- Qualitative research plan covers below-average digital literacy and accessibility, because telemetry cannot explain every failure.

## Open decisions

- Exact pilot north-star value state and legitimate/test exclusion.
- Analytics notice/choice/consent requirements by launch market and vendor.
- First-party event storage shape and whether an external warehouse/vendor is needed later.
- Pseudonym rotation/linkage rules for account, organization, installation and anonymous flow.
- Minimum group/reporting privacy threshold.
- Approved crash/performance tooling and low-end device/network cohorts.
- KPI baselines, targets and pilot stop/go thresholds.
- Research method for manual fallback, help dependency and task comprehension.
