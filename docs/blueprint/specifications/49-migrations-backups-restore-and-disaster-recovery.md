# Migrations, backups, restore and disaster recovery

Status: **Draft — recovery architecture proposed; provider, budgets and approved objectives unresolved**  
Depends on: deployment/database/storage, schema governance, retention/deletion and quality controls  
Blocks: provider procurement, migration tooling, production launch, recovery runbooks and continuity commitments

## Purpose

This document defines how Ubi changes schemas/data safely, protects authoritative records, restores verified service and continues essential booking operations through failures.

A backup is not successful because a provider dashboard says it exists. Ubi has a backup only when authorized people can restore the required database, objects, keys/configuration and deletion state within an approved objective and reconcile the result.

## Recovery vocabulary

- **RPO (recovery point objective):** maximum tolerable data loss measured backward from disruption.
- **RTO (recovery time objective):** target time to restore the defined service/data capability.
- **PITR:** restore PostgreSQL to an exact supported point from base backup plus transaction log/WAL.
- **HA failover:** provider/local resilience to an instance/zone fault; not the same as backup or region disaster recovery.
- **Backup:** protected recoverable copy with retention and access controls.
- **Restore:** technical recovery of bytes/database.
- **Recovery:** verified restoration of correct business service, security, integrations and reconciliation.
- **Business continuity:** how people continue essential work while technology is degraded.

No product/legal promise uses RPO/RTO before provider testing and leadership approval.

## Initial resilience direction

For an affordable pilot:

- one managed PostgreSQL primary with same-region multi-zone standby/automatic failover where the chosen service supports it;
- continuous transaction-log/PITR plus automated encrypted snapshots with bounded retention;
- private object storage with versioning/soft-delete and inventory where justified;
- backups/copies protected in a separate provider account/project/security boundary where feasible;
- infrastructure/config/schema/event contracts in version control and reproducible deployment;
- managed secret/KMS recovery design with strict separation and tested key availability;
- manual, rehearsed regional recovery rather than costly active-active multi-region writes;
- rebuildable search/cache/reporting/sync projections restored from PostgreSQL owner data;
- explicit merchant/customer degraded-mode communication and manual export/continuity procedure.

Multi-region active-active booking is rejected initially. It makes capacity consistency, payment effects, conflict resolution, data location and cost materially harder. Add it only from measured availability/regulatory need and a proven consistency model.

## Service/data recovery tiers

These are provisional selection/testing targets, not approved customer commitments.

| Tier | Capability/data | Provisional pilot direction | Recovery order |
|---|---|---|---|
| R0 | IAM/session authority, organization restrictions, booking/capacity owner state, idempotency/outbox, payment obligations/evidence | Seek tested PITR RPO no worse than minutes and essential write recovery within hours; exact values after provider rehearsal | First |
| R1 | Catalog/config, CRM needed for upcoming work, communications/process state, audit/privacy deletion/hold state | Same database recovery point; restore immediately after/with R0 | First/second |
| R2 | Required private media/object bytes and manifests | Versioned/protected object recovery matched to DB point; upload/processing can pause | Second |
| R3 | Search, cache, reporting projections, transformed public variants | Rebuild from owner records; stale/unavailable acceptable within declared degradation window | Later |
| R4 | Product analytics/debug telemetry beyond required evidence | Loss may be acceptable within approved retention; never delay booking recovery | Last/optional |

Before launch, replace “minutes/hours” with exact tested objectives, owner, dependency, business impact and cost. If affordable provider capability cannot meet them, scope/promise changes—documentation does not pretend otherwise.

## Backup scope

### PostgreSQL

Required:

- provider-managed encrypted continuous WAL/PITR and scheduled snapshots;
- retention longer than detection/recovery decision window but no longer than approved privacy/cost schedule;
- backup copy/control plane protected from ordinary application and migration credentials;
- monitoring of backup age, WAL archive gaps, snapshot failures and storage capacity;
- periodic portable logical/schema exports for portability/forensics where useful, never the only primary backup;
- exact PostgreSQL/extensions/collation/timezone/provider version inventory;
- roles/grants/RLS/functions/extensions included and verified in restore.

Read replicas and HA standbys are not backups because corruption/deletion can replicate.

### Object storage

Required design:

- versioning or bounded soft-delete against accidental overwrite/removal, if cost/legal review supports it;
- immutable inventory of object key/version/checksum/size/classification/tenant/lifecycle from PostgreSQL plus storage inventory;
- bucket policy, lifecycle, encryption and replication configuration in controlled infrastructure/configuration;
- protection from application credential deleting versions/lifecycle controls;
- sample and full-manifest reconciliation;
- cross-region/account copy only after data-location, cost and key recovery approval.

Object storage durability is not a complete application backup if metadata/keys/ownership/deletion state cannot be restored consistently.

### Secrets and keys

The recovery plan records:

- secret names/references, not plaintext exports;
- KMS key IDs/versions/policies and protected recovery administration;
- provider procedures for account/region/key disablement and restore;
- emergency rotation/revocation order;
- ability to decrypt retained database/object backups for their valid lifetime;
- rewrap/re-encrypt and destroy-old-key verification;
- break-glass credentials stored outside the failed primary path under dual control where practical.

A backup encrypted with an unavailable destroyed key is not recoverable; a broadly exported master key is not acceptable recovery.

### Code, infrastructure and contracts

Version control and release artefacts protect:

- application source/build manifests and dependency lockfiles;
- infrastructure/network/database configuration;
- migrations and schema snapshots;
- event/API/provider mapping versions;
- reference data and feature/config definitions that are not tenant secrets;
- runbooks, contact/escalation trees and vendor/account inventory.

Repository backup and access recovery are separated from the production cloud account to avoid one identity/provider failure removing both.

### External providers

Ubi cannot back up WhatsApp/Telegram/Meta/TikTok/payment/telecom service state as if it owns those systems. It preserves justified Ubi canonical facts, provider IDs/receipts/checkpoints and reconnection/reconciliation procedures. Provider export/replay capability and retention limits are documented per adapter.

## Backup security and governance

- Encrypt in transit and at rest; keys/access are environment separated.
- Backup operators cannot casually browse tenant content.
- Restore requires ticket/case, approved destination, purpose, scope, expiry and access audit.
- Production restoration into non-production is prohibited by default; exceptional forensic restore uses isolated environment, minimum people, no outbound integrations, masking/minimization and automatic destruction.
- Malware/ransomware/compromised-admin scenarios require immutable/separately controlled copies and MFA/dual approval for destructive backup-policy changes.
- Backup deletion follows retention/legal-hold policy and produces control evidence.
- Alert on disabled backup, retention reduction, mass deletion, key-policy change and restore activity.

## Schema migration principles

- All schema changes are version-controlled, reviewed, automated and reproducible from an empty database.
- Production is never the first environment to execute a migration.
- Mobile clients, queued commands, events and workers create a compatibility window; backend and schema cannot assume instant app upgrade.
- Prefer expand → migrate/backfill → verify → switch → contract.
- One migration has a clear owner, risk class, lock/runtime estimate, forward/rollback strategy and verification queries.
- DDL and application deployment order is explicit.
- Database migration role is JIT and not the application runtime.
- Destructive changes wait for retention, export, mobile/event/provider compatibility and recovery review.
- Schema drift from out-of-band SQL is detected and treated as an incident.

## Migration classes and controls

| Class | Example | Required controls |
|---|---|---|
| Additive/low lock | New nullable column/table/index | Compatibility review, representative timing and rollback/disable path |
| Constraint/index | FK/check/unique/GiST/RLS policy | Dirty-data precheck, concurrent/not-valid/validation plan, query/lock tests |
| Backfill/transformation | Populate tenant key or normalized contact | Bounded resumable batches, checkpoint, throttling, invalid-row quarantine, dual-read/write where needed |
| Semantic/state change | Split booking/payment status | Versioned domain/event/API migration, reconciliation and old-client behavior |
| Destructive/contract | Drop/rename/type narrow/remove enum | Evidence no readers/writers, retention/export approval, backup/restore proof, delayed cleanup |
| Provider/reference | New webhook/status/currency/timezone mapping | Adapter/schema version, unknown-state fail-safe, replay fixtures and reconciliation |
| RLS/privilege | Policy/grant/role change | Positive/negative isolation tests, plan/performance tests and clean-db grant diff |

## Expand–migrate–contract sequence

1. **Discover:** inventory readers/writers, mobile versions, jobs/events/providers, data volume/invalid rows, locks and legal implications.
2. **Expand:** add compatible schema/index/function/event version without removing old behavior.
3. **Dual capability:** deploy code that reads/writes safely across the transition; avoid indefinite dual-write ambiguity.
4. **Backfill:** process stable keys in bounded idempotent batches with checkpoint/rate/lag/error counters.
5. **Validate:** row counts, null/constraint candidates, hashes/sums, tenant integrity, owner/projection reconciliation and sample domain behavior.
6. **Switch:** use feature/config gate to move reads/owner behavior; monitor and retain fast disable/forward-fix path.
7. **Contract:** after every supported client/worker/event is beyond the boundary, validate constraints, remove old writes/data/indexes and update catalogue.
8. **Close:** record duration, result, deviations, metrics and recovery evidence.

DDL statements are checked against the selected PostgreSQL/provider version. Lock timeout/statement timeout are deliberate. A migration aborts safely rather than waiting indefinitely and freezing production.

## Data backfill and correction requirements

- Stable ordered cursor, not one huge transaction or deep offset.
- Idempotent result; restart from checkpoint after interruption.
- Per-batch transaction, bounded locks and production load throttling.
- Expected/processed/changed/skipped/invalid/error counters.
- Preserve original/provenance and transformation version where meaning changes.
- Never assign guessed tenant, currency, timezone, consent or payment state to make a constraint pass.
- Invalid/ambiguous rows quarantine into a controlled case with owner.
- Outbox/event emission strategy is explicit: suppress, synthesize versioned migration events or rebuild projection—never accidental notification floods.
- Backfill tooling has least privilege, dry-run/sample mode and auditable operator/release identity.

## Rollback versus forward recovery

Rollback is safe mainly for application code/config while expanded schema remains compatible. Data/state/provider effects are often not reversible.

- Never promise automatic rollback of a migration that sent messages, charged/refunded money, confirmed bookings or deleted data.
- Prefer disable/stop unsafe writers, preserve evidence, restore service with compatible prior code and issue an idempotent forward correction.
- PITR of the whole database is a disaster tool, not a routine rollback: it can erase valid concurrent work and desynchronize providers/objects/devices.
- Every high-risk release states the point after which rollback becomes forward-fix only.

## Restore scenarios

### Accidental row/table deletion or bad migration

1. Stop/limit unsafe writers and integrations.
2. Record incident time range, releases, migrations and external effects.
3. Select candidate recovery point before corruption.
4. Restore backup/PITR into isolated environment—never overwrite primary blindly.
5. Validate schema, constraints, counts, tenant distribution, audit/deletion/hold state and objects.
6. Decide targeted correction/import versus full cutover based on concurrent valid work and evidence.
7. Reconcile provider payments/messages, outbox/inbox, objects, sync and projections.
8. Resume gradually and monitor.

### Primary database/zone failure

Use managed HA failover, verify fencing prevents two primaries, validate application pools/DNS/read-write behavior and reconcile last transactions/idempotency/provider unknown outcomes.

### Regional/provider control-plane outage

Initial design uses a manual declared disaster process:

- confirm failure and expected recovery with provider;
- leadership/security/legal decide wait versus alternate-region recovery based on approved RTO/data location;
- provision isolated infrastructure from controlled code;
- restore latest approved cross-boundary backup plus keys/configuration;
- validate and fence old primary/queues/webhooks to avoid split brain/duplicate effects;
- update routing/DNS/provider callbacks carefully;
- reopen read-only/essential writes in stages;
- reconcile after old region returns.

Without approved cross-region backup and key availability, regional recovery cannot be claimed.

### Object deletion/corruption

Restore the exact version/checksum into quarantine, validate metadata/classification/owner and relink through media workflow. Purged-under-policy content is not recovered casually from backup.

### Credential/key compromise

Contain access, revoke/rotate affected service/provider credentials/keys, preserve audit evidence, assess encrypted data exposure, reissue scoped links/sessions as required and reconcile unauthorized provider effects. Restoring old credentials from backup is prohibited.

### Data integrity/tenant-isolation incident

Prioritize containment and evidence over rapid broad restore. Identify affected rows/objects/projections/exports/logs, stop leaks, preserve audit, notify/escalate under the approved legal runbook and perform reviewed tenant-by-tenant correction.

## Point-in-time consistency across systems

PostgreSQL, object storage, providers, Android devices and search do not share one atomic snapshot.

Recovery uses a manifest/cut line:

- database recovery point and transaction/outbox position;
- object metadata/version/checksum inventory;
- provider connection/checkpoint and outstanding unknown requests;
- search/report/sync projector checkpoints;
- deletion/revocation tombstone high-water mark;
- schema/app/adapter/reference-data versions.

After database restore:

1. keep external-effect workers/webhooks controlled;
2. replay deletion, revocation and legal-hold state first;
3. reconcile object existence/checksum;
4. rebuild disposable projections;
5. inspect outbox/inbox and deduplicate before provider work resumes;
6. reconcile payments and unknown messages/calls individually;
7. force client sync bootstrap where cursor history is incompatible;
8. validate canary tenants and critical workflows before general traffic.

## Tenant-specific restore

A shared-database PITR cannot rewind only one organization in place. For tenant-level accidental deletion:

- restore the whole database to an isolated temporary environment at the correct point;
- extract exact same-tenant graph through approved manifest/constraints;
- compare with current valid data and provider effects;
- create reviewed idempotent domain import/correction operations;
- preserve current/new legitimate work and audit every decision;
- reconcile downstream objects/projections/devices.

Do not copy raw rows around constraints/RLS into production. Tenant restore tooling is built only if incident likelihood/value justifies it; until then operations escalation and full isolated restore are the safe path.

## Android/offline recovery behavior

- Client command IDs/idempotency survive retry and server failover.
- After server restore, sync protocol detects incompatible cursor/epoch and requests a safe bootstrap.
- Pending commands created after recovered server point remain on devices and may replay; server applies current authority/state/idempotency and returns truthful conflict/rejection.
- Notifications/provider effects are not inferred from client local state.
- Emergency product communication explains pending/read-only status in simple localizable language.
- Business continuity provides a minimal export/print/manual booking log procedure with later controlled reconciliation; it must not create hidden double-booking.

## Recovery validation checklist

Before declaring recovered:

- schema/migration/version/extension/collation/timezone and roles/grants/RLS match approved release;
- expected database counts, constraints and checksums pass by tenant/domain;
- booking-reservation and payment-allocation invariants pass;
- audit, privacy cases, deletion tombstones and legal holds are present and applied;
- object manifest/checksum and required media samples pass;
- session/credential/restriction state is current enough and compromised credentials stay revoked;
- outbox/inbox/scheduled tasks are deduplicated and not about to flood providers;
- provider payments/messages/webhooks and unknown outcomes reconcile;
- search/cache/reporting/sync rebuild/checkpoints are correct;
- synthetic canary tenants prove cross-tenant isolation and critical booking workflows;
- observability/alerts/backups resume in the recovered environment;
- incident commander, security/data/domain owners sign staged reopen.

## Exercise cadence

Proposed minimum, subject to risk/provider/launch approval:

- automated daily verification of backup success/age/WAL continuity and alerts;
- scheduled isolated PostgreSQL restore at least monthly during build/pilot until reliable, then at least quarterly plus after material architecture/provider change;
- object manifest/sample restore with database restore exercise;
- table/tenant accidental-deletion exercise and migration recovery rehearsal before launch;
- regional recovery tabletop before launch and technical exercise when alternate-region capability is claimed;
- annual or major-change credential/key-compromise and business-continuity exercise;
- every exercise records achieved RPO/RTO, gaps, manual dependencies, cost and owned remediation.

An untested written runbook does not justify an availability claim.

## Monitoring and triggers

Alert on:

- PITR/WAL/snapshot failure or recovery window shrink;
- replication/failover lag and storage saturation;
- object versioning/lifecycle/bucket-policy/replication changes;
- KMS key disable/deletion/policy anomaly;
- unapproved schema drift or migration checksum mismatch;
- long migration locks, backfill error/lag and invalid-row growth;
- restore/export activity;
- outbox/inbox/projection/deletion replay lag after recovery;
- recovery exercise misses its objective.

## Roles and decision authority

- Incident commander controls recovery sequence and communication.
- Database/reliability owner performs approved restore/failover.
- Domain owners validate booking/payment/identity semantics.
- Security/privacy/legal assess compromise, data location, notification and deletion/hold duties.
- Provider/integration owners fence/reconcile external effects.
- Business operations decides continuity/manual workflow with technical capacity safeguards.
- No single support person can both authorize and invisibly perform a high-risk production restore/export.

## Release gates

- Exact RPO/RTO and outage behavior are approved per first-release capability from measured provider tests and business impact.
- Managed database HA/PITR/snapshot/retention/region/cost/key behavior is verified, not assumed from marketing.
- Clean database can be created from migrations and matches schema/grant/RLS snapshot.
- Expand/backfill/contract and old Android client compatibility are rehearsed on representative volume.
- Isolated PostgreSQL restore plus object/key/config recovery meets objective and passes validation checklist.
- Deletion/revocation/legal-hold replay executes before restored traffic.
- Unknown provider payment/message outcomes reconcile without duplicate effects.
- Backup/restore access is JIT/audited and destructive policy changes are protected.
- Merchant business-continuity and later reconciliation process is tested with pilot staff.
- Recovery runbooks 66–70 name people, contacts, commands, evidence and communications.
- Customer-facing availability/data-loss language matches tested capability.

## Open decisions

- Exact tier RPO/RTO/continuity promises and affordable loss/downtime tolerance.
- Managed PostgreSQL/object/KMS provider, version, region and separate-account backup options.
- Backup/PITR/object-version retention after privacy/legal/cost review.
- Whether/when cross-region copy is lawful, affordable and required.
- Restore automation/tooling and tenant-extraction/import investment.
- Mobile compatibility duration/minimum supported version and forced-upgrade emergency behavior.
- Recovery approval/dual-control staffing at pilot size.
- Portable logical export frequency and encrypted off-provider storage need.
- Objective evidence required before claiming disaster recovery beyond same-region HA.
