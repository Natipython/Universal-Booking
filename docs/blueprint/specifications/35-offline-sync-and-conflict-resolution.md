# Android offline synchronization and conflict resolution

Status: **Draft**  
Depends on: Android requirements, API governance, events/workflows, authority and domain invariants  
Blocks: Android architecture, local/server schemas, sync endpoints, user experience, tests and support tooling

## Purpose

Ubi Business Android must remain useful on unreliable or absent connectivity while never pretending that one phone has complete knowledge of other devices, customers or provider events. This document defines local truth, queued commands, server convergence and human conflict recovery.

## Meaning of local-first

For Ubi:

- the Android UI observes and reads a durable local relational store;
- valid cached data appears immediately without waiting for a network attempt;
- eligible work is saved locally before the app says it is saved;
- synchronization runs when allowed/possible without forcing the user to reopen each item;
- the server remains authoritative for tenant authority, capacity, price/policy, booking confirmation, payment verification and cross-device order;
- the app distinguishes local/pending, synchronizing, synchronized, rejected and conflict states in text/semantics;
- a network icon or spinner is never the only explanation.

Local-first does not mean “everything works offline” or “last device write wins.”

## Truth model

| Truth category | Local Android authority | Server authority |
|---|---|---|
| Cached schedule/catalog/customer projection | Best-known snapshot with `as_of`/version | Current authoritative state |
| Unsaved form draft | Authoritative only on that app install | None until command submission |
| Durable queued command | Authoritative evidence that the device saved the user’s intent | Final command acceptance/result after sync |
| Local tentative booking | Occupies a local shadow view to avoid same-device reuse | Not a confirmed/capacity commitment |
| Synced booking/request | Read projection of server record/version | Authoritative lifecycle/snapshot |
| Staff permission | Cached hint for navigation/offline eligibility | Current permission at execution |
| Availability | Cached candidate plus device-local pending work | Final capacity/hold/commitment |
| Payment report/receipt draft | Local unverified claim | Ubi accepted/verified/reconciled state |
| Provider/message delivery | Cached evidence | Ubi/provider current evidenced result |

User wording must follow the strongest actual evidence, not the optimistic intention.

## Sync components

### Android

- **Local relational store** — source observed by the UI.
- **Draft repository** — recoverable form/input state not yet queued.
- **Command outbox** — durable immutable command intent/envelope and dependencies.
- **Local projection updater** — applies optimistic/pending changes and server changes transactionally.
- **Sync coordinator** — schedules one logical sync per account/tenant scope, batches work and applies results.
- **Conflict repository** — preserves rejected/conflicting input plus server context/resolution options.
- **Secure session storage** — minimal credentials/keys separate from ordinary database fields.
- **OS work adapter** — persistent constrained background work, foreground refresh and network/battery policy.

### Server

- **Sync session/manifest endpoint** — supported protocol/client/scope and bootstrap negotiation.
- **Command batch endpoint** — authenticates and routes typed domain commands.
- **Idempotency/command-result store** — one logical effect and recoverable result.
- **Change-feed/projection service** — scoped incremental server changes and tombstones.
- **Cursor/retention service** — opaque resumable position and full-resync requirement.
- **Conflict serializer** — stable machine code, safe current context and allowed resolution commands.
- **Sync diagnostics** — privacy-safe lag/outcome/support evidence.

Sync is an envelope/orchestration concern. Booking, Scheduling, Customer, Catalog, Request and Organization modules still decide their own commands.

## Local data partition and scope

Every local row/command/projection carries:

- authenticated account ID;
- organization tenant ID;
- branch/operation/resource scope as appropriate;
- server/client resource ID;
- projection/aggregate/schema version;
- synchronization state and timestamps/evidence;
- classification/retention/cache policy.

Rules:

- Switching organization visibly switches the active local scope; data cannot appear merely because another membership exists.
- Cross-tenant joins/queries are prohibited by repository design and tested.
- One install may cache several authorized organizations, but each has independent cursor/commands/quota/wipe handling.
- Membership revocation removes active access after server knowledge; queued commands reject and protected cached data becomes inaccessible/deleted under policy.
- Only data needed for the user’s authorized work and offline horizon is synchronized.
- Customer/message/payment detail is minimized more aggressively than catalog/schedule metadata.

### Proposed offline dataset

Pilot candidate:

- active account/membership/role display and feature policy projection;
- selected organization/branch/operation identifiers and safe configuration;
- services needed to create manual bookings;
- assigned/scoped resources and bounded schedule horizon;
- future/recent bookings required for daily work;
- limited tenant customer lookup/records needed for manual booking;
- user drafts, queued commands, command results/conflicts;
- bounded notification/action queue;
- relevant tombstones/version/cursor data.

Not synchronized by default:

- organization-wide exports or unrestricted customer list;
- full conversation history/media;
- provider credentials or raw webhook payloads;
- settlement/bank details;
- Ubi operations/audit records;
- old records outside the approved working horizon.

Exact past/future horizons, customer limits and tenant quotas require pilot workflow/storage measurement.

## Local database lifecycle

### Creation/unlock

- Created for a verified app/account security context.
- Protected using platform-supported application storage and approved encryption/key management.
- Key material is bound to the app/security context where feasible and never hard-coded.
- Device/app restore cannot silently resurrect valid sessions or unreadable key material.

### Logout/session loss

Before logout, the app shows unsynchronized work and the consequence. Ordinary logout:

- revokes/clears usable session material as policy requires;
- locks or deletes protected local data according to account/device policy;
- never uploads/discards commands under another account;
- offers authorized sync/recovery before deletion when safe;
- cannot retain a readable database accessible from a new account.

Remote session revocation/organization removal invalidates access at next authoritative contact. High-risk local cache may require shorter offline access/reauthentication limits.

### Schema migration

- Local migrations are forward tested on supported prior app versions/data sizes.
- A failed migration backs up/rebuilds only rebuildable projections while preserving drafts/outbox commands in a compatible protected form.
- Destructive fallback cannot erase pending work silently.
- Command payload versions remain supported through the stated offline/client compatibility window.

## Synchronization phases

```text
0. Preflight/session and scope
1. Push eligible pending commands
2. Apply command results and returned authoritative projections
3. Pull incremental server changes/tombstones
4. Apply changes transactionally
5. Repeat while bounded work remains
6. Record cursor/health and schedule next need
```

The server may require pull-before-push for specific stale security/configuration states. The default prioritizes user-created pending commands without assuming they will succeed.

## Bootstrap protocol

Used for first login/scope, local reset, filter change or expired/invalid cursor.

1. App authenticates and selects an authorized organization/operation scope.
2. App sends protocol/client/schema capabilities and existing pending-command summary.
3. Server returns a sync manifest: authorized scopes, projection sets, server time, minimum client/protocol, paging/batch limits and bootstrap token/cursor.
4. App downloads bounded snapshot pages with checksums/version/order.
5. Pages enter staging tables/partition; partially downloaded snapshot does not replace the active usable dataset.
6. App validates completeness and transactionally activates snapshot plus cursor.
7. Local pending commands/drafts are re-applied as overlay or marked for rebase; they are never overwritten blindly.
8. Incremental synchronization begins from the issued cursor.

If bootstrap fails, the last valid dataset remains visible with `as of` status where safe.

## Incremental pull/change feed

Request inputs:

- opaque cursor;
- authorized tenant/scope/filter version resolved by server;
- supported projection schema versions;
- bounded page limit/byte capability.

Response contains:

- ordered page of projection upserts and tombstones;
- next cursor;
- `has_more`;
- server time;
- minimum client/protocol/action required;
- optional retry/backpressure guidance.

Rules:

- Cursor is opaque and scoped; client cannot modify it to query another tenant/scope.
- Cursor advances only in the same local transaction that applies the complete page.
- Reapplying a page is idempotent.
- Server projection/version prevents older change overwriting newer local synchronized state.
- Local pending overlay remains visible and linked while base server projection changes.
- Cursor has a documented retention window; outside it server returns `full_resync_required`.
- Feed may be operation/schedule-horizon scoped. Scope/filter change creates a new manifest/cursor.
- A missing version gap triggers retry/rebootstrap rather than guessing order.

## Command envelope

Proposed logical command:

```json
{
  "client_command_id": "opaque-stable-id",
  "command_type": "booking.submit_manual",
  "command_version": 1,
  "client_resource_id": "opaque-client-created-id",
  "tenant_id": "opaque-organization-id",
  "operation_id": "opaque-operation-id",
  "created_at_device": "2026-08-30T10:15:00+03:00",
  "base_resource_version": 12,
  "base_policy_versions": {
    "service": 5,
    "booking_policy": 3
  },
  "dependencies": [],
  "payload": {},
  "client_context": {
    "app_install_id": "opaque-install-id",
    "app_version": "version",
    "timezone": "Africa/Addis_Ababa"
  }
}
```

Server adds authenticated account/session/membership/assurance/IP/received time; client cannot assert those.

### Command rules

- Envelope and payload are immutable after queued. User correction cancels/supersedes the old unsent command where legal and creates a new command.
- `client_command_id` is stable across retry/restart and maps to server idempotency.
- `client_resource_id` allows offline references/dependencies; accepted server may adopt it as the canonical opaque ID or return a permanent mapping.
- Device time records user context only; it never proves expiry/order/authorization.
- Base versions identify what the user saw.
- Dependencies reference earlier client commands, not assumed upload order.
- Payload is typed/allow-listed by command version; no generic entity/field patch.
- The server evaluates each command under current account/membership/resource/domain state.

## Local command lifecycle

```text
draft
→ queued
→ in_flight
→ synchronized

in_flight → retry_wait → in_flight
queued/retry_wait → cancelled_before_send
in_flight → conflict
in_flight → rejected
in_flight → dependency_failed
in_flight → unknown_result → result_recovery → synchronized/conflict/rejected
```

### UI meanings

| Internal state | User-facing direction |
|---|---|
| draft | “Draft on this phone” |
| queued | “Saved on this phone. Waiting for connection.” |
| in_flight/retry_wait | “Sending to Ubi…” / “We’ll try again.” |
| synchronized | “Synchronized” plus actual booking/request/payment state |
| conflict | “Needs your attention” with specific reason/options |
| rejected | Specific safe reason such as permission ended/invalid value |
| dependency_failed | “Couldn’t send because an earlier change needs attention.” |
| unknown_result | “Checking whether Ubi received this—don’t create it again.” |

Green success is used only for an actual synchronized/confirmed outcome appropriate to the domain. Local `queued` is informational/amber/neutral.

## Server batch processing

1. Authenticate batch/session and check client/protocol compatibility.
2. Enforce one tenant per batch, size/count/time/rate limits.
3. For each command in dependency order:
   - validate envelope/schema;
   - resolve current principal/membership/tenant;
   - claim/recover idempotency record;
   - validate dependency results and client-resource mapping;
   - call owning application/domain command;
   - commit result, audit and outbox;
   - return typed result and authoritative projection/version.
4. Independent commands may complete separately. An explicitly documented atomic group is used only when the domain supports it.
5. HTTP retry of the batch replays per-command stored outcomes.

Per-command result:

```text
accepted/committed
duplicate_with_result
conflict
rejected
retry_later
dependency_failed
unsupported_command_version
```

An HTTP `200` for a valid batch envelope does not mean every command succeeded.

## Offline capability matrix

Legend:

- **R** — read cached/offline.
- **Q** — action may be durably queued but is visibly pending until server result.
- **D** — draft locally only; cannot be queued as an authoritative action yet.
- **O** — online authoritative interaction required before the action starts/completes.

| Action | Initial direction | Reason/behaviour |
|---|:---:|---|
| View Today/assigned schedule | R | Core offline promise with `as of` and pending overlay |
| View bounded services/resources/customers | R | Minimum data needed for manual booking |
| Create call/walk-in standard booking | Q | Saves local intent; server may confirm/request/conflict |
| Create minimal customer with booking | Q | Tenant-scoped client ID; duplicate suggestion/merge later |
| Add ordinary operational note/draft answers | Q/D | Typed visibility; sensitive/prohibited data controls still apply |
| Check in/start/complete/no-show assigned booking | Q | Server reauthorizes/state-checks; UI remains pending |
| Record customer-reported payment/receipt metadata | Q | Always local/unverified until server/provider review |
| Attach media/receipt | Q separately | Text command can proceed; upload waits and may fail/quarantine |
| Reschedule/cancel customer commitment | O initially | Capacity, policy, money and communication consequences; queued expansion only after tests |
| Accept/confirm booking request | O initially | Must establish current capacity/authority and customer-visible truth |
| Place/extend hold | O | Time/capacity authority cannot be local |
| Verify payment or execute refund | O | Provider/concurrency/financial authority |
| Publish price/service/policy/availability | D then O | Draft may be local; publication/impact review online |
| Invite/change staff/permissions/ownership | O | High-risk current authority/notification |
| Connect/change provider/integration | O | Official authorization and credential state |
| Export/delete/close organization | O | High-risk multi-system process |
| Native message send (later) | Q | Local queued not sent; block/relationship rechecked online |
| Ubi call | O | Real-time relationship/permission/network required |

This matrix can expand only after conflict/support evidence. The interface should prevent unsupported work before extensive input when possible.

## Local tentative capacity overlay

When a manual booking is queued:

- Android creates a local booking/client ID and pending schedule entry.
- The selected resource/time is marked tentatively occupied on that same local tenant dataset.
- Subsequent same-device choices warn/block according to resource rules.
- Pending overlay cannot know another device/customer has taken the time.
- Staff sees `Saved on this phone—not confirmed by Ubi` where confirmation matters.
- Server result replaces overlay with confirmed/requested/conflict/rejected projection without changing the client-visible reference unexpectedly.

If business policy makes manual staff bookings request-only, the label reflects that; offline cannot upgrade it to confirmed.

## Conflict taxonomy

| Code family | Example | Default resolution |
|---|---|---|
| `authority.*` | membership revoked, scope removed, step-up required | Stop retries; preserve input; sign in/ask owner/export permitted draft |
| `tenant.*` | organization/operation suspended/closed | Stop affected actions; show operational/support path |
| `capacity.*` | slot/resource taken, hold expired, capacity reduced | Show exact affected schedule and safe alternatives/manager override if authorized |
| `resource.*` | service/customer/resource archived/deleted/merged | Rebind/select current object; retain entered details/provenance |
| `version.*` | booking/catalog/policy changed | Show what changed and review/reapply selected fields |
| `price_policy.*` | amount/deposit/cancellation terms changed | Require explicit review; no automatic acceptance |
| `validation.*` | field now invalid/required | Focus field, preserve valid values and explain correction |
| `state.*` | already cancelled/completed/confirmed | Show current state; offer only legal next actions |
| `duplicate.*` | same booking/customer/action already exists | Link to existing candidate; authorized merge/keep-separate decision |
| `time.*` | action/quote/token expired or device date ambiguity | Recompute/reissue/review using server time |
| `payment.*` | online/provider verification required or state changed | Move to online/reconciliation flow; never mark verified locally |
| `dependency.*` | earlier offline-created customer/booking failed | Resolve parent first or rebind/cancel dependent command |
| `protocol.*` | command/client version unsupported | Safe app upgrade/migration/manual recovery; do not discard payload |

Conflict codes are context-qualified, documented API contracts and localized in the client.

## Conflict resolution rules

### Never silently auto-merge

- exclusive time/resource/capacity;
- confirmed/requested booking lifecycle;
- price, deposit, cancellation or customer-visible terms;
- customer identity/link/merge;
- permission/ownership/approval;
- payment/refund/reconciliation;
- external message recipient/content after takeover/block;
- deletion/retention/legal hold;
- provider/integration configuration.

### Potentially safe automatic convergence

Only owner-defined operations that are commutative/idempotent may merge, for example:

- marking a notification read when already read;
- adding a locally unique, authorized append-only operational tag/note where duplicate identity prevents repeat and policy allows;
- replacing an unchanged draft field when server base version still matches;
- applying a newer server projection/tombstone under version rules.

Each rule is explicit and tested. There is no platform-wide last-write-wins.

### Resolution screen contract

The user receives:

- what they attempted, with entered data preserved;
- what changed/current server state in ordinary language;
- whether anything already succeeded;
- exact affected business/customer/time/amount;
- recommended safe action first;
- alternatives such as select new time, update terms, link existing customer, ask manager, cancel local attempt or contact support;
- authority/connectivity requirement before choosing;
- resulting customer communication/payment consequence;
- no raw stack/HTTP/database error.

Resolution creates a new typed command referencing the conflict/original command. It never edits an already processed envelope.

## Multi-command dependencies

Offline example:

```text
C1 create customer record
C2 submit booking using C1 client_resource_id
C3 add customer-reported payment to C2 booking
```

Rules:

- C2 waits for C1 accepted/mapped result.
- C3 waits for C2 accepted result and uses its server mapping.
- C1 conflict may be resolved by linking an existing customer; mapping lets C2 proceed without re-entry.
- C2 capacity conflict makes C3 `dependency_failed`; payment evidence is preserved but not attached/verified incorrectly.
- Canceling an unsent parent previews/removes or rebinds dependents.
- Cycles and unbounded chains are rejected locally/server-side.

## Multi-device behaviour

- Each app install has an opaque install ID and independent local command sequence/cursor.
- Server aggregate versions/idempotency remain authoritative across devices.
- A device sees other-device changes through incremental pull/push notification hint.
- Push is only a hint to synchronize; it does not carry authoritative full state.
- Two devices can create conflicting offline bookings; both remain local until the server accepts one and returns conflict to the other.
- Semantic duplicate detection may suggest existing booking/customer but cannot merge/delete silently.
- Logging out/revoking one device does not remove another valid session unless policy says all sessions.
- Device-specific diagnostics do not expose another user’s content.

## Authorization expiry and long-offline use

Offline access policy considers:

- session/account risk and last successful online authentication;
- organization/member status and role risk;
- local data classification;
- device lock/security posture signals where approved;
- maximum offline duration by action/read dataset;
- need for emergency daily schedule access versus data exposure.

The launch policy may permit cached schedule viewing for a bounded period while blocking new queued commands after stronger expiry. Exact durations require threat/user research.

On reconnect every command uses current authority. Cached role/capability never overrides revocation.

## Time and clock handling

- Server time establishes received/commit/order/expiry.
- Device time is retained as contextual evidence and marked untrusted.
- Sync preflight returns server time; app may display significant clock mismatch warning.
- Device clock cannot extend a hold/quote/link or win a conflict.
- Locally scheduled views use stored instant/date/timezone/calendar rules, not current device timezone silently.
- If the phone changes timezone/calendar, existing commitment presentation remains tied to booking/branch context.

## Deletions and tombstones

A server tombstone includes minimum:

- resource/projection type and ID;
- tenant/scope;
- deletion/archive version/time;
- reason class needed for local behaviour;
- retention/visibility boundary.

Rules:

- Tombstone removes/hides the local projection and clears dependent sensitive cache.
- It remains long enough to reach supported offline clients or forces full resync after expiry.
- Local pending command against a tombstoned resource becomes explicit conflict/rebind/reject.
- Tombstone does not retain deleted personal payload unnecessarily.
- Archive, revoke, merge, restrict and legal delete are distinct outcomes where clients need different behaviour.

## Full resynchronization

Triggers:

- cursor expired/unknown/corrupt;
- scope/filter/membership materially changed;
- local schema/data integrity failure;
- server projection contract incompatible;
- security-initiated reset;
- mismatch/checksum/version gap that cannot reconcile incrementally.

Procedure:

1. Freeze new sync execution briefly; drafts may continue if safe.
2. Preserve and validate local drafts/pending commands/conflicts separately.
3. Obtain new manifest and staged snapshot.
4. Activate snapshot atomically.
5. Rebase/reapply pending overlay and send commands through current authorization.
6. Remove old protected projections/keys according to policy.
7. Record diagnostic outcome without personal payload.

“Clear app data and try again” is not an acceptable default when unsynchronized work exists.

## Background and foreground synchronization

### Triggers

- successful login/account/scope selection;
- app start/resume;
- eligible user pull-to-refresh/manual retry;
- new queued command;
- network becoming usable as an opportunistic hint;
- push notification hint;
- periodic OS-scheduled work within battery/data policy;
- high-priority conflict/security state.

### Rules

- Use OS persistent work scheduling; do not depend on a forever-running background service.
- Ensure one logical sync worker per account/tenant/scope; concurrent triggers coalesce.
- Connectivity signal does not guarantee internet/Ubi reachability.
- Backoff uses error class, jitter, max delay and user-visible manual retry.
- Foreground user action may accelerate but not create a second duplicate worker.
- Respect metered/data/battery policies for media and large bootstrap; critical small text commands take priority.
- Worker cancellation/process death leaves commands durable and leases recoverable.

Exact Android WorkManager/Room/DataStore choices follow the Android framework decision; the architecture requires equivalent guarantees.

## Network and payload efficiency

- Commands/changes are batched within measured size/time limits, not one network call per row.
- Compression is supported with decompression limits.
- Incremental projections omit unchanged/unused fields and media bytes.
- Media uploads resume separately and use checksums.
- Sync prioritizes user commands/results and today’s schedule before distant history/media.
- Pagination streams bootstrap to staging; do not hold entire tenant dataset in memory.
- Retry avoids re-downloading acknowledged pages using transactional cursor.
- Server can return backpressure/retry windows; Android follows without tight polling.
- Track bytes per sync/active org/booking and D1 battery/thermal cost.

## Local security and privacy

- Store minimum offline data; use explicit horizon/quota/eviction.
- Use platform-protected keys/secure storage; never store raw OTP/provider credentials.
- Sensitive values are redacted from Android logs/crash/analytics and notification payloads.
- Disable or securely control Android backup/restore for session-bound database material.
- Clipboard/share/screenshot policy is risk- and usability-reviewed; do not break accessibility casually.
- Media/receipts use private app storage, content URI/temporary access and prompt cleanup.
- Tenant/customer search caches and indexes follow the same encryption/deletion scope.
- Root/debug/tamper signals may trigger risk controls but cannot be the only authorization layer or silently destroy data.
- Remote wipe cannot be guaranteed while permanently offline; this limitation informs offline duration/data minimization and merchant policy.

## Local audit and provenance

Queued command captures:

- authenticated local account/session reference at creation;
- cached membership/operation scope/version;
- device/app-install/client command ID;
- device-created time/timezone and last known server offset;
- source such as phone/walk-in;
- base resource/config/policy versions;
- user-entered payload/provenance;
- local edits/supersession/cancellation;
- sync attempts and final server result.

The authoritative server audit records execution-time subject/current authorization plus original creation context. A local audit entry alone is not authoritative evidence that the server action occurred.

## Observability and support

Measure privacy-safe:

- last successful sync/pull/push by app version/network class;
- local pending age/count/bytes by command type;
- command accepted/duplicate/conflict/rejected/retry outcome;
- conflict type and resolution success/time;
- cursor expiry/full-resync frequency/duration/bytes;
- local migration/database integrity failure;
- worker duration/retry/cancellation/background restriction;
- server batch latency/size/rate/idempotency hit;
- projection lag/version gaps;
- D1 battery/data/performance impact.

Business UI shows a simple `All synchronized`, `Waiting for connection`, or `N items need attention` plus item-level state. Support can use an exported redacted diagnostic bundle only with user authorization; it excludes customer content/secrets by default.

## Testing matrix

### Deterministic domain/protocol tests

- command serialization/version compatibility;
- idempotency under concurrent replay;
- dependency graph/mapping;
- projection version/tombstone application;
- cursor atomicity/page replay;
- every conflict code/resolution;
- timezone/calendar/device-clock changes;
- revoked membership/operation closure;
- local quota/eviction preserving pending work.

### Device/process tests

- kill process after local transaction, during request, after server commit and during result application;
- reboot/device lock/key rotation/session expiration;
- database migration from each supported version with pending commands;
- organization/account switch and logout with pending work;
- network loss/flapping/captive portal/DNS/TLS/timeout/slow upload;
- low storage/memory/battery/background restriction;
- D1 performance and accessibility at every sync state.

### Multi-client/server tests

- two offline devices select same exclusive slot;
- customer web takes slot while staff device offline;
- catalog/price/policy changes while offline form exists;
- record deleted/merged/restricted while device offline;
- server deployment changes command/projection version;
- duplicate/out-of-order change pages/results;
- cursor expires after long offline period;
- restore server database while devices/provider events continue;
- 100 duplicate command submissions result in one logical effect.

### Usability tests

Target staff must be able to answer:

1. Is this booking only on the phone or in Ubi?
2. Is the customer actually confirmed?
3. Did payment get reported or verified?
4. Which item needs attention and why?
5. What information is saved if they go back/close the app?
6. What will happen when connection returns?
7. How do they choose another time without retyping everything?

If these answers are unclear, the sync implementation is not ready even when technically correct.

## Acceptance scenarios

1. On N0, staff opens the cached Today schedule and adds a walk-in booking; killing/restarting the app preserves it as local pending.
2. The same command is transmitted repeatedly after timeouts; one server booking/capacity effect results and the device recovers it.
3. Another device takes the slot first; original device receives capacity conflict, retains all customer/input data and offers valid alternatives.
4. Owner revokes receptionist before reconnect; queued booking is rejected with no server mutation and a respectful recovery path.
5. Customer record creation succeeds but booking conflicts; mapping/customer remains correct and payment-dependent command does not run.
6. Price/policy changed; device requires review and never submits old terms as silently accepted.
7. Cursor expires after a long offline period; staged full resync preserves/rebases pending commands and last valid data until activation.
8. Old tombstone removes a customer projection without retaining personal content; a dependent draft explains re-selection.
9. Receipt is recorded offline; every state says unverified and media can upload later without blocking text command.
10. Local DB migration fails on a device with queued work; app preserves/recovery-exports the pending commands rather than deleting them.
11. Network flaps and WorkManager triggers repeatedly; one logical sync drains with backoff and bounded battery/data.
12. Switching between two organizations never displays or submits data under the wrong tenant/cursor.

## Authoritative external baseline

- [Android offline-first architecture guidance](https://developer.android.com/topic/architecture/data-layer/offline-first) recommends immediate local reads, explicit online/queued/lazy write strategies, synchronization and persistent constrained work. Ubi adds stricter domain-specific authority and conflict rules for booking and money.
- [Android data-layer guidance](https://developer.android.com/topic/architecture/data-layer) supports a local data source as the source of truth for an offline-first client’s observed UI state.

## Open decisions

- Android framework and concrete local database/secure-storage/background-work libraries.
- Exact client-generated identifier format and whether server adopts or maps it.
- Sync change-feed mechanism, cursor retention and projection versioning.
- Offline schedule/customer horizon, data quota and maximum offline authorization duration.
- Initial queued-write matrix after pilot field observation.
- Local encryption/key/logout/backup policy by Android version/device capability.
- Batch byte/count/time limits and sync priority/fairness.
- Conflict/support diagnostic data and retention/notice.
- Whether push-to-sync hints are included in the first pilot.

