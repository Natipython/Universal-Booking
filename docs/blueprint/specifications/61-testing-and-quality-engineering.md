# Testing and quality-engineering architecture

Status: **Draft — test strategy and release evidence model**  
Depends on: requirements, domain/workflow/data/security policies, quality attributes and repository architecture  
Blocks: CI/CD approval, release contract, pilot and production deployment

## Purpose

Ubi must prove that a change preserves booking, money, authorization, offline, privacy and provider behavior under failure—not merely that a screen opens. This document defines the layered evidence required for every deployable and vertical slice.

Diagram: `diagrams/61-quality-evidence-pipeline.mmd`.

## Principles

- Tests are executable evidence linked to requirements, risks and invariants.
- Correctness at the authoritative transaction boundary outranks UI snapshots and coverage percentage.
- Every production defect adds the smallest reliable regression test at the correct layer.
- Mocks support fast tests but never replace database, device, browser, provider-sandbox and restore evidence.
- Failure, retry, concurrency, cancellation, process death, weak network and ambiguous external outcome are normal test cases.
- Test telemetry and data follow the same privacy/secret rules as production.
- A flaky red test is a product risk, not harmless noise.

## Evidence layers

| Layer | Main purpose | Examples | Normal cadence |
|---|---|---|---|
| Static/schema | Prevent known structural defects | formatting, lint, type checks, API/event/schema compatibility, SQL lint, dependency/secret/license checks | Every change |
| Unit/value | Prove pure rules cheaply | money, time/calendar, policy version, state transition, error mapping | Every change |
| Property/model | Explore large input/state spaces | price arithmetic, booking lifecycle, idempotency, sync convergence, parser invariants | PR plus deeper scheduled run |
| Database/integration | Prove real transaction/storage behavior | PostgreSQL constraints, RLS, locks, outbox, migration, restore queries | PR/merge/release by cost |
| Contract | Prove component/provider compatibility | OpenAPI consumers, webhook fixtures/signatures, OAuth/callback schemas | PR plus scheduled provider checks |
| Client component | Prove UI/application behavior | Flutter widget/state/local DB, React semantic forms/SSR/hydration | Every relevant change |
| End-to-end | Prove critical journeys | owner setup, manual booking, offline submit, customer link, conflict, reconciliation | Release candidate |
| Non-functional | Prove quality budgets | accessibility, performance, load, security, battery/data, resilience, recovery | Scheduled and release gates |
| Field/usability | Prove human success | Ethiopian merchants/customers on representative devices/languages | Before pilot and material workflow release |

## Traceability

Each release-critical test records:

```text
requirement/risk/control IDs
scenario and actor
preconditions and data classification
environment/build/schema/provider fixture version
expected domain/audit/user-visible outcome
failure/retry/recovery assertions
owner and evidence location
```

A dashboard pass without versioned test output is not release evidence.

## Domain and invariant suites

Mandatory invariant coverage includes:

- one exclusive resource cannot have overlapping held/committed reservations;
- hold-to-book conversion is one authoritative transaction and expiry is server-timed;
- one idempotency key plus request hash produces one compatible logical result;
- tenant/operation identity cannot change across an owned record relationship;
- permission/revocation is checked at command execution, not only route rendering;
- payment reported, processing, verified, failed, reversed and refunded remain distinct;
- exact currency arithmetic never uses binary floating point;
- committed booking survives notification/search/analytics failure;
- external unknown outcome remains action-required/reconcilable, never guessed success;
- policy/price/customer-visible terms are immutable snapshots for the transaction;
- sync replay, duplication and reordering converge or produce explicit conflict;
- deletion/hold/restriction propagates without silently deleting required evidence.

Use example tests for known rules and property/state-machine tests for combinations and long sequences.

## Concurrency and fault tests

Run deterministic barriers/failpoints around:

- two customers/agents selecting the last capacity;
- hold expiry racing confirmation;
- duplicate Android commands and webhook replays;
- worker lease expiry during provider call;
- database commit succeeds but response is lost;
- provider accepts request but Ubi times out;
- role revoked while queued offline command synchronizes;
- deployment overlap with old/new event/schema versions;
- connection-pool exhaustion, deadlock/serialization retry and failover;
- outbox publication before/after crash and inbox deduplication.

Assertions cover database state, external-call count/evidence, audit trail, client status and repair path.

## PostgreSQL, migration and tenant-isolation tests

- Run integration tests against supported real PostgreSQL, not an in-memory substitute.
- Exercise constraints, triggers/functions where approved, RLS under actual runtime roles and composite tenant foreign keys.
- Generate two or more tenants and attempt cross-tenant reads/writes/joins/exports/worker jobs through every route and direct repository boundary.
- Test migration from oldest supported production schema through current, including realistic volume/locks/timeouts and deploy overlap.
- Test expand/backfill/contract separately; destructive contract never runs in the same unsafe window as old code.
- Verify downgrade/roll-forward procedure and backup/restore checkpoint for risky migrations.
- Explain query plans and performance-regress critical queries against representative distributions.

## Android/Flutter test matrix

### Fast suites

- Dart unit/property tests for local use cases, state reducers, validation and sync classification;
- widget/golden tests for semantics, text scaling, themes, loading/offline/conflict/error states;
- local-database migration/outbox atomicity and kill-after-save tests;
- generated contract decoding, unknown fields and old/new server compatibility;
- platform-plugin fakes plus native adapter unit tests.

### Device suites

Run release builds on purchased D1/D2 devices and supported Android versions:

- cold/warm start, frame/ANR/crash/memory/storage/battery/network budgets;
- screen reader, font scaling, contrast, touch targets, rotation and larger displays where supported;
- N0/N1/N2 switching, captive/partial connectivity, DNS/TLS failure and metered mode;
- process kill, reboot, app update, local schema migration, background limits and low storage;
- permission denial/revocation, deep-link replay, notification behavior and secure local cleanup;
- offline create/edit/conflict/retry and session/role expiry.

Emulators widen coverage but cannot approve performance, OEM/background or accessibility behavior.

## Web/React test matrix

- server-rendered HTML contains business/service/price method/next action before hydration;
- critical customer paths remain meaningful with JavaScript delayed/failed where designed;
- semantic labels, focus order, keyboard, screen reader, zoom/reflow and error summary pass;
- scoped links do not leak through URL/query/referrer/analytics/cache;
- public/private route caching and authentication boundaries are adversarially tested;
- bundle/payload/Core Web Vitals budgets run on W1/N1 profiles;
- current approved browser/OS matrix and constrained Windows W2 are exercised;
- locale, Ethiopian/Gregorian calendar display and timezone/DST test vectors are fixed and independently reviewed.

## Native, C, C++ and Rust test requirements

Any native module adds:

- C ABI compatibility tests for every supported architecture/toolchain;
- address/undefined/thread/memory sanitizers where supported;
- fuzzing of every untrusted byte/parser boundary with retained corpus;
- ownership/lifetime/double-free/use-after-free/out-of-bounds tests;
- malicious size, integer overflow, cancellation, timeout and parallel-call tests;
- safe-wrapper tests in Dart/Kotlin/Swift/Go/Rust;
- load/unavailable/wrong-version fallback;
- reproducible release build, symbols, crash unwinding and CVE/OSS inventory.

Benchmark output alone cannot approve unsafe code.

## API, event and provider contracts

- Producer verifies schema and semantic invariants; each supported consumer verifies fixtures.
- Additive/breaking classification is automated; removal follows published compatibility window.
- Provider fixtures include signatures, timestamp/replay, missing/unknown fields, duplicates, pagination, rate limits, errors and policy-window rejection.
- Sandbox/live-canary checks use dedicated safe accounts and cannot send to real customers unexpectedly.
- Provider documentation examples are inputs, not substitutes for captured approved test cases.
- Contract drift disables unsupported capability and alerts rather than corrupting state.

## Payment quality gates

- reported evidence never becomes verified without licensed-provider evidence;
- duplicate/out-of-order/missing/forged callbacks and polling reconciliation;
- amount/currency/merchant/provider-reference mismatch quarantine;
- timeout before/after provider acceptance and idempotent retry;
- reversal/refund/partial/failed/unknown scenarios only where provider/model supports;
- cash/pay-at-location remains independent of digital-provider availability;
- daily reconciliation counts and sums prove ledger/provider alignment;
- no test uses live money without written, bounded approval and reconciliation owner.

## Security, privacy and abuse testing

- ASVS/MASVS/API-derived control matrix with evidence and justified exclusions;
- authentication/recovery/session/step-up and cross-tenant authorization adversarial suites;
- secret/PII/log/analytics/build-artifact scanning;
- webhook/OAuth/scoped-link replay, SSRF/file/media and rate/resource abuse;
- support JIT/break-glass and privileged-action audit tests;
- consent/withdrawal/marketing suppression, rights export/deletion and legal-hold interaction;
- dependency/container/IaC/mobile security scanning plus independent pre-launch assessment;
- moderation report/block/appeal and weak-signal fairness tests.

Security scanners do not replace threat-specific tests or review.

## Performance, capacity and resilience suites

- Microbenchmarks identify local regressions but do not prove user speed.
- End-to-end D1/W1/W2/N1 journeys validate real budgets.
- Load profiles include steady forecast, 3× peak, short 10× idempotent ingress burst, hot tenant, provider slowdown and database saturation.
- Soak tests expose leaks, timer/queue growth, connection churn and degraded cleanup.
- Fault injection covers process/zone/database/provider/cache/object/notification/telemetry loss within safe environments.
- Restore exercises validate RPO/RTO, constraints, tenant isolation, keys/configuration, objects and provider reconciliation.

## Test data and environments

- Generated synthetic people/businesses/messages/payments are the default.
- Production data never enters developer/CI/shared test environments by default.
- Sanitization is not assumed safe merely because names are changed; exceptional production-derived data requires privacy/security approval and isolated expiry.
- Fixtures include Ethiopian names/scripts/phone formats/timezones/calendars/currencies without copying real people.
- Test provider credentials/accounts are separate, least-privilege and rotated.
- Each run records seed so failures can be reproduced.

## Flaky-test and failure policy

- A failing required test blocks; rerunning until green without diagnosis is prohibited.
- Flake is recorded with owner, evidence, quarantine scope and short expiry.
- Quarantine cannot cover booking/money/authorization/migration/security critical gates.
- Track flake rate, time to repair, escaped defects and false-green incidents.
- Delete redundant/brittle tests only with replacement evidence; suite speed is a design concern.

## CI and release gates

### Pull request fast gate — target ≤15 minutes

- formatting/lint/type/schema checks;
- unit/property-selected/component tests;
- dependency/license/secret/static-security checks;
- affected contract/database verification;
- build of affected release artifacts where feasible.

### Merge/nightly gates

- broader property/concurrency/PostgreSQL/RLS/provider fixtures;
- device/browser matrix subsets;
- native fuzz/sanitizer and dependency updates;
- performance regression and migration rehearsal subsets.

### Release-candidate gate

- immutable candidate deployed to staging;
- critical end-to-end/offline/accessibility/security/payment scenarios;
- full migration/rollback, load and restore evidence proportional to change;
- D1/W1/W2 production-like performance;
- known-risk/exception list with accountable approvals;
- runbooks, dashboards, support/training and release communication verified.

## Quality reporting

Report by requirement/risk and release—not just pass count:

- critical scenario/invariant coverage;
- defects by severity/source/escape stage and recurrence;
- test flake/quarantine/median and p95 duration;
- device/browser/provider/schema coverage;
- performance/accessibility/security/recovery results;
- unresolved risk, owner, due date and release decision.

Code coverage highlights unexecuted areas but no universal percentage proves quality. Risk-critical modules may set local thresholds alongside mutation/property evidence.

## Open decisions

- Exact test libraries/runners and managed device/browser providers.
- Purchased D1/D2 models and supported Android/browser matrix.
- Field usability languages, sample sizes and task thresholds.
- Independent security/performance/accessibility reviewers and budget.
- Provider sandbox/live-canary availability.
- Release compatibility window and CI compute/cost budget.

