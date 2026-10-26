# Event, workflow and asynchronous processing architecture

Status: **Draft**  
Depends on: bounded contexts, platform components and API governance  
Blocks: outbox/inbox schema, workers, provider adapters, notifications, process managers, observability, recovery and testing

## Purpose

Ubi must coordinate committed bookings, notifications, payments, provider callbacks, expiries, projections and privacy work without pretending that a database and external provider share one transaction. This document defines the durable asynchronous model, delivery guarantees, event meanings and long-running workflow rules.

## Fundamental distinctions

| Artifact | Meaning | May fail/reject? | Primary owner |
|---|---|---:|---|
| Command | Request/intention to perform one domain action | Yes | Receiving domain |
| Command result | Accepted/committed, rejected, conflict, duplicate or retry-later outcome | N/A | Receiving domain |
| Domain event | Past-tense fact that committed inside an owning domain | No; it already happened | Producing domain |
| Integration event | Minimized stable fact published for other contexts/processes | Delivery may retry; fact remains true | Producing domain contract |
| Provider event | External provider’s claim/evidence, possibly duplicated/out of order | Authentication/normalization may reject/quarantine | Integration adapter |
| Notification intent | Purpose-bound instruction to communicate a domain fact/action | Delivery may fail | Communications |
| Audit event | Protected evidence of a significant action/decision/access | Recording is part of required action durability | Audit/source domain |
| Analytics event | Governed measurement for product/operations | May be sampled/delayed under approved policy | Analytics |
| Process-manager state | Durable coordination position across multiple commands/events/deadlines | May wait/retry/escalate | Workflow owner |

These must not share one generic `events` meaning or state model.

## Delivery model

### Direction

Ubi designs for **at-least-once delivery and exactly-once logical effect**.

- Networks, processes and providers can duplicate, delay or reorder work.
- Ubi does not claim infrastructure-level exactly-once delivery.
- Every command/event/side effect carries stable identity and is handled idempotently.
- The database commit plus durable outbox prevents “state changed but event disappeared.”
- Consumer inbox/deduplication prevents repeated logical effects.
- Reconciliation repairs ambiguous external outcomes.

### Initial transport

The authoritative relational database SHOULD initially provide:

- transactional outbox rows;
- scheduled work/timer rows;
- process-manager state;
- consumer/provider inbox/deduplication receipts;
- worker leases/attempts/dead-letter metadata.

Workers poll/claim bounded batches using safe database concurrency and backoff. A message broker is added only when measured latency, fan-out, volume, retention, provider isolation or operational needs exceed this design.

The domain contract is transport-independent so a later broker does not change event meaning.

## Transactional outbox pattern

### Producer transaction

```text
BEGIN
  verify idempotency and authorization evidence
  validate domain invariants
  update owner aggregate/current state
  append required domain history/audit evidence
  insert integration event/outbox envelope(s)
  persist command result/idempotency mapping
COMMIT
```

Only after commit may the application tell the client the authoritative action succeeded.

The transaction does not call SMS, social, payment, media, analytics or other external providers.

### Outbox lifecycle

```text
pending
→ leased
→ delivered

leased → pending          (retryable failure / lease expiry)
pending/leased → blocked  (dependency/configuration/action needed)
pending/leased → dead_lettered
dead_lettered/blocked → pending  (authorized replay after correction)
```

Fields include event/work ID, producer/tenant/aggregate, destination/consumer class, availability time, priority, attempts, lease owner/expiry, last safe error code, expiry, payload schema/version and created/delivered timestamps.

### Worker claiming

- Claim rows atomically with a short renewable lease.
- Expired worker lease makes work claimable again.
- Work processing is idempotent because a worker can die after external success but before marking delivered.
- Batch size and concurrency are bounded by queue/dependency/tenant fairness.
- High-priority security/payment/expiry work cannot permanently starve ordinary work; dedicated pools may be required.
- Database polling interval/load is monitored and becomes an extraction signal.

## Event envelope

Proposed logical envelope:

```json
{
  "event_id": "opaque-event-id",
  "event_type": "booking.confirmed",
  "event_version": 1,
  "occurred_at": "2026-08-30T10:15:30.123Z",
  "recorded_at": "2026-08-30T10:15:30.130Z",
  "producer": "booking",
  "tenant_id": "opaque-organization-id",
  "aggregate_type": "booking",
  "aggregate_id": "opaque-booking-id",
  "aggregate_version": 7,
  "actor": {
    "type": "account",
    "id": "opaque-account-id"
  },
  "correlation_id": "opaque-journey-id",
  "causation_id": "opaque-command-or-event-id",
  "idempotency_id": "opaque-logical-action-id",
  "data_classification": "confidential",
  "payload": {}
}
```

### Envelope rules

- `event_id` is globally unique/non-reused.
- `event_type` has an owner and stable semantic definition.
- `event_version` versions the payload/meaning, not the aggregate.
- `occurred_at` is the authoritative domain commit/effective time; `recorded_at` captures persistence/ingress difference where relevant.
- `tenant_id` is required for tenant-owned events; public/platform events use an explicit context, never null by accident.
- aggregate ID/version supports deduplication and per-aggregate ordering checks.
- actor is the accountable human/service/provider principal; background follow-up preserves original causation and adds executing service identity in processing metadata.
- correlation follows the business journey; causation points to the direct trigger; idempotency points to the logical effect.
- payload is minimum necessary and schema-defined. Secrets, OTPs, usable tokens and raw content are prohibited.
- trace context may be transport metadata rather than retained domain payload.

## Event types and ownership

Naming convention:

```text
<bounded_context>.<past_tense_fact>
```

Examples:

```text
identity.contact_verified
organization.membership_revoked
organization.operation_paused
catalog.service_version_published
scheduling.hold_placed
scheduling.hold_expired
scheduling.capacity_committed
booking.submitted
booking.confirmed
booking.rescheduled
booking.cancelled
booking.completed
request.submitted
quotation.issued
quotation.accepted
request.converted_to_booking
conversation.message_received
conversation.human_takeover_started
payment.payment_reported
payment.payment_verified
payment.refund_verified
payment.payment_reversed
communications.notification_delivery_failed
integration.connection_reauthorization_required
privacy.deletion_completed
```

Avoid vague types such as `booking.updated`, `status_changed`, `data_changed` or one enormous `entity_event`. A meaningful event states the business fact consumers are allowed to rely on.

## Domain versus integration event

The producing module may have richer internal domain events than it publishes.

Example:

- Internal Booking domain event contains policy-decision details used to update its own history.
- Published `booking.confirmed.v1` contains booking ID/version, operation/tenant, schedule snapshot reference, customer communication reference and minimum next-action data.
- Notification/Analytics do not receive full customer, price, policy or notes automatically.

The producer owns and documents the integration event. A consumer cannot demand a database-shaped payload that transfers ownership or unnecessary personal data.

## Consumer inbox and idempotency

Each consumer records:

- consumer name/version;
- event/work ID and producer;
- received/first-seen time;
- payload schema/version/hash where necessary;
- processing status/attempt/lease;
- resulting command/side-effect identity;
- terminal result or safe error;
- retention/replay checkpoint.

Processing rule:

```text
begin transaction
  claim/check (consumer, event_id)
  if already terminal: return stored outcome
  validate schema/tenant/authorization of service action
  apply one logical consumer effect
  record consumer outcome and any new outbox work
commit
```

For an external side effect, use provider idempotency/reference if supported and preserve ambiguous `submitted/unknown` outcomes for reconciliation. Do not mark consumer success before the external outcome required by that consumer is evidenced.

## Ordering

Ubi assumes no global total ordering.

- Aggregate version orders committed changes for one aggregate.
- Causation links order a specific workflow.
- Provider source order is retained but not trusted as complete.
- Workers may process different tenants/aggregates concurrently.
- Consumers reject/defer impossible version gaps where they require strict projection order, or rebuild from an owner query/change feed.
- An older event arriving late must not overwrite a newer projection/state.
- Timestamp alone cannot resolve domain conflict because device/provider clocks differ.

When a later broker is used, partition key SHOULD be immutable tenant + aggregate/connection key for consumers requiring local order. This is not a promise of cross-aggregate ordering.

## Event schema governance

Each published event definition records:

- owner/context and purpose;
- producers and authorized consumers;
- type/version and lifecycle;
- schema and examples;
- tenant/aggregate/order semantics;
- field classification, source and retention;
- idempotency/replay behaviour;
- compatibility promise;
- expected volume/burst and expiry;
- monitoring and dead-letter owner.

### Compatibility

- Additive optional fields are normally compatible.
- Existing field meaning/type/requiredness cannot change silently.
- Enum consumers handle unknown future values where the set can expand.
- Removing/renaming fields or changing meaning produces a new event version.
- Producer may dual-publish versions during migration with one logical event/causation relationship.
- Consumer upgrade, replay/backfill and removal evidence are required before retiring a version.
- Schema validation occurs at producer tests and consumer ingress; unknown versions quarantine/alert rather than guess.

## Replay and backfill

Replay is privileged operational work, not a normal retry button.

Rules:

- declare source range, tenant/context, consumer, reason, owner and approval;
- use original event IDs for deterministic consumer dedup or explicit replay IDs linked to originals according to purpose;
- default to dry-run/count/impact preview;
- external side effects (SMS, refund, message) are suppressed unless separately and explicitly authorized;
- analytics/projection rebuild uses a replay-safe sink/dataset/version;
- replay respects current privacy deletion/restriction and legal hold; deleted payload is not magically reconstructed;
- rate limits/backpressure protect live traffic;
- start/end/result/errors are audited.

Ubi is not initially event-sourced. Current operational tables plus explicit history are authoritative; the event stream alone is not promised to rebuild every aggregate forever.

## Retry classification

| Class | Examples | Behaviour |
|---|---|---|
| Transient Ubi | connection exhaustion, short lock timeout, temporary dependency | bounded exponential backoff with jitter |
| Transient provider | timeout, 429, temporary 5xx | provider-specific backoff/`Retry-After`, expiry and circuit policy |
| Unknown external outcome | timeout after submit, lost callback | do not repeat blindly; query/reconcile using same reference |
| Permanent input/schema | invalid payload, unsupported event version | quarantine/dead-letter and alert owner; no automatic retry |
| Authorization/configuration | revoked membership/connection, missing template | block/action-required; retry after controlled correction |
| Domain conflict | capacity taken, state/version changed, hold expired | terminal conflict requiring owning workflow/user resolution |
| Abuse/security | invalid signature, replay outside window, malicious media | reject/quarantine/rate/alert according to security policy |
| Expired/stale | late reminder, expired link/action | suppress/terminal with evidence; never send stale work |

Maximum attempts alone is insufficient. Work stops by error class, age/deadline, side-effect risk and operator policy.

## Dead-letter and blocked-work operations

Every durable queue has:

- owner and severity rules;
- age/depth/oldest-item/attempt/throughput metrics;
- blocked and dead-letter query with redacted context;
- safe payload/schema inspection under access policy;
- documented correct/retry/skip/compensate/escalate commands;
- batch-action guardrails and impact preview;
- retention and evidence;
- runbook and on-call/escalation.

“Dead letter” is not successful completion. Critical booking/payment/privacy/security work remains visible until resolved or explicitly risk-accepted.

## Scheduled work and timers

Timers cover holds, quotes, links, invites, reminders, payment deadlines, retry windows, retention and process escalation.

Timer record includes:

- owner context, tenant/resource/action;
- due time/timezone calculation evidence;
- expected aggregate/policy version;
- idempotency identity;
- state/lease/attempt;
- expiry/latest useful execution time;
- cancellation/supersession reference.

At execution the owner rechecks current authoritative state/version. A late timer is safe:

- consumed hold does not expire;
- rescheduled/cancelled booking does not receive the old reminder;
- accepted quote does not become expired;
- revoked member/link cannot be reactivated;
- retention job respects a newly created legal hold.

Scheduler availability is monitored by due-time lag, not only process uptime.

## Process-manager pattern

A process manager coordinates long-running multi-context work. It owns:

- process ID/type/version and tenant;
- initiating command/event/actor/correlation;
- current step and expected events/results;
- linked resource IDs/versions;
- deadlines/timers/attempts;
- external references/idempotency;
- compensation/reconciliation/manual action;
- terminal completed/failed/cancelled/action-required outcome.

It calls public commands on owner contexts. It never edits their tables or invents a global distributed transaction.

### Process state direction

```text
created
→ running
↔ waiting_for_event
↔ waiting_for_time
↔ waiting_for_human
→ completed

running/waiting_* → compensating → completed_with_compensation
running/waiting_* → failed
running/waiting_* → cancelled
running/waiting_* → action_required
```

Business-specific states may be richer; these are operational process summaries.

## Workflow catalogue

### WF-01 — Immediate/manual booking

Consistency: synchronous local Ubi transaction.

```text
SubmitBooking command
→ authorize + validate catalog/customer/policy
→ Scheduling commits capacity
→ Booking confirms/requests and snapshots terms
→ audit + outbox commit atomically
→ client receives result
→ Communications/Analytics projections asynchronously
```

No process manager is required unless external payment/approval enters the path.

### WF-02 — Payment-conditioned confirmation

Coordinator: Payment/Booking process manager (exact ownership decision pending; recommended application workflow owned by Payment condition policy).

```text
Booking/request + hold + obligation created
→ payment intent initiated
→ wait for authenticated/reconciled provider/manual evidence
→ Payment verifies allocation
→ process invokes ConfirmBooking using same condition/hold/version
→ success: booking confirmed, hold consumed, receipts scheduled
→ conflict/expired: paid-without-booking action-required case
→ authorized rebooking or refund workflow
```

The external payment and booking database cannot be atomic. The exception state is designed, monitored and staffed.

### WF-03 — Custom request to booking

Coordinator: Request context.

```text
request submitted/triaged
→ clarification/assessment
→ quotation version issued
→ customer accepts exact option/version
→ optional payment/hold condition
→ Request invokes idempotent Booking conversion
→ Booking + Scheduling hard transaction
→ Request records conversion/result
```

If conversion conflicts, accepted quote remains immutable and process waits for authorized revision/alternative/refund—not false confirmation.

### WF-04 — Notification/reminder

Coordinator: Communications.

```text
authoritative event
→ notification intent idempotently created
→ template/recipient/channel resolved
→ scheduled or ready
→ provider submission
→ delivery event/reconciliation
→ completed/failed action-required
```

Booking/payment never waits for delivery. Current state is rechecked before stale reminders/actions.

### WF-05 — Connected provider message

Coordinator: Integration ingress then Conversation.

```text
authenticated provider callback
→ receipt/dedup commit and prompt acknowledgement
→ canonical normalization
→ conversation/message idempotent append
→ routing/automation/human policy
→ optional structured draft/handoff
→ separate Request/Booking command only after authorized review/submit
```

### WF-06 — Provider refund

Coordinator: Payment.

```text
refund requested/approved
→ stable refund intent/provider key
→ submit once / wait or reconcile unknown
→ authenticated provider result
→ exact allocation/refund update
→ customer/business notification
→ later settlement reconciliation
```

Timeout after submission never creates a new refund blindly.

### WF-07 — Organization closure

Coordinator: Organization/Operations case.

```text
closure requested + strong authentication/approval
→ restrict unsafe new activity
→ inventory future bookings/holds/requests/payments/refunds/connections/exports
→ merchant review and customer fulfilment/cancellation plan
→ revoke/disconnect credentials at correct time
→ export/financial/privacy/retention tasks
→ close when blocking obligations resolved
→ deletion eligibility later under records policy
```

Closure is not a single `DELETE organization` transaction.

### WF-08 — Privacy deletion/export/restriction

Coordinator: Privacy & Records.

```text
verified request/case + responsibility/jurisdiction decision
→ plan owner systems/fields/exceptions/legal holds
→ idempotent tasks to operational/object/search/cache/analytics/providers
→ gather per-task completion/failure evidence
→ quality/identity/security review where required
→ truthful completed/partially completed/refused outcome
```

### WF-09 — Membership revocation with offline clients

Coordinator: Organization; Sync consumes current authority.

```text
membership revoked
→ authority version/session policy updated + audit/outbox
→ optional push/security notice
→ offline device later submits queued command
→ Sync reauthorizes against current membership
→ domain command rejected without effect
→ Android preserves user input and shows access-ended recovery/export path
```

## Compensation and reconciliation

Compensation is a new explicit business action, not rollback magic.

Examples:

- payment verified but booking unavailable → offer alternative/rebook or refund;
- notification failed → staff follow-up/fallback, not unconfirm booking;
- external message submitted with unknown status → provider query/manual decision, not automatic duplicate;
- media upload accepted but scan failed → quarantine/remove attachment while preserving text record;
- provider connection revoked → pause external work and direct users to Ubi/manual path;
- privacy propagation partially failed → retry/escalate and report partial state, not claim deletion.

Every compensation has authority, idempotency, audit, customer/business communication and financial/data consequence.

## Provider callback ingress workflow

1. Edge applies provider-route size/rate controls while preserving signed raw input.
2. Adapter validates signature/token/certificate/time/replay as the provider specifies.
3. Resolve live environment, provider account/connection and Ubi tenant from protected configuration.
4. Insert provider receipt/deduplication record with safe bounded raw evidence or hash/reference.
5. Return provider-required acknowledgement after the durable acceptance point.
6. Worker parses using the pinned adapter version and emits a canonical command/evidence object.
7. Owning domain validates current state/idempotency and commits result/outbox.
8. Unknown/invalid events are quarantined/alerted, never guessed.
9. Missing/out-of-order/ambiguous outcomes are reconciled by provider query/export where supported.

If a provider requires processing before acknowledgement, the adapter must still bound latency and document duplicate/recovery behaviour.

## Backpressure and fairness

- Separate queues/pools for core domain follow-up, provider ingress, communications, media, analytics and bulk privacy/export work.
- Prioritize by user harm/deadline, not merchant payment tier alone.
- Per-tenant/provider destination concurrency protects global capacity.
- Queue age/throughput/autoscaling limits are measured; autoscaling cannot overload the database/provider.
- Large exports/media/analytics yield to booking/payment/security work.
- When accepting more durable work would violate recovery capacity, reject before acceptance with `429/503` and safe retry rather than accept and lose it.
- Load shedding may drop only explicitly noncritical sampled analytics, never authoritative/audit/payment/accepted booking work.

## Security and authorization

- Worker identities are capability- and environment-scoped.
- Event possession is not authorization to access arbitrary current records; consumers have approved purpose/data scope.
- A delayed command is reauthorized when executed if it acts on behalf of a human/current membership.
- A domain event remains historical fact after membership revocation, but follow-up disclosure/action follows current policy.
- Provider secrets/signatures/tokens are not event payloads.
- Queue/replay/admin controls require strong workforce identity, case/reason and audit.
- Deserialization uses strict schemas/size/depth/type limits; events are untrusted at each boundary.
- Cross-tenant event/resource mismatches reject and alert as potential security incidents.

## Privacy and retention

- Prefer IDs/references over copying names, phones, messages, addresses or full booking/payment objects.
- Consumers fetch only authorized current/minimum data at processing time when feasible.
- Event field inventory includes purpose/classification/retention and deletion effect.
- Outbox/inbox/dead-letter payloads are production data with access/retention—not harmless logs.
- Analytics receives its own minimized event, not the richest domain/integration event.
- Deletion may redact/remove event personal payload while retaining minimum integrity/idempotency/audit tombstone according to approved policy.
- Replay cannot recreate data that should remain deleted.

## Observability

Per queue/consumer/process measure:

- depth, oldest age and due-time lag;
- enqueue/deliver/process rate;
- attempts/retry class and processing duration;
- blocked/dead-letter count/value/user impact;
- idempotency duplicate/conflict;
- provider latency/rate/unknown outcome;
- process step/time waiting/manual backlog;
- tenant fairness/saturation;
- schema/version rejection;
- end-to-end correlation from originating command to projection/communication.

Metrics use bounded event/consumer/error labels. Full payload, personal IDs and raw provider errors do not enter metric labels.

## Recovery and deployment

- Workers can be stopped/redeployed without losing committed outbox work.
- New worker versions support old/new event schemas through the deployment/replay window.
- Database migrations use expand/migrate/contract order for queues/process state.
- A rollback cannot make newly emitted events uninterpretable; dual readers/writers or roll-forward may be required.
- Restore tests include outbox/inbox leases, idempotency, timers, process managers and provider reconciliation—not only domain tables.
- After restore, leases expire safely and duplicates remain one logical effect.
- Provider callbacks during outage are recovered from provider retries/query/export or explicit reconciliation limitations.

## Event/workflow acceptance tests

1. Kill the application after domain commit but before dispatch; event later delivers once logically.
2. Kill worker after provider success but before marking success; retry reconciles/reuses reference without duplicate external effect.
3. Deliver an event 100 times and out of order; consumer state/effects remain correct.
4. Process two aggregates concurrently while preserving each aggregate’s version order where required.
5. Pause notification/analytics workers for hours; booking stays operational and queues recover within objective.
6. Expire a hold/reminder timer late after its state changed; no stale expiry/send occurs.
7. Restore database to an approved point; outbox/inbox/process workflows recover and reconcile external events after the restore point.
8. Replay a projection range with external side effects disabled; live provider/customer messages are not repeated.
9. Revoke membership after offline command creation; delayed execution rejects without domain mutation.
10. Create paid-without-booking, unknown-refund and partial-privacy-deletion cases; each remains visible with owner/deadline/resolution.

## Open decisions

- Database outbox polling/claim implementation and capacity thresholds for broker adoption.
- Event schema format/tooling/registry and compatibility enforcement.
- Exact process-manager owner for payment-conditioned confirmation.
- Per-queue priority, attempt, expiry, dead-letter and retention policy.
- Provider-specific acknowledgement and reconciliation designs.
- Audit event durability/protection relationship to source transaction.
- Cross-region/outage recovery and event replay source.
- Analytics durability/sampling boundary and warehouse timing.

