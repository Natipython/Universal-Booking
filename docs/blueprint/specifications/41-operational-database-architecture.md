# Operational database architecture

Status: **Draft — PostgreSQL direction proposed, provider/version not selected**  
Depends on: conceptual data model, platform/deployment architecture and quality requirements  
Blocks: executable schema, constraints/indexes/RLS, migrations, capacity, backup/restore and implementation

## Executive direction

Use one managed PostgreSQL primary database for initial authoritative Ubi operations, with:

- multi-zone high availability where the selected provider supports it affordably;
- point-in-time recovery and tested restores;
- connection pooling;
- domain-owned lowercase schemas/tables;
- explicit tenant columns and composite tenant integrity;
- server authorization plus database row-level/privilege defense in depth;
- hard booking/capacity constraints/transactions;
- transactional idempotency, outbox/inbox, timers and process state;
- object binaries in private object storage, not database large blobs;
- no sharding, dedicated search cluster or analytics warehouse until measured evidence.

This is one physical database initially, not one unstructured schema.

## Why PostgreSQL fits

Ubi needs:

- strong ACID transactions for booking and capacity;
- unique/check/foreign-key/exclusion constraints;
- exact numeric/money handling;
- range/time/date/timezone support;
- row locking and transaction isolation;
- relational joins across coherent modules;
- JSONB only for controlled extensibility;
- row-level security and granular roles;
- full-text/trigram/geospatial extensions where justified;
- mature backup/PITR/replication/monitoring ecosystem.

Another database requires proof it can improve the approved workload without weakening these guarantees or increasing operating burden.

## Logical schema namespaces

Proposed PostgreSQL schemas:

```text
iam              accounts, contacts, credentials, sessions, devices
org              organizations, brands, branches, operations, memberships, authority
public_presence  slugs and safe public projections
catalog          versioned services, forms, policies and price definitions
scheduling       resources, rules, holds and capacity commitments
booking          booking lifecycle, parties, snapshots, schedules and fulfilment
request          custom requests, answers, quotations, acceptances and amendments
crm              tenant-private customers, contacts, claims and merges
conversation     conversations, participants, messages, assignment and takeover
access_grant     scoped link/token metadata and redemption
payment          obligations, intents, evidence, allocations, refunds and reconciliation
communications   notification intents, templates, endpoints, reminders and delivery attempts
integration      provider connections, credentials metadata and event receipts
media            object metadata, links, variants and processing state
sync             client commands, results, cursors, changes and tombstones
operations       internal cases, assignments and controlled elevations
audit            protected audit events/checkpoints
privacy          inventories, requests, propagation, retention and holds
reporting        rebuildable operational projections/metric definitions
platform         idempotency, outbox/inbox, timers, process managers, schema metadata
```

These names describe ownership, not direct client access. Public-facing data lives in controlled views/projections, not a PostgreSQL `public` free-for-all schema.

## Database ownership and role model

### Roles

| Role class | Purpose | Restrictions |
|---|---|---|
| Cluster/provider administrator | Managed service/control-plane emergency operations | Not used by application or ordinary migrations |
| Schema owner (`no login`) | Owns objects by schema/context | Not a runtime credential |
| Migration role | Applies reviewed schema changes | JIT CI identity; no routine app traffic |
| API runtime role | Executes approved functions/queries/commands | No object ownership, no bypass RLS, no DDL/superuser |
| Worker runtime role | Claims/processes approved queues/provider work | Separate grants from API where practical |
| Reporting role | Reads approved views/projections | No operational writes or raw cross-tenant detail by default |
| Operations application role | Calls restricted operations/domain interfaces | No arbitrary SQL/table write |
| Backup/replication role | Backup/PITR/replication duties | No application mutation; protected separately |
| Monitoring role | Catalog/statistics/health views | No personal content or mutation |

### Privilege rules

- Revoke default broad `public` privileges.
- Runtime roles never own tables/functions/schemas and never use superuser/BYPASSRLS.
- Grant schema usage and exact table/sequence/function actions only.
- Default privileges are explicitly configured so new objects do not become accessible accidentally.
- Security-definer functions are exceptional, non-public, fixed `search_path`, fully qualified, least-privilege and have execute revoked by default.
- Database access is server/workload only; Android/web clients never connect directly.
- Production human database access is JIT, read-only by default, case/reason/audit controlled.

## Tenant isolation strategy

Defense uses several layers together.

### Layer 1 — Immutable tenant key

Every organization-owned operational row includes `organization_id`, even when derivable from its parent. It supports:

- safe repository signatures and filters;
- composite foreign keys that prove parent/child same tenant;
- tenant-leading indexes;
- RLS predicates;
- object/search/export/change-feed isolation;
- incident/restore/isolation testing.

### Layer 2 — Composite tenant referential integrity

Parent tables expose unique `(organization_id, id)` and child foreign keys reference both:

```text
(organization_id, booking_id)
→ booking.bookings (organization_id, id)
```

This prevents a valid child tenant ID from pairing with another tenant’s parent ID even if application validation fails.

Global entities such as Account are referenced through explicit relationship tables (Membership, customer-account claim) rather than adding tenant ownership to Account.

### Layer 3 — Application authorization

Server resolves authenticated principal, current membership/capability/scope, resource tenant and state policy. Repository/module interfaces require typed tenant context. Client-supplied tenant is never sufficient.

### Layer 4 — PostgreSQL RLS defense

Proposed direction: enable and force RLS on tenant-owned tables after performance/security prototype.

Within each database transaction, the application establishes trusted local context such as:

```text
app.current_account_id
app.current_organization_id
app.current_membership_id
app.current_authority_version
app.current_workforce_case_id (only operations path)
```

Use transaction-local settings (`SET LOCAL` or an approved transaction entry function) so pooled connections cannot leak context. Policies use indexed `organization_id` and small stable helper functions.

Important limitations:

- RLS does not replace application capability/resource/state authorization.
- If a compromised runtime can set arbitrary tenant context, RLS alone is not a complete boundary; database role separation, signed/validated context entry and application controls still matter.
- Table owner/BYPASSRLS must not be the runtime role; force RLS where applicable.
- Policies require positive/negative performance tests and explicit public/platform/workforce paths.
- Background cross-tenant jobs use a separate restricted role/function/queue scope and process one explicit tenant/work item at a time; they do not disable isolation globally casually.

Final RLS architecture is approved in `44-multitenancy-authorization-and-data-isolation.md`.

## Identifier strategy

Requirements:

- opaque, immutable, non-reassigned;
- client/offline creation safe for relevant aggregates;
- efficient PostgreSQL indexing/locality;
- string-safe API representation;
- no embedded personal/tenant/time/business meaning exposed as authority;
- distinct provider IDs, public codes/slugs and token secrets.

Proposed direction:

- time-ordered UUID-compatible IDs (for example UUIDv7 semantics) generated by trusted application/client library for offline-capable entities;
- PostgreSQL `uuid` storage, with generation implementation pinned/tested and no unreviewed extension dependency;
- `bigint generated ... identity` only for internal append/projection sequence positions where offline/public identity is unnecessary;
- high-entropy token secrets hashed/protected separately;
- provider identifiers stored as `text` within provider/environment/account/connection unique namespace.

Exact UUID/version/library/browser/mobile collision/clock behaviour is an ADR before schema DDL.

## Data type conventions

| Concept | PostgreSQL direction |
|---|---|
| Opaque domain ID | `uuid` after identifier ADR |
| Internal sequence/change position | `bigint generated ... as identity` or database sequence |
| Human/provider identifiers | `text` with semantic length/check/normalization constraints where needed |
| Current version | `bigint not null` incremented atomically |
| Server instant | `timestamptz` |
| Local date | `date` plus timezone/calendar context where needed |
| Local time/recurrence | structured `time`/weekday/rule/effective range, not only display string |
| Interval commitment | timestamp/date range or start/end with explicit non-empty/order constraint |
| Timezone/locale/currency | constrained `text` validated by application/reference data |
| Money | exact `bigint` minor units + currency when valid for approved currencies; `numeric` for exact rates/percent/exceptional decimal models |
| Quantity/capacity | `integer`/`bigint`/exact `numeric` with nonnegative checks |
| Boolean | `boolean`, not text |
| Status/type | lowercase `text` + check/reference/transition policy initially; PostgreSQL enum only after migration trade-off approval |
| Structured extension | bounded validated `jsonb` with schema/version and promoted indexed fields |
| Free text | `text`, with application/DB size and purpose constraints |
| Phone/email | protected `text` normalized/display fields; no numeric phone type |
| Network/address | PostgreSQL network/geospatial type only when purpose/extension approved |

Never use PostgreSQL `money`, binary floating point for financial truth, `timestamp without time zone` for an instant, or arbitrary `varchar(255)` as a universal habit.

## Naming conventions

- Lowercase `snake_case`, unquoted identifiers.
- Plural table names and singular column/reference names consistently.
- Primary key `id`; foreign key `<entity>_id`; tenant key `organization_id`.
- Timestamps `<verb>_at`; date-only `<meaning>_date`; booleans `is_`/`has_` only when it improves clarity.
- Unique/check/foreign-key/index names deterministic and table-prefixed.
- Context-qualified statuses in columns such as `booking_status`, `payment_status`; never one global `status` type.
- Avoid reserved/generic names such as `user`, `order`, `transaction`, `type`, `data` without domain qualification.

## Relational versus JSONB

### Relational first

Use columns/child tables for fields participating in:

- authorization/tenant constraints;
- foreign keys/uniqueness/checks;
- availability/capacity/money;
- lifecycle/version/idempotency;
- search/sort/report/deletion;
- provider routing and reconciliation.

### JSONB approved uses

- immutable command/integration event payload with schema version;
- provider raw/normalized evidence supplement;
- custom form value envelope tied to field definition/version;
- immutable snapshot rendering supplement;
- bounded diagnostics.

Every JSONB column has schema owner/version, maximum size, allowed fields, validation, indexing/query plan, retention/deletion and migration path. Do not create a generic `metadata jsonb` on every table.

## Transaction strategy

### Default

- PostgreSQL `read committed` for ordinary commands/queries with explicit conditional updates/locks/constraints.
- Short transactions; validate external data and call providers outside the DB transaction.
- Every transaction has bounded statement/lock/idle-in-transaction timeouts appropriate to route/job.
- Stable error mapping for unique/check/FK/exclusion/serialization/deadlock outcomes.

### Booking and capacity

Within one transaction:

1. claim/check idempotency;
2. revalidate tenant/authority/config/booking/hold versions;
3. lock affected existing booking/hold/resource/capacity rows in deterministic order by immutable key;
4. validate/insert/move/release capacity under constraints;
5. update Booking/current schedule/snapshot/history;
6. append audit/outbox and command result;
7. commit.

Use database exclusion/unique/check constraints where feasible for exclusive-resource overlaps. Pooled capacity may use locked bucket/ledger conditional update. Exactly one approved approach per archetype avoids competing algorithms.

### Deadlock/serialization handling

- Consistent lock order across resource/hold/booking/payment rows.
- Prefer one atomic statement or explicit ordered `FOR UPDATE` to check-then-update races.
- Detect and retry PostgreSQL serialization/deadlock errors only around idempotent complete command transactions, with bounded jitter/attempts.
- Never automatically retry an external payment/message side effect inside a DB transaction.
- Log safe conflict/deadlock diagnostics and test adversarial concurrency.

## Idempotency, outbox, inbox and timers

### Idempotency

Unique logical scope includes principal/tenant/command/action/key as defined. Store request hash, processing/terminal result and retention. Same key/different request is conflict.

### Outbox

Inserted in same owner transaction as domain change. Workers atomically claim ordered eligible rows using `FOR UPDATE SKIP LOCKED`/lease pattern. Partial indexes target pending/available rows.

### Inbox/provider receipt

Unique `(consumer/provider namespace, event_id)` prevents repeated logical effects. Raw provider evidence is protected/bounded; canonical processing is separately versioned.

### Timers/process managers

Rows include tenant/owner/resource/due time/expected version/status/lease/idempotency. Worker rechecks current state before expiry/reminder/action.

High-churn queue tables get tuned indexes/autovacuum and archival/partitioning only after volume measurement.

## Query and index strategy

Indexes are derived from approved access paths and measured query plans.

Rules:

- Every foreign key receives an appropriate index unless a documented query/size exception exists; PostgreSQL does not add it automatically.
- Tenant-owned common indexes generally lead with `organization_id`, followed by equality/scope/status and range/order columns.
- Composite index equality columns precede range/sort columns for the intended query.
- Cursor/keyset pagination uses the complete stable order key, often `(organization_id, scheduled_start, id)`.
- Partial indexes target active/pending/non-deleted states actually used by queries/workers.
- GiST/range/exclusion for scheduling; GIN only for approved JSONB/full-text/array queries; BRIN only for very large naturally ordered time tables.
- Avoid redundant/unused indexes because writes/storage/vacuum cost matters.
- Index RLS predicate/relationship columns.
- Query selects minimum columns; avoid unbounded `select *`, N+1 and deep `offset` pagination.
- `EXPLAIN (ANALYZE, BUFFERS)` on safe non-production/controlled data and `pg_stat_statements` guide tuning.

Detailed catalog is defined in `43-keys-constraints-indexes-and-partitioning.md`.

## Connection architecture

PostgreSQL connections are bounded resources.

### Pooling

- Application/worker connect through provider pooler or PgBouncer-equivalent.
- Transaction pooling is preferred for ordinary stateless transactions.
- Session-dependent features/temp tables/session `SET` are prohibited on transaction-pooled paths; use `SET LOCAL` inside transaction.
- Prepared statement/driver behaviour is tested with selected pool mode.
- Migration/administrative/streaming jobs use separate limited direct/session pool when required.

### Budget

Calculate:

```text
database max connections
- provider/admin/monitor/replication reserve
= available application connections

allocate by API replicas + worker pools + deployment overlap + emergency headroom
```

- Configure small per-process pools; do not multiply large defaults by replicas.
- Queue/admission/backpressure when pool wait rises.
- Monitor active/idle/waiting/long transaction and pool saturation.
- `idle_in_transaction_session_timeout`, statement/lock timeouts and maximum query duration by workload.
- Autoscaling cannot exceed the global DB connection/IO budget.

## Read/write topology

### Primary

All authoritative commands, authorization requiring current state, booking availability commit, payment allocation/idempotency and sync command results use primary/current transaction.

### Standby/failover

Managed synchronous/zone standby if chosen. Test failover endpoint, open transaction failure and pool reconnection/idempotency recovery.

### Read replica

Not required initially. Later permitted for lag-tolerant operational reporting/backfills with:

- explicit maximum lag/fallback;
- no current permission/availability/payment/booking decision;
- tenant/RLS equivalent controls;
- query workload/connection isolation;
- monitoring and failover semantics.

Analytics projections may move outside operational primary when measured.

## Extensions policy

Extensions expand attack, upgrade and portability surface. Maintain allow-list, owner, version/support, backup/restore/migration and provider availability.

Likely evaluation:

- `pg_stat_statements` — query performance visibility;
- `btree_gist` — scheduling exclusion constraints with tenant/resource equality plus ranges;
- `pg_trgm` — bounded tenant/name/code search if justified;
- cryptographic primitives where application/KMS design requires and security approves;
- PostGIS only when real service-area/geospatial queries justify it;
- time-ordered UUID generation only if provider/core version safely supports the approved strategy.

Do not install an extension because it might be useful later.

## Partitioning and archival

Do not partition ordinary business tables at launch. Partitioning adds key/unique/FK/migration/query/operations complexity.

Candidates only after volume/retention evidence:

- audit events;
- provider event receipts/raw envelopes;
- outbox/inbox history;
- notification delivery attempts;
- sync changes/tombstones;
- analytics ingestion;
- high-volume messages.

Use time-range partitioning when tables are very large, queries/retention include the partition key and dropping old partitions materially improves maintenance. Tenant hashing/sharding is not a substitute for proper indexes.

## Vacuum, statistics and maintenance

- Autovacuum remains enabled; tune per high-churn queue/idempotency/status tables based on dead tuples/rate.
- Monitor table/index bloat, freeze age, long transactions, last vacuum/analyze and database storage.
- Analyze after large backfill/import and raise statistics target only for demonstrated skew/plan issues.
- Reindex/vacuum-full/table rewrite are planned operations with space/lock/failover impact—not routine first response.
- Hot rows (global counters/single queue cursor) are avoided.
- Large deletion happens in bounded batches or partition retention, with vacuum/storage planning.

## Observability

Monitor:

- availability/failover/restart and recovery state;
- connections/pool wait/idle-in-transaction;
- transactions/commits/rollbacks/deadlocks/serialization/lock waits;
- query latency/calls/rows/shared/temp IO via normalized statements;
- slow query plans and index scan/usage/missing FK indexes;
- CPU/memory/storage/IOPS/WAL/checkpoint/cache hit;
- replication lag and backup/PITR freshness;
- vacuum/analyze/dead tuples/bloat/freeze age;
- table/index growth by context/tenant class without leaking content;
- outbox/inbox/timer/reconciliation age;
- constraint/tenant-isolation violations.

Query text/parameters/logs are redacted; no secrets/OTP/payment/message/customer payload.

## Backup and restore

- Managed automated backup plus continuous WAL/PITR targeting approved RPO.
- Encrypted backup access separated from runtime and protected from routine deletion.
- Retention/location/replication follow data governance/legal constraints.
- Schema/migration/tool version and object-storage/config dependencies captured.
- Restore to isolated environment before launch and quarterly thereafter.
- Verify row counts/checksums where meaningful, constraints, RLS/privileges, tenant isolation, booking-capacity integrity, idempotency/outbox/inbox/timers, audit evidence and provider reconciliation.
- A backup job success without application-level restore validation is not proof.

Detailed procedures are defined in `49-migrations-backups-restore-and-disaster-recovery.md`.

## Scaling path

1. Fix queries/indexes/N+1/pagination and pool sizing.
2. Tune schema/autovacuum/storage/instance based on evidence.
3. Isolate worker/report workload and add lag-tolerant replica if needed.
4. Archive/partition genuinely large append-only tables.
5. Extract specialized search/analytics/provider workloads.
6. Consider tenant/database sharding only after single-primary limits, residency or isolation tiers are evidenced and migration/routing/restore are designed.

Sharding too early makes booking transactions, tenant moves, analytics, backups and support harder.

## Database change governance

- Schema source/migrations reviewed and version-controlled; no production GUI drift.
- Expand/migrate/contract compatibility across app/worker/Android versions.
- Constraints/indexes added using lock-safe validation/concurrent strategy where PostgreSQL supports it.
- Every migration declares lock/rewrite/duration/disk/WAL/rollback/roll-forward and backup impact.
- Large backfills are resumable/idempotent/batched/observable and tenant-fair.
- Destructive column/table/partition removal after retention/client/event/replay window and verified backups.
- Migration role only during controlled deployment; app never auto-migrates on startup.

## Database fitness tests

1. Concurrently attempt conflicting exclusive-resource bookings; constraint/transaction permits one valid commitment.
2. Attempt child row with another tenant’s valid parent ID; composite FK/RLS rejects it.
3. Omit tenant filter in a runtime query; RLS/repository tests prevent cross-tenant rows.
4. Replay command/provider event 100 times; uniqueness/idempotency yields one logical effect.
5. Kill worker after claiming/outbound ambiguity; lease/inbox/provider reconciliation recovers.
6. Revoke membership then run an offline command; current authority rejects.
7. Saturate API replicas; pool/backpressure protects DB rather than exhausting connections.
8. Fail over database during command; same key recovers committed/failed outcome without duplicate.
9. Restore PITR and verify constraints/RLS/outbox/provider reconciliation within objectives.
10. Query each critical access path at forecast + 3× data using expected index/keyset plan.
11. Delete/restrict tenant customer data; tombstone/derived systems update without deleting required booking/financial evidence.
12. Old/new app/workers run during expand/contract migration without incompatible writes.

## Open decisions

- Managed PostgreSQL provider/version/region/HA/PITR and budget.
- Identifier strategy and database/application generation.
- RLS context-entry/policy/helper implementation and performance results.
- Booking exclusive/pooled capacity constraint/locking design by first archetype.
- Schema-per-context ownership and cross-schema FK policy.
- Pooler/driver/prepared statement and connection budgets.
- Approved extensions and portability.
- Audit storage/integrity/partition/retention.
- Encryption/field-level/tokenization for selected sensitive fields.
- Read replica/analytics/search extraction thresholds.

