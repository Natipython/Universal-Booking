# Observability, SLO and reliability architecture

Status: **Draft — internal objectives, not contractual SLA**  
Depends on: quality/performance model, data/privacy/audit architecture, deployment and incident processes  
Blocks: alerting/tool procurement, on-call readiness, pilot reliability approval and public status commitments

## Purpose

Ubi must detect whether people can safely book and operate—not merely whether servers return HTTP 200. This document defines service-level indicators/objectives, privacy-safe logs/metrics/traces/crash data, alerting and error-budget governance.

Diagram: `diagrams/63-observability-and-reliability-loop.mmd`.

## Reliability model

Measure four layers separately:

1. **User task:** customer/merchant reaches a truthful completed/pending/conflict/rejected outcome.
2. **Domain correctness:** no tenant leak, false payment, double booking or lost accepted command.
3. **Service delivery:** API/webhook/sync/worker/provider adapter latency/availability.
4. **Infrastructure:** process, database, queue/outbox, object, edge and vendor health.

Infrastructure green cannot override a failed task or integrity incident.

## Proposed SLIs and SLOs

| Capability | SLI | Proposed internal objective | Exclusions/notes |
|---|---|---|---|
| Booking command | good authoritative outcomes / eligible attempts | 99.9% monthly after general launch | Pilot measures; invalid/unauthorized requests excluded, internal errors/timeouts included |
| Booking integrity | Ubi-induced allocation conflicts/incorrect commitments | Zero in approved tests and production | Any confirmed invariant breach is severe regardless of percentage |
| Webhook ingestion | authenticated accepted-or-explicitly-rejected / valid arrivals | 99.95% monthly | Provider delivery outside Ubi measured separately |
| API reads | p95/p99 latency and success for bounded routes | p95 ≤500 ms under approved conditions | Segment provider-independent/dependent |
| Command acceptance | p95/p99 to durable accepted/rejected result | p95 ≤800 ms under approved conditions | Final async completion has separate indicator |
| Sync | accepted commands reaching terminal server/client-visible state | Objective set after pilot; no unbounded unknown | Track queue age/conflict/rejection separately |
| Notifications | eligible notifications reaching provider accepted/delivered evidence | Provider/channel-specific objective | Never equate accepted with human seen |
| Payment status | verified provider events reconciled to correct merchant/amount/currency | 100% correctness; timeliness provider-specific | Customer report is not provider verification |
| Android stability | crash-free sessions and user-perceived ANR | ≥99.8%; ANR <0.1% pilot target | Critical-flow known crash blocks release |
| Recovery | achieved RPO/RTO in exercises | proposed RPO ≤5 min/RTO ≤4 h core | Not promised until tested/approved |

SLO documents include window, population, good/bad/unknown definition, data source, delay, owner and review date.

## Error budgets

For percentage SLOs:

```text
error_budget = eligible_events × (1 - objective)
burn_rate = bad_events_over_window / allowed_bad_events_over_same_window
```

Rules:

- Tenant/privacy/security/money/booking-integrity violations have separate zero-tolerance release/incident policy; they are not traded for feature speed.
- Fast and slow burn alerts avoid waiting for month end.
- When budget is exhausted or repeatedly at risk, reliability work/traffic reduction/release restraint takes priority by accountable decision.
- Exclusions are version-controlled and cannot be added after an incident to improve the number.
- Provider-caused user impact is shown even when excluded from first-party SLO responsibility.

## Telemetry contract

Every signal records, where applicable:

- environment, service/module/version/region/instance;
- event/metric/trace schema version;
- privacy-safe request/correlation/idempotency/command/booking/provider-reference hash or opaque ID;
- tenant/operation only in controlled pseudonymous form needed for isolation/noisy-neighbor diagnosis;
- route/use-case/outcome/error class and retryability;
- timing/queue/provider phase and payload size class;
- device/app/OS/network profile in approved coarse form;
- sampling/redaction decision and data classification.

Do not emit message/customer content, passwords, OTPs, tokens, payment credentials, scoped links or unrestricted SQL/bodies.

## Logs

- Structured fields and controlled event names; no concatenated arbitrary payload dumps.
- Severity means operational action, not developer preference.
- Expected validation/authorization conflicts are metrics/audit as appropriate, not noisy stack traces.
- Sensitive values are allowlist-filtered/redacted at source, collector and sink where possible.
- Logs have bounded size/rate/retention; repeated errors aggregate/sample without hiding first/last/count.
- Security/audit records remain a separate controlled evidence system; ordinary logs cannot be the sole audit trail.
- Clock/source/correlation allow timeline reconstruction without trusting device time.

## Metrics

Required metric groups:

- request/task rate, good/bad/unknown outcome and latency distribution;
- booking/capacity/idempotency conflict and terminal-state age;
- PostgreSQL pool/query/lock/deadlock/replication/WAL/storage/backup;
- outbox/inbox/timer queue depth/age/attempt/poison/fairness;
- provider latency/error/rate-limit/token/webhook/reconciliation by provider/capability;
- Android startup/frame/ANR/crash/sync/local-store/storage/network/battery in approved aggregate;
- web Core Web Vitals/payload/cache/hydration/error by route/device/network;
- authorization denial/cross-tenant anomaly/support privilege/abuse signals;
- data deletion/export/hold and telemetry-pipeline health;
- usage/cost/forecast/headroom by causal unit.

Avoid unbounded labels such as raw user/booking/message IDs. High-cardinality investigation uses controlled logs/traces, not metric labels.

## Traces

- Propagate W3C-compatible or approved trace context across edge, Go API, database spans, outbox/worker and provider adapters.
- A client-generated correlation ID is untrusted metadata; server binds it to authenticated/tenant context.
- Trace provider calls by operation/outcome without leaking URL query/token/body.
- Tail/risk-aware sampling may retain errors/slow/rare flows while keeping cost bounded.
- Async links represent outbox/worker causality rather than pretending one continuous synchronous span.
- Trace absence/drop rate is itself observable.

## Crash and client diagnostics

- Release symbols/source maps/mapping files are access-controlled and retained by version.
- Crash/ANR reports default to technical context; screen text/form/message values are excluded.
- User-attached diagnostics are optional, previewed/minimized and securely uploaded.
- D1 startup/frame/storage/sync telemetry is sampled enough to detect regressions without persistent tracking.
- App offers visible sync/support reference useful to a user without exposing identifiers/secrets.

## Dashboards

Minimum operational views:

1. executive service/task health and active incidents;
2. booking integrity and command terminal-state health;
3. Android/web experience by D1/W1/network/release;
4. API/PostgreSQL/outbox/worker capacity and saturation;
5. provider/channel/payment health and reconciliation;
6. security/privacy/support privilege signals;
7. backups/restore readiness and data lifecycle jobs;
8. usage/cost/headroom and noisy tenants;
9. telemetry pipeline loss/delay/sampling configuration.

Every chart states definition, source, owner, timezone/window, freshness and drill-down route.

## Alerting

Alert only when a human action or imminent user/risk consequence exists.

| Class | Example trigger | Response |
|---|---|---|
| Page/critical | booking integrity, cross-tenant exposure, payment false verification, rapid SLO burn, DB unavailable, lost accepted work | Immediate incident process |
| Page/high | webhook ingestion outage, severe queue age, authentication outage, backup window threatened | On-call response by defined target |
| Ticket | slow burn, capacity/cost forecast threshold, dependency expiry, rising conflict/help rate | Owned work within due date |
| Dashboard/info | expected deployment shift, low-volume provider issue | Observe/correlate; no wake-up |

- Multi-window burn alerts combine fast impact and sustained degradation.
- Deduplicate/correlate symptoms around one causal incident.
- Alert includes user impact, dashboard, safe first checks, runbook, owner and silence/escalation rules.
- Test routing and acknowledgments; an alert existing in configuration is not readiness.

## Reliability ownership and review

- Capability owner owns SLI definition and corrective backlog.
- Reliability/on-call owner owns monitoring/alert mechanics and exercises.
- Product/CEO/CFO approve public SLA or risk/cost tradeoff.
- Security/privacy/data owners approve relevant telemetry and incident handling.
- Weekly pilot review covers task success, error budget, top regressions, queue/provider unknowns and cost.
- Monthly/quarterly review adjusts objectives only from evidence and business impact.

## Provider and telemetry failure

- Ubi continues core manual/cached workflows when observability vendor is degraded; local bounded buffering cannot exhaust service.
- Losing telemetry never converts a failed/unknown payment or booking into success.
- Critical local health/audit has an independent minimal path where feasible.
- Monitor collector/export/sampling/ingestion/query/alert delivery separately.
- Provider status pages inform diagnosis but never replace Ubi measurements/reconciliation.

## User/status communication

- In-product states distinguish offline, Ubi degraded, provider delayed, queued, conflict and action required.
- Public status communication describes known impact and workaround without exposing security details or unsupported cause.
- Updates have owner/cadence; resolution states remaining reconciliation and user action.
- Support receives the same incident truth and approved script.

## Release gates

- Every critical command has correlation from client/API/transaction/outbox/provider to terminal state.
- SLI definitions and dashboards are validated against known synthetic failures.
- Privacy/secret/cardinality/cost tests pass for logs, metrics, traces and crash data.
- Paging alerts and escalation are exercised; no critical alert lacks an owned runbook.
- Deployment annotations, version/symbol/schema/provider capability and feature-flag state are queryable.
- Pilot can measure specification 21 budgets by device/network/release and provider.
- Telemetry outage is bounded and does not break booking/payment correctness.
- No public SLA is published before pilot/provider/recovery evidence and executive/legal approval.

## Open decisions

- Observability/crash/status vendors and hosting/data location.
- Exact SLO eligibility definitions and pilot baseline volumes.
- On-call hours/staffing/escalation and public status domain.
- Retention/sampling/cardinality/cost budgets.
- Synthetic monitoring locations/carriers/devices in Ethiopia.

