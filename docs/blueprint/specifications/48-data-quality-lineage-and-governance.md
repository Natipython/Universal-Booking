# Data quality, lineage and governance

Status: **Draft — governance and control model proposed; owners, thresholds and tooling unresolved**  
Depends on: domain model, schema catalogue, events/metrics, retention and storage architecture  
Blocks: production datasets, dashboards, migration acceptance, provider reconciliation and data operations

## Purpose

Ubi’s booking promise depends on trustworthy availability, identities, customer contacts, prices, times, payment evidence and provider state. This document defines ownership, quality dimensions, contracts, lineage, reconciliation, incident handling and controlled change.

Data quality is not a cleanup job performed by a future data team. Product input, database constraints, domain transactions, provider adapters, migrations, projections and operations all share responsibility.

## Principles

- Prevent invalid data as close to creation as practical.
- Preserve original/provenance evidence while keeping a reviewed canonical interpretation.
- Never hide uncertainty by converting unknown into false, zero, empty or “paid.”
- One owner defines each fact; projections do not compete with owner state.
- Quality rules are executable, versioned, observable and linked to business impact.
- Tenant isolation is a data-quality dimension as well as a security boundary.
- Human correction is traceable, authorized and reversible where appropriate.
- Provider data is evidence with known semantics, not automatically truth.
- Metrics and ML/automation outputs publish lineage and confidence/limitations.
- A green pipeline is not enough: reconcile meaningful facts end to end.

## Quality dimensions

| Dimension | Meaning for Ubi | Example failure |
|---|---|---|
| Validity | Value conforms to approved type/domain/rule | End before start; invalid currency code |
| Completeness | Required fact is present for its lifecycle state | Confirmed booking lacks accepted schedule snapshot |
| Accuracy | Value reflects the real-world fact to the justified standard | Wrong customer contact or branch timezone |
| Consistency | Equivalent facts agree within documented model | Booking says confirmed but no active reservation |
| Integrity | Relationships, tenant and versions are structurally valid | Child references another organization |
| Uniqueness | One real/business event is not unintentionally duplicated | Provider webhook creates two payment transactions |
| Timeliness | Data is fresh enough for its decision | Search still shows unpublished service |
| Conformity | Value uses controlled format/semantics | Provider “success” mapped to Ubi “settled” incorrectly |
| Traceability | Source, version and transformation are knowable | Price total cannot be traced to line/version |
| Accessibility/interpretability | Authorized people can understand it | Ambiguous state/label causes staff error |

Quality targets differ by data product. Capacity integrity may require zero tolerated conflicts; a noncritical analytics projection may allow bounded freshness and later correction.

## Ownership model

| Role | Responsibility |
|---|---|
| Domain/data owner | Defines business meaning, quality objective, access and acceptable risk |
| Data steward | Maintains dictionary, rules, issues, reference data and user correction process |
| Producing engineering owner | Enforces schema/input/transaction/event contract and monitors production |
| Consuming owner | Declares dependencies, handles version/freshness and does not reinterpret silently |
| Security/privacy/legal | Reviews classification, purpose, access, retention and regulated meaning |
| Operations/support | Follows controlled correction/reconciliation playbooks and records evidence |
| Analytics owner | Defines metric semantics, lineage, reconciliation and caveats |

Initially one person may hold several roles, but every critical data product still names accountability. “The database team” is not a business owner.

## Data product and contract register

Treat important outputs as governed data products:

- availability/slot answer;
- booking and schedule snapshot;
- staff work queue;
- customer/contact view;
- quote and price total;
- payment/reconciliation status;
- public business/service projection;
- message/provider delivery state;
- sync change feed;
- audit trail and privacy export;
- merchant/platform metric.

Each contract records:

```text
stable product/contract ID and owner
business meaning and prohibited interpretations
authoritative sources and lineage
schema/version/compatibility policy
keys, tenant scope and granularity
required/optional fields and null semantics
quality rules/objectives and freshness
dedup/order/correction behavior
classification/purpose/retention
consumers and break/change notification
reconciliation and incident playbook
```

## Source and provenance model

Every mutable or externally derived high-value fact records sufficient provenance:

- source type: staff, customer, account, system rule, import, provider, calculation;
- source identifier/message/event/file where permissible;
- source timestamp and trusted receipt timestamp;
- producer/provider/schema/adapter/app version;
- actor and authority or service identity;
- transformation/rule/model/version;
- confidence/review state for suggestions, never fake certainty;
- superseded/corrected relationship and reason.

Preserve provider raw/digest evidence only as long as justified; canonical records retain exact reference and mapping version. A raw payload is not allowed to become an indefinite shadow database.

## Lineage levels

### Level 1 — system lineage

Shows that Booking PostgreSQL produces an outbox event, reporting projection, search document, sync change, notification and audit record.

### Level 2 — dataset/table lineage

Shows exact owner tables/events/versions feeding a projection/export/metric.

### Level 3 — field lineage for critical facts

Required for price/currency, payment state/value, booking schedule/timezone, customer contact destination, permission decision, legal/notice acceptance and KPI numerator/denominator.

Lineage is stored in version-controlled contracts/catalogue plus runtime checkpoints, not only a diagram that drifts from code.

## Validation boundaries

### Client

Provides immediate, understandable input help and prevents obvious errors. It is never trusted as the only validation because old/offline/modified clients exist.

### API/domain

Validates schema, normalization, actor permission, lifecycle policy, cross-row/domain rules and provider semantics. It returns field/action-specific errors that a below-average user can recover from.

### Database

Enforces types, nullability, checks, tenant-qualified relationships, uniqueness, capacity constraints/locked ledgers and immutable/version rules.

### Event/queue boundary

Validates producer identity, schema version, payload size, tenant/resource consistency, event ID and classification. Invalid messages quarantine with actionable diagnostics; they do not retry forever.

### Projection/consumer

Rejects stale aggregate versions, deduplicates, records checkpoint/schema version and reconciles with owner counts/samples.

### Provider/import boundary

Stores authenticated receipt/provenance, validates size/schema/known semantics, maps through a versioned adapter and isolates malformed/unknown data for review. Never treat a CSV column or provider label as a trusted internal enum automatically.

## Null, unknown and uncertainty semantics

Each field distinguishes relevant states explicitly:

- not collected/not applicable;
- unknown to source;
- pending verification;
- reported/unverified;
- verified;
- provider state unknown;
- intentionally withheld/restricted;
- deleted/anonymized.

Do not overload `NULL`, empty string, zero and false to mean all of these. APIs and reports must not present “unknown payment state” as unpaid/paid without context.

## Reference and master data

Controlled registries include capabilities, lifecycle codes, countries, currencies, timezone IDs, locales/scripts, calendar modes, provider types/capabilities, event/problem schemas, payment/notification result mappings and reason codes.

- Stable machine code is separate from localized display label.
- Registry changes are versioned/reviewed and backwards compatible where clients/events depend on them.
- ISO/IANA or regulator/provider identifiers use pinned source/version and local overrides only with evidence.
- Currency minor-unit/rounding rules are explicit by version; never assume two decimal places.
- Timezone database version is recorded where future schedule interpretation could change.
- A deleted/retired code remains interpretable in history.

## Customer/contact quality and deduplication

East African booking contacts may involve shared phones, multiple SIMs, reused numbers, informal names, different scripts/transliterations and limited addresses. Therefore:

- normalize phone numbers under selected country/numbering context while retaining display/source form;
- record contact purpose, provenance and verification status separately;
- do not assume one phone equals one human or global Account;
- fuzzy name/contact matches create candidates, never automatic destructive merges;
- merge requires same-tenant authority, evidence, preview and protected reversible history;
- conflicting booking snapshots remain historically accurate after CRM merge;
- staff can correct understandable fields without exposing hidden cross-tenant matches.

Quality research with local pilot merchants is required before selecting merge thresholds.

## Time, calendar and schedule quality

- Store UTC instants plus authoritative local timezone/context where domain meaning requires.
- Preserve local date/time/calendar entered/shown and conversion version in committed schedule snapshots.
- Validate DST gaps/folds even if the initial market rarely encounters them; providers/travel/future regions may.
- Availability and booking use the same published rule/version and rounding/buffer semantics.
- Timezone database updates trigger regression/recalculation review for future schedules, never silent history rewrite.
- Reports use documented organization/operation local day boundaries.

## Money and payment quality

- Exact amount plus currency; no float and no cross-currency sum without a dated source/version.
- Preserve obligation, reported evidence, provider transaction, allocation, refund, reversal, dispute and settlement as different facts.
- Provider status mapping is versioned; `submitted`/`successful`/`settled` are not synonyms.
- Allocation/refund totals reconcile under lock with source/obligation limits.
- Manual verification records reviewer, source, reason and approval requirement.
- Daily and event-driven reconciliation compares Ubi with authenticated provider/settlement evidence and creates a case for discrepancy.
- Unknown provider outcome remains unknown and is reconciled; retry does not manufacture a second charge.

## Core invariant/reconciliation catalogue

| Reconciliation | Example invariant | Response |
|---|---|---|
| Booking ↔ reservation | Every capacity-requiring confirmed booking version has matching active reservations; no unapproved overlap | Block/alert, isolate workflow, investigate before auto-fix |
| Hold ↔ expiry | Expired/released hold has no active reserving quantity | Idempotent release/rebuild and alert leakage |
| Booking ↔ snapshot | Current confirmed state points to immutable service/price/policy/schedule evidence | High-severity data case |
| Payment obligation ↔ allocation | Verified allocations/refunds/reversals balance within exact currency rules | Reconciliation case; never silently mark paid |
| Provider receipts ↔ canonical | Authenticated processed receipt has terminal mapping or explicit quarantined/unknown state | Retry/review adapter; preserve evidence |
| Database ↔ outbox/consumer | Aggregate change requiring event has durable outbox; consumer dedup/checkpoint progresses | Republish/rebuild through controlled tool |
| Owner ↔ search/public/sync | Projection version/count/sample matches owner and tombstones | Reindex/resync/purge; security urgency for stale access |
| Database ↔ object | Ready metadata maps to existing checksum-matching bytes; orphan bytes/rows bounded | Quarantine/delete/recover by lifecycle |
| Authority ↔ cache/device | Revoked scope/version no longer authorizes cached/server action; purge converges | Security alert and forced sync/session action |
| Metric ↔ domain | KPI counts/value reconcile to authoritative rows under same definition/window | Mark dashboard degraded, repair/backfill |

Automated correction is allowed only when the owner and safe deterministic rule are known. Otherwise quarantine/case/review protects evidence.

## Quality objectives and severity

Each rule defines measurement window, denominator and action threshold. Example severity direction:

- **Q0 critical:** cross-tenant mix, confirmed capacity conflict, money imbalance/loss, auth bypass, unrecoverable owner-data corruption. Immediate incident and workflow protection.
- **Q1 high:** missing booking snapshot, stale revocation/public restriction, provider duplicate effect, material reconciliation drift. Urgent response.
- **Q2 medium:** search/sync/report freshness outside objective, incomplete noncritical field, repeated correction burden. Scheduled incident/problem work.
- **Q3 low:** cosmetic label/reference enrichment issue without wrong decision. Backlog under trend review.

No percentage target is copied across all data. Critical invariants may demand zero tolerance; optional profile completeness should not pressure users to provide unnecessary data.

## Observability and scorecards

Quality scorecards include:

- rule pass/fail count and affected tenant/resources;
- freshness/checkpoint lag and last successful reconciliation;
- invalid/quarantined provider/event/import count;
- duplicate candidate/merge/reversal outcomes;
- unknown/unverified state age;
- correction and support reason patterns;
- owner/projection count/version/checksum drift;
- migration/backfill validation progress;
- privacy/security-sensitive stale copy or deletion-task breach.

Alerts route by business impact and owner. Tenant-facing status explains actionable impact without leaking other tenants/provider secrets. Platform aggregates do not expose customer content.

## Issue management and correction

1. Detect through constraint, reconciliation, support report, provider alert or analysis.
2. Assign severity, owner, affected scope/time/version and privacy/security impact.
3. Protect users: stop unsafe confirmation/payment/export, mark uncertainty or disable adapter narrowly.
4. Preserve evidence and create case/audit trail.
5. Diagnose source and lineage; do not edit symptoms in projections first.
6. Correct through versioned domain command/migration/rebuild with preview and idempotency.
7. Reconcile downstream copies and communicate affected parties as required.
8. Add preventive constraint/validation/monitor/test and review similar data.

Direct production SQL edits are emergency exceptions using JIT approval, reviewed script, before/after evidence, dry-run/backup/reconciliation and post-incident review.

## Schema/event/provider change governance

- Contracts are versioned and reviewed by producer and known consumers.
- Breaking changes use expand/migrate/contract with mobile/provider compatibility window.
- Schema registry records supported app/adapter/event versions.
- Unknown provider fields/statuses do not map to success by default.
- Migration includes data-quality preconditions, backfill counters, invalid-row quarantine and postcondition queries.
- Metrics receive a new definition version when meaning/denominator changes.
- Historical records retain enough registry/version context to remain interpretable.

## Test and non-production data

- Automated tests use synthetic/factory data designed for multilingual names/scripts, shared phones, missing landmarks, calendar/timezone edges, provider duplicates/out-of-order events, money/currency edges, offline conflicts and cross-tenant attacks.
- Production personal/business content is not copied into developer environments by default.
- If a rare production-derived dataset is essential, it requires approved minimization/de-identification, restricted environment/access, purpose, expiry and deletion proof.
- Test tenants/events carry explicit flags that cannot be confused with production KPI/payment/business truth.
- Seed/reference fixtures are versioned and reproducible.

## Automation/AI extraction quality

An extracted field is a suggestion until policy and confidence allow otherwise.

- Store source message reference, model/rule/version, confidence and reviewed result.
- Never auto-confirm price, time, identity, payment or acceptance from ambiguous text.
- Use local-language/realistic evaluation sets with consent/approved synthetic data.
- Measure per-field precision/recall and harmful error classes, not one average score.
- Low confidence/unsupported language routes to staff with clear source context.
- Staff correction feeds governed evaluation only under approved data-use purpose.
- Model/provider change is a contract/version change with regression and privacy review.

## Release gates

- Every first-release data product/fact has owner, steward, contract, quality objectives and incident route.
- Database, API, event and provider validations implement the same semantic rules without contradictory enums.
- Core booking/capacity/payment/tenant invariants pass property and concurrency tests.
- Reference data versions for country, currency, timezone, locale and providers are pinned/reviewed.
- Provider duplicates, out-of-order, malformed and unknown states are tested.
- Owner-to-search/object/sync/report reconciliation and rebuild are rehearsed.
- Quality scorecard uses synthetic canary tenants without leaking production content.
- Correction tooling is authorized/audited and ordinary support has no direct data rewrite path.
- Migration and backup-restore validation proves counts, constraints, checksums and tombstones.
- Local pilot users validate contact/name/address/search assumptions.

## Open decisions

- Named owners/stewards and on-call responsibility for each first-release data product.
- Numerical quality/freshness thresholds after representative capacity/pilot tests.
- Contract/catalogue implementation tooling versus version-controlled documents plus registry tables.
- Country phone/address/reference sources and update process.
- Customer duplicate-candidate algorithm and safe review experience.
- Field-level lineage scope/tooling for payments, schedule, authority and metrics.
- Data-quality dashboard/case integration implementation.
- Production-derived troubleshooting dataset exception process.
