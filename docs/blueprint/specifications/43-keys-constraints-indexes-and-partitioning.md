# Keys, constraints, indexes and partitioning

Status: **Draft — enforceable patterns proposed; executable DDL awaits release scope and ADRs**  
Depends on: data domain model, operational database architecture and logical table catalogue  
Blocks: production schema, repository contracts, concurrency tests, RLS performance and migrations

## Purpose

This document defines the database invariants that remain true even when an application bug, retry, offline client, worker race or provider duplicate reaches PostgreSQL. It does not replace domain validation; it establishes the smaller set of facts that must never depend on a screen or one service process behaving correctly.

The enforcement order is:

1. type and `NOT NULL`;
2. `CHECK` constraints for row-local truth;
3. primary, unique and foreign keys for identity and relationships;
4. exclusion or locked-ledger rules for scarce capacity;
5. short transactions with deterministic locking for multi-row truth;
6. idempotency/inbox/outbox for retries and distributed effects;
7. application policy for rules that legitimately change by tenant, version or jurisdiction.

## Non-negotiable conventions

- All identifiers are lowercase `snake_case` and never require quoted case-sensitive SQL.
- Every table has a primary key; association tables may use a meaningful composite primary key.
- Required facts are `NOT NULL`; an empty string is not a substitute for unknown.
- Foreign keys are real constraints unless a documented lifecycle or external-boundary ADR explains otherwise.
- Every foreign-key access path needed for delete, join or RLS is indexed deliberately; PostgreSQL does not create child indexes automatically.
- Tenant-owned relationships carry `organization_id` on both sides and use composite tenant foreign keys.
- Money is exact, carries a currency code and never uses floating point.
- Server instants are `timestamptz`; local dates/times retain timezone/calendar context where required.
- State fields use constrained values and legal transition commands; clients cannot write arbitrary labels.
- User/provider JSON is bounded and never replaces an invariant-bearing column.
- Cascade deletion is exceptional. Closing a business, account or customer must not silently destroy commercial, payment or audit evidence.

## Identifier and primary-key pattern

Subject to the identifier ADR, opaque domain identity uses:

```sql
id uuid primary key
```

Offline-creatable aggregates use a pinned UUIDv7-compatible generator in supported clients/server code. The database accepts valid IDs, but authorization and ownership are always server-established. An ID conveys no permission.

Internal ordered streams may use:

```sql
position bigint generated always as identity primary key
```

Use database-assigned positions for outbox/change/audit streams only where ordering is useful. Do not expose them as secrets or tenant authority.

### Stable identity versus human reference

- `id` is immutable technical identity.
- A booking/reference code is separately generated, displayable and searchable.
- Slugs and usernames are mutable public locators with reservation/history rules.
- Provider IDs are namespaced by provider, environment and connection.
- Token secrets are never primary keys and are stored only as verifiers/hashes or external secret references.

## Tenant-key and foreign-key pattern

Every tenant-owned parent exposes:

```sql
unique (organization_id, id)
```

Every tenant-owned child proves same-tenant ownership:

```sql
foreign key (organization_id, booking_id)
  references booking.bookings (organization_id, id)
  on delete restrict
```

Rules:

1. `organization_id` is set from trusted server context, not accepted as authority from a request body.
2. Tenant identity is immutable. Moving records between businesses is a controlled exceptional workflow, not an ordinary update.
3. A global `account_id` link exists only where Account is genuinely the subject. Tenant authority travels through an active Membership, claim or party relation.
4. Cross-tenant references are prohibited except explicit platform-governance records whose semantics require both sides and are separately authorized/audited.
5. Cyclic FKs are avoided. Current-version pointers that create a cycle require a safe deferred validation sequence or a derived-current rule.

## Deletion and update actions

Default FK action is `ON DELETE RESTRICT`/`NO ACTION` and `ON UPDATE NO ACTION`; IDs and tenant keys are immutable.

| Relationship | Default action | Reason |
|---|---|---|
| Aggregate-owned disposable child | `CASCADE` only before the aggregate becomes evidence | Safe draft cleanup |
| Commercial/history/evidence child | `RESTRICT` | Preserve booking, quote, payment and decision history |
| Optional operational pointer | `SET NULL` only when “formerly linked” remains truthful | Preserve record without pretending a current link |
| Projection/search/cache row | Projector-managed rebuild/delete | Projection is not owner truth |
| Identity/customer closure | Lifecycle, tombstone or approved anonymization workflow | Avoid unsafe graph cascade |
| Organization closure | Restrict + retention workflow | Never erase a tenant history through one delete |

Every cascade is listed in executable-schema review with expected maximum fan-out and lock impact.

## Core uniqueness catalogue

Exact predicates depend on final lifecycle values, but these intentions are mandatory.

### IAM and organization

| Subject | Required uniqueness/invariant |
|---|---|
| Account contact | Normalized contact is unique only under an approved authentication policy; shared business numbers require an explicit design |
| Authenticator | Provider/type/external credential ID unique within environment |
| Session token family | Hashed family identifier globally unique; revoked family cannot reactivate |
| Membership | At most one current `(organization_id, account_id)` membership; history is explicit |
| Role capability | Unique `(role_id, capability_id)` |
| Role assignment | No duplicate role/scope/effective assignment for one membership |
| Branch/operation link | Unique same-tenant pair and non-overlapping effective interval where historical |
| Approval decision | One effective decision per approver/request unless supersession is explicit |

### Public presence and catalogue

| Subject | Required uniqueness/invariant |
|---|---|
| Public slug/subdomain | Unique normalized locator across active/reserved namespace; confusable/reserved rules run first |
| Service | Unique stable service code within tenant/operation when codes are exposed |
| Versioned object | Unique `(owner_id, version_number)` and immutable after publication |
| Publication | One active publication per operation/environment/channel unless effective dating says otherwise |
| Intake field | Unique stable field key inside form version; display label is not a key |
| Field option | Unique stable option key inside field/version |
| Policy acceptance | Unique subject/purpose/version/acceptance episode; presentations may repeat without duplicate acceptance |

### CRM, booking and request

| Subject | Required uniqueness/invariant |
|---|---|
| Customer | No global phone/email uniqueness assumption; contact matches become reviewable candidates |
| Customer-account claim | One active verified tenant/customer/account relation; never auto-merge on contact match alone |
| Booking reference | Unique within the chosen organization/operation namespace; retries cannot consume a second reference outcome |
| Booking source | Unique source event/client command/provider message conversion |
| Schedule version | Unique `(organization_id, booking_id, version_number)` and at most one current marker |
| Booking snapshot | Unique `(organization_id, booking_id, version_number)`; committed snapshot immutable |
| State transition | Unique aggregate sequence and causation/idempotency key |
| Request conversion | One terminal conversion for an accepted quote/version/idempotency key |
| Quote version | Unique `(quotation_id, version_number)`; issuance freezes content |
| Quote acceptance | One effective response per version/option/authorized party; supersession explicit |

### Integration, payment and communications

| Subject | Required uniqueness/invariant |
|---|---|
| Provider connection | Provider environment + external business asset/account has one allowed organization relationship |
| Provider event receipt | Unique `(provider_application_id, provider_event_id)` or documented stable digest fallback |
| Provider request | Stable Ubi request/idempotency key unique per connection/operation |
| Payment obligation | Domain purpose/subject/version unique so retries cannot create another amount due |
| Provider payment object | Unique `(payment_connection_id, provider_object_type, provider_object_id)` |
| Payment allocation | Unique source-value to obligation allocation identity; totals checked under lock |
| Refund command | Unique idempotency key and refund intent |
| Notification intent | Purpose + recipient + domain version + episode unique where duplicates are harmful |
| Delivery callback | Provider/message/event identity unique |

### Platform and sync

| Subject | Required uniqueness/invariant |
|---|---|
| Idempotency record | Unique `(scope, principal_or_tenant, operation, idempotency_key)` |
| Outbox event | Globally unique event ID and aggregate-version semantic uniqueness where applicable |
| Consumer receipt | Unique `(consumer_name, event_id)` |
| Scheduled task | Unique actionable logical task key/version |
| Client command | Unique `(installation_id, client_command_id)` with immutable request hash |
| Sync change | Unique tenant/scope/position and resource/version emission where required |

## Row-local check constraints

The database, not only the UI, enforces patterns such as:

```sql
check (end_at > start_at)
check (expires_at > created_at)
check (quantity > 0)
check (amount_minor >= 0)
check (currency_code ~ '^[A-Z]{3}$')
check (version_number > 0)
check (percentage_basis_points between 0 and 10000)
check (latitude between -90 and 90)
check (longitude between -180 and 180)
```

Additional required checks:

- exactly one authorized target for typed links implemented with nullable typed FKs;
- outcome/value compatibility, such as failure reason required for failure only;
- published version has publication time and issuer;
- revoked/expired/redeemed grant cannot present as active;
- local schedule has timezone ID and valid start/end semantics;
- price line has legal quantity/rate/sign/type combination;
- payment currency compatibility is enforced by composite keys or a locked transaction;
- plaintext secret/token/OTP fields do not exist in ordinary tables;
- JSON size/shape is validated at the API boundary while its discriminator/version is typed and constrained.

Volatile “current time” checks do not belong in static constraints. Commands evaluate expiry and idempotent scheduled work advances lifecycle state.

## State-transition enforcement

The aggregate row keeps constrained current state and `version`; append-only rows record transitions.

A state-changing command must:

1. lock the aggregate;
2. verify expected version, authority and legal transition;
3. update state/version;
4. append actor/cause/reason evidence;
5. append outbox work;
6. commit once.

Controlled functions/triggers are allowed only when multiple trusted writers could otherwise bypass a durable invariant. Frequently changing business policy remains in versioned domain code/config. There is one owner write path, not competing client, trigger and service state machines.

## Booking and capacity correctness

### Critical correction: one physical reservation surface

An exclusion constraint cannot prevent overlaps across two independent tables. Therefore `hold_allocations` and `capacity_commitments` must be backed by one authoritative `scheduling.resource_reservations` table for exclusive resources.

Each row has tenant, resource, occupied interval including approved buffers, state and exactly one current owner relation: hold or booking commitment. Consuming a hold converts its locked reservation rows in the same transaction; it does not release and then race to reinsert.

The logical hold/commitment aggregates remain; `resource_reservations` is the physical collision boundary. The executable catalogue must include it before implementation.

For exclusive resources, proposed enforcement is conceptually:

```sql
exclude using gist (
  organization_id with =,
  resource_id with =,
  occupied_interval with &&
)
where (reservation_state in ('held', 'committed'))
```

This needs reviewed `btree_gist` availability or a tested discrete-slot/locking alternative. Open-ended booking ranges are prohibited. Cancellation/release changes current state while append-only history preserves evidence.

### Pooled capacity

Quantity capacity needs a locked ledger, not a simple overlap constraint. A command:

1. derives exact bucket keys from published service/resource rules;
2. locks buckets in ascending key/time order;
3. reconciles relevant stale holds idempotently;
4. verifies `committed + held + requested <= configured_capacity`;
5. appends ledger entries and updates guarded summaries;
6. creates/converts the hold or commitment;
7. writes outbox work and commits.

Ledger rows are evidence; bucket totals are guarded current state. A reconciliation job recomputes and alarms on mismatch.

### Concurrency acceptance tests

- 100+ concurrent attempts for the last slot produce exactly one success.
- Hold expiry racing confirmation produces one valid terminal result and no leak.
- A retried offline command replays its first result, not another booking.
- Moving a booking locks old/new resources deterministically and never loses both.
- Different organizations cannot affect one another’s capacity.
- Failure after commit but before response is replay-safe.
- DST gap/fold and local date boundary cases use stored timezone/calendar context.

Diagram: `diagrams/43-booking-inventory-transaction.mmd`.

## Index design policy

Indexes follow approved query shapes, tenant isolation and measured plans—not every column.

### Baseline patterns

| Query/constraint | Index direction |
|---|---|
| Tenant row by ID | Unique `(organization_id, id)` |
| Children of parent | `(organization_id, parent_id, created_at/id)` |
| Actionable work | Partial tenant/state/due-time index for active states |
| Staff inbox | `(organization_id, operation_or_queue_id, work_state, priority, next_action_at, id)` |
| Calendar window | `(organization_id, resource_id, start_at, id)` plus selected GiST range constraint/index |
| Customer contact lookup | Tenant + type + normalized/search-safe value after privacy review |
| Provider dedup | Provider/app/connection + external event key unique |
| Worker claim | Partial `(available_at, priority, id)` for pending rows |
| Sync stream | `(organization_id, scope_id, position)` |
| Audit review | Tenant + recorded time + stable tie-breaker, protected path |
| Public slug | Normalized locator unique under lifecycle/registry design |

### Rules

- Equality/security/tenant predicates generally lead, then equality filters, range/sort and tie-breaker.
- One correct composite index is not equivalent to several isolated column indexes.
- `INCLUDE` is for small stable projection fields only after measurement, not sensitive wide payloads.
- Partial indexes serve small actionable subsets such as pending outbox, active grants and unresolved cases; query and index predicates must align.
- PostgreSQL full-text/trigram is acceptable for an approved early search use case. Normalization must support local scripts and names rather than English-only assumptions.
- Sensitive contacts are not copied into broad search indexes without purpose and protection review.
- No default GIN index on every JSONB document. Index stable approved paths or promote them to columns.

### Review gate

Every critical query needs representative bind values/cardinality, `EXPLAIN (ANALYZE, BUFFERS)` evidence, latency and rows-scanned target, index cost, RLS plan validation, pagination contract and owner. Production review uses `pg_stat_statements` and slow-query evidence; no index is auto-dropped from a short observation period.

## Pagination and bounded reads

Operational lists use cursor/keyset pagination, normally not deep `OFFSET`.

```text
ORDER BY next_action_at ASC, id ASC
cursor = (last_next_action_at, last_id)
```

Exposed cursors are opaque, versioned/signed where needed, bound to tenant/filter/sort and may expire. Page sizes have server caps. Expensive exact counts come from governed projections or asynchronous computation.

## Worker claiming and deadlock prevention

Workers claim small batches with short transactions and `FOR UPDATE SKIP LOCKED` where appropriate. Claim/lease metadata commits before external network work; completion is an idempotent second transaction.

Multi-row writers use a published lock order:

```text
organization/operation
→ aggregate by ascending ID
→ resource/bucket by ascending key and time
→ payment obligation/source
→ outbox/idempotency
```

Deadlock or serialization failures retry with bounded jitter under the same idempotency key. User interaction and provider calls never hold database locks.

## Partitioning decision

### Initial decision

Do **not** partition ordinary bookings, customers, memberships, services, resources or payments initially. Partitioning adds constraint, migration, query and operational complexity and does not repair poor indexes.

### Candidates after evidence

- provider event receipts;
- audit/security streams;
- notification/delivery attempts;
- sync change/tombstone history;
- completed outbox/inbox history;
- very large message/event history.

Partition only when size measurably harms vacuum, backup, retention or critical SLOs; time retention needs predictable drops; or volume exceeds the tested unpartitioned envelope—and the team has rehearsed partition creation, parity, pruning, default-partition alarms and restore.

Default candidate is time range by recorded/created time at a measured month/week boundary. Tenant hash partitioning is not the default because it complicates small-tenant operations and uniqueness. Queries must include the partition key where pruning is expected.

## Sharding decision

No application-level tenant sharding in the first release. It requires measured capacity/regulatory/isolation need plus a tested routing directory, cross-shard workflow rule, restore/export path and low-downtime tenant migration plan. Potential future scale alone is not evidence.

## Constraint and migration rollout

For existing large tables:

1. add a compatible nullable column/constraint form;
2. deploy code that tolerates both shapes;
3. backfill in bounded resumable batches;
4. add FK/check as `NOT VALID` where supported;
5. validate separately;
6. make nullability/default changes after proof;
7. build indexes concurrently where safe;
8. remove compatibility after all clients/jobs pass the boundary.

Never combine an unbounded rewrite, index build and behavior cutover blindly. Verify exact PostgreSQL/provider-version behavior before approval.

## Initial executable-schema gate

The first migrations wait until:

- launch vertical/archetype and workflow boundary are selected;
- UUID/reference ADR is approved;
- booking/capacity concurrency is load-tested in a disposable prototype;
- lifecycle values/transitions are approved;
- every release table has exact columns/types/null/default semantics;
- every FK/deletion action and cross-tenant negative test is listed;
- every critical query has an index/pagination contract;
- document 44’s RLS/role design passes behavior and plan tests;
- retention/deletion does not contradict evidence duties;
- rollback/restore and forward compatibility are rehearsed.

## Open decisions

- UUIDv7 generator/library and database validation.
- Human booking reference namespace/format and offline behavior.
- `btree_gist` availability versus a discrete-slot/explicit-lock alternative.
- Exclusive-resource reservation conversion/history details.
- Pooled-capacity bucket granularity for the chosen archetype.
- Contact normalization/search/encryption/tokenization by launch country.
- Evidence threshold and period for future partitioning.
- Online index/constraint capabilities of the selected managed PostgreSQL.
