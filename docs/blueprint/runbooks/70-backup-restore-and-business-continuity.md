# Runbook 70 — Backup, restore and business continuity

Status: **Draft — provider-specific commands and approved RPO/RTO pending**  
Owners: Reliability/data recovery lead and business continuity lead  
Applies to: database/object/config/key loss or corruption, infrastructure/region/account outage, destructive change, ransomware and prolonged service/provider disruption

## Core rule

Recovery means restoring a safe, usable and reconciled booking service—not only starting a database. A restore must include schema/constraints, tenant isolation, objects, keys/configuration, idempotency/outbox/inbox/timers, audit/privacy deletion/hold state, devices and external-provider outcomes.

Diagram: `diagrams/70-backup-restore-continuity.mmd`.

## Proposed objectives

- Authoritative R0 data engineering target: RPO no worse than 5 minutes.
- Booking core target: RTO no worse than 4 hours.
- Supporting administration target: RTO no worse than 8 hours.

These are not SLA promises until chosen provider/topology, cost and exercises prove them. Security/privacy containment may deliberately extend recovery to prevent further harm.

## Recovery tiers

| Tier | Capability/data | Priority |
|---|---|---|
| R0 | Identity/session authority, organization restrictions, booking/capacity, idempotency/outbox, payment obligations/evidence | First |
| R1 | Catalog/config/upcoming CRM, communications state, audit/privacy deletion/restriction/hold | With/immediately after R0 |
| R2 | Objects/media/exports needed for active obligations | After metadata/key validation |
| R3 | Search/cache/reporting/analytics | Rebuild after authoritative service |
| External | Channel/payment/notification state | Reconcile; never assumed restored from Ubi backup |

## Invoke this runbook for

- primary database/zone/region/provider control-plane failure;
- accidental/bad migration/deletion/corruption affecting authoritative data;
- object/key/configuration/secret loss;
- compromised cloud/admin/backup control or ransomware;
- restore requested for investigation/tenant correction;
- prolonged outage requiring merchant manual-continuity mode;
- scheduled recovery exercise.

Security compromise invokes runbook 67 before connecting restored assets to production.

## Roles and authority

- Recovery commander chooses objective, recovery point/environment and traffic state.
- Database/reliability operator executes approved failover/restore.
- Security validates trust/control plane/credentials/backup integrity.
- Data/domain owners validate consistency, constraints, deletion/hold and reconciliation.
- Integration/payment owners reconcile provider outcomes.
- Operations/support coordinates merchant continuity and re-entry.
- Executive/legal/privacy approves regional/data-location/prolonged-outage tradeoffs.

No single support person authorizes and invisibly performs high-risk restore/export.

## First response

1. Declare incident, stop unsafe writes/deployments/jobs and preserve evidence.
2. Determine failure scope: service, data, control plane, credentials, region/account/provider.
3. Protect backup copies, key/config repositories and clean identities from suspected compromise.
4. Establish last known good transaction/time and accepted external/device work after it.
5. Decide HA failover versus isolated restore versus wait/alternate environment.
6. Announce merchant continuity mode and safe manual capture instructions.
7. Select recovery point based on data integrity/evidence, not simply newest timestamp.
8. Record approval, expected RPO/RTO and next update.

## Merchant business-continuity mode

When Ubi writes are unsafe/unavailable:

- merchants use approved printable/local manual booking sheet with stable temporary reference;
- capture minimum customer/contact/service/date/time/resource/deposit status and who recorded it;
- never label customer-reported payment as verified;
- check local cached schedule/as-of time but state it may be stale;
- do not issue overlapping confirmation without merchant's manual conflict check;
- direct customer/channel automation to truthful delayed/manual-contact message where possible;
- protect paper/device records and enter them once through controlled reconciliation, not bulk guess.

The pilot trains this mode before launch.

## Database failover versus restore

### HA failover

Use when managed standby is healthy and incident is instance/zone failure without logical corruption/compromise.

- verify provider failover state/writer endpoint and fencing of old writer;
- reconnect with bounded retry; never replay ambiguous command without idempotency lookup;
- validate transaction/replication point, constraints and queue processing;
- reconcile accepted commands during transition.

### PITR/backup restore

Use for logical corruption/deletion or unavailable primary/control plane after approval.

1. Create isolated recovery environment/account/network with outbound providers disabled.
2. Restore chosen PostgreSQL point/snapshot plus roles/extensions/RLS/functions.
3. Restore configuration/key access necessary to decrypt—never old compromised credentials.
4. Verify migration/schema/application compatibility.
5. Validate counts/checksums/constraints/tenant-isolation/idempotency/outbox/inbox/timers/audit/privacy state.
6. Restore object metadata/content versions by manifest and key; quarantine unknown/corrupt objects.
7. Replay approved deletion/revocation/restriction/legal-hold changes after recovery point before exposure.
8. Build provider/device reconciliation inventory after recovery point.
9. Approve traffic cutover, DNS/routing/fencing and phased service restoration.

Never overwrite the existing primary blindly; preserve it for evidence/comparison.

## Tenant-specific correction

Shared-database PITR cannot rewind one tenant in place. For scoped accidental deletion:

- restore whole database at target time into isolated environment;
- identify tenant records and dependency graph through approved export tool;
- compare current versus historical owner state and concurrent valid changes;
- create reviewed domain commands/correction migration preserving current unrelated work;
- verify constraints/RLS/audit and merchant result.

Do not copy raw rows around constraints/RLS directly into production.

## Objects, keys and configuration

- Database metadata and object bytes/version/checksum/key must agree.
- Purged content is not casually recovered; privacy/legal authority controls.
- Key recovery uses protected escrow/provider process; missing key may make backup unrecoverable and must be reported honestly.
- Credential compromise triggers rotation/reauthorization; restoring old secret is prohibited.
- IaC/config/domain/cert/provider app/webhook inventories restore the service control plane.
- Search/cache/analytics rebuild from authoritative sources after access controls/deletion filters are active.

## External provider reconciliation

For the outage/recovery interval, inventory:

- payment initiations/callbacks/refunds/reversals/unknowns via runbook 69;
- channel ingress/egress/delivery/unknowns via runbook 68;
- notifications accepted/delivered/expired/duplicate risk;
- OAuth/merchant credential actions and provider asset changes;
- files/uploads whose metadata/object commit diverged.

Provider state is not rewound by Ubi PITR. Reconcile before replaying or telling users the outcome.

## Android/offline client recovery

- Server exposes recovery epoch/minimum protocol and safe cursor behavior.
- Invalid/expired cursor triggers staged rebootstrap; last valid local data remains with `as of` status where safe.
- Pending commands keep stable IDs and query idempotency results before resend.
- Commands based on restored-away/changed versions may conflict and require user resolution.
- Revoked membership/session/links remain revoked through replayed security state.
- Client never overwrites local pending work blindly with restored server snapshot.

## Cutover validation

Before write traffic:

- trusted control plane/identity/keys and exact release/schema/config active;
- database consistency, constraints, RLS/negative tenant tests and clock/timezone correct;
- booking capacity/idempotency/money/audit/privacy deletion/hold checks pass;
- outbox/inbox/timers paused or safely positioned; no duplicate send storm;
- object references/sample downloads and key decrypt pass;
- backups/PITR/observability/alerts enabled in recovered environment;
- external sends/payments remain gated until reconciliation owner approves;
- capacity/headroom and rollback/fencing route ready.

Restore reads first where useful, then bounded writes/cohorts; monitor domain/SLO/security/cost.

## Re-entry of manual records

1. Freeze/collect each merchant's outage sheet/device records and verify recorder/time.
2. Sort by resource/date/time and check conflicts against recovered current schedule.
3. Enter with source `business_continuity_manual`, original observed time and operator.
4. Contact affected customers for ambiguous overlap/payment/confirmation.
5. Reconcile temporary to canonical reference and retain appropriate evidence.
6. Obtain merchant sign-off on unresolved/corrected items.

No bulk “mark all confirmed/paid.”

## Communication

State:

- affected capabilities and time;
- whether cached/manual booking is safe and its limitations;
- what merchants/customers must not repeat;
- recovery stage and next update;
- after restoration, recovery point, possible missing/duplicate/unknown work and reconciliation plan.

Never say “no data lost” until transaction/device/provider/object evidence supports it.

## Scheduled exercise

At minimum:

- daily automated backup age/WAL/snapshot/alert verification;
- isolated PostgreSQL restore monthly during build/pilot until reliable, then at least quarterly and after material change;
- object manifest/key/configuration plus application validation with restore;
- annual or pre-launch region/control-plane/business-continuity tabletop/technical exercise proportional to design;
- periodic tenant-specific correction and Android rebootstrap test.

Each exercise records actual RPO/RTO, people/manual steps, access failures, cost, reconciliation gaps and corrective owners.

## Closure/readiness gates

- Backups are monitored, separately protected and restorable with available keys/configuration.
- Restore meets approved objectives in an isolated exercise and passes tenant/domain/privacy checks.
- HA failover, logical corruption, compromised credential and provider-reconciliation paths are rehearsed.
- Merchant manual-continuity sheet/training/re-entry process works with representative users.
- Android clients safely rebootstrap/preserve pending commands after recovery epoch.
- Region/account recovery claims match actual approved data location/provider/cost design.
- Every recovered unknown outcome has owner; temporary environments/data/access are destroyed or retained under approved hold.

