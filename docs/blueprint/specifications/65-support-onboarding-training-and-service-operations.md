# Support, onboarding, training and service operations

Status: **Draft — Ethiopia-first pilot operating model**  
Depends on: actor/permission/Android/operations architecture, trust/privacy/legal policy, pilot and commercial decisions  
Blocks: merchant pilot, support-console access, service commitments and scalable onboarding

## Purpose

Ubi will fail if technically correct software still requires a highly skilled operator. This document defines assisted merchant onboarding, low-literacy training, support intake, escalation, service ownership and operational evidence for an Ethiopia-first pilot.

Diagram: `diagrams/65-merchant-onboarding-and-support-flow.mmd`.

## Operating principles

- Design the product to be understandable; training supplements rather than excuses confusing UX.
- Meet merchants in their current workflow: calls, walk-ins, notebooks, messaging and cash/payment evidence.
- One accountable case owner communicates; internal handoffs are invisible work, not a merchant burden.
- Support starts from redacted metadata and asks permission/step-up for scoped data/actions.
- Never request passwords, OTPs, full payment credentials, social-media sessions or copied private data in chat.
- Amharic/local-language/plain speech, screenshots and short demonstrations are preferred to technical manuals.
- Every repeated support problem feeds product/content/reliability improvement.

## Pilot service boundary

Proposed pilot:

- Ethiopia, preferably Addis Ababa after founder approval;
- approximately 10 photo studios/photographers for 8–12 weeks;
- founder/operations-led onboarding and named merchant contact;
- Ubi Business Android plus necessary customer link/manual booking workflow;
- official channel/payment capability only if its gate is passed; manual/link/cash fallback always available;
- limited published support hours with serious-incident escalation, not an unstaffed 24/7 promise.

The pilot is a learning service with explicit limits, not a disguised general public launch.

## Merchant qualification

Before onboarding, verify enough to operate safely without burdensome universal KYC:

- business/owner representative and reachable contact;
- actual service/location/workflow and authority to configure/receive bookings;
- supported first-vertical fit and excluded/restricted services;
- current booking sources, schedule/resources/staff, price/deposit/cancellation behavior;
- Android device/OS/storage/connectivity and shared-device reality;
- desired customer language/calendar/payment/channel behavior;
- willingness to run a bounded pilot, report issues and maintain schedule truth.

Enhanced evidence applies to risky claims/payment/channel ownership; Ubi verification is not a service-quality guarantee.

## Assisted onboarding journey

1. **Observe first:** watch representative call/walk-in/message/notebook booking without forcing Ubi terminology.
2. **Create organization:** verified owner/representative, operation/location and recovery contacts.
3. **Select template:** photo studio/photographer starter archetype, not blank universal configuration.
4. **Configure minimum:** services, duration/price method, resources, working hours, confirmation and cancellation/deposit choices.
5. **Import near-term truth:** upcoming bookings manually or through reviewed template; retain reconciliation checklist.
6. **Add staff:** individual memberships/roles; no shared owner password.
7. **Practice safely:** demo customer/manual/offline/conflict/cancel/payment-status scenarios in sandbox/demo mode.
8. **Publish carefully:** preview customer link/subdomain, policy and contact; owner explicitly activates.
9. **Go-live shadow:** compare notebook/current process and Ubi for a short agreed period without silent double truth.
10. **Adoption review:** correct configuration/training/product gaps and decide continue/pause/exit.

## Minimum onboarding record

- merchant/owner/operation/service/template and verification level;
- devices/connectivity/language/calendar and accessibility needs;
- staff roles and training completion;
- migration counts/reconciliation and unresolved items;
- configured policies/payment/channel status and exact limitations;
- acceptance/notice versions and pilot agreement;
- go-live owner/date/support route and success/exit criteria.

Do not duplicate customer content into onboarding notes unnecessarily.

## Training design

Train by tasks, each with “show → merchant does → recover from mistake”:

1. See Today and identify offline/as-of/sync state.
2. Add a phone/walk-in booking quickly.
3. Respond to request/quote and confirm only when availability is authoritative.
4. Edit/reschedule/cancel with customer-visible consequences.
5. Record cash/customer-reported payment versus provider-verified payment.
6. Resolve conflict/rejected/pending item.
7. Find conversation/customer/booking without exposing unrelated data.
8. Add/remove staff and protect owner account.
9. Use support/report/block/emergency paths.
10. Continue safely during Ubi/provider/network outage.

Materials:

- one-page visual quick start;
- short captioned/local-language videos/audio where useful;
- in-app contextual help and safe examples;
- printable outage/manual-booking sheet;
- owner/admin and staff-specific guides;
- searchable support articles with version/date.

Measure unassisted completion and error recovery—not whether the trainer spoke.

## Demo and sandbox accounts

- Visibly marked `DEMO`; cannot contact real customers/providers or produce real payment truth.
- Synthetic Ethiopian-localized data only.
- Resettable scenarios for offline, conflict, cancellation, reported/verified payment and support.
- Demo credentials are not production credentials and cannot be promoted.
- Training completion does not grant production roles automatically.

## Support intake channels

Initial channels may include in-app case, phone and an approved business messaging route. All create/attach to one case.

Capture:

- merchant/account/operation verified through safe method;
- affected task/booking reference, time and visible state;
- scope/number of affected users and workaround;
- app/version/device/network/provider;
- customer data only when necessary and through approved secure upload;
- consent/authority for scoped support action where required.

Social DMs are not used to exchange credentials or full personal/payment evidence.

## Case classification and targets

| Priority | Example | Initial operating target* |
|---|---|---|
| P0 | Safety/security/privacy/tenant leak/false payment or widespread booking integrity | Immediate incident escalation |
| P1 | Merchant cannot access/operate critical booking; no safe workaround | Rapid acknowledgment during covered hours and on-call escalation by policy |
| P2 | Degraded feature/provider issue with workaround; limited data/config error | Same business day target during pilot |
| P3 | How-to, cosmetic, request, planned configuration | Queued with clear expectation |

`*` Internal pilot target only. Exact hours/response/resolution/SLA are approved after staffing and demand evidence. Resolution time is not promised when provider/merchant/legal investigation controls it.

## Support workflow

```text
intake and identity/authority check
→ classify impact/risk/data
→ give safe workaround and case reference
→ inspect redacted telemetry/configuration
→ request scoped/JIT access or evidence only if necessary
→ resolve, escalate to runbook/engineering/provider or explain limitation
→ verify with merchant and reconcile affected records
→ revoke access, record outcome and identify product/content/root-cause work
```

No generic “log in as merchant,” raw database editor or permanent universal support role.

## Escalation routes

- product/configuration → product operations;
- account recovery/ownership → trained identity/security approval;
- booking/sync/data integrity → backend/data/reliability;
- channel/OAuth/provider → integration runbook 68;
- payment evidence/reconciliation → finance/payment runbook 69;
- privacy/security/safety → runbook 67/trust/legal;
- outage/system incident → runbook 66;
- backup/restore/continuity → runbook 70;
- legal/authority request → policy 57.

## Merchant communications

- Use plain language: what is affected, what is safe, what not to repeat, workaround, next update.
- Never say “paid,” “confirmed,” “sent” or “restored” beyond evidence.
- Provider failures are identified as provider-specific only when known.
- Status updates have time/cadence and case/status reference.
- Post-resolution message states reconciliation/user action and prevention work.

## Offboarding and exit

1. Verify authorized owner and effective date.
2. Explain customer-booking continuity, links, channels, payments, exports, retention and legal holds.
3. Export approved merchant data in documented format.
4. Disable new bookings/messages safely; preserve/manage existing obligations.
5. Disconnect provider tokens/webhooks and revoke staff/support access.
6. Process deletion/retention across stores/vendors under approved schedule.
7. Reserve/redirect slug safely and prevent impersonation.
8. Record completion/evidence and reasons for product learning.

## Operational metrics

- onboarding lead time and trainer/support minutes;
- first-value time: first accurate manual/customer booking;
- unassisted critical-task completion and repeated help rate;
- activation, schedule completeness and old-process duplication;
- cases per active merchant/booking by cause/priority/reopen/escalation;
- response, workaround and resolution distributions by covered hours;
- JIT/support-access use and policy violations;
- merchant/customer-impact minutes and provider dependency;
- retention/adoption and exit reason;
- product defects/content gaps prevented through feedback.

Do not optimize by closing cases prematurely or making merchants repeat contacts.

## Staffing and RACI direction

- Pilot operations owner: onboarding calendar, case ownership and merchant relationship.
- Product owner: configuration/templates/content and adoption decisions.
- Engineering on-call: technical incident/data/integration escalation.
- Security/privacy/trust owner: sensitive access, reports and incident classification.
- Finance/payment owner: reconciliation/refunds/provider financial cases.
- Founder/CEO: high-impact merchant communication, risk/continuation and commercial exceptions.

One person may hold multiple pilot roles, but must know which authority they are exercising; high-risk approvals retain separation where possible.

## Release/pilot gates

- Onboarding script/template/data-import/reconciliation and exit flow are rehearsed with synthetic data.
- Representative merchants complete critical tasks and recover errors with acceptable assistance.
- Support intake creates one auditable case and verifies authority proportionate to action.
- JIT/redaction/audit/revoke and break-glass controls work before support sees production data.
- P0/P1 escalation exercises reach named people and invoke correct runbooks.
- Training/demo cannot message real users, move money or grant production authority.
- Support hours/languages/limits and provider dependencies are stated truthfully.
- Pilot staffing/cost capacity supports the proposed merchant count.

## Open decisions needed from founder

- Confirm Addis Ababa, photo studios/photographers and 10-business 8–12 week pilot.
- Sole founder/cofounder ownership and legal entity formation timing.
- Initial support hours, languages and field/on-site availability.
- Free/paid pilot, onboarding fee and equipment/data assistance budget.
- Merchant eligibility/exclusions and exit/success thresholds.

