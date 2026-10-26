# Runbook 66 — General incident response

Status: **Draft — must be exercised before pilot**  
Owner: Reliability/operations lead  
Applies to: production/pilot service degradation, correctness failure and operational disruption not superseded by a more specific runbook

## Purpose

Restore safe service, protect people/data/money, communicate truthfully and learn without destroying evidence. Speed matters, but a fast unsafe recovery that duplicates bookings, leaks tenant data or invents payment status is failure.

Diagram: `diagrams/66-incident-lifecycle.mmd`.

## Invoke this runbook when

- critical task/SLO/error-budget alert indicates material user impact;
- accepted command may be lost/stuck/duplicated;
- booking/availability integrity is uncertain;
- widespread login/sync/web/API/database/worker failure occurs;
- support/provider/monitoring reports credible production impact;
- deployment/configuration/migration caused unexpected behavior;
- another runbook declares a coordinating incident.

Security/privacy, integration, payment and recovery specialists use runbooks 67–70 inside this common command structure.

## Severity

| Severity | Examples | Initial action |
|---|---|---|
| SEV-0 Critical | Cross-tenant exposure; false verified payment; Ubi-induced confirmed double booking at scale; destructive compromise; imminent safety harm | Immediate incident, executive/security/privacy/legal as relevant, stop unsafe actions |
| SEV-1 High | Booking core unavailable; accepted work at risk; widespread auth/sync/database/provider impact without safe workaround | Immediate on-call incident and merchant/status communication cadence |
| SEV-2 Medium | Degraded feature/limited cohort with safe workaround; growing backlog within recovery window | Named owner, active case, escalation if thresholds worsen |
| SEV-3 Low | Minor defect/no material task impact; routine operational issue | Ticket and normal priority |

Severity follows credible impact and uncertainty, not customer count alone. A single tenant privacy breach may be SEV-0.

## Roles

- **Incident commander (IC):** owns severity, priorities, cadence and decisions; does not perform every technical task.
- **Technical lead:** diagnosis, containment/recovery changes and verification.
- **Operations/support lead:** affected merchants, workarounds, case correlation and frontline scripts.
- **Communications lead:** internal/status/merchant updates approved by IC.
- **Scribe/evidence custodian:** timeline, hypotheses, actions, approvals, artifacts and decisions.
- **Security/privacy/legal/finance/provider liaison:** activated by scope.
- **Executive decision owner:** accepts high residual business risk or prolonged service tradeoff.

One person may cover roles at pilot scale, but IC and high-risk change approval should separate where feasible.

## First 15 minutes

1. Acknowledge alert/report; open incident ID, secure channel and timeline.
2. Name IC/technical lead/scribe; set provisional severity and next update time.
3. State observed facts, affected tasks/cohorts/time and uncertainty—do not guess cause.
4. Protect safety/integrity: disable risky feature/writes/provider route or reduce traffic if necessary.
5. Preserve logs/audit/deployment/config/provider evidence before rotation or cleanup.
6. Check recent release/migration/flag/provider/security changes and telemetry health.
7. Give support a safe workaround and prohibited actions.
8. Invoke 67/68/69/70 when their trigger is met.

## Stabilization priorities

1. Human safety and cross-tenant/privacy/security containment.
2. Prevent false payment/booking/authorization and further irreversible writes.
3. Preserve accepted commands/evidence and make unknown states visible.
4. Restore safe read/manual/offline/customer communication capability.
5. Restore normal writes/providers and drain backlog fairly.
6. Reconcile every uncertain external/internal outcome.
7. Restore secondary notification/search/analytics/reporting.

## Investigation method

Maintain a live table:

```text
time | fact/evidence | hypothesis | test | action | result | owner | rollback
```

- Correlate task/command/transaction/outbox/provider/release IDs.
- Compare healthy versus affected tenant/device/network/provider/version.
- Validate telemetry pipeline before trusting absence of errors.
- Bound queries/exports and use JIT case access; incident does not grant unrestricted browsing.
- Change one high-value variable at a time when possible.
- Mark facts, hypotheses and decisions distinctly.

## Containment options

- feature/capability/tenant/provider kill switch;
- stop deployment/backfill/worker pool or risky write route;
- force read-only/manual mode with truthful banner;
- reduce concurrency/retries and apply backpressure;
- revoke/rotate compromised credentials through runbook 67;
- isolate bad release/configuration/traffic cohort;
- disable outbound sends while preserving durable queue;
- extend no hold/authorization/payment validity merely to hide outage.

Every containment records scope, owner, expected effect, data impact, expiry and reversal check.

## Recovery decision

Choose using evidence:

- continue current version with safe mitigation;
- disable optional feature;
- compatible application/web rollback;
- forward fix/data correction;
- provider adapter pause/fallback;
- database failover/restore through runbook 70.

Never perform whole-database PITR as routine rollback; it can erase valid concurrent bookings and desynchronize devices/providers.

## Verification before reopening

- exact artifact/config/schema/flag/credential active;
- synthetic safe critical flows pass;
- authoritative booking/permission/payment constraints and audit intact;
- accepted/stuck/unknown commands inventoried and reconciliation owned;
- queue/database/provider/storage/telemetry healthy with headroom;
- representative merchant/device/client confirms user-visible recovery where appropriate;
- support/status message accurately states residual effects/actions;
- containment reversal cannot reintroduce cause.

Restore gradually and watch fast/slow burn plus domain integrity.

## Communication template

```text
Incident ID and status: Investigating / Identified / Recovering / Monitoring / Resolved
Started/observed time and affected capability/cohort
What users may see
What remains safe / what not to repeat
Safe workaround or no action required
Next update time
Known reconciliation/follow-up after recovery
```

Do not publish speculation, personal data, exploit detail, false provider blame or unsupported “no data lost.”

## Handoff and prolonged incident

- Incoming IC/lead receives live summary, timeline, current risk, access, actions, hypotheses and next decision.
- Fatigue triggers scheduled rotation; credentials are not shared.
- Maintain merchant/status cadence even without new cause.
- Executive decision records any extended read-only/outage/cost/risk tradeoff.

## Resolution and closure

“Resolved” requires safe service plus owned reconciliation, not merely alert recovery.

Record:

- final impact/time/cohorts/data/financial/provider effects;
- root/contributing/control-detection factors or explicitly unknown;
- all changes/credentials/flags/manual corrections;
- backlog/unknown/reconciliation and merchant follow-up;
- notification/legal/support obligations;
- evidence retention/access and temporary-resource cleanup.

## Learning review

Within a risk-appropriate period:

- blameless factual timeline and five contributing-system questions;
- why prevention/detection/response/recovery controls behaved as they did;
- customer/merchant and responder burden;
- tests/alerts/runbooks/architecture/training/vendor changes;
- action owner, priority, due date and verification evidence;
- risk register/SLO/error-budget/update.

Close actions only after verification, not ticket creation.

## Exercise and readiness gates

- Contacts/roles/escalation and access work without one unavailable person.
- Alert-to-declaration, feature disable, deployment rollback and status update are rehearsed.
- At least one accepted-command/unknown-provider/database scenario is simulated.
- Incident channel/timeline/evidence templates and secure storage exist.
- Support knows safe workarounds and credential/data prohibitions.
- Findings have owners; failed exercise blocks pilot for critical gaps.

