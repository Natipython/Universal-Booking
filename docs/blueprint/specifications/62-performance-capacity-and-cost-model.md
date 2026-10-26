# Performance, capacity and cost model

Status: **Draft — budgets are engineering targets, not public promises**  
Depends on: quality attributes, deployment/data architecture, code strategy, pilot scope and vendor shortlist  
Blocks: provider sizing, performance approval, pricing/unit economics and scale/release decisions

## Purpose

Ubi's “fast and smooth” goal must be measured on Ethiopian devices/networks and paid for sustainably. This model ties user latency, system capacity and per-merchant/per-booking cost together so optimization does not move a problem from screen to server, server to provider or performance to an unaffordable bill.

Diagram: `diagrams/62-performance-capacity-cost-loop.mmd`.

## Measurement hierarchy

1. **Task outcome:** booking/manual entry/conflict/payment status completes correctly.
2. **User experience:** input response, usable content, frames, clear pending/offline state.
3. **End-to-end latency:** device/web edge → Ubi → database/provider → visible result.
4. **Component latency/resource:** route/query/queue/native function CPU, memory, I/O.
5. **Capacity and cost:** headroom, saturation, vendor units and operator effort.

Optimizing layer 4 while worsening task success, battery, data usage, correctness or total cost is a regression.

## Provisional release budgets

The detailed source is specification 21. Stage 8 must continuously measure at least:

| Area | Provisional target/gate |
|---|---|
| D1 Android cold usable screen | ≤3.0 s p75; cached content first |
| D1 Android warm usable screen | ≤1.0 s p75 |
| Local interaction acknowledgment | ≤100 ms |
| D1 common frame behavior | p99 ≤33.3 ms and missed-frame ratio <5%; no >700 ms freeze |
| Normal sync/booking command data | ≤200 KiB incremental sync with no attachments; ≤25 KiB typical command |
| Android artifact | aim ≤35 MB; mandatory review >50 MB |
| Customer web | LCP ≤2.5 s, INP ≤200 ms, CLS ≤0.1 at p75 field target |
| Customer critical route | ≤250 KiB compressed excluding media; initial JS aim ≤150 KiB |
| First-party API | read p95 ≤500 ms; command acceptance p95 ≤800 ms under approved load |
| Availability query | p95 ≤1 s for bounded launch workload |
| Booking-path availability | 99.9% monthly objective after general launch; pilot measures, does not contractually promise |
| Android stability | crash-free sessions ≥99.8%; user-perceived ANR <0.1%; no known critical-flow defect |
| Recovery | proposed authoritative RPO ≤5 min; booking-core RTO ≤4 h, only after tested approval |

Percentiles are segmented by release/device/network/route/provider. A fast global average cannot hide D1/N1 failure.

## Workload model

Forecast and measure:

- organizations, operations/branches, active staff and concurrent sessions;
- services/resources/availability rules and schedule horizon;
- manual, customer-link, walk-in/call and channel-originated enquiries/bookings;
- command retry/offline backlog and sync change volume;
- quotes, holds, timers, reminders and scheduled jobs;
- provider webhooks/messages/templates/media and rate-limit windows;
- reported/verified payments and reconciliation records;
- public page/link traffic, abusive scans and seasonal/event bursts;
- images/files/object transformations and egress;
- logs/metrics/traces/audit/analytics volume;
- support/onboarding cases and human minutes.

Track average, p95 tenant and largest expected tenant. Aggregate totals alone miss noisy-neighbor risk.

## Forecast equations

Use transparent inputs rather than invented “users”:

```text
daily_bookings = active_merchants × bookings_per_active_merchant_day
peak_booking_rate = daily_bookings × peak_fraction / peak_window_seconds
webhook_ingress = channel_events + payment_events + expected_retries
sync_egress = active_devices × sync_sessions × mean_delta_bytes
database_write_rate = commands × writes_per_command + workers + audit/outbox
monthly_variable_cost = provider_units + compute + database + storage + egress
                        + observability + support labor + tax/FX contingency
contribution_per_merchant = collected_recurring_revenue
                            - merchant-attributable variable cost
```

Every input records source, date, confidence, low/base/high case and owner. Pilot measurements replace assumptions.

## Capacity envelopes

Each component maintains:

| Dimension | Limit/evidence |
|---|---|
| Application | request concurrency, CPU, memory, open files, graceful replicas/start time |
| PostgreSQL | connections, TPS, locks, query p95/p99, WAL, IOPS, storage/growth, failover |
| Worker/outbox | enqueue/dequeue, queue age, attempts, lease/concurrency, poison backlog |
| Public edge | requests, bytes, rate controls, cache hit, TLS/connection limits |
| Object/media | upload/transform/egress, size/count, lifecycle and malware queue |
| Providers | per-account/merchant/template/API rate, timeout, retry, price, quota |
| Android | startup/frame/heap/local DB/storage/network/battery/background quota |
| Web | HTML/JS/CSS/font/media bytes, request count, CPU/hydration and cache |

Autoscaling cannot exceed downstream capacity silently. Application replica limits are calculated with database/provider budgets.

## Load and stress profiles

- **Baseline:** forecast normal day and peak hour.
- **Release:** sustained forecast peak ×3 without SLO/invariant loss.
- **Burst:** short ×10 idempotent webhook/public ingress with bounded backpressure.
- **Hot tenant:** one tenant at p99 expected activity while ordinary tenants retain service.
- **Retry storm:** provider/network recovery releases queued retries without duplicate effects.
- **Offline return:** many devices reconnect and synchronize after outage.
- **Degraded provider:** slow/time-out/rate-limited channel/payment while core manual booking remains healthy.
- **Soak:** long enough to expose leaks, backlog and storage/telemetry growth.
- **Limit:** controlled test identifies safe saturation and rejection behavior, never in production by surprise.

Successful test means correctness/audit/recovery and fair rejection, not merely no crash.

## Client performance engineering

### Flutter/Dart Android

- profile release/profile builds on D1/D2, never infer from debug;
- instrument cold/warm start phases, first useful local render and sync separately;
- minimize startup initialization/plugin registration; defer media/analytics/noncritical configuration;
- query/index local data for bounded screens; paginate and virtualize long lists;
- isolate CPU-heavy parsing/image work without spawning unbounded isolates;
- bound caches/logs/media and verify low-storage behavior;
- measure native bridge crossings, allocations, jank, ANR, battery and radio/data use;
- compare framework/native prototype if a target fails before reaching for C++.

### React/TypeScript web

- render critical content on server and stream/cache only within privacy policy;
- route-split customer/business/operations bundles;
- ship no unnecessary component/analytics/media code to booking links;
- responsive image/font strategy and declared dimensions;
- measure field Core Web Vitals plus W1/N1 laboratory journeys;
- protect server rendering from slow provider/data fan-out.

## Server and database performance

- profile end-to-end first; optimize queries/data access before language/native rewrites;
- bound every query/page/window/batch and require keyset pagination where appropriate;
- inspect query plans with representative skew and table size;
- keep transactions short, lock order deterministic and external calls outside DB transaction;
- connection pool has global budget across deploy overlap, workers and support tools;
- cache only where measured and correctness-safe; invalidate/version explicitly;
- outbox workers use fair tenant/provider concurrency and queue-age scaling;
- reject/backpressure before memory/goroutine/connection explosion.

## C/C++/Rust optimization gate

Native optimization is considered only after a reproducible profile shows a material D1/task/cloud-cost bottleneck and an algorithm/data-layout/framework solution is insufficient.

Approval compares:

- before/after task percentile, CPU, memory, battery, binary size and cost;
- FFI copy/call overhead and concurrency;
- fallback behavior and supported ABIs;
- fuzz/sanitizer/security/maintenance burden;
- total build/CI/store artifact and staffing cost.

Handwritten assembly requires an additional architecture-specific exception; compiler optimization/intrinsics are preferred.

## Cost ledger

Track fixed, step and variable costs separately:

- company/domain/certificates/developer-store and compliance costs;
- cloud edge/compute/database/backup/object/egress/KMS/secrets;
- SMS/OTP/push/email and messaging conversation/template costs;
- payment provider setup/transaction/refund/reconciliation/settlement costs;
- analytics/crash/log/trace/security/support tools;
- maps/media/call/AI costs only if approved;
- taxes, FX, minimum commitments and payment method limitations;
- engineering/on-call/support/onboarding/training/moderation minutes;
- failed/retried/abusive units that still incur vendor cost.

Allocate costs by causal unit where possible: active merchant, booking, message, GB, OTP, payment, support case—not arbitrary equal division.

## Cost controls

- Budget and forecast alerts at multiple thresholds; alert before provider suspension.
- Hard caps only where rejection is safe and clearly communicated.
- Tenant/provider/action quotas protect shared capacity and unit economics.
- Debug/trace sampling and retention respond to value/risk, not unlimited ingestion.
- Expire ephemeral environments, stale media, abandoned exports and test accounts.
- Detect retry loops, duplicate notifications, webhook storms and oversized payloads.
- Finance sees payable currency, tax and FX exposure; engineering sees unit consumption.
- A cheaper provider is not selected if onboarding, failures or support erase the savings.

## Performance/cost regression process

Every material release compares a stable baseline on identical conditions. Regression requires owner and disposition:

- correct before release;
- accept temporarily with quantified user/cost impact and expiry;
- change the budget only with product/engineering/finance evidence, never to make a failed test green.

## Pilot outputs

The pilot must produce distributions, not anecdotes:

- active merchants/staff/devices and task frequency;
- D1/D2/N1 user latency/stability/offline behavior;
- API/query/queue/provider percentiles and saturation;
- bytes, messages, storage, payment attempts and retries per completed booking;
- onboarding/support minutes and repeat help rate;
- cost per active merchant/completed booking and contribution scenarios;
- forecast error and revised low/base/high 12-month cases.

## Release gates

- Purchased representative devices and reproducible N0–N3 profiles exist.
- Critical budgets in specification 21 pass on immutable release candidate.
- Forecast peak ×3 and approved burst/degradation tests preserve invariants and fair service.
- Capacity dashboard exposes database, worker, provider, client and cost limits/headroom.
- Per-merchant/per-booking cost model includes vendor and human operations.
- No unowned cost/usage dimension or uncontrolled retry loop remains.
- Native optimization, if any, demonstrates net user/cost benefit and passes ADR-002/61 gates.
- Public pricing/SLA/performance claims do not exceed pilot/production evidence.

## Open decisions

- Real pilot demand/seasonality and D1/D2 devices/carriers.
- Hosting region/provider prices and payable currency/FX/tax.
- SMS/OTP/channel/payment vendor tariffs and minimums.
- Approved product/operations cost ceilings and gross-margin targets.
- Public SLA timing and any merchant plan differentiation.

